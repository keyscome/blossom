# 使用官方的Ubuntu镜像作为基础镜像
FROM ubuntu:latest

# 安装软件包，为安装Ansible做准备
RUN apt-get update && apt-get install -y \
    software-properties-common

# 添加Ansible的APT仓库并安装Ansible
RUN add-apt-repository --yes --update ppa:ansible/ansible \
    && apt-get install -y ansible

# 设置工作目录
WORKDIR /ansible

# 默认运行ansible --version检查版本
# ENTRYPOINT [ "ansible" ] # 执行playbook等操作时，使用ansible-playbook等命令，因此不固定ENTRYPOINT
CMD ["ansible" "--version"]
