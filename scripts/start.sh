#!/bin/bash

# 启动脚本 - 适用于麒麟操作系统
# 此脚本用于启动Star LLM Backend服务

# 定义变量
PROJECT_NAME="star_llm_backend"
BASE_DIR=$(dirname $(dirname $(readlink -f "$0")))
BIN_DIR="${BASE_DIR}/bin"
CONF_DIR="${BASE_DIR}/conf"
LOG_DIR="${BASE_DIR}/logs"
PID_FILE="${BASE_DIR}/${PROJECT_NAME}.pid"

# 创建日志目录
mkdir -p ${LOG_DIR}

# 显示启动信息
echo "=== 启动 ${PROJECT_NAME} 服务 ==="
echo "工作目录: ${BASE_DIR}"

# 检查服务是否已经运行
if [ -f "${PID_FILE}" ]; then
    PID=$(cat ${PID_FILE})
    if ps -p ${PID} > /dev/null; then
        echo "服务已经在运行中，PID: ${PID}"
        exit 1
    else
        echo "发现过期的PID文件，将被删除"
        rm -f ${PID_FILE}
    fi
fi

# 检查可执行文件是否存在
if [ ! -f "${BIN_DIR}/${PROJECT_NAME}" ]; then
    echo "错误: 可执行文件不存在 ${BIN_DIR}/${PROJECT_NAME}"
    exit 1
fi

# 检查配置文件是否存在
if [ ! -f "${CONF_DIR}/config.yml" ]; then
    echo "错误: 配置文件不存在 ${CONF_DIR}/config.yml"
    exit 1
fi

# 启动服务
echo "正在启动服务..."
cd ${BASE_DIR}
nohup ${BIN_DIR}/${PROJECT_NAME} > ${LOG_DIR}/server.log 2>&1 &

# 保存PID
echo $! > ${PID_FILE}

# 检查服务是否成功启动
sleep 2
if ps -p $(cat ${PID_FILE}) > /dev/null; then
    echo "服务启动成功，PID: $(cat ${PID_FILE})"
    echo "日志文件: ${LOG_DIR}/server.log"
    echo "使用以下命令查看日志: tail -f ${LOG_DIR}/server.log"
    echo "使用以下命令停止服务: ${BASE_DIR}/scripts/stop.sh"
    exit 0
else
    echo "服务启动失败，请检查日志文件: ${LOG_DIR}/server.log"
    exit 1
fi