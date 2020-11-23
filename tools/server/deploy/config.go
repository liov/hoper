package main

import (
	"log"

	"github.com/spf13/viper"
)

var config = struct {
	Flow string `json:"flow"`
	Env  string `json:"env"`

	Version       string `json:"version"`
	ExportPort    string `json:"exportPort"`
	Author        string `json:"author"`
	Namespace     string `json:"namespace"`
	LogPath       string `json:"logPath"`
	LogTargetPath string `json:"logTargetPath"`
	Domain        string `json:"domain"`
	ServicePort   string `json:"servicePort"`
	CpuLimit      string `json:"cpuLimit"`
	MemoryLimit   string `json:"memoryLimit"`
	CpuRequest    string `json:"cpuRequest"`
	MemoryRequest string `json:"memoryRequest"`
	Url           string `json:"url"`
	CmdArgs       string `json:"cmdArgs"`
	ProjectPath   string `json:"projectPath"`
	ProxyStatus   string `json:"proxyStatus"`
	Type          string `json:"type"`

	Name          string `json:"name"`          //系统名称
	ConfigPath    string `json:"configPath"`    //配置文件路径
	ConfigMapName string `json:"configMapName"` //configMap名称
	Conf          string `json:"conf"`
	BranchName    string `json:"branchName"` //系统git分支名称
	UserName      string `json:"userName"`   //系统git用户名称
	CommitNote    string `json:"commitNote"` //系统git发布特性内容
	Replicas      string `json:"replicas"`   //系统pod数量
	Port          string `json:"port"`
	ConfigName    string `json:"configName"` //配置文件名称

	ConfigFile             string //发布脚本目录下的配置文件路径
	ImageTag               string
	dockerfile, deployment string
}{
	Author:        "线上发布",
	Namespace:     "default",
	LogPath:       "/data/logs",
	LogTargetPath: "/data/logs",
	ServicePort:   "8080",
	CpuLimit:      "500m",
	MemoryLimit:   "512Mi",
	CpuRequest:    "50m",
	MemoryRequest: "64Mi",
	ProxyStatus:   "nginx",
	Type:          "go",
	Replicas:      "1",
	Port:          "80",
	ConfigFile:    "./build/config.json",
}

func SetupConfig(systemName string, filePath string) {
	if systemName == "" {
		panic("需要发布的系统名称错误")
	}

	if config.Domain == "" {
		panic("项目缺少 domain配置")
	}
	if config.Version == "" {
		panic("未输入发布版本，版本号定义规则：项目版本_年月日时分秒，如：v3.4_20180207121201")
	}
	config.Name = systemName
	config.ConfigPath = filePath + "/config/"
	config.ConfigMapName = systemName + "-config"
	config.Conf = config.ConfigPath + config.Env

	//重置viper，读取新的配置文件
	viper.Reset()
	viper.SetConfigName("deploy-" + config.Env)
	viper.AddConfigPath(config.ConfigPath)
	log.Println("正在从", config.ConfigPath+"deploy-"+config.Env, "读取配置文件")
	err := viper.ReadInConfig()
	if err != nil {
		panic("配置读取失败:" + err.Error())
	}
	//从配置文件读取配置
	viper.Unmarshal(&config)
	config.ImageTag = config.Url + "/" + config.Author + "/" + config.Name + ":" + config.Version
	//先获取最新代码
	GitPull(filePath)
	if config.Version == "" {
		//获取git版本号，每个系统都获取
		config.Version = GetGitTag(filePath)
		log.Println("获取到的tag版本是: ", config.Version)
	}

	log.Println("发布分支是: ", config.BranchName)
	//获取发布内容
	config.CommitNote = GetBranchNote()
	log.Println("发布内容特性: ", config.CommitNote)
	//获取发布者名称
	config.UserName = "线上发布程序"
	log.Println("发布人: ", config.UserName)

	if config.ConfigName == "" {
		config.ConfigName = config.Namespace + "-" + config.Env + ".json"
	}
}
