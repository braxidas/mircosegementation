package mstype

type Application struct {
	Server ServerConfig `yaml:port`
	Spring SpringConfig `yaml:spring`
}

type ServerConfig struct {
	Port string `yaml:port`
}

type SpringConfig struct {
	Application ApplicationConfig `yaml:application`
	Cloud       CloudConfig       `yaml:cloud`
}

type ApplicationConfig struct {
	Name string `yaml:name`
}

type CloudConfig struct {
	Nacos NacosConfig `yaml:nacos`
}

type NacosConfig struct {
	Discovery DiscoveryConfig `yaml:discovery`
}

type DiscoveryConfig struct {
	ServerAddr string `yaml:server-addr`
}
