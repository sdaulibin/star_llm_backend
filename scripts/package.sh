#!/bin/bash

# 打包脚本 - 适用于麒麟操作系统
# 此脚本用于编译和打包Star LLM Backend项目

set -e

# 定义变量
PROJECT_NAME="star_llm_backend"
PACKAGE_DIR="./dist"
BIN_DIR="${PACKAGE_DIR}/bin"
CONF_DIR="${PACKAGE_DIR}/conf"
SCRIPTS_DIR="${PACKAGE_DIR}/scripts"
DOC_DIR="${PACKAGE_DIR}/doc"

# 显示打包信息
echo "=== 开始打包 ${PROJECT_NAME} ==="
echo "打包目录: ${PACKAGE_DIR}"

# 创建打包目录结构
rm -rf ${PACKAGE_DIR}
mkdir -p ${BIN_DIR} ${CONF_DIR} ${SCRIPTS_DIR} ${DOC_DIR}

# 设置Go环境变量
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 编译项目
echo "正在编译项目..."
go build -o ${BIN_DIR}/${PROJECT_NAME} ./cmd/api/main.go
if [ $? -ne 0 ]; then
    echo "编译失败，请检查错误信息"
    exit 1
fi
echo "编译完成: ${BIN_DIR}/${PROJECT_NAME}"

# 复制配置文件
echo "正在复制配置文件..."
cp -r ./conf/* ${CONF_DIR}/
echo "配置文件已复制到 ${CONF_DIR}"

# 复制启动脚本
echo "正在复制启动脚本..."
cp ./scripts/start.sh ${SCRIPTS_DIR}/
chmod +x ${SCRIPTS_DIR}/start.sh
echo "启动脚本已复制到 ${SCRIPTS_DIR}"

# 复制启动脚本
echo "正在复制停止脚本..."
cp ./scripts/stop.sh ${SCRIPTS_DIR}/
chmod +x ${SCRIPTS_DIR}/stop.sh
echo "启动脚本已复制到 ${SCRIPTS_DIR}"

# 复制文档
echo "正在复制文档..."
cp ./README.md ${DOC_DIR}/
cp -r ./doc/* ${DOC_DIR}/
echo "文档已复制到 ${DOC_DIR}"

# 创建压缩包
echo "正在创建压缩包..."
cd ${PACKAGE_DIR}/../
tar -czvf ${PROJECT_NAME}.tar.gz dist/
echo "压缩包已创建: ${PROJECT_NAME}.tar.gz"

echo "=== ${PROJECT_NAME} 打包完成 ==="
echo "打包文件位置: ${PROJECT_NAME}.tar.gz"