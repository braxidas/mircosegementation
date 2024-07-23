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
	Fdfs   FdfsConfig   `yaml:"fdfs"`
	Minio  MinioConfig  `yaml:"minio"`
}

// #
type ServerConfig struct {
	Port string `yaml:"port"`
}

// #
type SpringConfig struct {
	Application ApplicationConfig `yaml:"application"`
	Cloud       CloudConfig       `yaml:"cloud"`
	Redis       RedisConfig       `yaml:"redis"`
	DataSource  DataSourceConfig  `taml:"datasource"`
}

// ##
type ApplicationConfig struct {
	Name string `yaml:"name"`
}

// ##
type CloudConfig struct {
	Nacos NacosConfig `yaml:"nacos"`
	Gateway GatewayConfig `yaml:"gateway"`
}

// ##
type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type DataSourceConfig struct {
	Url     string        `yaml:"url"`
	Dynamic DynamicConfig `yaml:"dynamic"`
}

// ###
type NacosConfig struct {
	Discovery DiscoveryConfig `yaml:"discovery"`
	Config    ConfigConfig    `yaml:"config"`
}

type GatewayConfig struct{
	Routes []*RouteConfig `yaml:"routes"`
}

// ###
type DynamicConfig struct {
	DataSource DataSourceLiConfig `yaml:"datasource"`
}

// ####
type RouteConfig struct{
	Id string `yaml:"id"`
	Uri string  `yaml:"uri"`
}


// ####
type DataSourceLiConfig struct {
	Master DataSourceElemConfig `yaml:"master"`
	Slave  DataSourceElemConfig `yaml:"slave"`
}

// ####
type DiscoveryConfig struct {
	ServerAddr string `yaml:"server-addr"`
}
type ConfigConfig struct {
	ServerAddr string `yaml:"server-addr"`
	Group      string `yaml:"group"`
	NameSpace string `yaml:"namespace"`
}

// #####
type DataSourceElemConfig struct {
	Url string `yaml:"url"`
}

// #
type DubboConfig struct {
	Application ApplicationConfig `yaml:"application"`
	Registry    RegistryConfig    `yaml:"registry"`
}

type RegistryConfig struct {
	Address string `yaml:"address"`
}

// #
type FdfsConfig struct {
	Domain      string `yaml:"domain"`
	TrackerList string `yaml:"trackerList"`
}

// #
type MinioConfig struct {
	Url string `yaml:"url"`
}

func (application *Application) GetApplicationName() (string, error) {
	if len(application.Spring.Application.Name) != 0 {
		return application.Spring.Application.Name, nil
	}
	if len(application.Dubbo.Application.Name) != 0 {
		return application.Dubbo.Application.Name, nil
	}
	return "", fmt.Errorf("Application Name not Found")
}

// func (application *Application) ReplaceEnv(env map[string]string)error{




// 	return nil
// }