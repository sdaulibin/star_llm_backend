# Star LLM Backend 部署指南（麒麟操作系统）

本文档提供了在麒麟操作系统上打包、安装和部署 Star LLM Backend 服务的详细步骤。

## 1. 打包项目

在开发环境中执行以下步骤来打包项目：

```bash
# 进入项目目录
cd /path/to/star_llm_backend_n

# 赋予打包脚本执行权限
chmod +x scripts/package.sh

# 执行打包脚本
./scripts/package.sh
```

打包完成后，将在项目根目录生成 `star_llm_backend.tar.gz` 文件，这是可以部署到麒麟操作系统的完整安装包。

## 2. 部署到麒麟操作系统

### 2.1 准备工作

- 确保麒麟操作系统已安装
- 确保有足够的磁盘空间（建议至少 200MB）
- 确保有 root 或 sudo 权限

### 2.2 安装步骤

1. 将生成的 `star_llm_backend.tar.gz` 文件复制到麒麟操作系统服务器上

2. 赋予安装脚本执行权限并运行：

```bash
# 解压安装包
tar -xzf star_llm_backend.tar.gz

# 进入解压后的目录
cd dist

# 赋予安装脚本执行权限
chmod +x scripts/install.sh

# 以root用户或sudo运行安装脚本
sudo ./scripts/install.sh
```

### 2.3 配置服务

安装完成后，需要修改配置文件以适应您的环境：

```bash
# 编辑配置文件
sudo vi /opt/star_llm_backend/conf/config.yml
```

主要配置项：
- `api.base_url`: Dify API 的基础 URL
- `api.key`: Dify API 的认证密钥
- `server.port`: 服务监听端口
- `database`: 数据库连接信息

## 3. 服务管理

### 3.1 启动服务

```bash
# 使用systemd启动服务
sudo systemctl start star_llm_backend

# 或直接使用启动脚本
sudo /opt/star_llm_backend/scripts/start.sh
```

### 3.2 停止服务

```bash
# 使用systemd停止服务
sudo systemctl stop star_llm_backend

# 或直接使用停止脚本
sudo /opt/star_llm_backend/scripts/stop.sh
```

### 3.3 查看服务状态

```bash
# 查看服务状态
sudo systemctl status star_llm_backend

# 查看服务日志
sudo journalctl -u star_llm_backend -f
# 或直接查看日志文件
sudo tail -f /opt/star_llm_backend/logs/server.log
```

## 4. 故障排除

如果服务无法正常启动或运行，请检查：

1. 配置文件是否正确
2. 数据库连接是否可用
3. 查看日志文件中的错误信息
4. 确保服务器防火墙允许指定端口的访问

## 5. 安全建议

1. 定期更新服务和依赖项
2. 不要使用默认的API密钥
3. 限制对配置文件的访问权限
4. 考虑使用HTTPS进行API通信

## 6. 联系支持

如有任何问题或需要技术支持，请联系系统管理员或开发团队。