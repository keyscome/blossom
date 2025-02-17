package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config 定义 YAML 配置结构
type Config struct {
	Monitor struct {
		Disk []string `yaml:"disk"`
	} `yaml:"monitor"`
}

var cfg Config

// loadConfig 从配置文件中加载监控配置
func loadConfig() error {
	configFile := viper.GetString("config")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &cfg)
}

// monitorData 保存监控信息
type monitorData struct {
	CPUUsage  float64
	MemUsage  float64
	DiskStats []diskStat
	SwapUsage float64
	Load1     float64
	Load5     float64
	Load15    float64
	NetRecv   uint64
	NetSent   uint64
	OpenFDs   int32
	Uptime    string
	Procs     int
}

// diskStat 保存单个挂载点的磁盘使用情况
type diskStat struct {
	Mount       string
	UsedPercent float64
}

// getMonitorData 采集各项系统监控数据
func getMonitorData() monitorData {
	// CPU 使用率（整体）
	cpuPercents, _ := cpu.Percent(time.Second/5, false)
	// 内存
	vm, _ := mem.VirtualMemory()
	// 磁盘：根据配置中的挂载点采集数据
	var dStats []diskStat
	for _, mount := range cfg.Monitor.Disk {
		usage, err := disk.Usage(mount)
		if err != nil {
			dStats = append(dStats, diskStat{Mount: mount, UsedPercent: 0})
		} else {
			dStats = append(dStats, diskStat{Mount: mount, UsedPercent: usage.UsedPercent})
		}
	}
	// 交换区
	swap, _ := mem.SwapMemory()
	// 系统负载
	avg, _ := load.Avg()
	// 网络 I/O（取第一个网卡）
	netIO, _ := net.IOCounters(false)
	var netRecv, netSent uint64
	if len(netIO) > 0 {
		netRecv = netIO[0].BytesRecv
		netSent = netIO[0].BytesSent
	}
	// 当前进程打开的 FD 数
	proc, _ := process.NewProcess(int32(os.Getpid()))
	fds, _ := proc.NumFDs()
	// 系统运行时长
	uptimeSec, _ := host.Uptime()
	uptime := fmt.Sprintf("%.2f hrs", float64(uptimeSec)/3600)
	// 系统进程数
	pids, _ := process.Pids()

	return monitorData{
		CPUUsage:  cpuPercents[0],
		MemUsage:  vm.UsedPercent,
		DiskStats: dStats,
		SwapUsage: swap.UsedPercent,
		Load1:     avg.Load1,
		Load5:     avg.Load5,
		Load15:    avg.Load15,
		NetRecv:   netRecv,
		NetSent:   netSent,
		OpenFDs:   fds,
		Uptime:    uptime,
		Procs:     len(pids),
	}
}

// 消息类型定义
type tickMsg time.Time
type cmdOutputMsg string

// model 定义 TUI 状态
type model struct {
	cmdOutput string           // 命令输出内容
	monitor   monitorData      // 监控数据
	cpuProg   progress.Model   // CPU 进度条
	memProg   progress.Model   // 内存进度条
	diskProgs []progress.Model // 每个挂载点对应一个磁盘进度条
}

// Init 启动时返回批量命令：定时刷新和执行命令
func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), runCmd())
}

// tick 每秒发送 tickMsg
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// runCmd 执行命令 "ls ~"（模拟 watch 命令）
func runCmd() tea.Cmd {
	return func() tea.Msg {
		home := os.Getenv("HOME")
		out, err := exec.Command("ls", home).CombinedOutput()
		if err != nil {
			return cmdOutputMsg(fmt.Sprintf("Command error: %v", err))
		}
		return cmdOutputMsg(string(out))
	}
}

// Update 根据消息更新模型状态
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// 更新监控数据
		m.monitor = getMonitorData()
		_ = m.cpuProg.SetPercent(m.monitor.CPUUsage / 100)
		_ = m.memProg.SetPercent(m.monitor.MemUsage / 100)
		for i, ds := range m.monitor.DiskStats {
			if i < len(m.diskProgs) {
				_ = m.diskProgs[i].SetPercent(ds.UsedPercent / 100)
			}
		}
		// 每秒刷新，同时重新执行命令
		return m, tea.Batch(tick(), runCmd())
	case cmdOutputMsg:
		m.cmdOutput = string(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View 渲染 TUI 界面
func (m model) View() string {
	// 标题
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Render("🌸 Blossom DevOps Tool")
	// 命令输出面板
	cmdTitle := lipgloss.NewStyle().Bold(true).Underline(true).Render("Command Output")
	cmdPanel := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(
		fmt.Sprintf("%s\n%s", cmdTitle, m.cmdOutput),
	)
	// 监控信息面板：CPU、内存、磁盘（每个磁盘一个进度条）
	labelStyle := lipgloss.NewStyle().Bold(true).Width(6)
	cpuLine := fmt.Sprintf("%s %s %.2f%%", labelStyle.Render("CPU:"), m.cpuProg.View(), m.monitor.CPUUsage)
	memLine := fmt.Sprintf("%s %s %.2f%%", labelStyle.Render("MEM:"), m.memProg.View(), m.monitor.MemUsage)
	diskLines := ""
	for i, ds := range m.monitor.DiskStats {
		if i < len(m.diskProgs) {
			diskLines += fmt.Sprintf("DISK(%s): %s %.2f%%\n", ds.Mount, m.diskProgs[i].View(), ds.UsedPercent)
		}
	}
	monInfo := fmt.Sprintf("%s\n%s\n%s", cpuLine, memLine, diskLines)
	monTitle := lipgloss.NewStyle().Bold(true).Underline(true).Render("Monitoring Info")
	monPanel := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(
		fmt.Sprintf("%s\n%s", monTitle, monInfo),
	)
	// 左右并排显示
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Top, cmdPanel, monPanel)
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press 'q' to exit.")
	return lipgloss.JoinVertical(lipgloss.Center, title, "", mainPanel, "", footer)
}

// monitorCmd 定义启动监控 TUI 的 Cobra 子命令
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "启动监控界面",
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置文件
		err := loadConfig()
		if err != nil {
			fmt.Println("加载配置文件失败，使用默认配置：", err)
			cfg.Monitor.Disk = []string{"/"}
		}
		// 为配置中每个磁盘挂载点创建一个进度条
		var diskProgs []progress.Model
		for range cfg.Monitor.Disk {
			diskProgs = append(diskProgs, progress.New(progress.WithDefaultGradient()))
		}
		// 初始化 TUI 模型
		m := model{
			cpuProg:   progress.New(progress.WithDefaultGradient()),
			memProg:   progress.New(progress.WithDefaultGradient()),
			diskProgs: diskProgs,
		}
		// 启动 TUI 程序
		p := tea.NewProgram(m)
		if err := p.Start(); err != nil {
			fmt.Println("TUI 运行错误:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
