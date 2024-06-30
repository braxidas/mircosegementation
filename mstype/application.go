package mstype

import (
	"fmt"
)
/*
此类用于描述一个微服务的的配置文件
*/
type Application struct {
	Server ServerConfig `yaml:"server"`
	Spring SpringConfig `yaml:"spring"`
	Dubbo  DubboConfig  `yaml:"dubbo"`
}

// #
type ServerConfig struct {
	Port string `yaml:"port"`
}

// #
type SpringConfig struct {
	Application ApplicationConfig `yaml:"application"`
	Cloud       CloudConfig       `yaml:"cloud"`
}

// ##
type ApplicationConfig struct {
	Name string `yaml:"name"`
}

// ##
type CloudConfig struct {
	Nacos NacosConfig `yaml:"nacos"`
}

// ###
type NacosConfig struct {
	Discovery DiscoveryConfig `yaml:"discovery"`
}

// ####
type DiscoveryConfig struct {
	ServerAddr string `yaml:"server-addr"`
}

// #
type DubboConfig struct {
	Application ApplicationConfig `yaml:"application"`
	Registry    RegistryConfig    `yaml:"registry"`
}

type RegistryConfig struct {
	Address string `yaml:"address"`
}

func (application Application) GetApplicationName()(string, error){
	if(len(application.Spring.Application.Name) != 0){
		return application.Spring.Application.Name,nil
	}
	if(len(application.Dubbo.Application.Name) != 0){
		return application.Dubbo.Application.Name,nil
	}
	return "", fmt.Errorf("Application Name not Found")
}