# 阿里云 ECS 部署指南

本文档提供将 xiaohongshu-mcp 服务部署到阿里云 ECS 的完整操作指南。

## 📋 部署任务清单

- [ ] 1. 准备阿里云 ECS 服务器
- [ ] 2. 配置 ECS 基础环境
- [ ] 3. 安装 Go 运行环境
- [ ] 4. 安装 Chrome/Chromium 浏览器依赖
- [ ] 5. 部署 xiaohongshu-mcp 代码
- [ ] 6. 配置持久化存储
- [ ] 7. 完成小红书账号登录
- [ ] 8. 配置 systemd 服务
- [ ] 9. 配置防火墙和安全组规则
- [ ] 10. 配置反向代理（可选）
- [ ] 11. 配置日志管理和监控
- [ ] 12. 测试 MCP 服务

---

## 📐 服务器配置要求

### 最低配置
- **操作系统**: Ubuntu 22.04 LTS 或 CentOS 8+
- **CPU**: 2 核
- **内存**: 4GB
- **带宽**: 3Mbps
- **磁盘**: 系统盘 40GB + 数据盘 20GB

### 推荐配置
- **CPU**: 2-4 核
- **内存**: 8GB
- **带宽**: 5Mbps
- **磁盘**: 系统盘 40GB + 数据盘 50GB

---

## 1️⃣ 准备阿里云 ECS 服务器

### 1.1 创建 ECS 实例

1. 登录 [阿里云控制台](https://ecs.console.aliyun.com/)
2. 点击"创建实例"
3. 选择配置：
   - **地域**: 根据需求选择（建议选择离用户较近的地域）
   - **实例规格**: 2核4GB 或以上
   - **镜像**: Ubuntu 22.04 64位 或 CentOS 8.x 64位
   - **网络**: 选择或创建 VPC 和安全组
   - **公网 IP**: 分配公网 IP（或使用弹性 IP）
   - **带宽**: 按固定带宽，至少 3Mbps

### 1.2 配置密钥对或密码

- **密钥对方式**（推荐）: 创建或选择现有密钥对，下载 .pem 文件
- **密码方式**: 设置 root 用户密码

### 1.3 连接到服务器

**使用密钥对连接：**

```bash
# 设置密钥文件权限
chmod 400 your-key.pem

# SSH 连接
ssh -i your-key.pem root@<ECS公网IP>
```

**使用密码连接：**

```bash
ssh root@<ECS公网IP>
# 输入密码
```

---

## 2️⃣ 配置 ECS 基础环境

### 2.1 更新系统

**Ubuntu 系统：**

```bash
# 更新软件包列表
sudo apt update

# 升级已安装的软件包
sudo apt upgrade -y

# 安装基础工具
sudo apt install -y wget curl git build-essential ca-certificates vim htop
```

**CentOS 系统：**

```bash
# 更新系统
sudo yum update -y

# 安装基础工具
sudo yum install -y wget curl git gcc make ca-certificates vim htop
```

### 2.2 配置时区

```bash
# 查看当前时区
timedatectl

# 设置为中国时区
sudo timedatectl set-timezone Asia/Shanghai

# 验证时区
date
```

### 2.3 配置主机名（可选）

```bash
# 设置主机名
sudo hostnamectl set-hostname xiaohongshu-mcp

# 编辑 /etc/hosts
sudo vim /etc/hosts
# 添加: 127.0.0.1 xiaohongshu-mcp
```

---

## 3️⃣ 安装 Go 运行环境

### 3.1 下载并安装 Go

```bash
# 下载 Go 1.24.0 (或访问 https://go.dev/dl/ 获取最新版本)
cd /tmp
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# 解压到 /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# 清理安装包
rm go1.24.0.linux-amd64.tar.gz
```

### 3.2 配置环境变量

```bash
# 配置当前用户环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc

# 使配置生效
source ~/.bashrc

# 验证安装
go version
# 输出: go version go1.24.0 linux/amd64
```

### 3.3 配置 Go 代理（国内访问加速）

```bash
# 配置 GOPROXY（三选一）

# 1. 七牛 CDN（推荐）
go env -w GOPROXY=https://goproxy.cn,direct

# 2. 阿里云
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 3. 官方
go env -w GOPROXY=https://goproxy.io,direct

# 验证配置
go env | grep GOPROXY
```

---

## 4️⃣ 安装 Chrome/Chromium 浏览器依赖

xiaohongshu-mcp 依赖无头浏览器运行，需要安装浏览器及相关依赖。

### 4.1 Ubuntu 系统安装

```bash
# 安装 Chromium 和依赖
sudo apt install -y \
    chromium-browser \
    fonts-wqy-zenhei \
    fonts-wqy-microhei \
    xvfb \
    libx11-dev \
    libxcomposite-dev \
    libxcursor-dev \
    libxdamage-dev \
    libxext-dev \
    libxfixes-dev \
    libxi-dev \
    libxrandr-dev \
    libxrender-dev \
    libxss-dev \
    libxtst-dev \
    libgbm-dev \
    libnss3 \
    libatk-bridge2.0-0 \
    libgtk-3-0 \
    libasound2

# 验证安装
chromium-browser --version
```

### 4.2 CentOS 系统安装

```bash
# 启用 EPEL 仓库
sudo yum install -y epel-release

# 安装 Chromium 和依赖
sudo yum install -y \
    chromium \
    wqy-zenhei-fonts \
    libX11 \
    nss \
    atk \
    at-spi2-atk \
    libXcomposite \
    libXcursor \
    libXdamage \
    libXext \
    libXi \
    libXrandr \
    libXrender \
    libXtst \
    libgbm \
    gtk3 \
    alsa-lib

# 验证安装
chromium-browser --version
```

### 4.3 中文字体支持

```bash
# Ubuntu
sudo apt install -y fonts-noto-cjk fonts-noto-cjk-extra

# CentOS
sudo yum install -y google-noto-sans-cjk-fonts google-noto-serif-cjk-fonts
```

---

## 5️⃣ 部署 xiaohongshu-mcp 代码

### 5.1 创建工作目录

```bash
# 创建项目目录
mkdir -p ~/xiaohongshu-mcp
cd ~/xiaohongshu-mcp
```

### 5.2 方式一：从源码编译（推荐）

```bash
# 克隆代码仓库
git clone https://github.com/xpzouying/xiaohongshu-mcp.git .

# 查看可用版本
git tag

# 切换到稳定版本（可选）
# git checkout v1.x.x

# 下载依赖
go mod download

# 编译主程序
go build -o xiaohongshu-mcp .

# 编译登录工具
go build -o xiaohongshu-login cmd/login/main.go

# 设置可执行权限
chmod +x xiaohongshu-mcp xiaohongshu-login

# 验证编译结果
ls -lh xiaohongshu-mcp xiaohongshu-login
./xiaohongshu-mcp -h
```

### 5.3 方式二：下载预编译二进制

```bash
# 设置版本号（访问 https://github.com/xpzouying/xiaohongshu-mcp/releases 查看最新版本）
VERSION="v1.0.0"  # 替换为实际版本号

# 下载主程序
wget https://github.com/xpzouying/xiaohongshu-mcp/releases/download/${VERSION}/xiaohongshu-mcp-linux-amd64

# 下载登录工具
wget https://github.com/xpzouying/xiaohongshu-mcp/releases/download/${VERSION}/xiaohongshu-login-linux-amd64

# 重命名
mv xiaohongshu-mcp-linux-amd64 xiaohongshu-mcp
mv xiaohongshu-login-linux-amd64 xiaohongshu-login

# 设置可执行权限
chmod +x xiaohongshu-mcp xiaohongshu-login

# 验证
./xiaohongshu-mcp -h
```

---

## 6️⃣ 配置持久化存储

### 6.1 创建数据目录

```bash
cd ~/xiaohongshu-mcp

# 创建必要的目录
mkdir -p data      # 存储 cookies
mkdir -p images    # 存储发布的图片
mkdir -p logs      # 存储日志文件
mkdir -p videos    # 存储视频文件（可选）

# 设置权限
chmod 755 data images logs
```

### 6.2 目录结构

```
~/xiaohongshu-mcp/
├── xiaohongshu-mcp          # 主程序
├── xiaohongshu-login        # 登录工具
├── data/
│   └── cookies.json         # 登录状态（自动生成）
├── images/                  # 图片目录
├── videos/                  # 视频目录
└── logs/
    ├── service.log          # 服务日志
    └── error.log            # 错误日志
```

---

## 7️⃣ 完成小红书账号登录

### 7.1 方式一：本地登录后上传（推荐）

**适用场景**: 服务器无 GUI 界面或无法运行图形化程序

**步骤：**

1. **在本地电脑运行登录工具：**

```bash
# 在本地电脑上
./xiaohongshu-login
```

2. **完成登录操作**
   - 程序会打开浏览器窗口
   - 手动登录小红书账号
   - 登录成功后，cookies 会保存到 `data/cookies.json`

3. **上传 cookies 文件到服务器：**

```bash
# 在本地电脑上执行
scp data/cookies.json root@<ECS公网IP>:~/xiaohongshu-mcp/data/

# 如果使用密钥对
scp -i your-key.pem data/cookies.json root@<ECS公网IP>:~/xiaohongshu-mcp/data/
```

4. **在服务器上验证文件：**

```bash
# SSH 到服务器
cd ~/xiaohongshu-mcp
ls -lh data/cookies.json
cat data/cookies.json  # 检查内容格式
```

### 7.2 方式二：服务器直接登录

**适用场景**: 服务器支持 GUI 或 X11 转发

**使用 X11 转发：**

```bash
# 从本地连接时启用 X11 转发
ssh -X root@<ECS公网IP>

# 在服务器上运行登录工具
cd ~/xiaohongshu-mcp
./xiaohongshu-login
```

**使用 VNC：**

如果服务器安装了 VNC 服务，可以通过 VNC 连接后运行登录工具。

### 7.3 验证登录状态

```bash
# 启动 MCP 服务（临时测试）
cd ~/xiaohongshu-mcp
./xiaohongshu-mcp -headless=true &

# 等待服务启动（约5秒）
sleep 5

# 测试登录状态
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"check_login_status","arguments":{}},"id":1}'

# 停止测试服务
pkill xiaohongshu-mcp
```

---

## 8️⃣ 配置 systemd 服务

将 xiaohongshu-mcp 配置为系统服务，实现开机自启和自动重启。

### 8.1 创建 systemd 服务文件

```bash
sudo tee /etc/systemd/system/xiaohongshu-mcp.service <<'EOF'
[Unit]
Description=Xiaohongshu MCP Service
Documentation=https://github.com/xpzouying/xiaohongshu-mcp
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/xiaohongshu-mcp
ExecStart=/root/xiaohongshu-mcp/xiaohongshu-mcp -headless=true
Restart=always
RestartSec=10
StandardOutput=append:/root/xiaohongshu-mcp/logs/service.log
StandardError=append:/root/xiaohongshu-mcp/logs/error.log

# 环境变量
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
Environment="DISPLAY=:99"

# 资源限制
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOF
```

### 8.2 重载并启动服务

```bash
# 重载 systemd 配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start xiaohongshu-mcp

# 查看服务状态
sudo systemctl status xiaohongshu-mcp

# 设置开机自启
sudo systemctl enable xiaohongshu-mcp
```

### 8.3 服务管理命令

```bash
# 启动服务
sudo systemctl start xiaohongshu-mcp

# 停止服务
sudo systemctl stop xiaohongshu-mcp

# 重启服务
sudo systemctl restart xiaohongshu-mcp

# 查看服务状态
sudo systemctl status xiaohongshu-mcp

# 查看实时日志
journalctl -u xiaohongshu-mcp -f

# 查看最近100行日志
journalctl -u xiaohongshu-mcp -n 100

# 查看服务日志文件
tail -f ~/xiaohongshu-mcp/logs/service.log
tail -f ~/xiaohongshu-mcp/logs/error.log
```

---

## 9️⃣ 配置防火墙和安全组规则

### 9.1 阿里云安全组配置

**操作步骤：**

1. 登录 [阿里云 ECS 控制台](https://ecs.console.aliyun.com/)
2. 选择实例 → 点击实例 ID
3. 进入"安全组" → 点击"配置规则"
4. 添加"入方向"规则：

**配置参数：**

| 字段       | 值                                        | 说明                         |
| ---------- | ----------------------------------------- | ---------------------------- |
| 授权策略   | 允许                                      |                              |
| 优先级     | 1                                         | 数字越小优先级越高           |
| 协议类型   | 自定义 TCP                                |                              |
| 端口范围   | 18060/18060                               | MCP 服务端口                 |
| 授权对象   | 见下方说明                                | 根据安全需求选择             |
| 描述       | xiaohongshu-mcp service                   |                              |

**授权对象选择：**

- **特定 IP 访问**（推荐）: `你的IP地址/32`
  - 例如: `123.123.123.123/32`
  - 最安全，仅允许特定 IP 访问

- **IP 段访问**: `192.168.1.0/24`
  - 允许特定网段访问

- **所有 IP 访问**: `0.0.0.0/0`
  - ⚠️ 不推荐，存在安全风险
  - 如需公开访问，建议配合 Nginx + Basic Auth

### 9.2 服务器防火墙配置

**Ubuntu 系统（使用 UFW）：**

```bash
# 检查 UFW 状态
sudo ufw status

# 允许 SSH（必须！防止锁死）
sudo ufw allow 22/tcp

# 允许 MCP 服务端口
sudo ufw allow 18060/tcp

# 启用防火墙
sudo ufw enable

# 查看规则
sudo ufw status numbered

# 如需删除规则
# sudo ufw delete <规则编号>
```

**CentOS 系统（使用 firewalld）：**

```bash
# 检查 firewalld 状态
sudo systemctl status firewalld

# 如未安装，先安装
sudo yum install -y firewalld
sudo systemctl start firewalld
sudo systemctl enable firewalld

# 允许 SSH
sudo firewall-cmd --permanent --add-port=22/tcp

# 允许 MCP 服务端口
sudo firewall-cmd --permanent --add-port=18060/tcp

# 重载配置
sudo firewall-cmd --reload

# 查看开放的端口
sudo firewall-cmd --list-ports

# 查看所有规则
sudo firewall-cmd --list-all
```

### 9.3 验证端口开放

```bash
# 在服务器上检查端口监听
sudo netstat -tulpn | grep 18060
# 或
sudo ss -tulpn | grep 18060

# 从本地测试连接
telnet <ECS公网IP> 18060
# 或
nc -zv <ECS公网IP> 18060
```

---

## 🔟 配置反向代理（可选）

使用 Nginx 作为反向代理，可以实现：
- 域名访问
- HTTPS 加密
- 访问控制
- 负载均衡

### 10.1 安装 Nginx

**Ubuntu 系统：**

```bash
sudo apt install -y nginx
```

**CentOS 系统：**

```bash
sudo yum install -y nginx
```

### 10.2 配置 Nginx

**方式一：HTTP 访问**

```bash
sudo tee /etc/nginx/sites-available/xiaohongshu-mcp <<'EOF'
# HTTP 配置
server {
    listen 80;
    server_name mcp.yourdomain.com;  # 替换为你的域名

    # 日志配置
    access_log /var/log/nginx/xiaohongshu-mcp.access.log;
    error_log /var/log/nginx/xiaohongshu-mcp.error.log;

    # MCP 服务代理
    location /mcp {
        proxy_pass http://localhost:18060;
        proxy_http_version 1.1;

        # WebSocket 支持
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';

        # 请求头
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 300s;
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;

        # 缓存配置
        proxy_cache_bypass $http_upgrade;
    }

    # 健康检查接口（可选）
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF
```

**方式二：HTTPS 访问（推荐）**

```bash
sudo tee /etc/nginx/sites-available/xiaohongshu-mcp <<'EOF'
# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name mcp.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name mcp.yourdomain.com;

    # SSL 证书配置
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    # SSL 安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # 日志配置
    access_log /var/log/nginx/xiaohongshu-mcp.access.log;
    error_log /var/log/nginx/xiaohongshu-mcp.error.log;

    # MCP 服务代理
    location /mcp {
        proxy_pass http://localhost:18060;
        proxy_http_version 1.1;

        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_connect_timeout 300s;
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF
```

### 10.3 配置 Basic Auth（可选，增强安全性）

```bash
# 安装密码工具
sudo apt install -y apache2-utils  # Ubuntu
# 或
sudo yum install -y httpd-tools     # CentOS

# 创建密码文件（用户名: admin）
sudo htpasswd -c /etc/nginx/.htpasswd admin
# 输入密码

# 修改 Nginx 配置，在 location /mcp 块中添加：
    auth_basic "Restricted Access";
    auth_basic_user_file /etc/nginx/.htpasswd;
```

### 10.4 启用配置

**Ubuntu 系统：**

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/xiaohongshu-mcp /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx

# 查看状态
sudo systemctl status nginx
```

**CentOS 系统：**

```bash
# CentOS 直接编辑主配置或在 /etc/nginx/conf.d/ 目录创建配置

# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx
```

### 10.5 配置 SSL 证书

**使用 Let's Encrypt 免费证书（推荐）：**

```bash
# 安装 Certbot
sudo apt install -y certbot python3-certbot-nginx  # Ubuntu
# 或
sudo yum install -y certbot python3-certbot-nginx  # CentOS

# 申请证书
sudo certbot --nginx -d mcp.yourdomain.com

# 自动续期（Certbot 会自动添加 cron 任务）
sudo certbot renew --dry-run
```

### 10.6 开放 HTTP/HTTPS 端口

```bash
# Ubuntu UFW
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# CentOS firewalld
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload
```

---

## 1️⃣1️⃣ 配置日志管理和监控

### 11.1 配置日志轮转

防止日志文件无限增长占满磁盘。

```bash
# 创建日志轮转配置
sudo tee /etc/logrotate.d/xiaohongshu-mcp <<'EOF'
/root/xiaohongshu-mcp/logs/*.log {
    daily                  # 每天轮转
    rotate 7               # 保留7天
    compress               # 压缩旧日志
    delaycompress          # 延迟压缩（保留最近一天的日志不压缩）
    missingok              # 日志文件不存在不报错
    notifempty             # 空文件不轮转
    create 0644 root root  # 创建新文件的权限
    sharedscripts
    postrotate
        systemctl reload xiaohongshu-mcp > /dev/null 2>&1 || true
    endscript
}
EOF

# 测试配置
sudo logrotate -d /etc/logrotate.d/xiaohongshu-mcp

# 手动执行轮转（测试）
sudo logrotate -f /etc/logrotate.d/xiaohongshu-mcp
```

### 11.2 配置服务监控脚本

**创建监控脚本：**

```bash
tee ~/xiaohongshu-mcp/monitor.sh <<'EOF'
#!/bin/bash

# 配置
LOG_FILE="/root/xiaohongshu-mcp/logs/monitor.log"
SERVICE_NAME="xiaohongshu-mcp"

# 日志函数
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
}

# 检查服务状态
if ! systemctl is-active --quiet $SERVICE_NAME; then
    log "服务已停止，正在重启..."
    systemctl restart $SERVICE_NAME

    # 等待5秒后检查是否启动成功
    sleep 5
    if systemctl is-active --quiet $SERVICE_NAME; then
        log "服务重启成功"
    else
        log "服务重启失败，请检查"
    fi
else
    # 检查端口监听
    if ! netstat -tulpn | grep -q ":18060"; then
        log "服务运行中但端口未监听，正在重启..."
        systemctl restart $SERVICE_NAME
    fi
fi

# 检查磁盘空间
DISK_USAGE=$(df -h /root | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 80 ]; then
    log "警告: 磁盘使用率过高 - $DISK_USAGE%"
fi

# 检查内存使用
MEM_USAGE=$(free | awk 'NR==2 {printf "%.0f", $3/$2*100}')
if [ "$MEM_USAGE" -gt 90 ]; then
    log "警告: 内存使用率过高 - $MEM_USAGE%"
fi
EOF

# 设置可执行权限
chmod +x ~/xiaohongshu-mcp/monitor.sh

# 测试脚本
~/xiaohongshu-mcp/monitor.sh
cat ~/xiaohongshu-mcp/logs/monitor.log
```

**配置定时任务：**

```bash
# 编辑 crontab
crontab -e

# 添加以下行（每5分钟检查一次）
*/5 * * * * /root/xiaohongshu-mcp/monitor.sh

# 查看已配置的定时任务
crontab -l
```

### 11.3 配置阿里云监控（可选）

**安装阿里云监控插件：**

```bash
# 下载安装脚本
wget https://cloudmonitor-agent.oss-cn-hangzhou.aliyuncs.com/release/install.sh

# 执行安装
sudo bash install.sh

# 查看状态
sudo systemctl status argusagent
```

**在阿里云控制台配置告警规则：**

1. 登录 [云监控控制台](https://cloudmonitor.console.aliyun.com/)
2. 创建告警规则：
   - **CPU 使用率** > 80%
   - **内存使用率** > 85%
   - **磁盘使用率** > 80%
   - **进程存活监控**: xiaohongshu-mcp

### 11.4 实时监控命令

```bash
# 查看服务状态
sudo systemctl status xiaohongshu-mcp

# 实时查看服务日志
journalctl -u xiaohongshu-mcp -f

# 实时查看应用日志
tail -f ~/xiaohongshu-mcp/logs/service.log

# 查看系统资源使用
htop

# 查看进程详情
ps aux | grep xiaohongshu-mcp

# 查看端口监听
sudo netstat -tulpn | grep 18060
```

---

## 1️⃣2️⃣ 测试 MCP 服务

### 12.1 本地服务测试

**测试 1: 检查服务启动**

```bash
# 查看服务状态
sudo systemctl status xiaohongshu-mcp

# 查看端口监听
sudo netstat -tulpn | grep 18060

# 查看进程
ps aux | grep xiaohongshu-mcp
```

**测试 2: 测试 MCP 连接**

```bash
# 测试初始化接口
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'

# 预期响应: 包含 "result" 字段和服务信息
```

**测试 3: 检查登录状态**

```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "check_login_status",
      "arguments": {}
    },
    "id": 2
  }'

# 预期响应: 显示登录状态
```

**测试 4: 列出可用工具**

```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 3
  }'

# 预期响应: 列出所有 MCP 工具
```

### 12.2 远程服务测试

**从本地机器测试：**

```bash
# 替换 <ECS公网IP> 为实际 IP
ECS_IP="<ECS公网IP>"

# 测试连接
curl -X POST http://$ECS_IP:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'
```

**如果配置了 Nginx 反向代理：**

```bash
# HTTP 测试
curl -X POST http://mcp.yourdomain.com/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'

# HTTPS 测试
curl -X POST https://mcp.yourdomain.com/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'

# 如果配置了 Basic Auth
curl -u admin:password -X POST https://mcp.yourdomain.com/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'
```

### 12.3 使用 MCP Inspector 测试

**在本地机器运行：**

```bash
# 安装并运行 MCP Inspector
npx @modelcontextprotocol/inspector

# 浏览器会自动打开，连接地址:
# 直连: http://<ECS公网IP>:18060/mcp
# 或反向代理: http://mcp.yourdomain.com/mcp
```

**测试步骤：**

1. 在 Inspector 界面输入 MCP 服务地址
2. 点击 "Connect" 按钮
3. 测试 "List Tools" - 应该显示所有工具
4. 测试 "check_login_status" - 验证登录状态
5. 测试其他功能

### 12.4 客户端接入测试

**Claude Code 接入：**

```bash
# 添加 MCP 服务器
claude mcp add --transport http xiaohongshu-mcp http://<ECS公网IP>:18060/mcp

# 或使用域名
claude mcp add --transport http xiaohongshu-mcp http://mcp.yourdomain.com/mcp

# 查看 MCP 列表
claude mcp list

# 测试功能
claude "使用 xiaohongshu-mcp 检查登录状态"
```

**Cursor/VSCode 接入：**

参考主 README.md 中的客户端配置部分。

---

## 🔍 故障排查

### 问题 1: 服务无法启动

**症状:** `systemctl status xiaohongshu-mcp` 显示 failed

**排查步骤:**

```bash
# 查看详细错误日志
journalctl -u xiaohongshu-mcp -n 50

# 查看应用日志
cat ~/xiaohongshu-mcp/logs/error.log

# 手动启动测试
cd ~/xiaohongshu-mcp
./xiaohongshu-mcp -headless=true

# 检查依赖
ldd ./xiaohongshu-mcp
```

**常见原因:**
- Go 环境变量未配置
- 浏览器依赖缺失
- cookies.json 文件格式错误
- 端口被占用

### 问题 2: 端口无法访问

**症状:** 远程无法连接到 18060 端口

**排查步骤:**

```bash
# 1. 检查服务是否监听
sudo netstat -tulpn | grep 18060

# 2. 检查本地是否可访问
curl http://localhost:18060/mcp

# 3. 检查防火墙
sudo ufw status      # Ubuntu
sudo firewall-cmd --list-all  # CentOS

# 4. 检查阿里云安全组
# 登录控制台检查入方向规则

# 5. 测试端口连通性
nc -zv <ECS公网IP> 18060
```

### 问题 3: 登录状态丢失

**症状:** check_login_status 返回未登录

**解决方法:**

```bash
# 1. 检查 cookies 文件
cat ~/xiaohongshu-mcp/data/cookies.json

# 2. 重新上传 cookies（从本地）
scp data/cookies.json root@<ECS公网IP>:~/xiaohongshu-mcp/data/

# 3. 重启服务
sudo systemctl restart xiaohongshu-mcp

# 4. 验证登录
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"check_login_status","arguments":{}},"id":1}'
```

### 问题 4: 浏览器无法启动

**症状:** 日志显示 browser launch failed

**解决方法:**

```bash
# 检查 Chromium 是否安装
chromium-browser --version

# 重新安装依赖（Ubuntu）
sudo apt install -y --reinstall \
    chromium-browser \
    libnss3 \
    libgbm-dev

# 检查字体
fc-list | grep -i chinese

# 测试浏览器启动
chromium-browser --headless --no-sandbox --disable-gpu --dump-dom https://www.baidu.com
```

### 问题 5: 内存不足

**症状:** 服务频繁重启，日志显示 OOM

**解决方法:**

```bash
# 1. 查看内存使用
free -h
ps aux --sort=-%mem | head

# 2. 增加 swap（临时方案）
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 3. 升级 ECS 配置（推荐）
# 在阿里云控制台升级实例规格
```

---

## 📚 附录

### A. 常用命令速查

```bash
# 服务管理
systemctl start xiaohongshu-mcp     # 启动
systemctl stop xiaohongshu-mcp      # 停止
systemctl restart xiaohongshu-mcp   # 重启
systemctl status xiaohongshu-mcp    # 状态
journalctl -u xiaohongshu-mcp -f    # 日志

# 防火墙
ufw status                          # UFW 状态
firewall-cmd --list-all             # firewalld 状态

# 端口检查
netstat -tulpn | grep 18060
ss -tulpn | grep 18060

# 进程查看
ps aux | grep xiaohongshu-mcp
pgrep -fa xiaohongshu-mcp

# 磁盘空间
df -h
du -sh ~/xiaohongshu-mcp/*

# 内存使用
free -h
top
htop
```

### B. 配置文件位置

```
/etc/systemd/system/xiaohongshu-mcp.service  # systemd 服务
/etc/nginx/sites-available/xiaohongshu-mcp   # Nginx 配置（Ubuntu）
/etc/nginx/conf.d/xiaohongshu-mcp.conf       # Nginx 配置（CentOS）
/etc/logrotate.d/xiaohongshu-mcp             # 日志轮转
~/xiaohongshu-mcp/data/cookies.json          # 登录 cookies
~/xiaohongshu-mcp/logs/                      # 日志目录
```

### C. 安全建议

1. **最小权限原则**
   - 不要使用 root 用户运行服务（可创建专用用户）
   - 限制安全组访问范围

2. **定期更新**
   - 定期更新系统软件包
   - 关注 xiaohongshu-mcp 项目更新

3. **备份策略**
   - 定期备份 `data/cookies.json`
   - 备份配置文件

4. **监控告警**
   - 配置服务监控
   - 设置告警通知

5. **日志审计**
   - 定期查看访问日志
   - 监控异常行为

### D. 性能优化

1. **调整 Go 运行参数**

```bash
# 修改 systemd 服务文件
[Service]
Environment="GOGC=50"           # 更频繁的 GC，降低内存占用
Environment="GOMAXPROCS=2"      # 限制 CPU 核心数
```

2. **Nginx 优化**

```nginx
# worker 进程数
worker_processes auto;

# 连接数
worker_connections 1024;

# 启用 gzip
gzip on;
gzip_types text/plain application/json;
```

3. **系统参数优化**

```bash
# 增加文件描述符限制
sudo tee -a /etc/security/limits.conf <<EOF
* soft nofile 65536
* hard nofile 65536
EOF
```

---

## 🎉 部署完成

恭喜！你已经成功将 xiaohongshu-mcp 部署到阿里云 ECS。

**后续步骤：**

1. 测试所有 MCP 功能
2. 配置客户端接入（Claude Code、Cursor、VSCode 等）
3. 设置监控告警
4. 定期备份重要数据
5. 关注项目更新

**遇到问题？**

- 查看 [GitHub Issues](https://github.com/xpzouying/xiaohongshu-mcp/issues)
- 查阅 [疑难杂症文档](https://github.com/xpzouying/xiaohongshu-mcp/issues/56)
- 加入交流群（参考主 README.md）

**项目地址：**
https://github.com/xpzouying/xiaohongshu-mcp

---

**文档版本**: v1.0
**更新日期**: 2025-12-05
**适用版本**: xiaohongshu-mcp v1.x
