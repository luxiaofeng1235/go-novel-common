# Repository Guidelines

_仓库指南_

## 项目结构与模块组织

- `go.mod`：模块 `go-novel`（Go `1.19`）。
- 根目录 `*.go`：包含多个独立 `package main` 入口（API/后台/source/抓取器等）。由于根目录存在多个 `func main()`，请按“单文件入口”方式构建/运行。
- `app/`：按分层组织的业务代码：
  - `app/controller/`：Gin 处理器（如 `app/controller/api`、`app/controller/admin`、`app/controller/source`）
  - `app/service/`：业务逻辑（`app/service/...`）
  - `app/models/`：数据模型（`app/models/...`）
- `routers/`：Gin 路由注册与分组（`routers/*_routes/*_route.go`）。
- `db/`：DB/日志/缓存/NSQ/WS 等初始化与启动编排（见 `db/bootstrap.go`）。
- `config/`：基于 Viper 的配置加载；默认在仓库根目录读取 `config.yml`（示例见 `config.yml.dev`）。
- `public/`：静态资源与 Casbin 配置（`public/casbin_conf/...`）。
- `gofound-1.4.1/`：内嵌组件（自带代码/测试）；除非必要，避免改动。

## 开发规范（Router → Controller → Service → Model）

- 路由（Router）：只做“路径/方法/中间件/分组”注册，代码放在 `routers/*_routes/`（例如 `routers/api_routes/user_route.go`、`routers/source_routes/route.go`）。新增接口先加路由，再补 controller/service。
- 控制器（Controller）：只做参数绑定、轻量校验、调用 service、统一返回；放在 `app/controller/<api|admin>/`。推荐统一使用 `utils.SuccessEncrypt()` / `utils.FailEncrypt()` 输出 JSON（保持前端兼容）。
- 服务（Service）：承载业务逻辑与数据访问；放在 `app/service/<api|admin|common>/...`。DB 统一通过 `global.DB` 访问，避免在 controller 内直接写 SQL/GORM；尽量不要在 service 里直接写 HTTP 响应。
- 模型（Model）：GORM 模型与请求结构体集中在 `app/models/`（如 `app/models/mc_user.go` 已包含 `LoginReq/RegisterReq/GuestLoginReq`）。新增表时：模型字段需有 `gorm:"column:..."`，并实现 `TableName()` 返回真实表名。
- 强约束（必须遵守）：严禁在 `app/controller/...` 或 `app/models/...` 直接编写业务逻辑（包括：权限判断、状态机、落库/事务、复杂校验、跨表查询、第三方调用等）；所有业务处理必须“转发/下沉”到 `app/service/...` 统一实现与复用。
- 依赖方向：`routers` → `app/controller` → `app/service` → `app/models`（单向依赖）；禁止反向引用（例如 models 引用 service/routers）。
- 示例（新增一个 API 接口）：1) `routers/api_routes/*.go` 注册 `POST /api/xxx`；2) `app/controller/api/xxx.go` `ShouldBind` 到 `models.*Req`；3) `app/service/api/xxx_service/*.go` 实现核心逻辑；4) 需要落库则补 `app/models/*.go`。

## 构建、测试与本地开发命令

- 初始化本地配置：
  - `cp config.yml.dev config.yml`（项目级配置；不要提交 `config.yml`）
  - `cp config/upload.yml.dev config/upload.yml`（业务级上传配置；不要提交 `config/upload.yml`）
- 数据库与缓存连接统一从 `config.yml` 读取：`mysql.*`、`redis.*`（示例见 `config.yml.dev`）。
- 入口规范（底层改造约定）：`api.go` / `admin.go` 入口文件保持“最小化”，只负责调用 `db.StartApiServer()` / `db.StartAdminServer()` 等启动函数；不要在入口里堆叠复杂参数解析、服务编排或业务初始化逻辑，统一收敛到 `db/` 的启动编排代码中。
- 端口与模块约定（默认）：`api=8006`、`admin=8005`、`source=8007`；API 监听优先读取 `api.host/api.port`，若未配置则回退到 `server.host/server.port`；其他模块使用 `admin.host/admin.port`、`source.host/source.port`。
- 运行服务（在仓库根目录）：
  - `go run ./api.go`（API 服务）
  - `go run ./admin.go`（后台服务）
  - `go run ./source.go`（source 服务）
- 构建单个入口：
  - `go build -o novel-api ./api.go`
  - `go build -o novel-admin ./admin.go`
- 脚本：`./startsource.sh`（依赖 `novel-source` 二进制，日志输出到 `source.log`）。

## WSL/Windows 联调（重要）

- WSL2 下如果服务只监听 `::1:PORT`（IPv6 回环），Windows 访问 `http://127.0.0.1:PORT` 会失败（常见现象：Windows `netstat` 只看到 `[::1]:PORT`）。
- 本仓库默认采用 IPv4 监听（`tcp4`）保证 Windows/WSL 直连可用；若你自行调整启动/监听逻辑，请确保至少监听到 `127.0.0.1:PORT`（或 `0.0.0.0:PORT`）以支持 IPv4 访问。
- Windows 侧联调建议使用 `curl.exe`（PowerShell 的 `curl` 可能是别名）：`curl.exe -i http://127.0.0.1:8006/api/user/guest`

## 代码风格与命名约定

- Go 代码统一使用 `gofmt`（标准 Go 格式；默认 tab 缩进）。如团队使用 `goimports` 可优先。
- 目录/包名使用小写（如 `routers/api_routes`）；导出标识符使用 `CamelCase`。
- 路由文件放在 `routers/*_routes/`，处理器实现放在 `app/controller/...`，避免交叉堆放。
- 文件头注释统一使用以下模板；新增文件 `@Author` / `@LastEditors` 统一写 `red`：
  - `/* @Descripttion: ... @Author: ... @Date: ... @LastEditors: ... @LastEditTime: ... */`

## 测试指南

- 测试文件命名为 `*_test.go`，尽量与被测包同目录（优先表驱动测试）。
- 现有测试主要位于 `gofound-1.4.1/`：`go test ./gofound-1.4.1/...`
- 多数启动逻辑依赖根目录 `config.yml`；运行涉及初始化的测试/命令前请确保该文件存在。

## 提交与 Pull Request 指南

- 历史提交信息较随意（如 `修改`、`12`），当前未形成强制规范。
- 建议新提交采用“范围 + 动词”的祈使句：`api: fix startup crash`，并保持提交聚焦。
- PR 请包含：变更目的、如何运行/测试（`go run ...` / `go test ...`）、以及新增/变更的配置项（尤其是 `config.yml` 字段）；涉及 `public/` 静态资源变更时请附截图。
- 权限约定：贡献者未经仓库持有者允许，不得私自执行 `merge`、`commit` 到受保护分支、`reset`/`rebase` 改写历史、`push --force`、打 `tag`/发版等操作；上述操作统一由仓库持有者执行。
