# 在线代码评测系统 MVP 设计

## 背景

当前项目是一个轻量 Go 服务骨架，已经存在 `internal/model/base.go`，其中基础模型约定如下：

- 主键字段使用 `uint`
- 时间字段使用 `int64`
- 软删除使用 `gorm.DeletedAt`

本设计基于这些约定，先定义一个可快速落地的 OJ MVP 数据模型，覆盖以下范围：

- 用户邮箱注册、登录
- 题库/题单
- 题目标签
- 题目
- 用户每次提交记录

不在本次 MVP 范围内：

- 邮箱验证码、激活邮件
- 真正的异步判题任务拆分
- 测试用例管理
- 代码运行沙箱
- 讨论区、题解、排行榜

## 设计结论

采用标准关系型建模，核心实体如下：

- `users`
- `problem_sets`
- `problems`
- `tags`
- `problem_set_problems`
- `problem_tags`
- `submissions`
- `user_problem_stats`

其中：

- 题库和题目是多对多关系
- 题目和标签是多对多关系
- 提交记录按“每次提交一条”存储
- 用户题目聚合状态单独建表，避免列表页频繁聚合 `submissions`

## 为什么选这个方案

相比把标签塞进 JSON 或只保留提交表的极简设计，这个方案的优点是：

- 查询路径稳定，GORM 映射清晰
- 后续做“题库页 / 标签筛选 / 我的做题状态”时不需要大改
- 既保留每次提交明细，也保留用户题目聚合状态，兼顾可追溯性和查询效率

相比一开始就做测试用例、判题任务、结果明细等复杂表，这个方案足够轻，适合 MVP。

## 表设计

### 1. `users`

用途：注册、登录、账号状态维护。

建议字段：

- `id`
- `email`：唯一索引
- `password_hash`
- `nickname`
- `status`：`active` / `disabled`
- `last_login_at`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一索引：`email`
- 普通索引：`status`

建议模型：

```go
type User struct {
	BaseModel
	Email        string `gorm:"size:128;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Nickname     string `gorm:"size:64;not null" json:"nickname"`
	Status       string `gorm:"size:32;not null;default:active;index" json:"status"`
	LastLoginAt  int64  `gorm:"not null;default:0" json:"last_login_at"`
}
```

### 2. `problem_sets`

用途：题库 / 题单，例如 `hot100`、`字节必刷`、`腾讯必刷`。

建议字段：

- `id`
- `name`
- `slug`：唯一索引，用于 URL 或编码
- `description`
- `status`
- `sort_order`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一索引：`slug`
- 普通索引：`status`
- 排序索引可按实际列表需求补充

建议模型：

```go
type ProblemSet struct {
	BaseModel
	Name        string `gorm:"size:128;not null" json:"name"`
	Slug        string `gorm:"size:128;not null;uniqueIndex" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"size:32;not null;default:active;index" json:"status"`
	SortOrder   int    `gorm:"not null;default:0" json:"sort_order"`
}
```

### 3. `problems`

用途：题目主表。

建议字段：

- `id`
- `title`
- `slug`
- `source`
- `difficulty`：`easy` / `medium` / `hard`
- `description`
- `input_description`
- `output_description`
- `sample_input`
- `sample_output`
- `hint`
- `status`：`draft` / `published` / `offline`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一索引：`slug`
- 普通索引：`difficulty`
- 普通索引：`status`
- 普通索引：`title`

建议模型：

```go
type Problem struct {
	BaseModel
	Title             string `gorm:"size:255;not null;index" json:"title"`
	Slug              string `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Source            string `gorm:"size:128" json:"source"`
	Difficulty        string `gorm:"size:32;not null;index" json:"difficulty"`
	Description       string `gorm:"type:longtext;not null" json:"description"`
	InputDescription  string `gorm:"type:text" json:"input_description"`
	OutputDescription string `gorm:"type:text" json:"output_description"`
	SampleInput       string `gorm:"type:text" json:"sample_input"`
	SampleOutput      string `gorm:"type:text" json:"sample_output"`
	Hint              string `gorm:"type:text" json:"hint"`
	Status            string `gorm:"size:32;not null;default:draft;index" json:"status"`
}
```

### 4. `tags`

用途：统一标签体系，当前主要承载题型标签，如：

- 数组
- 字符串
- 哈希表
- 动态规划

可选扩展类型：

- `topic`
- `company`
- `language`

建议字段：

- `id`
- `name`
- `slug`
- `type`
- `color`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一索引：`slug`
- 普通索引：`type`

建议模型：

```go
type Tag struct {
	BaseModel
	Name  string `gorm:"size:64;not null" json:"name"`
	Slug  string `gorm:"size:64;not null;uniqueIndex" json:"slug"`
	Type  string `gorm:"size:32;not null;index" json:"type"`
	Color string `gorm:"size:32" json:"color"`
}
```

说明：

- `difficulty` 建议保留在 `problems` 表独立字段，不强制进入标签体系
- “按语言分标签”不建议作为第一版主设计。更合理的做法是把语言放到 `submissions.language`

### 5. `problem_tags`

用途：题目和标签的多对多关系。

建议字段：

- `id`
- `problem_id`
- `tag_id`

建议索引：

- 唯一联合索引：`(problem_id, tag_id)`
- 普通索引：`problem_id`
- 普通索引：`tag_id`

建议模型：

```go
type ProblemTag struct {
	BaseModel
	ProblemID uint `gorm:"not null;uniqueIndex:uk_problem_tag;index" json:"problem_id"`
	TagID     uint `gorm:"not null;uniqueIndex:uk_problem_tag;index" json:"tag_id"`
}
```

### 6. `problem_set_problems`

用途：题库和题目的多对多关系，同时支持题单内排序。

建议字段：

- `id`
- `problem_set_id`
- `problem_id`
- `sort_order`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一联合索引：`(problem_set_id, problem_id)`
- 普通索引：`problem_set_id`
- 普通索引：`problem_id`

建议模型：

```go
type ProblemSetProblem struct {
	BaseModel
	ProblemSetID uint `gorm:"not null;uniqueIndex:uk_set_problem;index" json:"problem_set_id"`
	ProblemID    uint `gorm:"not null;uniqueIndex:uk_set_problem;index" json:"problem_id"`
	SortOrder    int  `gorm:"not null;default:0" json:"sort_order"`
}
```

### 7. `submissions`

用途：保存每次提交记录，一次提交一条。

建议字段：

- `id`
- `user_id`
- `problem_id`
- `language`
- `code`
- `status`：`pending` / `accepted` / `wrong_answer` / `runtime_error` / `compile_error` / `time_limit_exceeded`
- `score`
- `runtime_ms`
- `memory_kb`
- `submit_at`
- `judged_at`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- `user_id`
- `problem_id`
- `status`
- `language`
- `submit_at`
- 组合索引可按后续查询热点补充，如 `(user_id, problem_id, submit_at)`

建议模型：

```go
type Submission struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	ProblemID uint   `gorm:"not null;index" json:"problem_id"`
	Language  string `gorm:"size:32;not null;index" json:"language"`
	Code      string `gorm:"type:longtext;not null" json:"code"`
	Status    string `gorm:"size:32;not null;default:pending;index" json:"status"`
	Score     int    `gorm:"not null;default:0" json:"score"`
	RuntimeMS int    `gorm:"not null;default:0" json:"runtime_ms"`
	MemoryKB  int    `gorm:"not null;default:0" json:"memory_kb"`
	SubmitAt  int64  `gorm:"not null;index" json:"submit_at"`
	JudgedAt  int64  `gorm:"not null;default:0" json:"judged_at"`
}
```

### 8. `user_problem_stats`

用途：保存用户在题目维度上的聚合状态，提升题目列表和个人中心查询效率。

建议字段：

- `id`
- `user_id`
- `problem_id`
- `status`：`unsolved` / `attempted` / `solved`
- `submit_count`
- `accepted_count`
- `first_accepted_at`
- `last_submit_at`
- `last_submission_id`
- `created_at`
- `updated_at`
- `deleted_at`

建议索引：

- 唯一联合索引：`(user_id, problem_id)`
- 普通索引：`status`

建议模型：

```go
type UserProblemStat struct {
	BaseModel
	UserID           uint `gorm:"not null;uniqueIndex:uk_user_problem;index" json:"user_id"`
	ProblemID        uint `gorm:"not null;uniqueIndex:uk_user_problem;index" json:"problem_id"`
	Status           string `gorm:"size:32;not null;default:attempted;index" json:"status"`
	SubmitCount      int    `gorm:"not null;default:0" json:"submit_count"`
	AcceptedCount    int    `gorm:"not null;default:0" json:"accepted_count"`
	FirstAcceptedAt  int64  `gorm:"not null;default:0" json:"first_accepted_at"`
	LastSubmitAt     int64  `gorm:"not null;default:0" json:"last_submit_at"`
	LastSubmissionID uint   `gorm:"not null;default:0" json:"last_submission_id"`
}
```

## 核心关系

- 一个用户可以有多次提交：`users 1 -> n submissions`
- 一个题目可以有多次提交：`problems 1 -> n submissions`
- 一个用户在一个题目上有一条聚合状态：`users 1 -> n user_problem_stats`
- 一个题目对应多条用户聚合状态：`problems 1 -> n user_problem_stats`
- 一个题库包含多个题目，一个题目可属于多个题库：`problem_sets n <-> n problems`
- 一个题目可有多个标签，一个标签可关联多个题目：`problems n <-> n tags`

## 查询场景覆盖

该设计可直接支持以下 MVP 查询：

- 用户邮箱注册、登录
- 题库列表、题库详情、题库内题目排序
- 题目列表按难度筛选
- 题目列表按标签筛选
- 我的提交记录
- 某题的提交历史
- 用户是否做过某题、是否通过某题

## 关于“按语言分标签”的处理建议

MVP 不建议把编程语言作为普通题目标签主建模，原因如下：

- 语言更天然属于提交记录维度，而不是算法题内容维度
- 把语言混到标签会导致语义不清，例如“这道题是 Go 标签”到底表示推荐解法还是限制语言

更合理的做法：

- 当前阶段：只在 `submissions.language` 中记录提交语言
- 后续如果需要限制题目支持语言，再增加 `problem_supported_languages` 表

## 暂不进入 MVP 的表

以下表建议后续再加：

- `test_cases`
- `judge_tasks`
- `submission_case_results`
- `problem_supported_languages`
- `refresh_tokens` 或 `user_sessions`

原因：这些都不是当前最小可用产品的必要前提。

## 实施建议

第一阶段只落以下模型即可：

- `User`
- `ProblemSet`
- `Problem`
- `Tag`
- `ProblemTag`
- `ProblemSetProblem`
- `Submission`
- `UserProblemStat`

接口优先级建议如下：

1. 注册 / 登录
2. 题库列表 / 题库详情
3. 题目列表 / 题目详情
4. 提交代码 / 查询提交记录
5. 我的做题状态

## 自检结论

已检查以下问题：

- 无 `TODO` / `TBD` 占位
- 表关系与功能需求一致
- 范围仍然聚焦于单个 MVP
- 对“语言标签”和“聚合状态表”的边界已明确
