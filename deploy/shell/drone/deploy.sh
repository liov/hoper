#!/bin/bash

set -e

TPL_DIR="/tpl"

# 默认配置变量
DEPLOY_DIR="${PLUGIN_DEPLOY_DIR:-./deploy}"
CERT_DIR="${PLUGIN_CERT_DIR:-${DEPLOY_DIR}/cert}"
GROUP="${PLUGIN_GROUP}"
NAME="${PLUGIN_NAME}"
FULL_NAME="${PLUGIN_FULL_NAME}"
DEPLOY_KIND="${PLUGIN_DEPLOY_KIND}"
BUILD_TYPE="${PLUGIN_BUILD_TYPE}"


# Docker配置
DOCKERFILE_PATH="${PLUGIN_DOCKERFILE_PATH}"
DOCKER_USERNAME="${PLUGIN_DOCKER_USERNAME}"
DOCKER_PASSWORD="${PLUGIN_DOCKER_PASSWORD}"
DOCKER_CMD="${PLUGIN_DOCKER_CMD}"
IMAGE_TAG="${PLUGIN_IMAGE_TAG}"

# Kubernetes配置
DEPLOY_PATH="${PLUGIN_DEPLOY_PATH}"
DATA_DIR="${PLUGIN_DATA_DIR}"
CONF_DIR="${PLUGIN_CONFIG_DIR}"
SCHEDULE="${PLUGIN_SCHEDULE}"
CLUSTER="${PLUGIN_CLUSTER}"
CACRT="${PLUGIN_CA_CRT}"
DEVCRT="${PLUGIN_DEV_CRT}"
DEVKEY="${PLUGIN_DEV_KEY}"

# 其他环境变量
COMMIT_TAG="${DRONE_TAG:-${DRONE_COMMIT_TAG}}"

# 后处理配置
if [ -z "$DEPLOY_DIR" ]; then
    DEPLOY_DIR="./deploy"
fi

# 移除路径末尾的斜杠
DEPLOY_DIR=$(echo "$DEPLOY_DIR" | sed 's|/$||')
CERT_DIR="$DEPLOY_DIR/cert"

if [ -z "$NAME" ]; then
    FULL_NAME="$GROUP"
else
    FULL_NAME="${GROUP}-${NAME}"
fi

if [ -z "$IMAGE_TAG" ]; then
    # 尝试从COMMIT_TAG中提取版本信息
    if [[ "$COMMIT_TAG" =~ ^$FULL_NAME-v(.+) ]]; then
        image_tag="${BASH_REMATCH[1]}"
        IMAGE_TAG="$DOCKER_USERNAME/$FULL_NAME:$image_tag"
    elif [[ "$COMMIT_TAG" =~ ^v(.+) ]]; then
        IMAGE_TAG="$DOCKER_USERNAME/$FULL_NAME:${BASH_REMATCH[1]}"
    elif [[ "$COMMIT_TAG" =~ .*-v(.+) ]]; then
        # 使用awk分割字符串
        IMAGE_TAG="$DOCKER_USERNAME/$FULL_NAME:$(echo $COMMIT_TAG | awk -F'-v' '{print $2}')"
    else
        IMAGE_TAG="$DOCKER_USERNAME/$FULL_NAME:$(date '+%Y-%m-%d %H:%M:%S')"
    fi
else
    IMAGE_TAG="$DOCKER_USERNAME/$FULL_NAME:$(date '+%Y-%m-%d %H:%M:%S')"
fi

DOCKERFILE_PATH="$TPL_DIR/Dockerfile-$BUILD_TYPE"
DEPLOY_PATH="$TPL_DIR/deploy-$DEPLOY_KIND.yaml"

# 创建部署目录
if [ ! -d "$DEPLOY_DIR" ]; then
    mkdir -p "$DEPLOY_DIR"
fi

# 读取并处理Dockerfile模板
if [ ! -f "$TPL_DIR/Dockerfile-$BUILD_TYPE" ]; then
    echo "错误: 找不到Dockerfile模板 $TPL_DIR/Dockerfile-$BUILD_TYPE"
    exit 1
fi

dockerfile_content=$(cat "$TPL_DIR/Dockerfile-$BUILD_TYPE")
# 使用sed批量替换Dockerfile中的变量
processed_dockerfile=$(echo "$dockerfile_content" | \
    sed "s|\${app}|$FULL_NAME|g" | \
    sed "s|\${cmd}|$(echo $DOCKER_CMD | sed 's/","/","/g')|g")

# 写入新的Dockerfile
dockerfile_path="$DEPLOY_DIR/$FULL_NAME-Dockerfile"
echo "$processed_dockerfile" > "$dockerfile_path"

# 构建Docker镜像
echo "构建Docker镜像..."
docker build -f "$dockerfile_path" -t "$IMAGE_TAG" "$DEPLOY_DIR"

# 登录Docker仓库
echo "登录Docker仓库..."
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"

# 推送Docker镜像
echo "推送Docker镜像..."
docker push "$IMAGE_TAG"

# 处理Kubernetes部署文件
if [ ! -f "$TPL_DIR/deploy-$DEPLOY_KIND.yaml" ]; then
    echo "错误: 找不到Kubernetes部署模板 $TPL_DIR/deploy-$DEPLOY_KIND.yaml"
    exit 1
fi

deploy_content=$(cat "$TPL_DIR/deploy-$DEPLOY_KIND.yaml")

# 使用sed批量替换部署文件中的变量
processed_deploy=$(echo "$deploy_content" | \
    sed "s|\${app}|$FULL_NAME|g" | \
    sed "s|\${image}|$IMAGE_TAG|g" | \
    sed "s|\${group}|$GROUP|g" | \
    sed "s|\${datadir}|$DATA_DIR|g" | \
    sed "s|\${confdir}|$CONF_DIR|g")

# 如果是cronjob类型，还需要替换schedule
if [ "$DEPLOY_KIND" = "cronjob" ]; then
    processed_deploy=$(echo "$processed_deploy" | sed "s|\${schedule}|$SCHEDULE|g")
fi
# 写入新的部署文件
deploy_path="$DEPLOY_DIR/$FULL_NAME-$DEPLOY_KIND.yaml"
echo "$processed_deploy" > "$deploy_path"

# 处理证书文件
cacrt_path="$CERT_DIR/$CLUSTER/ca.crt"
if [ -n "$CACRT" ]; then
    mkdir -p "$CERT_DIR/$CLUSTER"
    echo "$CACRT" | base64 -d > "$cacrt_path"
fi

devcrt_path="$CERT_DIR/$CLUSTER/dev.crt"
if [ -n "$DEVCRT" ]; then
    echo "$DEVCRT" | base64 -d > "$devcrt_path"
fi

devkey_path="$CERT_DIR/$CLUSTER/dev.key"
if [ -n "$DEVKEY" ]; then
    echo "$DEVKEY" | base64 -d > "$devkey_path"
fi

# 设置KUBECONFIG


# 根据集群设置API服务器地址
if [ "$CLUSTER" = "tot" ]; then
    server="https://192.168.1.212:6443"
else
    server="https://hoper.xyz:6443"
fi

kubeconfig="--kubeconfig=/root/.kube/config"
# 配置kubectl
../kubectl_auth.sh $server k8s


# 删除旧资源（如果是job或cronjob）
if [ "$DEPLOY_KIND" = "job" ] || [ "$DEPLOY_KIND" = "cronjob" ]; then
    kubectl $kubeconfig delete --ignore-not-found -f "$deploy_path"
fi

# 应用新资源
kubectl $kubeconfig apply -f "$deploy_path"

echo "部署完成!"
