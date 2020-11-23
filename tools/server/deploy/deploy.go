package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"
)

func GitPull(filePath string) {
	log.Println("正在获取当前项目git最新代码")
	err, _ := execCommand("cd " + filePath + "\ngit fetch --all && git reset --hard origin/master && git pull")

	err, _ = execCommand("cd " + filePath + "\ngit checkout -b " + config.BranchName + " origin/" + config.BranchName)
	if err != nil {
		execCommand("cd " + filePath + "\ngit checkout " + config.BranchName)
	}

	err, _ = execCommand("cd " + filePath + "\ngit pull")
	if err != nil {
		log.Println("获取代码失败,错误原因是: ", err.Error())
		return
	}
}

func GetGitTag(filePath string) string {
	log.Println("正在获取当前项目git最新tag")
	err, result := execCommand("cd " + filePath + "\ngit tag")
	if err != nil {
		log.Println("获取tag失败,错误原因是: ", err.Error())
		return ""
	}
	arr := strings.Split(result[:], "\n")
	if len(arr) < 2 {
		return ""
	}
	return arr[len(arr)-2]
}

func execCommand(command string) (error, string) {
	log.Println("当前正在执行: \n", command)
	cmd := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in
	go func() {
		in.WriteString(command)
	}()
	result, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return err, ""
	}
	return nil, string(result[:])
}

func BuildDockerfile() {
	var dockerfile string
	switch config.Type {
	case "go":
		dockerfile = Dockerfile_Go
	case "jar":
		dockerfile = Dockerfile_Jar
	}
	t := template.Must(template.New("dockerfile").Parse(dockerfile))
	file := WriteFile("./Dockerfile")
	defer file.Close()
	err := t.Execute(file, &config)
	log.Println(err)
}

func BuildIngress() {
	t := template.Must(template.New("ingress").Parse(Ingress))
	file := WriteFile("./k8s/ing.yaml")
	defer file.Close()
	err := t.Execute(file, &config)
	log.Println(err)
}

func BuildService() {
	t := template.Must(template.New("msg").Parse(Service))
	file := WriteFile("./k8s/svc.yaml")
	defer file.Close()
	err := t.Execute(file, &config)
	log.Println(err)

}

func BuildDeployment() {
	log.Println("正在构建deployment文件......")

	if config.CmdArgs == "" {
		config.CmdArgs += "&& cp ../../config/default.json default.json && nohup ./{{name}} -conf default"
	}

	t := template.Must(template.New("deployment").Parse(Deployment))
	file := WriteFile("./k8s/dep.yaml")
	defer file.Close()
	err := t.Execute(file, &config)
	log.Println(err)
}

func Build(filePath string) {
	switch config.Type {
	case "go":
		{
			log.Println("正在打包go项目......")
			execAndPrint("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-o", "./main", filePath+"/main.go")
		}
	case "web":
		{
			log.Println("正在打包web项目......")
			execAndPrint("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-o", "./build/main", filePath+"/main.go")
			pwd := os.Getenv("PWD")
			os.Chdir(filePath)
			execAndPrint("npm", "i")
			execAndPrint("npm", "run-script", "build")
			os.Chdir(pwd)
			execAndPrint("cp", "-r", filePath+"/build", "./build/")
		}
	case "jar":
		{
			log.Println("正在打包java项目jar包......")
			execAndPrint("gradle", "-b", filePath+"/build.gradle", "clean", "jar", "-Pprofile="+config.Env)
			execAndPrint("cp", "-r", filePath+"/build/libs/"+config.Name+".jar", "./build/"+config.Name+".jar")
		}
	}
	time.Sleep(time.Second * 2)

}

func DockerRun() {
	log.Println("本地测试运行go项目......")
	execAndPrint("docker", "run", "-ti", "--rm", config.ImageTag)
}

func CopyConfig() {
	log.Println("正在复制配置", config.Type, "项目文件到发布系统中......")
	err1 := os.RemoveAll("./build/")
	if err1 != nil {
		log.Println("删除原配置文件失败：" + "./build/config/")
		return
	}
	err := os.MkdirAll("./build/config/", os.ModePerm)
	if err != nil {
		log.Println("创建配置文件目录失败：" + "./build/config/")
		return
	}
	execAndPrint("pwd")
	switch config.Type {
	case "go":
		execAndPrint("cp", "-f", config.Conf+".json", config.ConfigFile)
	case "java":
		execAndPrint("cp", "-r", config.ConfigPath+config.Env, "./build/config")
	case "war":

	case "web":
		execAndPrint("cp", "-f", config.Conf+".json", config.ConfigFile)
	case "jar":

	}

}

func BuildDockerImage(path string) {
	os.Chdir(path)
	log.Println("正在构建docker镜像......")
	execAndPrint("docker", "build", "-t", config.ImageTag, ".")
}

func PushDockerImage() {
	log.Println("正在推送docker镜像......")
	execAndPrint("docker", "push", config.ImageTag)
}

func ChangeEnv() {
	log.Println("正在切换k8s环境......")
	pwd, _ := os.Getwd()
	os.Chdir("/home/crm/k8s")
	command := exec.Command("/bin/bash", "-c", "./kube.sh "+config.Env)
	command.Start()
	command.Wait()
	execAndPrint("kubectl", "config", "use-context", config.Env)
	log.Println("已切换到" + config.Env)
	os.Chdir(pwd)
}

func ApplyDeployment() {
	log.Println("正在发布dep......")
	execAndPrint("kubectl", "apply", "-f", "./k8s/dep.yaml", "--record")
	backupConfig()
}

func ApplyService() {
	log.Println("正在发布svc......")
	execAndPrint("kubectl", "apply", "-f", "./k8s/svc.yaml", "--record")
}

func ApplyIngress() {
	log.Println("正在发布ing......")
	execAndPrint("kubectl", "apply", "-f", "./k8s/ing.yaml", "--record")
}

func ApplyConfigMap() {
	if config.Type == "go" || config.Type == "web" {
		log.Println("正在发布配置文件.....")
		execAndPrint("kubectl", "create", "configmap", config.ConfigMapName, "--from-file="+config.ConfigPath)
	}
}

func DelDeployment() {
	log.Println("正在删除dep......")
	execAndPrint("kubectl", "delete", "-f", "./k8s/dep.yaml")
}

func DelReplicationController() {
	log.Println("正在删除rc......")
	execAndPrint("kubectl", "delete", "-f", "./k8s/rc.yaml")
}

func DelService() {
	log.Println("正在删除svc......")
	execAndPrint("kubectl", "delete", "-f", "./k8s/svc.yaml")
}

func DelIngress() {
	log.Println("正在删除ing......")
	execAndPrint("kubectl", "delete", "-f", "./k8s/ing.yaml")
}

func DelConfigMap() {
	if config.Type == "go" || config.Type == "web" {
		log.Println("正在发布配置文件.....")
		execAndPrint("kubectl", "delete", "configmap", config.ConfigMapName)
	}
}

func ChangeNameSpace() {
	execAndPrint("kubectl", "config", "set-context", config.Env, "--namespace="+config.Namespace)
	log.Println(config.Env, "环境的命名空间已切换成", config.Namespace)
}

func execAndPrint(commandName string, params ...string) string {
	endLine := ""
	cmd := exec.Command(commandName, params...)
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("执行出错", err)
		return endLine
	}

	errStart := cmd.Start()
	if errStart != nil {
		log.Println("执行出错", errStart)
		return endLine
	}
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		endLine = line
		fmt.Printf(line)
	}

	err = cmd.Wait()
	if err != nil {
		if params[0] != "delete" {
			panic("脚本执行出错,请仔细检查docker是否启动,k8s客户端是否安装,版本是否正确" + err.Error())
		}
	}
	return endLine
}

func WriteFile(fileName string) *os.File {
	var file *os.File
	var err error
	if checkFileIsExist(fileName) {
		//如果文件存在
		file, err = os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY, 0666) //打开文件
		if err != nil {
			log.Println("写入文件失败1:", fileName, err)
		}
	} else {
		err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
		if err != nil {
			log.Println("写入文件失败2:", fileName, err)
		}
		file, err = os.Create(fileName) //创建文件
		defer file.Close()
		if err != nil {
			log.Println("写入文件失败3:", fileName, err)
		}
	}
	return file
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(fileName string) bool {
	var exist = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//发送钉钉信息
func SendDingTalk() {
	log.Println("Env:" + config.Env)
	//if Env != "prod" {
	//	return
	//}
	//发送钉钉信息
	hookUrl := "https://oapi.dingtalk.com/robot/send?access_token=dd01e86251264d6c5d089d295ffdc3c3d4a7ac59cce2e7c7d66782adb4f8ec5f"
	noteArray := strings.Split(config.CommitNote, "\n")
	note := strings.Join(noteArray[0:len(noteArray)-1], " \n - ")

	text := `#### 系统：{{Name}}.{{Namespace}} \n
		#### 环境：{{Env}} \n 
		#### 时间：` + time.Now().Format("2006年01月02日 15:04:05") + ` \n
		#### 模块：{{Name}}\n 
		#### 分支：{{BranchName}}\n
		#### 版本：{{Version}}\n 
		#### 特性：\n 
		- ` + note + `\n
		#### 发布人：{{UserName}} \n `
	t := template.Must(template.New("dingding").Parse(text))
	buf := new(bytes.Buffer)
	t.Execute(buf, &config)
	str := `
		{
		     "msgtype": "markdown",
		     "markdown": {
			 "title":"系统发布",
			 "text": "` + buf.String() + `"
		     },
		    "at": {
			"atMobiles": [],
			"isAtAll": false
		    }
		 }
		`
	body := strings.NewReader(str)
	http.Post(hookUrl, "application/json", body)
}

/**
回滚读取配置文件
*/
func RollBackConfig() {
	if (config.Type == "go" || config.Type == "web") && config.Env == "prod" {
		dir := "./history/" + config.Name + "/" + config.Namespace + "/" + config.Env
		log.Println("读取历史版本发布配置config")
		//备份当前配置文件config
		execAndPrint("mkdir", "build")
		execAndPrint("cp", "-f", dir+"/"+config.Version+".json", config.ConfigFile)
	}
}

/**
备份发布历史配置文件config.json
*/
func backupConfig() {
	if (config.Type == "go" || config.Type == "web") && config.Env == "prod" {
		dir := "./history/" + config.Name + "/" + config.Namespace + "/" + config.Env
		//创建备份目录
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("创建备份文件失败：" + dir)
			return
		}
		//java项目不复制配置文件
		log.Println("备份当前发布配置config")
		//备份当前配置文件config
		execAndPrint("cp", "-f", config.ConfigFile, dir+"/"+config.Version+".json")
	}
	//备份完后删除build目录
	os.RemoveAll("./build/")
}

func GetBranchNote() string {
	log.Println("正在获取当前分支的注释")
	//err, result := execCommand("cd ../" + Name + " && git branch -v|awk '{print $NF}'")
	err, result := execCommand("cd ../" + config.Name + " && git log --pretty=format:'%s' -n 3")
	if err != nil {
		log.Println("获取注释异常,错误原因是: ", err.Error())
		return ""
	}
	return result
}

func GetBranchName() string {
	log.Println("正在获取当前发布的分支名称")
	err, result := execCommand("cd ../" + config.Name + " && git symbolic-ref --short -q HEAD")
	if err != nil {
		log.Println("获取分支失败,错误原因是: ", err.Error())
		return ""
	}
	return result
}

func Help() {
	helpInfo := `	 ----------------------------------------------------------
	              项目发布脚本使用指南
	 ----------------------------------------------------------
	 可选参数如下:
	 flow [必须] 指定脚本执行流程 all,ing,svc,dep,config,build,push,delAll,delSvc,delIng,delDep
		--all:正常发布流程，但是不会发布svc,ing，如果版本号一样需要手动重启项目
		--config:只更新配置文件，需要项目支持配置文件热更新
		--dep:根据版本号滚动更新deployment和configMap
		--svc:只更新svc,一般上线时执行一次即可
		--ing:只更新ingress,一般上线时执行一次即可
		--build:只执行docker打包镜像,很少用
		--push:只执行docker Push命令,很少用
	 name [必须] 需要操作的的项目名称, 多个,隔开
	 env [必须] 指定发布环境 dev[测试环境], stage[预正式环境], prod[正式环境]
	 ns [必须] 切换项目namespace
	 ver [必须] 指定项目版本
	 path [非必须] 项目相对路径
	 -----------------------------------------------------------
	 常用命令示例：
	 发布项目: make flow=all ns={ns} name={project} env=dev ver=v1.0
	 停止项目: make flow=delAll ns={ns} name={project} env=dev
	 发布dep: make flow=dep ns={ns} name={project} env=dev ver=v1.0
	 切换namespace和发布环境: make flow=env ns={ns} name={project} env=dev
	 强制删除po:  kube delete po xx --grace-period=0 --force=true
	 示例:  make flow=all ns={ns} name={project} env=dev
	 -----------------------------------------------------------
	 ps:每次更新脚本后,需要执行 make init 或 go build -o ./main ./main.go`
	log.Println(helpInfo)
}
