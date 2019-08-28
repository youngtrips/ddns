package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type PluginConfig struct {
	Name   string            `yaml:"name"`
	Params map[string]string `yaml:"params"`
}

type BaseConfig struct {
	Domain     string        `yaml:"domain"`
	CheckIpUrl string        `yaml:"check_ip_url"`
	Interval   int32         `yaml:"check_interval"`
	Records    []string      `yaml:"records"`
	Plugin     *PluginConfig `yaml:"dns_plugin"`
}

//////////////
var (
	_cfg BaseConfig
)

func init() {
}

func Load(cfgfile string) error {
	data, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal([]byte(data), &_cfg); err != nil {
		return err
	}
	return nil
}

func Config() *BaseConfig {
	return &_cfg
}
