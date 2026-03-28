# OJ 后台账号与权限体系设计

## 1. 背景与目标

当前项目已经具备前台用户注册登录能力，核心用户表为 `users`，主要服务于做题用户。随着后台管理端进入设计阶段，需要补齐后台账号体系与权限控制方案，满足以下目标：

- 前台做题用户与后台管理人员身份边界清晰
- 后台权限能够按模块进行控制，而不是只靠一个管理员布尔值
- 权限模型可随后台模块增长而扩展，避免后续频繁重构
- 与现有后台路由规划 `/admin/*` 对齐，便于接口和前端菜单统一落地

本设计只覆盖后台账号与权限体系，不涉及后台页面视觉设计与具体接口实现细节。

## 2. 设计结论

本期采用以下总体方案：

- 前后台身份体系分离
- 前台继续使用 `users` 表
- 后台新增 `admin_users` 表
- 后台权限使用 Casbin 实现 RBAC
- 首期后台角色固定为 `admin` 和 `assistant`
- Casbin 存储采用标准 `casbin_rule` 表结构
- 首期引入简版 `operation_logs` 记录后台关键操作

## 3. 身份边界

### 3.1 前台用户

前台身份包括：

- `guest`
- `user`

其核心能力包括：

- 注册登录
- 浏览题目和题单
- 提交代码
- 查看提交记录
- 查看个人做题进度

前台权限以“是否登录”和“是否访问本人数据”为主，不引入 Casbin。

### 3.2 后台用户

后台身份统一为 `admin_user`，仅用于后台管理端登录和授权。

后台首期角色包括：

- `admin`
- `assistant`

其中：

- `admin` 拥有后台全部模块权限
- `assistant` 拥有内容维护与部分运营能力，不具备高风险配置权限

### 3.3 边界原则

必须明确以下原则：

- `users` 不承载后台角色字段
- `admin_users` 不承载前台做题数据
- 前后台账号即使邮箱相同，也视为两套独立身份体系
- 后台最终授权结果以 Casbin 策略为准，而不是表中的单一 `role` 字段

## 4. 权限模型

### 4.1 采用 Casbin RBAC with domain

后台使用 Casbin 的 RBAC with domain 模型，请求维度定义为：

- `sub`：后台用户主体，使用 `admin_user_id`
- `dom`：权限域，首期固定为 `admin`
- `obj`：后台资源模块
- `act`：动作

建议模型如下：

```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

### 4.2 资源与动作定义

首期资源 `obj` 与后台导航一一对应：

- `dashboard`
- `problems`
- `problem_sets`
- `tags`
- `test_cases`
- `judge_configs`
- `submissions`
- `users`
- `settings`

首期动作 `act` 固定为：

- `read`
- `write`

动作语义统一如下：

- `read`：页面访问、列表、详情、查询
- `write`：新增、编辑、删除、启停、配置修改

此设计刻意避免首期引入更细的 `create/update/delete/publish` 粒度，以降低理解和接入成本。

### 4.3 角色权限矩阵

`admin`：

- 对全部后台资源拥有 `read` 和 `write`

`assistant`：

- `dashboard:read`
- `problems:read/write`
- `problem_sets:read/write`
- `tags:read/write`
- `test_cases:read/write`
- `submissions:read`

`assistant` 默认不具备以下权限：

- `judge_configs:read/write`
- `users:read/write`
- `settings:read/write`

这里采用偏保守策略，优先控制高风险管理能力。

## 5. 数据表设计

### 5.1 `users`

`users` 继续作为前台用户表存在，现有字段语义保持不变。

保留字段：

- `email`
- `password_hash`
- `nickname`
- `status`
- `last_login_at`

可选增强字段：

- `avatar_url`
- `email_verified_at`
- `bio`

明确不新增以下字段：

- `role`
- `is_admin`
- `admin_type`
- `permissions`

### 5.2 `admin_users`

新增后台账号表 `admin_users`，建议字段如下：

- `id`
- `email`
- `password_hash`
- `display_name`
- `status`
- `last_login_at`
- `created_at`
- `updated_at`
- `deleted_at`

状态值首期固定为：

- `active`
- `disabled`

索引建议：

- 唯一索引：`email`
- 普通索引：`status`

### 5.3 `casbin_rule`

Casbin 策略表采用通用结构：

- `id`
- `ptype`
- `v0`
- `v1`
- `v2`
- `v3`
- `v4`
- `v5`

字段映射约定如下：

- `p, role, domain, resource, action`
- `g, admin_user_id, role, domain`

索引建议：

- 唯一索引：`ptype, v0, v1, v2, v3, v4, v5`
- 可选辅助索引：`v0`、`v1`

策略示例：

```text
p, admin, admin, problems, read
p, admin, admin, problems, write
p, admin, admin, settings, write

p, assistant, admin, problems, read
p, assistant, admin, problems, write
p, assistant, admin, problem_sets, write
p, assistant, admin, submissions, read

g, 1, admin, admin
g, 2, assistant, admin
```

### 5.4 `operation_logs`

首期建议新增简版后台操作日志表 `operation_logs`，用于记录高价值管理操作。

建议字段：

- `id`
- `admin_user_id`
- `resource`
- `action`
- `target_type`
- `target_id`
- `request_id`
- `detail_json`
- `ip`
- `user_agent`
- `created_at`

首期优先记录：

- 题目新增与编辑
- 题单新增与编辑
- 标签维护
- 测试用例维护
- 判题配置修改
- 用户状态变更
- 系统配置修改

## 6. 后台请求授权流程

后台请求统一走以下链路：

1. 后台用户登录成功后获得后台 JWT
2. 中间件解析 JWT，得到 `admin_user_id`
3. 查询后台用户状态，要求账号存在且 `status = active`
4. 将请求映射为 `(sub, dom, obj, act)`
5. 调用 Casbin 执行 `Enforce`
6. 通过则进入业务处理，拒绝则返回 `403`

状态异常时的处理约定：

- 未登录或 token 无效：返回 `401`
- 后台账号被禁用：返回 `403`
- 权限不足：返回 `403`

## 7. 路由与权限映射

与当前 PRD 路由规划对应，首期建议映射如下：

- `GET /admin` -> `dashboard:read`
- `GET /admin/problems` -> `problems:read`
- `POST /admin/problems` -> `problems:write`
- `PUT /admin/problems/:id` -> `problems:write`
- `GET /admin/problem-sets` -> `problem_sets:read`
- `POST /admin/problem-sets` -> `problem_sets:write`
- `PUT /admin/problem-sets/:id` -> `problem_sets:write`
- `GET /admin/tags` -> `tags:read`
- `POST /admin/tags` -> `tags:write`
- `PUT /admin/tags/:id` -> `tags:write`
- `GET /admin/test-cases` -> `test_cases:read`
- `POST /admin/test-cases` -> `test_cases:write`
- `PUT /admin/test-cases/:id` -> `test_cases:write`
- `GET /admin/judge-configs` -> `judge_configs:read`
- `PUT /admin/judge-configs/:id` -> `judge_configs:write`
- `GET /admin/submissions` -> `submissions:read`
- `GET /admin/users` -> `users:read`
- `PATCH /admin/users/:id/status` -> `users:write`
- `GET /admin/settings` -> `settings:read`
- `PATCH /admin/settings` -> `settings:write`

后台前端菜单是否展示可复用同一权限映射，从而让前端展示策略与后端鉴权保持一致。

## 8. 初始化与演进策略

### 8.1 首期初始化

数据库迁移时完成以下初始化动作：

- 创建 `admin_users` 表
- 创建 `casbin_rule` 表
- 创建 `operation_logs` 表
- 创建默认超级管理员账号
- 写入 `admin` 和 `assistant` 的初始策略

### 8.2 首期不做的内容

为控制范围，以下内容不在本期实现：

- 后台角色管理页面
- 后台权限点可视化配置页面
- 更细粒度动作拆分
- 前后台账号自动绑定
- 多域或多租户权限模型

### 8.3 后续扩展方向

当后台复杂度提升后，可在不推翻现有设计的前提下增加：

- 新角色，例如 `problem_editor`、`ops_admin`
- 更细粒度动作，例如 `delete`、`publish`
- 权限元数据表，用于角色与权限的后台可视化管理
- 审计日志检索与导出能力

## 9. 测试要求

至少覆盖以下验证场景：

- `admin` 可以访问全部后台模块
- `assistant` 可以操作题目、题单、标签、测试用例
- `assistant` 无法修改判题配置、用户管理、系统配置
- 被禁用的 `admin_user` 无法登录后台
- 未认证后台请求返回 `401`
- 已认证但未授权请求返回 `403`
- Casbin 策略初始化后能够正确完成角色授权判断

## 10. 最终决策摘要

本设计最终确认如下：

- 前后台账号与权限体系分开
- 前台使用 `users`
- 后台使用 `admin_users`
- 后台权限由 Casbin 统一裁决
- 后台首期角色为 `admin` 与 `assistant`
- `users` 表不增加后台角色语义字段
- 首期同步引入 `operation_logs`

该方案兼顾了首期实现成本、后续扩展能力和权限边界清晰度，适合作为 OJ 后台管理端的基础权限设计。
