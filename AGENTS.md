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
- 环境（`server.env`）：当前脚手架不做 `local/dev/prod` 分支逻辑，尽量避免“按环境写代码”；如需差异化请通过不同的 `config.yml` 管理。
- `public/`：静态资源与 Casbin 配置（`public/casbin_conf/...`）。
- `gofound-1.4.1/`：内嵌组件（自带代码/测试）；除非必要，避免改动。

## 开发规范（Router → Controller → Service → Model）

- 路由（Router）：只做“路径/方法/中间件/分组”注册，代码放在 `routers/*_routes/`（例如 `routers/api_routes/user_route.go`、`routers/source_routes/route.go`）。新增接口先加路由，再补 controller/service。
- 路由鉴权约定：若某个模块/一组接口都需要登录，优先使用 `r.Group(...).Use(middleware.ApiJwt())` 统一鉴权；若只有少数接口需要登录，则在单条路由上添加 `middleware.ApiJwt()`（例如 `user.POST("/follow", middleware.ApiJwt(), userApi.Follow)`）。
- 控制器（Controller）：只做参数绑定、轻量校验、调用 service、统一返回；放在 `app/controller/<api|admin>/`。推荐统一使用 `utils.SuccessEncrypt()` / `utils.FailEncrypt()` 输出 JSON（保持前端兼容）。
- 控制器（Controller）补充：避免在 Controller 里做“配置读取/文件大小限制/复杂判断/权限分支”等易膨胀逻辑；优先通过 `middleware`（通用拦截）或 `service`（业务处理）承接，Controller 保持“接参 → 转发 → 返回”。
- 服务（Service）：承载业务逻辑与数据访问；放在 `app/service/<api|admin|common>/...`。DB 统一通过 `global.DB` 访问，避免在 controller 内直接写 SQL/GORM；尽量不要在 service 里直接写 HTTP 响应。
- Service 查询单条（约定）：按 `userID` 等主键查询时，优先使用 `First` 并显式处理“未找到”，返回时清空敏感字段（如 `Passwd`），例如：`func GetUserInfoByUserID(userID int64) (user *models.McUser, err error) { ... }`。
- 模型（Model）：GORM 模型与请求结构体集中在 `app/models/`（如 `app/models/mc_user.go` 已包含 `LoginReq/RegisterReq/GuestLoginReq`）。新增表时：模型字段需有 `gorm:"column:..."`，并实现 `TableName()` 返回真实表名。
- 强约束（必须遵守）：严禁在 `app/controller/...` 或 `app/models/...` 直接编写业务逻辑（包括：权限判断、状态机、落库/事务、复杂校验、跨表查询、第三方调用等）；所有业务处理必须“转发/下沉”到 `app/service/...` 统一实现与复用。
- 依赖方向：`routers` → `app/controller` → `app/service` → `app/models`（单向依赖）；禁止反向引用（例如 models 引用 service/routers）。
- 示例（新增一个 API 接口）：1) `routers/api_routes/*.go` 注册 `POST /api/xxx`；2) `app/controller/api/xxx.go` `ShouldBind` 到 `models.*Req`；3) `app/service/api/xxx_service/*.go` 实现核心逻辑；4) 需要落库则补 `app/models/*.go`。

## 配置文件说明

### 主配置文件（config.yml）
核心配置项包括：
- `server.host/port`：API 服务监听地址（默认 0.0.0.0:8006）
- `server.apiUrl`：对外访问地址
- `server.encrypt`：是否加密请求/响应（默认 false）
- `jwt.secret`：JWT 签名密钥（必配）
- `auth.passwordSalt`：密码盐值（可选）
- `mysql.*`：MySQL 连接配置（host、port、database、user、password、pool）
- `redis.*`：Redis 连接配置（host、port、password、db）
- `source.*`：静态资源服务配置（host、port、publicBaseUrl）

### 上传配置文件（config/upload.yml，可选）
- `baseDir`：本地存储根目录（默认 `./public/upload`）
- `maxSizeMB`：最大文件大小（默认 50MB）
- `allowedExts`：允许的文件后缀白名单（.jpg/.png/.mp4 等）
- `allowedMimePrefixes`：允许的 MIME 类型前缀（image/、video/ 等）

## 全局变量说明（global/global.go）

项目通过 `global` 包统一管理全局变量，核心实例如下：

### 常用全局变量
- `global.DB`：MySQL 数据库连接（GORM）
- `global.Redis`：Redis 客户端
- `global.WsHubManager`：WebSocket HubManager 实例（分片架构，根据 CPU 核心数自动优化并发能力）
- `global.Errlog`、`global.Sqllog`：错误日志、SQL 日志（Zap SugaredLogger）

### 使用规范
- **DB 访问**：所有 Service 层通过 `global.DB` 访问数据库，禁止在 Controller 中直接使用
- **日志记录**：选择对应模块的日志实例，如 `global.Errlog.Error("错误信息")`
- **Redis 操作**：通过 `global.Redis` 执行缓存操作

## 构建、测试与本地开发命令

- 初始化本地配置：
  - `cp config.yml.dev config.yml`（项目级配置；不要提交 `config.yml`）
  - `cp config/upload.yml.dev config/upload.yml`（业务级上传配置；不要提交 `config/upload.yml`）
- 数据库与缓存连接统一从 `config.yml` 读取：`mysql.*`、`redis.*`（示例见上方配置说明）。
- 入口规范（底层改造约定）：`api.go` / `admin.go` 入口文件保持“最小化”，只负责调用 `db.StartApiServer()` / `db.StartAdminServer()` 等启动函数；不要在入口里堆叠复杂参数解析、服务编排或业务初始化逻辑，统一收敛到 `db/` 的启动编排代码中。
- 端口与模块约定（默认）：`api=8006`、`admin=8005`、`source=8007`；API 监听统一读取 `server.host/server.port`（`api.host/api.port` 已废弃）；其他模块使用 `admin.host/admin.port`、`source.host/source.port`。
- 运行服务（在仓库根目录）：
  - `go run ./api.go`（API 服务；同进程启动 `source` 静态服务，默认监听 `source.host/source.port`）
  - `go run ./api.go -host 0.0.0.0 -port 18016`（覆盖监听地址/端口，便于本地多实例/避免端口占用）
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

## 数据库与缓存初始化详解

### MySQL 连接配置（db/mysql.go）
- 连接参数读取优先级：`mysql.address` > `mysql.host:mysql.port`（默认 `127.0.0.1:3306`）
- 数据库名优先级：`mysql.database` > `mysql.dbname` > `mysql.name`（默认 `novel`）
- 连接池配置（可在 `config.yml` 的 `mysql.pool` 节点配置）：
  - `MaxIdleConns`：默认 25
  - `MaxOpenConns`：默认 100
  - `ConnMaxLifetime`：默认 600 秒
  - `ConnMaxIdleTime`：可选配置
- GORM 配置：
  - `TablePrefix: "mc_"`：所有表名自动加 `mc_` 前缀
  - `SingularTable: true`：使用单数表名（如 `mc_user` 而非 `mc_users`）
  - `Logger`：生产环境建议 `Silent`，开发环境可设为 `Info`

### Redis 连接配置（db/redis.go）
- 连接参数读取优先级：`redis.addr` > `redis.host:redis.port`（默认 `127.0.0.1:6379`）
- 密码配置：`redis.password`（可选）
- 数据库编号：`redis.db`（默认 0）
- 连接实例存储在 `global.Redis`，全局访问

### 日志系统（db/zaplog.go）
- 日志目录：`./public/logs/{模块名}/`
- 日志分类（按模块）：
  - `global.Sqllog`：SQL 日志
  - `global.Errlog`：错误日志
  - `global.Paylog`：支付日志
  - `global.Wslog`：WebSocket 日志
  - `global.Collectlog`：采集日志
  - 其他：`Zssqlog`、`Nsqlog`、`Updatelog`、`Requestlog` 等
- 日志配置（`config.yml` 的 `logs` 节点）：
  - `level`：日志级别（-1=Debug, 0=Info, 1=Warn, 2=Error）
  - `max-size`：单文件最大大小（MB）
  - `max-backups`：保留备份文件数量
  - `max-age`：保留天数
  - `compress`：是否压缩

## JWT 认证详解

### JWT 配置与实现（utils/jwt.go）
- **密钥配置**：`config.yml` 的 `jwt.secret`（必配，不能为空）
- **Token 有效期**：24 小时 × 90 = 2160 小时（约 90 天）
- **Claims 结构**：
  ```go
  type Claims struct {
    UserID    int64  // 用户 ID
    Username  string // 用户名
    Authority int    // 权限级别
    jwt.StandardClaims
  }
  ```
- **签发函数**：`utils.GenerateToken(userID, username, authority)`
- **验证函数**：`utils.ParseToken(token)` - 防止 alg 混淆攻击（仅允许 HMAC 系列）
- **刷新函数**：`utils.RefreshToken(tokenString)` - 保留原 userID/username，签发新 token

### 密码加密机制（app/service/api/user_service/user_api.go）
- **加密算法**：MD5 + 可配置盐值
- **盐值配置**：`config.yml` 的 `auth.passwordSalt`（可选，为空则仅 MD5）
- **加密函数**：`hashLoginPasswd(plain string)`
  - 若配置了 `passwordSalt`，则：`MD5(plain + salt)`
  - 若未配置，则：`MD5(plain)`
- **密码比对**：登录时对输入密码执行相同加密，与数据库 `Passwd` 字段比对

## 中间件详解

### JWT 鉴权中间件（middleware/auth_user.go）
- **函数名**：`middleware.ApiJwt()` 或 `middleware.AuthUser()`
- **Token 提取**：
  1. 优先从 `Authorization: Bearer <token>` 头提取
  2. 兼容 `Authorization: <token>`（直接传 token）
  3. 兼容历史 `Token: <token>` 头
- **校验流程**：
  1. 提取 token，为空则返回 403 "缺少token"
  2. 调用 `utils.ParseToken()` 验证签名和过期时间
  3. 校验 `claims.UserID > 0` 和 `claims.Username` 非空
  4. 将 `user_id` 和 `username` 注入到 `gin.Context`
  5. 校验失败时返回 403 "token无效" 并中止请求
- **Controller 获取用户信息**：
  ```go
  userIDVal, ok := c.Get("user_id")
  userID, ok := userIDVal.(int64)
  ```

### 请求解密中间件（middleware/api_req_decrypt.go）
- **函数名**：`middleware.ApiReqDecrypt()`
- **配置开关**：`config.yml` 的 `server.encrypt`（默认 false）
- **解密流程**（当 `encrypt=true` 时）：
  1. 读取 POST body 中的 `{data: <加密内容>}` 格式
  2. 使用 AES-CFB 模式解密（密钥：`utils.ApiAesKey = "WB0nMZHXlxNndORe"`）
  3. 将解密后的 JSON 重新绑定到 `Request.Body`
- **响应加密**：`utils.SuccessEncrypt()` / `utils.FailEncrypt()` 根据同一开关决定是否加密响应

### CORS 中间件（middleware/cors.go）
- **函数名**：`middleware.Cors()`
- **允许来源**：`Access-Control-Allow-Origin: *`（所有来源）
- **允许请求头**：`Content-Type`、`Authorization`、`Token`、`X-Token`、`X-User-Id` 等
- **允许方法**：`POST`、`GET`、`PUT`、`DELETE`、`OPTIONS`、`PATCH`
- **预检请求**：OPTIONS 请求直接返回 204

## 用户业务说明（app/models/mc_user.go）

### 用户模型（表名：mc_user）
- 核心字段：用户名、密码（MD5+盐）、昵称、手机、邮箱、头像
- 游客支持：`IsGuest`、`Deviceid`、`Oaid`、`Imei` 等设备标识
- 推荐体系：`ParentId`（上级 ID）、`Invitation`（邀请码）、`ParentLink`（推荐链条）
- VIP 体系：`Vip`、`Viptime`、`Rmb`（余额）、`Cion`（金币）
- 状态管理：`Status`（0-锁定 1-正常）、渠道号、包名、IP、登录时间等

### 请求结构体（同文件）
- `LoginReq`：登录（支持 username/tel/email）
- `RegisterReq`：注册（username、passwd、nickname、referrer）
- `GuestLoginReq`：游客登录（deviceid、referrer、sex）
- `EditUserReq`：编辑用户信息（tel/email/nickname/sex/pic/passwd/book_type）

## 文件上传服务

### 上传接口（POST /api/common/upload）
- **路由**：`routers/api_routes/common_route.go`
- **Controller**：`app/controller/api/common.go:Upload()`
- **Service**：`app/service/common/file_service/file_service.go:LocalUpload()`
- **表单字段**：`file`（必需）、`dir`（可选子目录，如 `avatar`、`book_cover`）
- **配置文件**：`config/upload.yml`（可选，包含 baseDir、maxSizeMB、白名单等）

### 安全机制
- 文件大小限制（默认 50MB）
- 后缀白名单校验
- MIME 类型白名单校验
- 路径安全（禁止 `..`、绝对路径）
- 随机文件名（防止冲突和路径遍历）

## WebSocket 通信

### 架构设计（pkg/ws/）
- **Hub 模式**：单 goroutine 主循环 + 多 goroutine 读写泵，避免 map 并发竞态
- **连接管理**：支持同账号多端在线（按 userID 归档连接）
- **核心文件**：`pkg/ws/hub.go`（Hub）、`pkg/ws/client.go`（客户端）、`pkg/ws/protocol.go`（协议）

### 连接与消息
- **路由**：`GET /api/ws?token=<JWT>`（token 可选，传入后支持鉴权功能）
- **心跳**：服务端 54 秒发送 Ping，客户端回复 Pong；读超时 60 秒
- **消息类型**：
  - `ping/pong`：心跳检测
  - `chat`：群聊广播（需 token，发送给所有已鉴权用户）
  - `dm`：单聊私信（需 token，支持多端同步和回显）

## 代码风格与命名约定

- Go 代码统一使用 `gofmt`（标准 Go 格式；默认 tab 缩进）。如团队使用 `goimports` 可优先。
- 目录/包名使用小写（如 `routers/api_routes`）；导出标识符使用 `CamelCase`。
- 路由文件放在 `routers/*_routes/`，处理器实现放在 `app/controller/...`，避免交叉堆放。
- 避免硬编码：可复用的业务常量/枚举统一维护在 `utils/parame.go`；环境相关（端口/地址/开关/密钥等）统一走 `config.yml` / `config/*.yml`。
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
