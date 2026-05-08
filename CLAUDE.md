# CLAUDE.md - OmniEvent 指南

## 📋 项目概览

**OmniEvent** 是一个通用的活动全生命周期管理平台。通过 **Schema 驱动** 的核心引擎，支持高度自定义的活动创建、动态表单报名、多模式签到及多维数据统计。

* **核心理念**：灵活配置 (Flexibility)、流程自动化 (Flow)、协议驱动 (Protocol-first)。
* **开发环境**：macOS / UTC+8 / `zzhao`。

---

## 🛠 技术栈

### 后端 (Backend)

* **框架**: Gin (Golang)
* **数据库**: PostgreSQL 15+ (关键使用 `JSONB` 存储 Schema 和动态数据)
* **缓存/锁**: Redis
* **ORM**: XORM
* **配置**: INI 配置文件 + 环境变量覆盖
* **规范**: 符合 Clean Architecture，逻辑层与持久层分离。

### 前端 (Frontend)

* **框架**: Vue 3 (Composition API) + TypeScript
* **UI 库**: Vuetify 3
* **表单引擎**: Formily 2 (核心：`@formily/core`, `@formily/vue`)
* **状态管理**: Pinia
* **构建**: Vite

---

## 💻 常用命令

### 后端

| 命令 | 说明 |
|------|------|
| `go run cmd/server/main.go` | 运行服务 |
| `go build -o bin/server cmd/server/main.go` | 构建项目 |
| `go mod tidy` | 整理依赖 |
| `go build ./...` | 编译检查 |

### 前端

| 命令 | 说明 |
|------|------|
| `npm run dev` 或 `pnpm dev` | 启动开发环境 |
| `npm run build` | 构建发布 |
| `npm run type-check` 或 `vue-tsc --noEmit` | 类型检查 |

---

## ⚙️ 配置系统

### 配置文件

配置文件路径查找顺序：
1. `conf/omnievent.ini` (项目目录)
2. `conf/omnievent.conf`
3. `/etc/omnievent/omnievent.ini`
4. 环境变量 `OMNIEVENT_CONFIG`

### 环境变量覆盖

环境变量格式：`OMNIEVENT_SECTION_KEY` (全大写)

示例：
- `OMNIEVENT_DB_HOST` → database.host
- `OMNIEVENT_JWT_SECRET` → jwt.secret
- `OMNIEVENT_SERVER_PORT` → server.http_port

### 配置节

| 节 | 说明 |
|----|------|
| `global` | 应用全局配置 (app_name, mode, version) |
| `server` | 服务器配置 (protocol, http_port, enable_gzip) |
| `database` | 数据库配置 (host, port, sslmode) |
| `redis` | Redis 配置 (host, port, password) |
| `jwt` | JWT 配置 (secret, expire_hours, refresh_days) |
| `security` | 安全配置 (secret_key, token_expire_seconds) |
| `user` | 用户配置 (enable_register, min/max_username_length) |

---

## 📏 编码规范

### 通用规则

* **命名**: 变量与函数名使用 `camelCase`，组件名使用 `PascalCase`，后端 Go 代码使用 `PascalCase` (公开) 或 `camelCase` (私有)。
* **错误处理**: 严禁忽略后端错误；前端使用 Vuetify 的 `v-alert` 或消息条展示业务异常。

### 错误码系统

| 范围 | 类别 |
|------|------|
| 100000+ | System 错误 |
| 2000-2999 | User 错误 |
| 3000-3999 | Token 错误 |
| 4000-4999 | Validation 错误 |

### 活动引擎 (Formily & JSONB)

* **协议一致性**: 前端生成的 `JSON Schema` 必须符合 Formily 规范，后端 `JSONB` 字段严禁手动修改结构，需通过接口更新。
* **桥接逻辑**: 所有的 Vuetify 组件适配逻辑应统一存放在 `src/components/FormilyBridge`。
* **联动设计**: 优先使用 Formily 的 `Reaction` 实现字段联动（Linkages），减少业务组件内的 `v-if` 硬编码。

---

## 📂 项目结构

```text
OmniEvent/
├── CLAUDE.md
│
├── omnievent-backend/           # Gin 后端
│   ├── cmd/
│   │   ├── server/main.go      # 服务入口、路由、中间件
│   │   ├── config.go           # INI配置加载、环境变量覆盖
│   │   ├── database.go         # XORM 数据库连接
│   │   └── redis.go            # Redis 连接
│   ├── conf/
│   │   └── omnievent.ini       # 配置文件模板
│   ├── internal/
│   │   ├── api/                # API 处理器
│   │   │   ├── user_handler.go # 用户资料 API
│   │   │   ├── auth_handler.go # 登录/注册/密码重置 API
│   │   │   └── token_handler.go # Token 管理 API
│   │   ├── service/            # 业务逻辑层
│   │   │   ├── user_service.go # 用户业务逻辑
│   │   │   └── token_service.go # Token 业务逻辑
│   │   ├── repository/         # 数据访问层 (XORM)
│   │   │   └── user_repository.go
│   │   ├── model/              # 数据模型
│   │   │   └── user.go         # User, TokenRecord 模型
│   │   └── middleware/         # 中间件
│   │       └── auth.go         # JWT 认证中间件
│   └── pkg/
│       ├── context/            # WebContext 封装
│       │   └── web.go         # Token Claims、请求上下文
│       ├── errs/               # 错误码定义
│       │   ├── error.go       # 错误基础结构 (Category, Code)
│       │   ├── user.go        # 用户错误 (2001-2007)
│       │   ├── token.go       # Token 错误 (3001-3005)
│       │   ├── database.go    # 数据库错误
│       │   └── validation.go   # 验证错误
│       ├── response/           # 统一 API 响应
│       │   └── response.go    # Success/Error 响应封装
│       └── utils/             # 工具函数
│           ├── time.go        # 时间工具
│           ├── string.go      # 字符串工具
│           └── crypto.go      # 加密工具 (SHA256, Salt)
│
└── omnievent-frontend/        # Vue 3 前端
    ├── src/
    │   ├── main.ts            # 入口文件
    │   ├── App.vue            # 根组件
    │   ├── router/index.ts     # 路由配置
    │   ├── stores/             # Pinia 状态管理
    │   │   ├── index.ts
    │   │   └── user.ts         # 用户状态 (login, register, logout)
    │   ├── lib/
    │   │   ├── api.ts          # API 错误处理、错误码常量
    │   │   ├── services.ts     # API 客户端 (userService, tokenService)
    │   │   ├── userstate.ts    # Token 状态管理
    │   │   ├── logger.ts       # 日志工具 (debug/info/warn/error)
    │   │   ├── utils/          # 前端工具
    │   │   │   ├── datetime.ts # 时间格式化
    │   │   │   ├── string.ts   # 字符串处理
    │   │   │   └── index.ts
    │   │   └── validators/     # 表单验证
    │   │       └── index.ts   # required, email, minLength 等
    │   ├── models/             # TypeScript 类型定义
    │   │   └── user.ts
    │   └── views/
    │       └── desktop/        # 桌面端视图
    │           ├── LoginPage.vue
    │           ├── SignupPage.vue
    │           └── ProfilePage.vue
    ├── index.html
    ├── vite.config.ts
    ├── tsconfig.json
    └── package.json
```

---

## 🎯 当前任务 (Milestone 1)

* [x] ~~完成 PostgreSQL 基础表结构设计（含 Activity 及其 JSONB 字段）~~ → 当前为用户模块
* [ ] 用户注册/登录 API
* [ ] JWT Token 认证
* [ ] 用户资料管理
* [ ] 前端页面集成

---

## Project Rules

*   **语言规则**: 始终使用中文进行对话。
*   **Documentation Location**: All design documents, architecture specs, and planning files MUST be created or updated in the **current root directory** (`./`).
*   **No Subdirectories for Docs**: Do not create or use a `docs/` or `plans/` folder unless explicitly instructed for a specific task.
*   **Naming Convention**: Design documents should follow the pattern `DESIGN_*.md` and planning documents should follow `PLAN_*.md`.
