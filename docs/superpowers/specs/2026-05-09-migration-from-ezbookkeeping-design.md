# OmniEvent 迁移设计

## 概述

将 ezbookkeeping 项目的可复用模块迁移到 OmniEvent，按层顺序迁移，确保架构清晰、每阶段可独立验证。

**迁移原则：**
- 先基础设施，后业务模块
- 每波迁移后验证可用性
- 保持 OmniEvent 的目录结构不变，源文件按功能复制到目标位置
- 记账特有代码直接跳过，不做无用迁移

---

## 第一波：后端基础设施

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `pkg/errs/` | `pkg/errs/` | 错误码系统，直接复制 |
| `pkg/middlewares/auth.go` | `internal/middleware/auth.go` | JWT 认证，直接复制 |
| `pkg/validators/` | `pkg/validators/` | 验证器集合，直接复制 |
| `pkg/log/` | `pkg/log/` | 日志系统，直接复制 |
| `pkg/utils/` | `pkg/utils/` | 工具函数，对比合并 |
| `pkg/uuid/` | `pkg/uuid/` | UUID 生成器，直接复制 |
| `pkg/response/` | `pkg/response/` | 对比 OmniEvent 后决定保留哪个 |

**适配说明：**
- `pkg/core/context.go` / `pkg/core/token_claims.go` - OmniEvent 已有类似实现 (`pkg/context/web.go`)，需对比差异后决定
- `pkg/response/response.go` - OmniEvent 已有简单实现，需对比合并

---

## 第二波：后端用户模块

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `pkg/models/user.go` | `internal/model/user.go` | 用户模型 |
| `pkg/models/token_record.go` | `internal/model/token_record.go` | Token 记录 |
| `pkg/services/users.go` | `internal/service/user_service.go` | 用户业务逻辑 |
| `pkg/services/tokens.go` | `internal/service/token_service.go` | Token 业务逻辑 |
| `pkg/services/forget_passwords.go` | `internal/service/forget_password_service.go` | 密码重置 |
| `pkg/api/authorizations.go` | `internal/api/auth_handler.go` | 登录/注册 API |
| `pkg/api/tokens.go` | `internal/api/token_handler.go` | Token 管理 API |
| `pkg/api/forget_passwords.go` | `internal/api/forget_password_handler.go` | 密码重置 API |
| `pkg/api/users.go` | `internal/api/user_handler.go` | 用户资料 API |

**需要适配的部分：**
- 数据库表名和字段（ezbookkeeping 用 `t_users`，OmniEvent 需要新的表结构）
- 密码加密方式（第一波已迁移 `pkg/utils/crypto.go`）
- 邮件发送集成（第五波已迁移 `pkg/mail/`）

---

## 第三波：前端基础设施

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `src/lib/services.ts` | `src/lib/services.ts` | API 客户端封装，对比合并 |
| `src/lib/logger.ts` | `src/lib/logger.ts` | 日志工具，直接复制 |
| `src/lib/userstate.ts` | `src/lib/userstate.ts` | 用户状态，对比合并 |
| `src/core/api.ts` | `src/lib/api.ts` | API 错误处理，对比合并 |
| `src/core/datetime.ts` | `src/lib/utils/datetime.ts` | 时间格式化，已有需对比 |
| `src/core/cache.ts` | `src/lib/cache.ts` | 缓存，复制 |
| `src/consts/api.ts` | `src/consts/api.ts` | API 路径常量，复制 |
| `src/models/user.ts` | `src/models/user.ts` | TypeScript 用户模型 |
| `src/models/auth_response.ts` | `src/models/auth.ts` | 认证响应模型 |
| `src/stores/user.ts` | `src/stores/user.ts` | Pinia 用户状态，对比合并 |

---

## 第四波：前端用户页面

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `src/views/desktop/LoginPage.vue` | `src/views/desktop/LoginPage.vue` | 登录页 |
| `src/views/desktop/SignupPage.vue` | `src/views/desktop/SignupPage.vue` | 注册页 |
| `src/views/desktop/ForgetPasswordPage.vue` | `src/views/desktop/ForgetPasswordPage.vue` | 忘记密码页 |
| `src/views/desktop/ResetPasswordPage.vue` | `src/views/desktop/ResetPasswordPage.vue` | 重置密码页 |
| `src/views/desktop/VerifyEmailPage.vue` | `src/views/desktop/VerifyEmailPage.vue` | 邮箱验证页 |

---

## 第五波：邮件模块

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `pkg/mail/` | `pkg/mail/` | 邮件发送模块 |
| `pkg/settings/` | `pkg/settings/` | 配置管理（邮件配置节） |

**需要适配的部分：**
- 邮件模板内容（ezbookkeeping 的模板是记账相关的，需替换为 OmniEvent 的活动相关文案）

---

## 第六波：存储模块

| 源 (ezbookkeeping) | 目标 (OmniEvent) | 说明 |
|------|------|------|
| `pkg/storage/` | `pkg/storage/` | 对象存储模块 |
| `pkg/api/avatars.go` | 新建 `internal/api/avatar_handler.go` | 头像上传 API |

**需要适配的部分：**
- 存储路径和 bucket 名称

---

## 验收标准

每波迁移后需满足：
1. 编译检查通过 `go build ./...` 或 `npm run type-check`
2. 不破坏现有功能
3. 迁移文件数与计划一致

---

## 不迁移的模块

以下模块为记账特有，不适合 OmniEvent：
- `transactions/` - 账目流水
- `accounts/` - 账户
- `explorers/` - 统计图表
- `exchange_rates/` - 汇率
- 移动端 UI 组件
- `models/transaction*.go` - 记账数据模型
