# mNas - Claude 开发上下文

## 项目概述

mNas 是基于 alist 项目扩展开发的网络附加存储解决方案，目标是在不破坏现有 alist 功能的前提下，增加 NFS 和 SMB 协议支持，让用户可以通过标准的网络文件系统协议挂载和访问 alist 管理的多种存储后端。

## Communication
- 永远使用简体中文进行思考和对话

## Documentation
- 编写 .md 文档时，也要用中文
- 正式文档写到项目的 docs/ 目录下
- 用于讨论和评审的计划、方案等文档，写到项目的 discuss/ 目录下
- 非必要不要创建文档，如果一定要创建文档也需要用户同意或用户主动要求

## Code Architecture
- 编写代码的硬性指标，包括以下原则：
  （1）对于 Python、JavaScript、TypeScript 等动态语言，尽可能确保每个代码文件不要超过 300 行
  （2）对于 Java、Go、Rust 等静态语言，尽可能确保每个代码文件不要超过 400 行
  （3）每层文件夹中的文件，尽可能不超过 8 个。如有超过，需要规划为多层子文件夹
- 除了硬性指标以外，还需要时刻关注优雅的架构设计，避免出现以下可能侵蚀我们代码质量的「坏味道」：
  （1）僵化 (Rigidity): 系统难以变更，任何微小的改动都会引发一连串的连锁修改。
  （2）冗余 (Redundancy): 同样的代码逻辑在多处重复出现，导致维护困难且容易产生不一致。
  （3）循环依赖 (Circular Dependency): 两个或多个模块互相纠缠，形成无法解耦的“死结”，导致难以测试与复用。
  （4）脆弱性 (Fragility): 对代码一处的修改，导致了系统中其他看似无关部分功能的意外损坏。
  （5）晦涩性 (Obscurity): 代码意图不明，结构混乱，导致阅读者难以理解其功能和设计。
  （6）数据泥团 (Data Clump): 多个数据项总是一起出现在不同方法的参数中，暗示着它们应该被组合成一个独立的对象。
  （7）不必要的复杂性 (Needless Complexity): 用“杀牛刀”去解决“杀鸡”的问题，过度设计使系统变得臃肿且难以理解。
- 【非常重要！！】无论是你自己编写代码，还是阅读或审核他人代码时，都要严格遵守上述硬性指标，以及时刻关注优雅的架构设计。
- 【非常重要！！】无论何时，一旦你识别出那些可能侵蚀我们代码质量的「坏味道」，都应当立即询问用户是否需要优化，并给出合理的优化建议。

## Run & Debug
- 必须首先在项目的 scripts/ 目录下，维护好 Run & Debug 需要用到的全部 .sh 脚本
- 对于所有 Run & Debug 操作，一律使用 scripts/ 目录下的 .sh 脚本进行启停。永远不要直接使用 npm、pnpm、uv、python 等等命令
- 如果 .sh 脚本执行失败，无论是 .sh 本身的问题还是其他代码问题，需要先紧急修复。然后仍然坚持用 .sh 脚本进行启停
- Run & Debug 之前，为所有项目配置 Logger with File Output，并统一输出到 logs/ 目录下

## 核心设计原则

1. **非破坏性扩展**：不修改现有 alist 核心功能和 API
2. **协议层并列**：NFS/SMB 与现有 FTP/SFTP/WebDAV 服务并列
3. **存储层复用**：通过适配器模式访问 alist 的所有存储驱动
4. **统一配置管理**：复用现有配置系统和 Web 管理界面

## 技术架构

### 目录结构
```
server/
├── nfs/           # NFS 协议实现
│   ├── filesystem.go  # billy.Filesystem 适配器
│   ├── file.go        # billy.File 实现
│   └── server.go      # NFS 服务器（待实现）
├── smb/           # SMB 协议实现
│   ├── vfs.go         # VFS 适配器（待实现）
│   └── server.go      # SMB 服务器（待实现）
└── [existing protocols...]
```

### 关键技术选型
- **NFS 实现**：基于 `github.com/willscott/go-nfs`
- **文件系统抽象**：使用 `github.com/go-git/go-billy/v5` 接口
- **SMB 实现**：参考 `references/go-smb2` 项目
- **存储适配**：通过 `internal/fs` 包访问 alist 存储层

## 当前实现状态

### ✅ 已完成
1. **项目结构设置**
   - 创建 `server/nfs/` 和 `server/smb/` 目录
   - 添加 go-nfs 相关依赖到 go.mod

2. **NFS 文件系统适配器**
   - `AlistFS` 实现 `billy.Filesystem` 接口
   - `AlistFile` 实现 `billy.File` 接口
   - `AlistFileInfo` 实现 `os.FileInfo` 接口
   - 基础的文件读取、目录列表功能

### 🚧 进行中
1. **完善 NFS 适配器**
   - HTTP 流读取实现
   - 错误处理优化
   - 性能优化

### ⏳ 待实现
1. **NFS 服务器**
   - 基于 go-nfs 的服务器实现
   - 认证和权限管理
   - 配置集成

2. **SMB 协议支持**
   - VFS 接口适配器
   - SMB 服务器实现
   - 与 alist 存储层集成

3. **服务集成**
   - 在 `cmd/server.go` 中集成启动逻辑
   - 配置选项支持（端口、认证等）
   - 生命周期管理

4. **测试和优化**
   - 协议兼容性测试
   - 性能优化
   - 并发安全

## 开发约束

### 必须遵守
- **不能修改**：`internal/driver/`、`drivers/`、现有 API 路由
- **不能破坏**：现有 Web 界面、存储配置、用户管理
- **必须保持**：向后兼容性，现有功能完整性

### 建议遵循
- 使用现有的日志系统（logrus）
- 复用现有的配置结构
- 遵循现有的代码风格和命名约定
- 优先使用现有的工具函数和中间件

## 关键文件位置

### 核心接口
- `internal/driver/driver.go` - 存储驱动接口
- `internal/fs/` - 文件系统操作包
- `internal/model/` - 数据模型定义

### 服务器相关
- `cmd/server.go` - 主服务器启动逻辑
- `server/router.go` - 路由设置
- `internal/conf/` - 配置管理

### 参考实现
- `server/ftp.go` - FTP 服务器实现参考
- `server/sftp.go` - SFTP 服务器实现参考
- `server/webdav.go` - WebDAV 服务器实现参考

## 开发提示

### 测试命令
```bash
# 编译项目
go build

# 运行服务器
./alist server

# 测试 NFS 挂载（开发完成后）
mount -o port=2049,mountport=2049,nfsvers=3,tcp -t nfs localhost:/mount ./testmount
```

### 配置样例
```yaml
# 预期的配置格式
scheme:
  nfs:
    enable: true
    port: 2049
    bind: "0.0.0.0"
  smb:
    enable: true
    port: 445
    bind: "0.0.0.0"
```

## 当前待办事项

1. 创建 NFS 服务器实现
2. 实现 alist 到 SMB VFS 的适配器
3. 创建 SMB 服务器实现
4. 在 cmd/server.go 中集成 NFS/SMB 服务启动
5. 添加配置选项支持 NFS/SMB 端口等设置

## 注意事项

- NFS 协议需要 root 权限或特殊配置才能监听 2049 端口
- SMB 协议需要 445 端口，可能与系统服务冲突
- 考虑使用非标准端口并通过挂载参数指定
- 文件权限和所有者信息的处理需要特别注意
- 大文件传输和流式读取的性能优化很重要