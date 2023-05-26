package config

import (
	"cloud-maque-common/utils/netutils"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/**
配置与服务注册
@param  ConfigSetCallback 回调参数设置
*/
func InitConfigServiceRegister(params AppStartParams, ConfigSetCallback func()) {

	fmt.Println("[初始化参数]......")

	var MqConfig AppConfig

	if params.YamlPath == "" {
		params.YamlPath = "./"
	}

	//系统环境变量
	//viper.SetEnvPrefix("maque") //将自动大写
	viper.BindEnv("config_enabled")
	viper.BindEnv("discovery_enabled")
	viper.BindEnv("config_type")
	viper.BindEnv("config_ip")
	viper.BindEnv("config_port")
	viper.BindEnv("service_ip") //当前服务IP
	//os.Setenv("SPF_ID", "13") // 通常在应用以外完成
	config_enabled := viper.GetString("config_enabled") // 13

	println(config_enabled)
	config_ip := viper.GetString("config_ip")
	config_port := viper.GetInt("config_port")

	if config_enabled != "" && strings.ToLower(config_enabled) == "true" {

		if params.RemoteConfigFile == "" {
			panic("远程配置文件名不能为空")
		}

		config_url := ""
		config_type := viper.GetString("config_type")
		if config_type == "" {
			config_type = "consul"
		}
		if config_ip == "" {
			config_url = "http://" + config_type + ".gisquest.com:8500"
			config_ip = config_type + ".gisquest.com"
			config_port = 8500
		} else {
			config_url = "http://" + config_ip + ":" + strconv.Itoa(config_port)
		}

		viper.AddRemoteProvider(config_type, config_url, params.RemoteConfigFile)
		viper.SetConfigType("yaml") // Need to explicitly set this to json
		err := viper.ReadRemoteConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config consul: %s \n", err))
			return
		}

		// 监控服务端配置变化
		go func() {
			for {
				time.Sleep(time.Second * 5) // delay after each request

				// currently, only tested with etcd support
				err2 := viper.WatchRemoteConfig()
				if err2 != nil {
					println("unable to read remote config: %v", err2)
					continue
				}

				MqConfig.Name = viper.GetString("application.server.name")
			}
		}()
	} else {
		Home, _ := os.UserHomeDir()
		fmt.Println("Home Directory:", Home)

		viper.SetConfigType("yaml")
		viper.SetConfigName("application") // name of config file (without extension)
		viper.AddConfigPath(params.YamlPath)
		//viper.AddConfigPath("./")          // 配置文件路径，多次使用可以查找多个目录
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
			return
		}
	}

	MqConfig.Name = viper.GetString("application.server.name")
	MqConfig.Port = viper.GetInt("application.server.port")

	discovery_enabled := viper.GetBool("discovery_enabled")
	if discovery_enabled {
		//注册到服务中心
		service_ip := viper.GetString("service_ip")

		if params.ServiceName == "" {
			panic("服务名称不能为空")
		}

		serviceRegister(params, config_ip, config_port, service_ip, MqConfig.Port)
	}

	//回调参数赋值
	ConfigSetCallback()
}

/**
服务注册
*/
func serviceRegister(params AppStartParams, ip string, port int, service_ip string, servciePort int) {

	serviceIp := service_ip
	if serviceIp == "" {
		serviceIp = netutils.GetLocalIpV4()
	}
	config := api.DefaultConfig()
	config.Address = ip + ":" + strconv.Itoa(port)
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	hostname := netutils.GetOsHostname()
	registration := new(api.AgentServiceRegistration)
	registration.Name = params.ServiceName
	registration.ID = params.ServiceName + "_" + hostname + "_" + strconv.Itoa(servciePort)
	registration.Address = serviceIp
	registration.Port = servciePort
	registration.Check = &api.AgentServiceCheck{
		CheckID:  registration.ID + "_Check",
		TCP:      serviceIp + ":" + strconv.Itoa(servciePort),
		Timeout:  "5s",
		Interval: "5s",
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

}

func unRegister(params AppStartParams, client *api.Client) {

	// 监听退出信号并注销服务
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
		log.Println("Shutting down...")
		err := client.Agent().ServiceDeregister(params.ServiceName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
