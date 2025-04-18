#!/bin/bash

# 停止脚本 - 适用于麒麟操作系统
# 此脚本用于停止Star LLM Backend服务

# 定义变量
PROJECT_NAME="star_llm_backend"
BASE_DIR=$(dirname $(dirname $(readlink -f "$0")))
PID_FILE="${BASE_DIR}/${PROJECT_NAME}.pid"

# 显示停止信息
echo "=== 停止 ${PROJECT_NAME} 服务 ==="

# 检查PID文件是否存在
if [ ! -f "${PID_FILE}" ]; then
    echo "服务未运行或PID文件不存在"
    exit 0
fi

# 读取PID并检查服务是否运行
PID=$(cat ${PID_FILE})
if ! ps -p ${PID} > /dev/null; then
    echo "服务未运行，删除过期的PID文件"
    rm -f ${PID_FILE}
    exit 0
fi

# 停止服务
echo "正在停止服务，PID: ${PID}"
kill ${PID}

# 等待服务停止
echo "等待服务停止..."
for i in {1..30}; do
    if ! ps -p ${PID} > /dev/null; then
        echo "服务已成功停止"
        rm -f ${PID_FILE}
        exit 0
    fi
    sleep 1
done

# 如果服务未能正常停止，尝试强制终止
echo "服务未能在30秒内停止，尝试强制终止"
kill -9 ${PID}
sleep 2

if ! ps -p ${PID} > /dev/null; then
    echo "服务已被强制终止"
    rm -f ${PID_FILE}
    exit 0
else
    echo "无法终止服务，请手动检查进程 ${PID}"
    exit 1
fi