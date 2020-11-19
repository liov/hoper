package main

import (
	"log"

	"github.com/spf13/viper"
)

var (
	Flow string
	Env  string

	Version       string
	ExportPort    string
	Author        string
	Namespace     string
	LogPath       string
	LogTargetPath string
	Domain        string
	ServicePort   string
	CpuLimit      string
	MemoryLimit   string
	CpuRequest    string
	MemoryRequest string
	Url           string
	CmdArgs       string
	ProjectPath   string
	ProxyStatus   string
	Type          string

	Name          string //系统名称
	ConfigPath    string //配置文件路径
	ConfigMapName string //configMap名称
	Conf          string
	BranchName    string //系统git分支名称
	UserName      string //系统git用户名称
	CommitNote    string //系统git发布特性内容
	Replicas      string //系统pod数量
	Port          string
	ConfigName    string //配置文件名称

	ConfigFile string //发布脚本目录下的配置文件路径

)

func SetupConfig(systemName string, filePath string) {
	ConfigFile = "./build/config.json"
	Name = systemName
	if Name == "" {
		panic("需要发布的系统名称错误")
	}
	if Version == "" {
		panic("未输入发布版本，版本号定义规则：项目版本_年月日时分秒，如：v3.4_20180207121201")
	}
	ConfigPath = filePath + "/config/"
	ConfigMapName = Name + "-config"
	Conf = ConfigPath + Env
	//重置viper，读取新的配置文件
	viper.Reset()
	viper.SetConfigName("deploy-" + Env)
	viper.AddConfigPath(ConfigPath)
	log.Println("正在从", ConfigPath+"deploy-"+Env, "读取配置文件")
	err := viper.ReadInConfig()
	if err != nil {
		panic("配置读取失败:" + err.Error())
	}

	//先获取最新代码
	GitPull(filePath)
	if Version == "" {
		//获取git版本号，每个系统都获取
		Version = GetGitTag(filePath)
		log.Println("获取到的tag版本是: ", Version)
	}

	log.Println("发布分支是: ", BranchName)
	//获取发布内容
	CommitNote = GetBranchNote()
	log.Println("发布内容特性: ", CommitNote)
	//获取发布者名称
	UserName = "线上发布程序"
	log.Println("发布人: ", UserName)

	//从配置文件读取配置
	GetAllConfigValue()

	//空值的情况下，设置默认参数值
	if Replicas == "" {
		Replicas = "1"
	}
	if Port == "" {
		Port = "80"
	}
	if ConfigName == "" {
		ConfigName = Namespace + "-" + Env + ".json"
	}
	if ServicePort == "" {
		ServicePort = "8080"
	}
	if Author == "" {
		Author = "erp"
	}
	if Namespace == "" {
		Namespace = "default"
	}
	if CpuLimit == "" {
		CpuLimit = "500m"
	}
	if MemoryLimit == "" {
		MemoryLimit = "512Mi"
	}
	if CpuRequest == "" {
		CpuLimit = "50m"
	}
	if MemoryRequest == "" {
		MemoryLimit = "64Mi"
	}
	if LogPath == "" {
		LogPath = "/data/logs"
	}
	if LogTargetPath == "" {
		LogTargetPath = "/data/logs"
	}
	if ProxyStatus == "" {
		ProxyStatus = "nginx"
	}
	if Env == "dev" {
		if Url == "" {
			Url = "reg.hoper.xyz"
		}
		if CmdArgs == "" {
			CmdArgs = "echo 192.168.1.210 keycloak.erp >> /etc/hosts"
		}
	} else if Env == "vpc" {
		if Url == "" {
			Url = "reg.hoper.com"
		}
		ExportPort = ""
	} else if Env == "classic" {
		if Url == "" {
			Url = "creg.hoper.com"
		}
		ExportPort = ""
	}
	if Type == "" {
		Type = "go"
	}
}

func GetAllConfigValue() {
	//获取配置内容
	Port = viper.GetString("port")             //svc内部端口
	ExportPort = viper.GetString("exportPort") //svc对外暴露端口
	Author = viper.GetString("author")
	Url = viper.GetString("url")
	LogPath = viper.GetString("log.path")
	LogTargetPath = viper.GetString("log.targetPath")
	CpuLimit = viper.GetString("cpuLimit")
	MemoryLimit = viper.GetString("memoryLimit")
	CpuRequest = viper.GetString("cpuRequest")
	MemoryRequest = viper.GetString("memoryRequest")
	Domain = viper.GetString("domain")
	CmdArgs = viper.GetString("cmdArgs")
	ServicePort = viper.GetString("servicePort") //pod内部程序端口
	Replicas = viper.GetString("replicas")
	if Namespace == "" {
		Namespace = viper.GetString("namespace")
	}
	ConfigName = viper.GetString("configName")
	Type = viper.GetString("type")

}
