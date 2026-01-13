#!/bin/bash

# 检查环境变量
if [ -z "$DINGTALK_TOKEN" ]; then
  echo "Error: $DINGTALK_TOKEN 环境变量未设置"
  exit 1
fi


# 钉钉机器人 Webhook 地址（从环境变量或直接替换）
WEBHOOK_URL="https://oapi.dingtalk.com/robot/send?access_token=${$DINGTALK_TOKEN}"


# 使用方法提示
usage() {
    echo "Usage: $0 [-t text|markdown] [-c 'content'] [-a '@手机号1,@手机号2'] [-A] [-T '标题']"
    echo "Example:"
    echo "  $0 -t text -c '服务已上线' -a '13812345678'"
    echo "  $0 -t markdown -T '告警' -c '### 问题\nCPU负载过高' -A"
    exit 1
}

# 解析参数
while getopts "t:c:a:T:A" opt; do
    case $opt in
        t) MSG_TYPE="$OPTARG" ;;
        c) CONTENT="$OPTARG" ;;
        a) AT_MOBILES="$OPTARG" ;;
        T) TITLE="$OPTARG" ;;
        A) IS_AT_ALL=true ;;
        *) usage ;;
    esac
done

# 检查必要参数
if [[ -z "$MSG_TYPE" || -z "$CONTENT" ]]; then
    usage
fi

# 构造 JSON 数据
generate_json() {
    local at_json=""
    if [[ -n "$AT_MOBILES" || "$IS_AT_ALL" = true ]]; then
        at_json="\"at\": {"
        if [[ -n "$AT_MOBILES" ]]; then
            at_json+="\"atMobiles\": [$(echo "$AT_MOBILES" | sed 's/,/", "/g')],"
        fi
        at_json+="\"isAtAll\": ${IS_AT_ALL:-false}}"
    fi

    case "$MSG_TYPE" in
        text)
            local json_data=$(cat <<EOF
{
    "msgtype": "text",
    "text": {
        "content": "$CONTENT"
    }$( [[ -n "$at_json" ]] && echo ",$at_json" )
}
EOF
)
            echo "$json_data"
            ;;
        markdown)
            local json_data=$(cat <<EOF
{
    "msgtype": "markdown",
    "markdown": {
        "title": "${TITLE:-通知}",
        "text": "$CONTENT"
    }$( [[ -n "$at_json" ]] && echo ",$at_json" )
}
EOF
)
            echo "$json_data"
            ;;
        *)
            echo "错误：不支持的消息类型 '$MSG_TYPE'"
            exit 1
            ;;
    esac
}

# 发送消息（带重试逻辑）
send_message() {
    local json="$1"
    local max_retries=3
    local retry_count=0
    local success=false

    while [[ $retry_count -lt $max_retries && "$success" = false ]]; do
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$WEBHOOK_URL" \
            -H 'Content-Type: application/json' \
            -d "$json" 2>/dev/null)

        local http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d':' -f2)
        local result=$(echo "$response" | sed '/HTTP_CODE:/d')

        if [[ "$http_code" == "200" ]]; then
            echo "✅ 消息发送成功"
            success=true
        else
            echo "⚠️ 尝试 $((retry_count+1))/$max_retries 失败: $result"
            sleep 2
            ((retry_count++))
        fi
    done

    if [[ "$success" = false ]]; then
        echo "❌ 消息发送失败，请检查网络或Webhook配置"
        exit 1
    fi
}

# 主流程
main() {
    echo "📤 正在发送钉钉通知..."
    local json=$(generate_json)
    send_message "$json"
}

main