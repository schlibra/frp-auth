# Frp-Auth
Frps认证服务，可以当作frps插件直接对连接进行认证，可以通过web UI管理Token认证、端口范围

## 编译
### 前端
```bash
cd /path-to-project
cd frontend
yarn install
yarn build-only
```
### 后端
```bash
cd /path-to-project
go get
go build -o frp-auth
```
## 运行服务
> 初始化数据库
```bash
./frp-auth init
```
> 临时运行
```bash
./frp-auth serve
```
> 运行后台服务
```bash
./frp-auth start
```
> 停止后台服务
```bash
./frp-auth stop
```
> 重启后台服务
```bash
./frp-auth restart
```
> 检查后台服务状态
```bash
./frp-auth status
```
## FrpAuth配置文件
> 需要将`example-config.toml`改名为`config.toml`
```toml
[mysql]
# MySQL数据库配置
host = "127.0.0.1"
port = 3306
user = "root"
pass = "123456"
name = "frp-auth"

[server]
# FrpAuth服务监听地址
host = "0.0.0.0"
port = 8080
# 服务Debug模式
debug = false

[frps-web]
# frps web配置，用于获取在线客户端和映射
host = "127.0.0.1"
port = 7500
# Basic Auth Token, 生成方法：base64encode(username + ":" + password)
token = "VXNlcm5hbWU6UGFzc3dvcmQ="

[frps]
# frps实例配置，用于在web UI中生成frpc配置时读取的地址
host = "127.0.0.1"
port = 7000

[jwt]
# 随机密钥，生成方法: openssl rand -hex 32
key = "<JWT-Key>"

[dangerous]
# 危险操作配置
# 允许清空数据库，用于重置数据
allowClean = false
```
## Frps配置
```toml
[[httpPlugins]]
name = "user-login"
addr = "127.0.0.1:8080"
path = "/api/frp/login"
ops = ["Login"]

[[httpPlugins]]
name = "user-proxy"
addr = "127.0.0.1:8080"
path = "/api/frp/proxy"
ops = ["NewProxy"]
```