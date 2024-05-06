package configs

import (
	// "io/ioutil"
	// "yaml"
	// "log"
)

type ConfigBlossom struct {
	ConfigCM
	ConfigRT
	ConfigWF
	ConfigTest
	ConfigKanban
	ConfigMeasure
}

type ConfigCM struct {

}

type ConfigRT struct {
}

type ConfigWF struct {}
type ConfigTest struct {}
type ConfigKanban struct {}
type ConfigMeasure struct {}
