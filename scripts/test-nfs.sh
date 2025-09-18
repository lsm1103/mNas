#!/bin/bash

# NFS 客户端测试脚本

set -e

# 配置
NFS_SERVER="${1:-localhost}"
NFS_PORT="${2:-2049}"
MOUNT_POINT="${3:-./nfs-mount}"

echo "=== mNas NFS 客户端测试 ==="
echo "NFS 服务器: $NFS_SERVER:$NFS_PORT"
echo "挂载点: $MOUNT_POINT"
echo ""

# 检查挂载点
if [ ! -d "$MOUNT_POINT" ]; then
    echo "创建挂载点目录: $MOUNT_POINT"
    mkdir -p "$MOUNT_POINT"
fi

# 检查是否已经挂载
if mountpoint -q "$MOUNT_POINT" 2>/dev/null; then
    echo "警告: $MOUNT_POINT 已经挂载，先卸载..."
    sudo umount "$MOUNT_POINT" || echo "卸载失败，请手动处理"
    sleep 1
fi

echo "开始挂载 NFS..."

# 根据操作系统选择挂载命令
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    echo "检测到 Linux 系统"
    sudo mount -t nfs -o port=$NFS_PORT,mountport=$NFS_PORT,nfsvers=3,tcp,timeo=900,retrans=3 \
        $NFS_SERVER:/mount "$MOUNT_POINT"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    echo "检测到 macOS 系统"
    sudo mount -o port=$NFS_PORT,mountport=$NFS_PORT,nfsvers=3,tcp \
        -t nfs $NFS_SERVER:/mount "$MOUNT_POINT"
else
    echo "错误: 不支持的操作系统 $OSTYPE"
    echo "请手动执行挂载命令:"
    echo ""
    echo "Linux:"
    echo "sudo mount -t nfs -o port=$NFS_PORT,mountport=$NFS_PORT,nfsvers=3,tcp $NFS_SERVER:/mount $MOUNT_POINT"
    echo ""
    echo "macOS:"
    echo "sudo mount -o port=$NFS_PORT,mountport=$NFS_PORT,nfsvers=3,tcp -t nfs $NFS_SERVER:/mount $MOUNT_POINT"
    exit 1
fi

# 检查挂载状态
if mountpoint -q "$MOUNT_POINT"; then
    echo "✅ NFS 挂载成功!"
    echo ""
    echo "挂载信息:"
    df -h "$MOUNT_POINT"
    echo ""
    echo "内容预览:"
    ls -la "$MOUNT_POINT" | head -10
    echo ""
    echo "测试完成。使用以下命令卸载:"
    echo "sudo umount $MOUNT_POINT"
else
    echo "❌ NFS 挂载失败"
    exit 1
fi