# go-novel

小说阅读平台后端，基于 Go、Gin、GORM，提供用户 API、管理后台、书源接入和采集能力。详细目录、开发规范和提交约定见 [AGENTS.md](./AGENTS.md)。

## 项目流程

1. 书源服务和采集程序接入书籍、章节数据，通过 NSQ 异步处理更新。
2. MySQL 保存书籍及业务数据，章节目录和正文也使用文件存储，Redis 提供缓存。
3. 用户 API 提供书籍浏览、搜索、阅读、书架和会员等功能；管理后台维护书籍、用户、运营配置及权限。

接口调用链为 `routers → app/controller → app/service → app/models`，公共业务放在 `app/service/common/`，基础组件的初始化与启动编排位于 `db/`。

## 本地运行

首次运行时，将 `config.yml.example` 复制为 `config.yml`，按本地环境调整配置。`config.yml` 已被 Git 忽略。

本地 MySQL 使用 `127.0.0.1:3306`，数据库名为 `novel`，账号和密码均为 `root`；Redis 使用 `127.0.0.1:6379`，无密码，数据库为 `0`。需准备好数据库及业务表；书源和采集流程还需要 NSQ。

在仓库根目录按需启动，各服务分别运行：

| 服务 | 启动命令 | 默认端口 | 路由前缀 |
| --- | --- | --- | --- |
| 用户 API | `go run ./api.go` | 8005 | `/api` |
| 管理后台 | `go run ./admin.go` | 8006 | `/system` |
| 书源服务 | `go run ./source.go` | 8007 | `/source` |

## 接口临时调试

在控制器或业务代码中导入 `go-novel/internal/debugger`，直接调用：

```go
import "go-novel/internal/debugger"

debugger.Hang(111)
debugger.Hang(result)
debugger.Hang("查询结果", result)
```

`Hang` 只向服务终端打印，不依赖 Gin 中间件，也不会自动停止执行或生成 HTTP 响应。需要结束当前控制器时，显式写 `return`：

```go
debugger.Hang("调试数据", result)
return
```

`return` 只结束它所在的函数；如果放在业务方法里，上层控制器仍会继续执行。当前 `goTest` 测试接口使用上述方式，终端输出调试内容，HTTP 响应为空。

输出包含时间、调用文件及行号、变量类型和内容。变量在调用时格式化，终端由后台协程统一输出；容量为 128 条的队列满时跳过输出，可通过 `debugger.Dropped()` 查看累计丢弃数量。

程序正常退出前可按需调用 `debugger.Flush(ctx)`，等待此前入队的消息完成写入尝试；建议给 `ctx` 设置超时，不要在每次接口调用时执行。

## 构建与平台切换

以下命令在 **Windows CMD** 中执行。根目录存在多个 `main()`，构建时需指定入口文件。

切换到 Linux amd64 后构建：

```bat
set GOARCH=amd64
go env -w GOARCH=amd64
set GOOS=linux
go env -w GOOS=linux

go build -o novel-api ./api.go
go build -o novel-admin ./admin.go
go build -o novel-source ./source.go
```

生成的二进制文件用于 Linux 运行，运行目录需提供 `config.yml` 和所需资源文件。

切回 Windows amd64：

```bat
set GOARCH=amd64
go env -w GOARCH=amd64
set GOOS=windows
go env -w GOOS=windows

go build -o novel-api.exe ./api.go
go build -o novel-admin.exe ./admin.go
go build -o novel-source.exe ./source.go
```

`set` 设置当前 CMD 窗口的环境变量，`go env -w` 保存 Go 的默认值。当前窗口的环境变量优先，因此切回 Windows 时也要同步执行 `set GOOS=windows`。
