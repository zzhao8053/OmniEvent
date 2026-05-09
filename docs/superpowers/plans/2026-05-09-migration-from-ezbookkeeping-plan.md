# OmniEvent Migration from ezbookkeeping Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 ezbookkeeping 的可复用模块迁移到 OmniEvent，按层顺序执行，确保每阶段可独立验证。

**Architecture:** 按波次迁移：后端基础设施 → 后端用户模块 → 前端基础设施 → 前端用户页面 → 邮件/存储。按功能垂直迁移，每波迁移后验证编译通过。

**Tech Stack:** Go (Gin), Vue 3 (Composition API), PostgreSQL (JSONB), Redis

---

## 第一波：后端基础设施

### Task 1: 迁移 errs 错误码系统

**Files:**
- Create: `omnievent-backend/pkg/errs/error.go`
- Create: `omnievent-backend/pkg/errs/user.go`
- Create: `omnievent-backend/pkg/errs/token.go`
- Create: `omnievent-backend/pkg/errs/database.go`
- Create: `omnievent-backend/pkg/errs/validation.go`
- Create: `omnievent-backend/pkg/errs/category.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/`

- [ ] **Step 1: 读取源文件**

Run: `cat /home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/error.go`

- [ ] **Step 2: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/errs`

- [ ] **Step 3: 复制 error.go**

Copy file content to `omnievent-backend/pkg/errs/error.go`

- [ ] **Step 4: 复制 user.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/user.go` to `omnievent-backend/pkg/errs/user.go`

- [ ] **Step 5: 复制 token.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/token.go` to `omnievent-backend/pkg/errs/token.go`

- [ ] **Step 6: 复制 database.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/database.go` to `omnievent-backend/pkg/errs/database.go`

- [ ] **Step 7: 复制 validation.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/validation.go` to `omnievent-backend/pkg/errs/validation.go`

- [ ] **Step 8: 复制 category.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/errs/category.go` to `omnievent-backend/pkg/errs/category.go`

- [ ] **Step 9: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/errs/...`
Expected: PASS

- [ ] **Step 10: 提交**

```bash
git add omnievent-backend/pkg/errs/ && git commit -m "chore(backend): migrate errs error code system"
```

---

### Task 2: 迁移 validators 验证器集合

**Files:**
- Create: `omnievent-backend/pkg/validators/validator.go`
- Create: `omnievent-backend/pkg/validators/user.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/validators/`

- [ ] **Step 1: 读取源文件**

Run: `ls /home/zzhao/zzhao/githubs/ezbookkeeping/pkg/validators/`

- [ ] **Step 2: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/validators`

- [ ] **Step 3: 复制所有验证器文件**

Copy all files from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/validators/` to `omnievent-backend/pkg/validators/`

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/validators/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/validators/ && git commit -m "chore(backend): migrate validators"
```

---

### Task 3: 迁移 log 日志系统

**Files:**
- Create: `omnievent-backend/pkg/log/logger.go`
- Create: `omnievent-backend/pkg/log/file.go`
- Create: `omnievent-backend/pkg/log/rotate.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/log/`

- [ ] **Step 1: 读取源文件**

Run: `ls /home/zzhao/zzhao/githubs/ezbookkeeping/pkg/log/`

- [ ] **Step 2: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/log`

- [ ] **Step 3: 复制所有日志文件**

Copy all files from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/log/` to `omnievent-backend/pkg/log/`

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/log/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/log/ && git commit -m "chore(backend): migrate log system"
```

---

### Task 4: 迁移 utils 工具函数

**Files:**
- Create: `omnievent-backend/pkg/utils/crypto.go`
- Create: `omnievent-backend/pkg/utils/string.go`
- Create: `omnievent-backend/pkg/utils/time.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/utils/`

- [ ] **Step 1: 对比现有文件**

Run: `diff omnievent-backend/pkg/utils/crypto.go /home/zzhao/zzhao/githubs/ezbookkeeping/pkg/utils/crypto.go` 或阅读两者内容

- [ ] **Step 2: 比较后决定合并策略**

如果 ezbookkeeping 版本功能更全，覆盖 OmniEvent 现有文件；否则保留 OmniEvent 版本

- [ ] **Step 3: 复制 string.go 和 time.go**

Copy from ezbookkeeping if they have additional useful functions not in OmniEvent version

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/utils/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/utils/ && git commit -m "chore(backend): migrate utils (compare and merge)"
```

---

### Task 5: 迁移 uuid 生成器

**Files:**
- Create: `omnievent-backend/pkg/uuid/uuid.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/uuid/`

- [ ] **Step 1: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/uuid`

- [ ] **Step 2: 复制 uuid.go**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/uuid/uuid.go` to `omnievent-backend/pkg/uuid/uuid.go`

- [ ] **Step 3: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/uuid/...`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add omnievent-backend/pkg/uuid/ && git commit -m "chore(backend): migrate uuid generator"
```

---

### Task 6: 迁移 middlewares 认证中间件

**Files:**
- Create: `omnievent-backend/internal/middleware/auth.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/middlewares/authorization.go`

**Note:** OmniEvent 已有一个 `internal/middleware/auth.go`，需对比后合并

- [ ] **Step 1: 读取 OmniEvent 现有 auth.go**

Read: `omnievent-backend/internal/middleware/auth.go`

- [ ] **Step 2: 读取 ezbookkeeping authorization.go**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/middlewares/authorization.go`

- [ ] **Step 3: 合并两个版本的功能**

保留 OmniEvent 的结构，融入 ezbookkeeping 的双因素认证等扩展功能

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./internal/middleware/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/internal/middleware/ && git commit -m "chore(backend): migrate and merge auth middleware"
```

---

### Task 7: 迁移 response 响应封装

**Files:**
- Modify: `omnievent-backend/pkg/response/response.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/response/response.go`

- [ ] **Step 1: 读取 OmniEvent response.go**

Read: `omnievent-backend/pkg/response/response.go`

- [ ] **Step 2: 读取 ezbookkeeping response.go**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/response/response.go`

- [ ] **Step 3: 比较差异，决定合并策略**

如果 ezbookkeeping 版本有 OmniEvent 缺少的功能（如分页封装），则融入；否则保留 OmniEvent 版本

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/response/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/response/ && git commit -m "chore(backend): compare and merge response package"
```

---

## 第二波：后端用户模块

### Task 8: 迁移 user 模型

**Files:**
- Create: `omnievent-backend/internal/model/user.go`
- Create: `omnievent-backend/internal/model/token_record.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/models/user.go` 和 `token_record.go`

- [ ] **Step 1: 读取源文件**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/models/user.go`
Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/models/token_record.go`

- [ ] **Step 2: 创建目标目录**

Run: `mkdir -p omnievent-backend/internal/model`

- [ ] **Step 3: 复制并适配模型**

基于 ezbookkeeping 模型，创建 OmniEvent 版本（表名、字段适配）

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./internal/model/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/internal/model/ && git commit -m "feat(backend): migrate user and token models"
```

---

### Task 9: 迁移 user service

**Files:**
- Create: `omnievent-backend/internal/service/user_service.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/users.go`

- [ ] **Step 1: 读取源文件**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/users.go`

- [ ] **Step 2: 创建目标目录**

Run: `mkdir -p omnievent-backend/internal/service`

- [ ] **Step 3: 复制并适配 service**

适配到 OmniEvent 的 model 和 repository

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./internal/service/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/internal/service/user_service.go && git commit -m "feat(backend): migrate user service"
```

---

### Task 10: 迁移 token service

**Files:**
- Create: `omnievent-backend/internal/service/token_service.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/tokens.go`

- [ ] **Step 1: 复制 token service**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/tokens.go` to `omnievent-backend/internal/service/token_service.go`

- [ ] **Step 2: 适配**

适配到 OmniEvent 的 model

- [ ] **Step 3: 编译检查**

Run: `cd omnievent-backend && go build ./internal/service/...`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add omnievent-backend/internal/service/token_service.go && git commit -m "feat(backend): migrate token service"
```

---

### Task 11: 迁移 forget password service

**Files:**
- Create: `omnievent-backend/internal/service/forget_password_service.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/forget_passwords.go`

- [ ] **Step 1: 复制 forget password service**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/services/forget_passwords.go` to `omnievent-backend/internal/service/forget_password_service.go`

- [ ] **Step 2: 适配**

适配到 OmniEvent 的 model 和邮件模块

- [ ] **Step 3: 编译检查**

Run: `cd omnievent-backend && go build ./internal/service/...`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add omnievent-backend/internal/service/forget_password_service.go && git commit -m "feat(backend): migrate forget password service"
```

---

### Task 12: 迁移 auth handler (登录/注册)

**Files:**
- Create: `omnievent-backend/internal/api/auth_handler.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/authorizations.go`

- [ ] **Step 1: 读取 OmniEvent 现有 handler 结构**

Read: `omnievent-backend/cmd/server/main.go` 了解现有路由结构

- [ ] **Step 2: 复制 auth handler**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/authorizations.go` to `omnievent-backend/internal/api/auth_handler.go`

- [ ] **Step 3: 适配**

适配到 OmniEvent 的 service 和 model

- [ ] **Step 4: 在 main.go 中注册路由**

在 main.go 添加 auth handler 的路由注册代码

- [ ] **Step 5: 编译检查**

Run: `cd omnievent-backend && go build ./...`
Expected: PASS

- [ ] **Step 6: 提交**

```bash
git add omnievent-backend/internal/api/auth_handler.go omnievent-backend/cmd/server/main.go && git commit -m "feat(backend): migrate auth handler and register routes"
```

---

### Task 13: 迁移其他 API handlers

**Files:**
- Create: `omnievent-backend/internal/api/token_handler.go`
- Create: `omnievent-backend/internal/api/forget_password_handler.go`
- Create: `omnievent-backend/internal/api/user_handler.go`

**Source:** ezbookkeeping 的 `pkg/api/tokens.go`, `forget_passwords.go`, `users.go`

- [ ] **Step 1: 复制 token_handler**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/tokens.go` to `omnievent-backend/internal/api/token_handler.go`

- [ ] **Step 2: 复制 forget_password_handler**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/forget_passwords.go` to `omnievent-backend/internal/api/forget_password_handler.go`

- [ ] **Step 3: 复制 user_handler**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/users.go` to `omnievent-backend/internal/api/user_handler.go`

- [ ] **Step 4: 注册所有路由到 main.go**

在 main.go 添加所有新 handler 的路由

- [ ] **Step 5: 编译检查**

Run: `cd omnievent-backend && go build ./...`
Expected: PASS

- [ ] **Step 6: 提交**

```bash
git add omnievent-backend/internal/api/ omnievent-backend/cmd/server/main.go && git commit -m "feat(backend): migrate remaining API handlers"
```

---

## 第三波：前端基础设施

### Task 14: 迁移 lib/services.ts (API 客户端)

**Files:**
- Modify: `omnievent-frontend/src/lib/services.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/services.ts`

- [ ] **Step 1: 读取 OmniEvent 现有 services.ts**

Read: `omnievent-frontend/src/lib/services.ts`

- [ ] **Step 2: 读取 ezbookkeeping services.ts**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/services.ts`

- [ ] **Step 3: 对比合并**

融入 ezbookkeeping 的 API 封装、错误处理、认证逻辑到 OmniEvent 版本

- [ ] **Step 4: 类型检查**

Run: `cd omnievent-frontend && npm run type-check` 或 `vue-tsc --noEmit`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-frontend/src/lib/services.ts && git commit -m "feat(frontend): migrate and merge API services"
```

---

### Task 15: 迁移 lib/logger.ts

**Files:**
- Create: `omnievent-frontend/src/lib/logger.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/logger.ts`

- [ ] **Step 1: 复制 logger.ts**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/logger.ts` to `omnievent-frontend/src/lib/logger.ts`

- [ ] **Step 2: 编译检查** (logger 不涉及编译，仅确保文件存在)

- [ ] **Step 3: 提交**

```bash
git add omnievent-frontend/src/lib/logger.ts && git commit -m "feat(frontend): migrate logger"
```

---

### Task 16: 迁移 lib/userstate.ts

**Files:**
- Modify: `omnievent-frontend/src/lib/userstate.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/userstate.ts`

- [ ] **Step 1: 对比现有文件**

Read: `omnievent-frontend/src/lib/userstate.ts`
Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/src/lib/userstate.ts`

- [ ] **Step 2: 合并功能**

- [ ] **Step 3: 提交**

```bash
git add omnievent-frontend/src/lib/userstate.ts && git commit -m "feat(frontend): migrate and merge userstate"
```

---

### Task 17: 迁移 core/api.ts 到 lib/api.ts

**Files:**
- Create: `omnievent-frontend/src/lib/api.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/core/api.ts`

- [ ] **Step 1: 复制 api.ts**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/core/api.ts` to `omnievent-frontend/src/lib/api.ts`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/lib/api.ts && git commit -m "feat(frontend): migrate core api"
```

---

### Task 18: 迁移前端缓存

**Files:**
- Create: `omnievent-frontend/src/lib/cache.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/core/cache.ts`

- [ ] **Step 1: 复制 cache.ts**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/core/cache.ts` to `omnievent-frontend/src/lib/cache.ts`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/lib/cache.ts && git commit -m "feat(frontend): migrate cache"
```

---

### Task 19: 迁移 consts/api.ts

**Files:**
- Create: `omnievent-frontend/src/consts/api.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/consts/api.ts`

- [ ] **Step 1: 创建目标目录**

Run: `mkdir -p omnievent-frontend/src/consts`

- [ ] **Step 2: 复制 api.ts**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/consts/api.ts` to `omnievent-frontend/src/consts/api.ts`

- [ ] **Step 3: 适配 API 路径**

修改为 OmniEvent 的 API 前缀和路径

- [ ] **Step 4: 提交**

```bash
git add omnievent-frontend/src/consts/api.ts && git commit -m "feat(frontend): migrate API constants"
```

---

### Task 20: 迁移 models/user.ts

**Files:**
- Modify: `omnievent-frontend/src/models/user.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/models/user.ts` 和 `auth_response.ts`

- [ ] **Step 1: 读取 OmniEvent 现有 user.ts**

Read: `omnievent-frontend/src/models/user.ts`

- [ ] **Step 2: 复制并补充模型**

基于 ezbookkeeping 补充 OmniEvent 缺少的字段

- [ ] **Step 3: 创建 auth.ts**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/models/auth_response.ts` to `omnievent-frontend/src/models/auth.ts`

- [ ] **Step 4: 提交**

```bash
git add omnievent-frontend/src/models/ && git commit -m "feat(frontend): migrate user models"
```

---

### Task 21: 迁移 stores/user.ts

**Files:**
- Modify: `omnievent-frontend/src/stores/user.ts`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/stores/user.ts`

- [ ] **Step 1: 读取 OmniEvent 现有 user store**

Read: `omnievent-frontend/src/stores/user.ts`

- [ ] **Step 2: 对比 ezbookkeeping user store**

Read: `/home/zzhao/zzhao/githubs/ezbookkeeping/src/stores/user.ts`

- [ ] **Step 3: 合并功能**

融入登录/注册/登出/刷新 token 等逻辑

- [ ] **Step 4: 提交**

```bash
git add omnievent-frontend/src/stores/user.ts && git commit -m "feat(frontend): migrate and merge user store"
```

---

## 第四波：前端用户页面

### Task 22: 迁移 LoginPage.vue

**Files:**
- Modify: `omnievent-frontend/src/views/desktop/LoginPage.vue`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/LoginPage.vue`

- [ ] **Step 1: 读取 OmniEvent 现有 LoginPage.vue**

Read: `omnievent-frontend/src/views/desktop/LoginPage.vue`

- [ ] **Step 2: 复制 ezbookkeeping LoginPage.vue 内容到 OmniEvent 版本**

融入 ezbookkeeping 的双因素认证支持等扩展功能

- [ ] **Step 3: 编译检查**

Run: `cd omnievent-frontend && npm run type-check`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add omnievent-frontend/src/views/desktop/LoginPage.vue && git commit -m "feat(frontend): migrate login page"
```

---

### Task 23: 迁移 SignupPage.vue

**Files:**
- Modify: `omnievent-frontend/src/views/desktop/SignupPage.vue`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/SignupPage.vue`

- [ ] **Step 1: 复制注册页**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/SignupPage.vue` to `omnievent-frontend/src/views/desktop/SignupPage.vue`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/views/desktop/SignupPage.vue && git commit -m "feat(frontend): migrate signup page"
```

---

### Task 24: 迁移 ForgetPasswordPage.vue

**Files:**
- Create: `omnievent-frontend/src/views/desktop/ForgetPasswordPage.vue`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/ForgetPasswordPage.vue`

- [ ] **Step 1: 复制忘记密码页**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/ForgetPasswordPage.vue` to `omnievent-frontend/src/views/desktop/ForgetPasswordPage.vue`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/views/desktop/ForgetPasswordPage.vue && git commit -m "feat(frontend): migrate forget password page"
```

---

### Task 25: 迁移 ResetPasswordPage.vue

**Files:**
- Create: `omnievent-frontend/src/views/desktop/ResetPasswordPage.vue`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/ResetPasswordPage.vue`

- [ ] **Step 1: 复制重置密码页**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/ResetPasswordPage.vue` to `omnievent-frontend/src/views/desktop/ResetPasswordPage.vue`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/views/desktop/ResetPasswordPage.vue && git commit -m "feat(frontend): migrate reset password page"
```

---

### Task 26: 迁移 VerifyEmailPage.vue

**Files:**
- Create: `omnievent-frontend/src/views/desktop/VerifyEmailPage.vue`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/VerifyEmailPage.vue`

- [ ] **Step 1: 复制邮箱验证页**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/src/views/desktop/VerifyEmailPage.vue` to `omnievent-frontend/src/views/desktop/VerifyEmailPage.vue`

- [ ] **Step 2: 提交**

```bash
git add omnievent-frontend/src/views/desktop/VerifyEmailPage.vue && git commit -m "feat(frontend): migrate verify email page"
```

---

## 第五波：邮件模块

### Task 27: 迁移 mail 邮件模块

**Files:**
- Create: `omnievent-backend/pkg/mail/mail.go`
- Create: `omnievent-backend/pkg/mail/template.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/mail/`

- [ ] **Step 1: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/mail`

- [ ] **Step 2: 复制所有邮件文件**

Copy all files from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/mail/` to `omnievent-backend/pkg/mail/`

- [ ] **Step 3: 适配模板内容**

将邮件模板内容修改为 OmniEvent 相关的文案（活动管理平台）

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/mail/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/mail/ && git commit -m "feat(backend): migrate mail module"
```

---

## 第六波：存储模块

### Task 28: 迁移 storage 存储模块

**Files:**
- Create: `omnievent-backend/pkg/storage/storage.go`
- Create: `omnievent-backend/pkg/storage/local.go`
- Create: `omnievent-backend/pkg/storage/minio.go`
- Create: `omnievent-backend/pkg/storage/webdav.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/storage/`

- [ ] **Step 1: 创建目标目录**

Run: `mkdir -p omnievent-backend/pkg/storage`

- [ ] **Step 2: 复制所有存储文件**

Copy all files from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/storage/` to `omnievent-backend/pkg/storage/`

- [ ] **Step 3: 适配存储路径和 bucket 名称**

- [ ] **Step 4: 编译检查**

Run: `cd omnievent-backend && go build ./pkg/storage/...`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add omnievent-backend/pkg/storage/ && git commit -m "feat(backend): migrate storage module"
```

---

### Task 29: 创建 avatar API handler

**Files:**
- Create: `omnievent-backend/internal/api/avatar_handler.go`

**Source:** `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/avatars.go`

- [ ] **Step 1: 复制 avatar handler**

Copy from `/home/zzhao/zzhao/githubs/ezbookkeeping/pkg/api/avatars.go` to `omnievent-backend/internal/api/avatar_handler.go`

- [ ] **Step 2: 注册路由**

在 main.go 添加头像上传路由

- [ ] **Step 3: 编译检查**

Run: `cd omnievent-backend && go build ./...`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add omnievent-backend/internal/api/avatar_handler.go omnievent-backend/cmd/server/main.go && git commit -m "feat(backend): add avatar upload API"
```

---

## 验收总检查

### 编译验证

- [ ] **Step 1: 后端编译**

Run: `cd omnievent-backend && go build ./...`
Expected: PASS，无报错

- [ ] **Step 2: 前端类型检查**

Run: `cd omnievent-frontend && npm run type-check`
Expected: PASS，无类型错误

- [ ] **Step 3: 提交所有遗留更改**

检查 git status，确保所有迁移文件已提交

