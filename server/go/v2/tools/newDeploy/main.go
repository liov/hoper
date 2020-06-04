package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
)

type mypage struct {
	Projects []string
}

func New() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}
				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)
				ctx.Application().Logger().Warn(logMessage)

				ctx.StatusCode(500)
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}

func main() {
	app := iris.New()
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		//关闭所有主机
		app.Shutdown(ctx)
	})

	app.Use(New())

	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html"))
	app.Get("/{project}", handle)
	app.Get("/{project}/project/{mod}", branches)
	app.Post("/{project}/deploy", deploy)
	if err := app.Run(iris.Addr(":9090")); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}
}

func handle(ctx iris.Context) {
	ctx.Gzip(true)
	GOPATH := os.Getenv("GOPATH")
	rd, _ := ioutil.ReadDir(GOPATH + "/src/hoper.com/" + ctx.Params().Get("project"))
	var page mypage
	for _, fi := range rd {
		if fi.IsDir() {
			page.Projects = append(page.Projects, fi.Name())
		}
	}
	ctx.ViewData("", page)
	ctx.View("mypage.html")
}

func branches(ctx iris.Context) {
	ctx.Gzip(true)
	GOPATH := os.Getenv("GOPATH")
	path := GOPATH + "/src/maizuo.com/" + ctx.Params().Get("project") + "/" + ctx.Params().Get("mod")
	cmd := exec.Command("/bin/bash", "-c", "cd "+path+"\ngit branch -r")
	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		log.Println("获取分支失败,错误原因是: ", err)
	}
	arr := strings.Split(string(out), "\n")
	ctx.JSON(map[string]interface{}{"data": arr[:len(arr)-1]})
}

type Dep struct {
	Project string `json:"project"`
	Flow    string `json:"flow"`
	Env     string `json:"env"`
	Version string `json:"version"`
	Branch  string `json:"branch"`
}

var dep Dep
var mutexMap = make(map[string]*sync.Mutex)

func deploy(ctx iris.Context) {
	ctx.ReadJSON(&dep)
	fmt.Println(dep)
	Env = dep.Env
	Flow = dep.Flow
	BranchName = dep.Branch
	Namespace = ctx.Params().Get("project")
	m, ok := mutexMap[Namespace]
	if ok {
		m.Lock()
	} else {
		mutexMap[Namespace] = &sync.Mutex{}
		mutexMap[Namespace].Lock()
	}
	defer func() {
		mutexMap[Namespace].Unlock()
	}()
	Version = "v" + dep.Version + "-" + time.Now().Format("20060102150405")
	SystemName := dep.Project
	GOPATH := os.Getenv("GOPATH")
	os.Chdir(GOPATH + "/src/hoper.com/" + ctx.Params().Get("project") + "/" + dep.Project + "/")
	ProjectPath = GOPATH + "/src/hoper.com/" + ctx.Params().Get("project") + "/" + dep.Project

	if Flow != "help" {
		if Namespace == "" || Flow == "" || Env == "" || SystemName == "" {
			ctx.Write([]byte("缺少必要参数 flow,name,env,ns"))
			Help()
			return
		}
	}

	names := strings.Split(SystemName, ",")
	paths := strings.Split(ProjectPath, ",")
	if ProjectPath != "" {
		if len(paths) != len(names) {
			ctx.Write([]byte("参数关系不对应 name,path"))
			return
		}
	}

	for i, sName := range names {
		filePath := "../" + sName
		if ProjectPath != "" && paths[i] != "" {
			filePath = paths[i]
		}
		//读取配置文件
		SetupConfig(sName, filePath)
		switch Flow {
		case "all":
			//拷贝项目配置文件
			CopyConfig()
			//生成Dockerfile文件
			BuildDockerfile()
			//编译go可执行文件
			Build(filePath)
			//编译代码生成docker镜像
			BuildDockerImage(filePath)
			//推送docker镜像到k8s
			PushDockerImage()
			//生成k8s文件
			BuildDeployment()
			//生成svc文件
			//BuildService()
			//切换k8s环境
			ChangeEnv()
			ChangeNameSpace()
			//使用k8s发布项目
			DelConfigMap()
			ApplyConfigMap()
			ApplyDeployment()
			//ApplyService()
			SendDingTalk()

		case "rollback":
			//读取配置文件
			RollBackConfig()
			//切换k8s环境
			ChangeEnv()
			ChangeNameSpace()
			BuildDeployment()
			DelConfigMap()
			ApplyConfigMap()
			DelDeployment()
			ApplyDeployment()
			SendDingTalk()

		case "ing":
			//切换k8s环境
			ChangeEnv()
			BuildIngress()
			DelIngress()
			ApplyIngress()

		case "svc":
			//切换k8s环境
			ChangeEnv()
			BuildService()
			DelService()
			ApplyService()

		case "dep":
			//切换k8s环境
			ChangeEnv()
			ChangeNameSpace()
			//拷贝项目配置文件
			CopyConfig()
			DelConfigMap()
			ApplyConfigMap()
			BuildDeployment()
			DelDeployment()
			ApplyDeployment()
			SendDingTalk()

		case "config":
			//切换k8s环境
			ChangeEnv()
			ChangeNameSpace()
			DelConfigMap()
			ApplyConfigMap()

		case "build":
			//生成Dockerfile文件
			BuildDockerfile()
			//拷贝项目配置文件
			CopyConfig()
			//编译go可执行文件
			Build(filePath)
			//编译代码生成docker镜像
			BuildDockerImage(filePath)

		case "test":
			//本地运行
			DockerRun()

		case "push":
			//推送docker镜像到k8s
			PushDockerImage()

		case "delAll":
			ChangeEnv()
			//生成k8s文件
			BuildDeployment()
			//生成svc文件
			BuildService()
			//生成ing文件
			BuildIngress()
			DelReplicationController()
			DelDeployment()
			DelService()
			DelIngress()

		case "delSvc":
			ChangeEnv()
			//生成svc文件
			BuildService()
			DelService()
		case "delIng":
			ChangeEnv()
			//生成ing文件
			BuildIngress()
			DelIngress()
		case "delDep":
			ChangeEnv()
			//生成k8s文件
			BuildDeployment()
			DelDeployment()
		default:
			ctx.Write([]byte("不存在的flow参数"))
		}
	}

	ctx.Write([]byte("成功"))
}
