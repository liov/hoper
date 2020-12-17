
# 常量定义 每个项目处理
SERVICE_NAME= business-sp-mng
NAME_SPACE= hoper



MAIN_GO = "./main.go"
DEPLOY_ENV =
DEPLOY_ENV_DEV = dev
DEPLOY_ENV_PROD = prod
DEPLOY_ENV_STAGE = stage

CONFIG_PATH =
CONFIG_LOCAL = "./conf/local.json"
CONFIG_DEV   = "./conf/dev.json"
CONFIG_PROD  = "./conf/prod.json"
CONFIG_STAGE  = "./conf/stage.json"

DEPLOY_TMP_DIR = tmp
CONFIG_VALUE = ""

# 镜像仓库地址
DEV_DOCKER_PREFIX = reg.hoper.so
PROD_DOCKER_PREFIX = reg.hoper.xyz

# 镜像仓库项目
DOCKER_REG_SPACE = hoper

# docker config文件占位替换
DOCKER_ENV_PROD = prod
DOCKER_ENV_DEV = dev
DOCKER_ENV_STAGE = stage

# K8S上下文
K8S_CONTEXT =
K8S_CONTEXT_DEV = dev
K8S_CONTEXT_PROD = prod
K8S_CONTEXT_STAGE = stage

# K8S命名空间
K8S_NAMESPACE = hoper

VERSION=`git describe --tags 2>/dev/null || git rev-parse --abbrev-ref HEAD 2>/dev/null || echo v1.0.0`


ifeq (prod, ${env})
	DOCKER_REG = $(PROD_DOCKER_PREFIX)
	CONFIG_PATH = $(CONFIG_PROD)
	DOCKER_ENV = $(DOCKER_ENV_PROD)
	K8S_CONTEXT = $(K8S_CONTEXT_PROD)
	DEPLOY_ENV = $(DEPLOY_ENV_PROD)
	CONFIG_VALUE = $(shell cat conf/prod.json)
else ifeq (stage, ${env})
    DOCKER_REG = $(PROD_DOCKER_PREFIX)
    CONFIG_PATH = $(CONFIG_STAGE)
    DOCKER_ENV = $(DOCKER_ENV_STAGE)
    K8S_CONTEXT = $(K8S_CONTEXT_STAGE)
    DEPLOY_ENV = $(DEPLOY_ENV_STAGE)
    CONFIG_VALUE = $(shell cat conf/stage.json)
else
	DOCKER_REG = $(DEV_DOCKER_PREFIX)
	CONFIG_PATH = $(CONFIG_DEV)
	DOCKER_ENV = $(DOCKER_ENV_DEV)
	K8S_CONTEXT = $(K8S_CONTEXT_DEV)
	DEPLOY_ENV = $(DEPLOY_ENV_DEV)
	CONFIG_VALUE = $(shell cat conf/dev.json)
endif


DOCKER_TAG=${tag}-dev
ifeq (prod, ${env})
	DOCKER_TAG=${tag}-prod
else ifeq (stage, ${env})
	DOCKER_TAG=${tag}-stage
endif

# 定义命令包

define build-linux
GOOS=linux GOARCH=amd64 go build -o main
endef

define go
go run main.go -conf $(CONFIG_LOCAL)
endef

# 清理
.PHONY : clean
clean:
	git checkout Dockerfile
	rm $(SERVICE_NAME)

# ----------------------------------------
# 项目初始化
.PHONY : init
init:
	@echo creating deploy directory: $(DEPLOY_TMP_DIR)/$(DEPLOY_ENV)
	@mkdir -p $(DEPLOY_TMP_DIR)/$(DEPLOY_ENV)

# ----------------------------------------
# 本地运行项目
.PHONY :  go
go:
	@echo starting go project
	$(go)

# ----------------------------------------
# 构建执行文件
.PHONY : package
package:
	@echo building go file
	$(build-linux)
	@echo finished building

# ----------------------------------------
# 打包docker
.PHONY : docker
docker: init check package
	@echo building docker image tag: ${tag}
	docker build -t $(DOCKER_REG)/$(DOCKER_REG_SPACE)/$(SERVICE_NAME):$(DOCKER_TAG) .
	docker push $(DOCKER_REG)/$(DOCKER_REG_SPACE)/$(SERVICE_NAME):$(DOCKER_TAG)
	@echo Image has been pushed !


# ----------------------------------------
# 检查参数
.PHONY : config
config: init
	@echo Set ConfigMap: $(K8S_CONTEXT)
	@head -n 3  deploy/k8s/configmap.yaml > tmp/$(DEPLOY_ENV)/configmap.yaml
	@echo '    $(CONFIG_VALUE)' >> tmp/$(DEPLOY_ENV)/configmap.yaml
	@tail -n 4 deploy/k8s/configmap.yaml >>  tmp/$(DEPLOY_ENV)/configmap.yaml
	kubectl apply -f tmp/$(DEPLOY_ENV)/configmap.yaml -n ${NAME_SPACE}

# ----------------------------------------
# 检查镜像标签参数
.PHONY : check
check:
ifeq (""a, ${tag}a)
	@echo  miss tag value
	exit 1
else ifeq (  , ${tag})
	@echo  miss tag value
	exit 1
endif
	@echo "tag : ${tag} "

# ----------------------------------------
# k8s发布
.PHONY : k8s
k8s:
	kubectl config use-context $(K8S_CONTEXT)
	sed  's/#tag/$(DOCKER_TAG)/g' deploy/k8s/$(DEPLOY_ENV)/deployment.yaml > tmp/$(DEPLOY_ENV)/deployment_$(DOCKER_TAG).yaml
	kubectl apply -f tmp/$(DEPLOY_ENV)/deployment_$(DOCKER_TAG).yaml -n ${NAME_SPACE}
	kubectl get po -n aurayou | grep $(SERVICE_NAME) -n ${NAME_SPACE}

# ----------------------------------------
# 完整发布
.PHONY : deploy
deploy: init package check docker k8s
	@echo deploy complete
