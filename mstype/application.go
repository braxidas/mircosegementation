package mstype

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
