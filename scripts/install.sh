#!/bin/bash

# 安装脚本 - 适用于麒麟操作系统
# 此脚本用于安装Star LLM Backend服务

set -e

# 定义变量
PROJECT_NAME="star_llm_backend"
INSTALL_DIR="/opt/${PROJECT_NAME}"
TAR_FILE="${PROJECT_NAME}.tar.gz"

# 显示安装信息
echo "=== 开始安装 ${PROJECT_NAME} ==="
echo "安装目录: ${INSTALL_DIR}"

# 检查是否以root用户运行
if [ "$(id -u)" -ne 0 ]; then
    echo "错误: 请使用root用户或sudo运行此脚本"
    exit 1
fi

# 检查安装包是否存在
if [ ! -f "${TAR_FILE}" ]; then
    echo "错误: 安装包不存在 ${TAR_FILE}"
    exit 1
fi

# 创建安装目录
echo "创建安装目录..."
mkdir -p ${INSTALL_DIR}

# 解压安装包
echo "解压安装包..."
tar -xzf ${TAR_FILE} -C ${INSTALL_DIR} --strip-components=1

# 设置权限
echo "设置权限..."
chmod +x ${INSTALL_DIR}/bin/${PROJECT_NAME}
chmod +x ${INSTALL_DIR}/scripts/*.sh

# 创建日志目录
mkdir -p ${INSTALL_DIR}/logs

# 配置服务
echo "配置服务..."
cat > /etc/systemd/system/${PROJECT_NAME}.service << EOF
[Unit]
Description=Star LLM Backend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/scripts/start.sh
ExecStop=${INSTALL_DIR}/scripts/stop.sh
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 重新加载systemd配置
systemctl daemon-reload

# 启用服务开机自启
systemctl enable ${PROJECT_NAME}.service

echo "=== ${PROJECT_NAME} 安装完成 ==="
echo "使用以下命令启动服务: systemctl start ${PROJECT_NAME}"
echo "使用以下命令停止服务: systemctl stop ${PROJECT_NAME}"
echo "使用以下命令查看服务状态: systemctl status ${PROJECT_NAME}"
echo "使用以下命令查看服务日志: journalctl -u ${PROJECT_NAME} -f"