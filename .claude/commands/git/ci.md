---
name: git 提交代码
description: 分析变更代码内容，git commit 提交代码
---

# Git 操作技能

## /git:ci

总结当前暂存区（staged）的代码变更，生成符合项目规范的提交信息，并执行 `git commit`。

### 执行步骤

1. 运行 `git status --porcelain` 查看暂存区文件
2. 运行 `git diff --cached --stat` 查看暂存的具体变更
3. 分析变更内容，生成简洁的提交信息（遵循 Conventional Commits 规范）
4. 执行 `git commit -m "<提交信息>"`
5. 运行 `git status` 确认提交成功

### 提交信息规范

- 格式：`type(scope): description`
- type 可选：feat, fix, docs, style, refactor, test, chore, perf, ci
- description 简短描述，不超过 50 字，采用列表形式


