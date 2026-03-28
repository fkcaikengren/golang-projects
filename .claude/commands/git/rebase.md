---
name: git rebase
description: 开发分支 git rebase 主(master)分支代码
---


## /git:rebase

拉取最新的 master 分支代码，并用当前分支 rebase master 分支。

### 执行步骤

1. 运行 `git fetch origin` 拉取最新远程分支
2. 运行 `git rebase origin/master` 执行 rebase
3. 如果有冲突，提示用户解决冲突后运行 `git rebase --continue`
4. 运行 `git log --oneline -3` 显示最近提交，确认 rebase 成功

### 注意事项

- 如果当前分支与 master 差异较大，可能产生较多冲突
- rebase 会重写提交历史，适合保持线性提交历史