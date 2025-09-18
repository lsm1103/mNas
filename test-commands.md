# mNas NFS 测试指南

## 1. 启动服务器

```bash
# 使用脚本启动（推荐）
./scripts/start.sh

# 或手动启动
go build -o mnas
./mnas server --data data
```

**注意**: alist 使用 `--data` 参数指定数据目录，配置文件位于 `data/config.json`

启动后应该看到类似输出：
```
INFO[0000] NFS 服务器启动成功，监听地址: 0.0.0.0:2049
INFO[0000] start HTTP server @ 0.0.0.0:5244
```

## 2. 验证服务状态

```bash
# 检查端口监听
netstat -tlnp | grep :2049
netstat -tlnp | grep :5244

# 检查进程
ps aux | grep mnas
```

## 3. 客户端挂载测试

### 方法一：使用测试脚本（推荐）
```bash
# 自动检测系统并挂载
./scripts/test-nfs.sh

# 指定参数
./scripts/test-nfs.sh localhost 2049 ./my-mount-point
```

### 方法二：手动挂载

#### Linux 系统：
```bash
# 创建挂载点
mkdir -p ./nfs-mount

# 挂载 NFS
sudo mount -t nfs \
  -o port=2049,mountport=2049,nfsvers=3,tcp,timeo=900,retrans=3 \
  localhost:/mount ./nfs-mount

# 验证挂载
df -h ./nfs-mount
ls -la ./nfs-mount

# 卸载
sudo umount ./nfs-mount
```

#### macOS 系统：
```bash
# 创建挂载点
mkdir -p ./nfs-mount

# 挂载 NFS
sudo mount -o port=2049,mountport=2049,nfsvers=3,tcp \
  -t nfs localhost:/mount ./nfs-mount

# 验证挂载
df -h ./nfs-mount
ls -la ./nfs-mount

# 卸载
sudo umount ./nfs-mount
```

## 4. 功能测试

```bash
# 挂载成功后，测试基本操作
cd nfs-mount

# 查看目录内容
ls -la

# 查看文件内容（如果有文件）
cat filename.txt

# 检查文件属性
stat filename.txt

# 测试目录遍历
find . -type f | head -10
```

## 5. 故障排除

### 常见错误及解决方案

1. **Permission denied**
   ```bash
   # 确保以 root 权限运行挂载命令
   sudo mount ...
   ```

2. **Connection refused**
   ```bash
   # 检查服务器是否启动
   netstat -tlnp | grep :2049

   # 检查防火墙设置
   sudo ufw status
   ```

3. **Protocol not supported**
   ```bash
   # 确保系统支持 NFS 客户端
   # Ubuntu/Debian:
   sudo apt-get install nfs-common

   # CentOS/RHEL:
   sudo yum install nfs-utils

   # macOS: 内置支持
   ```

4. **Timeout errors**
   ```bash
   # 增加超时时间
   sudo mount -t nfs -o port=2049,mountport=2049,nfsvers=3,tcp,timeo=1800,retrans=5 \
     localhost:/mount ./nfs-mount
   ```

### 调试信息

```bash
# 查看服务器日志
tail -f data/log/log.log

# 查看系统日志
# Linux:
sudo dmesg | grep -i nfs
sudo journalctl -u nfs* -f

# macOS:
sudo log stream --predicate 'process == "mount_nfs"' --info
```

## 6. 配置说明

在 `config-example.json` 中的 NFS 相关配置：

```json
{
  "nfs": {
    "enable": true,      // 启用 NFS 服务
    "address": "0.0.0.0", // 监听地址，0.0.0.0 表示所有网卡
    "port": 2049         // NFS 端口，标准端口是 2049
  }
}
```

**注意事项：**
- 端口 2049 通常需要 root 权限才能监听
- 如果遇到权限问题，可以尝试使用非标准端口如 20049
- 确保防火墙允许相应端口的访问