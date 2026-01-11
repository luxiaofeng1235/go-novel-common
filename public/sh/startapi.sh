#!/bin/bash

PROJECT_DIR="/www/wwwroot/code/go-novel"      # 项目代码目录
EXEC_NAME="novel-api"                       # 可执行文件名称
LOG_FILE="$PROJECT_DIR/api.log"             # 日志文件路径
BRANCH="master"                             # Git 分支名称
NAME="api"

# 编译环境变量
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 在脚本开头添加
set -e

# 日志函数
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# 错误处理函数
error_exit() {
    log "$1"
    exit 1
}

# 进入项目目录
log "进入项目目录: $PROJECT_DIR"
cd "$PROJECT_DIR" || error_exit "无法进入项目目录: $PROJECT_DIR"
log "当前目录: $(pwd)"  # 调试输出

# 拉取最新代码
log "拉取最新代码..."
git pull || error_exit "拉取代码失败"

# 编译项目
log "编译项目..."
go build -o "$EXEC_NAME" "$NAME".go || error_exit "编译项目失败"
log "编译完成"

# 检查是否存在旧进程并杀掉
log "检查并杀掉旧的 $EXEC_NAME 进程（如果存在）"
pkill -f "$EXEC_NAME" || log "没有找到正在运行的进程"  # 如果没有进程也不报错

# 删除旧的备份文件（如果存在）
BACKUP_FILE="$PROJECT_DIR/${EXEC_NAME}-backup"
if [ -f "$BACKUP_FILE" ];then
    log "删除旧的备份文件: $BACKUP_FILE"
    rm "$BACKUP_FILE" || error_exit "删除备份文件失败"
fi

# 复制当前可执行文件为 backup
CURRENT_EXEC="$PROJECT_DIR/$EXEC_NAME"
if [ -f "$CURRENT_EXEC" ]; then
    log "复制当前可执行文件为: $BACKUP_FILE"
    cp "$CURRENT_EXEC" "$BACKUP_FILE" || error_exit "复制可执行文件失败"
fi

# 启动新的进程
log "启动新的 $EXEC_NAME 进程"
nohup "$PROJECT_DIR/$EXEC_NAME" > "$PROJECT_DIR/output$NAME.log" 2>&1 &
log "新进程启动完成，进程 ID: $!"

log "部署完成"