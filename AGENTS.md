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

## 配置文件完整说明

### 主配置文件（config.yml）
```yaml
# 服务配置
server:
  env: prod                              # 环境标识
  debug: false                           # 调试模式
  host: 0.0.0.0                          # API 监听地址
  port: 8006                             # API 监听端口
  apiUrl: http://127.0.0.1:8006          # API 对外访问地址
  downUrl: http://127.0.0.1:8006         # 下载域名
  encrypt: false                         # 是否加密请求/响应

# 后台配置
admin:
  host: 0.0.0.0
  port: 8005

# 静态资源服务
source:
  host: 0.0.0.0
  port: 8007
  apiUrl: http://127.0.0.1              # 静态资源服务地址
  publicBaseUrl: http://127.0.0.1:8007  # 对外访问地址（优先级高于 apiUrl）

# JWT 配置（必配）
jwt:
  secret: "your-secret-key-here"        # JWT 签名密钥（不能为空）

# 密码加密配置
auth:
  passwordSalt: "your-salt-here"        # 密码盐值（可选）

# MySQL 配置
mysql:
  address: ""                           # 完整地址（优先级最高，如 "127.0.0.1:3306"）
  host: "127.0.0.1"                     # 主机地址
  port: 3306                            # 端口
  database: "novel"                     # 数据库名（优先级：database > dbname > name）
  user: "root"                          # 用户名
  password: "root"                      # 密码
  params: "charset=utf8mb4&parseTime=True&loc=Local"  # 连接参数
  pool:                                 # 连接池配置
    maxIdleConns: 25                    # 最大空闲连接数
    maxOpenConns: 100                   # 最大打开连接数
    connMaxLifetimeSeconds: 600         # 连接最大生命周期（秒）
    connMaxIdleTimeSeconds: 300         # 连接最大空闲时间（秒）

# Redis 配置
redis:
  addr: ""                              # 完整地址（优先级最高，如 "127.0.0.1:6379"）
  host: "127.0.0.1"                     # 主机地址
  port: 6379                            # 端口
  password: ""                          # 密码（可选）
  db: 0                                 # 数据库编号

# 日志配置
logs:
  level: -1                             # 日志级别（-1=Debug, 0=Info, 1=Warn, 2=Error）
  path: "logs"                          # 日志路径
  max-size: 50                          # 单文件最大大小（MB）
  max-backups: 100                      # 保留备份文件数量
  max-age: 30                           # 保留天数
  compress: false                       # 是否压缩

# Casbin 权限配置
casbin:
  modelFile: ./public/casbin_conf/rbac_model.conf
  policyFile: ./public/casbin_conf/rbac_policy.csv
```

### 上传配置文件（config/upload.yml，可选）
```yaml
upload:
  baseDir: "./public/upload"           # 本地存储根目录
  publicPathPrefix: "/public/upload"   # 公共访问路径前缀
  maxSizeMB: 50                        # 最大文件大小（MB）
  allowedExts:                         # 允许的文件后缀
    - .jpg
    - .jpeg
    - .png
    - .gif
    - .webp
    - .bmp
    - .mp4
    - .mp3
    - .wav
    - .pdf
    - .doc
    - .docx
    - .xls
    - .xlsx
  allowedMimePrefixes:                 # 允许的 MIME 前缀
    - "image/"
    - "video/"
    - "audio/"
    - "application/pdf"
    - "application/msword"
    - "application/vnd.openxmlformats-officedocument"
```

## 全局变量说明（global/global.go）

项目通过 `global` 包统一管理全局变量，核心实例如下：

### 常用全局变量
- `global.DB`：MySQL 数据库连接（GORM）
- `global.Redis`：Redis 客户端
- `global.WsHub`：WebSocket Hub 实例
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

## 用户模型完整字段说明（app/models/mc_user.go）

### McUser 结构体（表名：mc_user）
| 字段名 | 类型 | JSON 标签 | 说明 |
|--------|------|-----------|------|
| Id | int64 | id | 用户 ID（主键） |
| ParentId | int64 | parent_id | 上级用户 ID（推荐人） |
| Username | string | username | 登录账号 |
| Passwd | string | passwd | 密码（MD5+盐） |
| Nickname | string | nickname | 昵称 |
| Tel | string | tel | 手机号 |
| Pic | string | pic | 头像 URL |
| Email | string | email | 邮箱 |
| Sex | int | sex | 性别（0-未知 1-男 2-女） |
| Text | string | text | 个人简介 |
| Referrer | string | referrer | 上级邀请码 |
| Invitation | string | invitation | 本人邀请码（6 位随机） |
| ParentLink | string | parent_link | 推荐链条（`,id1,id2,...`） |
| Vip | int64 | vip | VIP 标记（0-否 1-是） |
| Rmb | float64 | rmb | 账户余额（人民币） |
| Cion | int64 | cion | 金币余额 |
| Viptime | int64 | viptime | VIP 到期时间（Unix 时间戳） |
| Status | int | status | 账户状态（0-锁定 1-正常） |
| IsGuest | int | is_guest | 游客标记（1-是 0-否） |
| Deviceid | string | deviceid | 游客匿名设备 ID |
| Mark | string | mark | 渠道号 |
| Package | string | package | 应用包名 |
| Imei | string | imei | IMEI 设备号 |
| Oaid | string | oaid | OAID 设备标识 |
| LastLoginTime | int64 | last_login_time | 上次登录时间（Unix 时间戳） |
| Ip | string | ip | 最后登录 IP |
| RegistId | string | regist_id | 推送设备 ID（极光推送等） |
| IsCheckinRemind | int | is_checkin_remind | 签到提醒开关（0-关闭 1-开启） |
| LastRemindTime | int64 | last_remind_time | 最后签到提醒时间 |
| ReportStatus | int | report_status | 上报状态 |
| BookType | int | book_type | 书籍类型偏好 |
| Addtime | int64 | addtime | 注册时间（Unix 时间戳） |
| Uptime | int64 | uptime | 更新时间（Unix 时间戳） |

### 请求结构体（同文件）
- `LoginReq`：登录请求（支持 username/tel/email 三种方式）
- `RegisterReq`：注册请求（username、passwd、nickname、referrer）
- `GuestLoginReq`：游客登录请求（deviceid、referrer、sex）
- `EditUserReq`：编辑用户信息请求（支持 tel/email/nickname/sex/pic/passwd/book_type）
- `LogoffReq`：注销请求
- `UserInfoReq`：查询用户信息请求

## 文件上传服务详解

### 上传接口规范（POST /api/common/upload）
- **路由位置**：`routers/api_routes/common_route.go`
- **Controller**：`app/controller/api/common.go:Upload()`
- **Service**：`app/service/common/file_service/file_service.go:LocalUpload()`
- **表单字段**：
  - `file`：上传文件（必需，multipart/form-data）
  - `dir`：相对子目录（可选，如 `avatar`、`book_cover`）

### 上传响应格式
```json
{
  "code": 0,
  "data": {
    "localPath": "/path/to/public/upload/avatar/xxx.jpg",
    "publicPath": "/public/upload/avatar/xxx.jpg",
    "url": "http://127.0.0.1:8007/public/upload/avatar/xxx.jpg",
    "filename": "20260113_abc123.jpg",
    "size": 102400,
    "contentType": "image/jpeg"
  },
  "msg": "ok"
}
```

### 安全校验机制
1. **文件大小限制**：默认 50MB（可配 `upload.maxSizeMB`）
2. **后缀白名单**：`upload.allowedExts`（如 `.jpg`、`.png`、`.gif`、`.mp4`）
3. **MIME 类型白名单**：`upload.allowedMimePrefixes`（如 `image/`、`video/`）
4. **路径安全**：禁止 `..`、绝对路径、特殊字符
5. **随机文件名**：`{timestamp}_{random}.{ext}`，防止文件名冲突和路径遍历

### 上传配置（config/upload.yml）
```yaml
upload:
  baseDir: "./public/upload"           # 本地存储根目录
  publicPathPrefix: "/public/upload"   # 公共访问路径前缀
  maxSizeMB: 50                        # 最大文件大小（MB）
  allowedExts:                         # 允许的文件后缀
    - .jpg
    - .jpeg
    - .png
    - .gif
    - .webp
    - .mp4
    - .mp3
    - .pdf
  allowedMimePrefixes:                 # 允许的 MIME 前缀
    - "image/"
    - "video/"
    - "audio/"
    - "application/pdf"
```

## WebSocket 架构详解

### Hub 中心架构（pkg/ws/hub.go）
- **设计模式**：单 goroutine 主循环 + 多 goroutine 读写泵
- **核心结构**：
  ```go
  type Hub struct {
    register     chan *Client           // 客户端注册队列
    unregister   chan *Client           // 客户端注销队列
    broadcastAll chan []byte            // 全局广播队列
    chat         chan chatRequest       // 群聊请求队列
    dm           chan dmRequest         // 单聊请求队列

    clients map[*Client]struct{}        // 所有连接（集合）
    users   map[int64]map[*Client]struct{}  // 按 userID 归档（同账号多端）
  }
  ```
- **Hub.Run() 主循环**：
  - 单 goroutine 处理所有状态变更（注册/注销/广播/聊天）
  - 避免 map 并发读写竞态
  - 使用 `select` 多路复用 channel

### 客户端生命周期（pkg/ws/client.go）
1. **连接建立**（`ws.HandleRequest()`）：
   - 升级 HTTP 连接为 WebSocket
   - 解析 URL 参数中的 `token`（可选）
   - 若有 token，验证并提取 `userId` 和 `username`
   - 创建 `Client` 实例，并启动读写泵
2. **读泵**（`client.readPump()`）：
   - 单 goroutine 持续读取客户端消息
   - 设置读超时（60 秒）和最大消息大小（64KB）
   - 解析 JSON 消息，分发到不同处理函数（ping/chat/dm）
3. **写泵**（`client.writePump()`）：
   - 单 goroutine 持续发送消息给客户端
   - 从 `client.send` channel 读取待发送消息
   - 定时发送 Ping frame（54 秒间隔，= 60 * 9/10）
   - 写超时保护（10 秒）
4. **断开连接**：
   - 读泵/写泵任一退出时，触发 `hub.unregister`
   - Hub 从 `clients` 和 `users` 中移除该连接
   - 关闭连接和 `send` channel

### 消息协议详解（pkg/ws/protocol.go）
| 消息类型 | 入站格式 | 出站格式 | 说明 |
|---------|---------|---------|------|
| ping | `{"type":"ping"}` | `{"type":"pong","data":{"ts":1234567890}}` | 心跳检测 |
| chat | `{"type":"chat","data":{"text":"hello"}}` | `{"type":"chat","data":{"text":"hello","userId":123,"username":"user","ts":1234567890}}` | 群聊（需 token） |
| dm | `{"type":"dm","data":{"toUserId":2,"text":"hi"}}` | `{"type":"dm","data":{"text":"hi","fromUserId":1,"fromUsername":"user1","ts":1234567890}}` | 单聊（需 token） |
| error | 无 | `{"type":"error","msg":"错误信息"}` | 错误消息 |

### 群聊实现（pkg/ws/chat.go）
- **鉴权要求**：发送方必须携带有效 token（`client.UserID > 0`）
- **广播规则**：仅发送给所有已鉴权用户（`user.UserID > 0`）
- **消息字段**：包含发送方的 `userId`、`username`、`text`、`ts`（时间戳）

### 单聊实现（pkg/ws/dm.go）
- **鉴权要求**：发送方必须携带有效 token
- **目标查找**：通过 `hub.users[toUserId]` 获取目标用户的所有连接
- **多端同步**：消息发送给目标用户的所有在线设备
- **回显机制**：同时回显给发送方（便于 UI 渲染）
- **消息字段**：包含发送方的 `fromUserId`、`fromUsername`、`text`、`ts`

### 超时与心跳机制
- **读超时**：60 秒（收到任何消息时重置，包括客户端的 Pong frame）
- **写超时**：10 秒
- **Ping 周期**：54 秒（= 60 * 9/10，确保在读超时前发送 Ping）
- **协议层心跳**：除 WebSocket 协议的 Ping/Pong frame 外，还支持 JSON 格式的 `{"type":"ping"}`，便于调试

## WebSocket（基础通信）

- 路由：`GET /api/ws`（使用 `gorilla/websocket`，最小可用：广播/心跳/房间）。
- 连接示例：`ws://127.0.0.1:18016/api/ws?token=<JWT>`（token 可选；传入时会解析并记录 `userId/username`，便于后续扩展鉴权/权限/房间等）。
- 消息协议（JSON）：
  - `{"type":"ping"}`：服务端回 `{"type":"pong","data":{"ts":<unix>}}`
  - `{"type":"chat","data":{"text":"hello"}}`：群聊广播（当前实现无房间；仅对"携带 token 建立连接"的用户开放，返回 `chat`，带 `userId/username/ts`）
  - `{"type":"dm","data":{"toUserId":2,"text":"hi"}}`：单聊（需要在连接时携带有效 token 才有 `userId`）；服务端会把 `dm` 发给目标用户的所有连接，并回显给发送方
- 心跳：服务端会定时发送 WS Ping frame，客户端正常回复 Pong 即可；协议层也支持 `ping/pong` 便于调试。

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
