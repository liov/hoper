#!/bin/bash

# 获取配置参数
REPO="${DRONE_REPO:-}"
COMMIT="${DRONE_COMMIT:-}"
COMMIT_TAG="${DRONE_TAG:-}"
COMMIT_LINK="${DRONE_COMMIT_LINK:-}"
COMMIT_REF="${DRONE_COMMIT_REF:-}"
COMMIT_MESSAGE="${DRONE_COMMIT_MESSAGE:-}"
COMMIT_BRANCH="${DRONE_COMMIT_BRANCH:-}"
COMMIT_AUTHOR="${DRONE_COMMIT_AUTHOR:-}"
COMMIT_AUTHOR_NAME="${DRONE_COMMIT_AUTHOR_NAME:-}"
DING_TOKEN="${PLUGIN_DING_TOKEN:-}"
DING_SECRET="${PLUGIN_DING_SECRET:-}"
BUILD_LINK="${DRONE_BUILD_LINK:-}"

# 检查是否设置了钉钉 Token
if [ -z "$DING_TOKEN" ]; then
    exit 0
fi

# 构建消息内容
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
MESSAGE_BODY="
# 发布通知
### 项目: $REPO
### 操作人: $COMMIT_AUTHOR
### 参考: $COMMIT_REF
### 分支: $COMMIT_BRANCH
### 标签: $COMMIT_TAG
### 时间: $TIMESTAMP
### 提交: $COMMIT
### 提交信息: $COMMIT_MESSAGE
### 发布详情: $BUILD_LINK
"

# 发送钉钉消息的函数
send_dingtalk_message() {
    local url="$1"
    local message_body="$2"

    # 发送请求
    curl -X POST \
        -H "Content-Type: application/json" \
        -d "{\"msgtype\": \"markdown\", \"markdown\": {\"title\": \"发布通知\", \"text\": \"$message_body\"}}" \
        "$url"
}

# 计算带签名的 URL
if [ -n "$DING_SECRET" ]; then
    # 计算 timestamp 和 sign
    TIMESTAMP=$(date +%s)000
    STRING_TO_SIGN=$(printf '%s\n%s' "$TIMESTAMP" "$DING_SECRET" | openssl dgst -sha256 -hmac "$DING_SECRET" | cut -d' ' -f2)

    REQUEST_URL="https://oapi.dingtalk.com/robot/send?access_token=$DING_TOKEN&timestamp=$TIMESTAMP&sign=$SIGN"
else
    REQUEST_URL="https://oapi.dingtalk.com/robot/send?access_token=$DING_TOKEN"
fi

# 发送消息
send_dingtalk_message "$REQUEST_URL" "$MESSAGE_BODY"
# test notify
# docker run --rm -e PLUGIN_DING_TOKEN=xxx -e PLUGIN_DING_SECRET=xxx jybl/notify