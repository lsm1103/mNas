#!/bin/bash

# mNas 启动脚本

set -e

# 检查是否已编译
if [ ! -f "./mnas" ]; then
    echo "正在编译 mNas..."
    go build -o mnas
    echo "编译完成"
fi

# 设置数据目录
DATA_DIR="${1:-data}"

# 确保数据目录存在
mkdir -p "$DATA_DIR/temp" "$DATA_DIR/bleve" "$DATA_DIR/log"

# 如果 data/config.json 不存在，则从示例配置复制
if [ ! -f "$DATA_DIR/config.json" ]; then
    if [ -f "config-example.json" ]; then
        echo "复制示例配置到 $DATA_DIR/config.json"
        cp config-example.json "$DATA_DIR/config.json"
    else
        echo "警告: 配置文件 $DATA_DIR/config.json 不存在，将使用默认配置"
    fi
fi

echo "启动 mNas 服务器..."
echo "数据目录: $DATA_DIR"
echo "配置文件: $DATA_DIR/config.json"
echo "Web 界面: http://localhost:5244"
echo "NFS 服务: localhost:2049"
echo ""
echo "按 Ctrl+C 停止服务"

# 启动服务器
./mnas server --data "$DATA_DIR"