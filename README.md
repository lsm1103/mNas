# mNas - 多协议网络附加存储

基于 [alist](https://github.com/alist-org/alist) 扩展开发，支持 NFS、SMB 协议的网络文件系统服务。

## 🚀 项目特性

### 核心功能
- **多协议支持**：在 alist 现有协议基础上增加 NFS 和 SMB 支持
- **存储统一**：通过单一服务访问多种云存储和本地存储
- **原生挂载**：支持操作系统原生的网络文件系统挂载
- **配置复用**：完全复用 alist 的存储配置和 Web 管理界面

### 协议支持
- ✅ **HTTP/HTTPS** - Web 界面和 API 访问
- ✅ **WebDAV** - 网盘客户端支持
- ✅ **FTP/FTPS** - 传统文件传输协议
- ✅ **SFTP** - SSH 文件传输协议
- 🚧 **NFS v3** - 网络文件系统（开发中）
- ⏳ **SMB/CIFS** - Windows 文件共享（计划中）

### 存储后端
继承 alist 的所有存储驱动支持，包括但不限于：
- 本地存储
- 阿里云盘、百度网盘、OneDrive 等云盘
- Amazon S3、阿里云 OSS 等对象存储
- FTP、SFTP 等远程存储
- 更多 70+ 存储驱动...

## 📦 安装

### 环境要求
- Go 1.23.4 或更高版本
- Linux/macOS/Windows 系统

### 从源码构建
```bash
# 克隆仓库
git clone https://github.com/your-org/mNas.git
cd mNas

# 安装依赖
go mod tidy

# 构建
go build -o mnas

# 运行
./mnas server
```

### Docker 部署
```bash
# 构建镜像
docker build -t mnas .

# 运行容器
docker run -d \
  --name mnas \
  -p 5244:5244 \
  -p 2049:2049 \
  -p 445:445 \
  -v /path/to/data:/opt/alist/data \
  mnas
```

## 🔧 配置

### 基础配置
mNas 完全兼容 alist 的配置格式，在此基础上扩展了协议配置：

```yaml
# config.yml
scheme:
  address: "0.0.0.0"
  http_port: 5244
  https_port: -1

  # NFS 配置
  nfs:
    enable: true
    port: 2049
    bind: "0.0.0.0"

  # SMB 配置（计划中）
  smb:
    enable: false
    port: 445
    bind: "0.0.0.0"
```

### 存储配置
通过 Web 界面 `http://localhost:5244` 配置存储驱动，与原版 alist 完全相同。

## 📖 使用方法

### Web 管理
访问 `http://localhost:5244` 进行存储配置和文件管理，功能与 alist 完全一致。

### NFS 挂载
```bash
# Linux
sudo mount -t nfs -o port=2049,mountport=2049,nfsvers=3,tcp localhost:/mount /mnt/mnas

# macOS
sudo mount -o port=2049,mountport=2049 -t nfs localhost:/mount /mnt/mnas

# 卸载
sudo umount /mnt/mnas
```

### SMB 挂载（开发中）
```bash
# Linux
sudo mount -t cifs //localhost/share /mnt/mnas -o port=445

# Windows
net use Z: \\localhost\share

# macOS
mount -t smbfs //localhost/share /mnt/mnas
```

## 🔧 开发状态

### 当前进度
- ✅ 项目基础架构搭建
- ✅ NFS 文件系统适配器实现
- 🚧 NFS 服务器集成
- ⏳ SMB 协议适配器
- ⏳ SMB 服务器实现
- ⏳ 配置系统集成

### 技术架构
```
┌─────────────────┐    ┌──────────────────┐
│   Web/API       │    │   NFS/SMB        │
│   (alist 原生)   │    │   (mNas 扩展)     │
└─────────────────┘    └──────────────────┘
          │                       │
          └───────────┬───────────┘
                      │
        ┌─────────────────────────────┐
        │      alist 存储抽象层         │
        └─────────────────────────────┘
                      │
        ┌─────────────────────────────┐
        │     70+ 存储驱动支持          │
        │  (本地/云盘/对象存储/...)      │
        └─────────────────────────────┘
```

## 🤝 贡献

### 开发环境
```bash
# 安装开发工具
go install github.com/air-verse/air@latest

# 热重载开发
air
```

### 贡献指南
1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目基于 AGPL-3.0 许可证开源，与上游 alist 项目保持一致。

## 🙏 致谢

- [alist](https://github.com/alist-org/alist) - 优秀的文件列表程序
- [go-nfs](https://github.com/willscott/go-nfs) - NFS 协议 Go 实现
- [go-smb2](https://github.com/hirochachacha/go-smb2) - SMB 协议 Go 实现

## 📞 支持

- 🐛 问题反馈：[GitHub Issues](https://github.com/your-org/mNas/issues)
- 💬 讨论交流：[GitHub Discussions](https://github.com/your-org/mNas/discussions)
- 📧 邮件联系：your-email@example.com

---

**注意**：本项目目前处于早期开发阶段，不建议在生产环境使用。