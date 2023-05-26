package config

/**
配置或服务注册的参数
*/
type AppStartParams struct {
	RemoteConfigFile string //远程配置文件名 ,比如 config/maque-cloud
	YamlPath         string //Yaml文件所在的目录
	ServiceName      string //注册到注册中心的服务名称
}
