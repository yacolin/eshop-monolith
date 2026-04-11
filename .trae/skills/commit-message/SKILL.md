---
name: "commit-message"
description: "生成符合项目规范的commit message。当用户需要提交代码时调用，帮助用户编写规范的commit message。"
---

# Commit Message 生成器

## 技能信息

- **技能名称**: Commit Message 生成器
- **技能ID**: commit-message
- **描述**: 生成符合项目规范的commit message
- **版本**: 1.0.0
- **作者**: System

## 适用场景

- 代码提交前生成规范的commit message
- 帮助团队成员统一commit message格式
- 确保commit message清晰明了，便于代码审查和版本管理

## Commit Message 规范

### 格式

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### 类型 (Type)

| 类型     | 描述               |
|---------|-------------------|
| feat    | 新功能             |
| fix     | 修复bug           |
| refactor| 代码重构           |
| docs    | 文档更新           |
| style   | 代码风格调整        |
| test    | 测试相关           |
| chore   | 构建或依赖更新      |
| Merge   | 合并分支           |

### 范围 (Scope)

范围应该是一个简短的标识符，描述了更改的范围，例如：

- api: API相关
- routes: 路由相关
- rbac: 权限管理相关
- middleware: 中间件相关
- service: 服务层相关
- model: 模型相关
- config: 配置相关
- db: 数据库相关
- mq: 消息队列相关

### 描述 (Description)

- 简洁明了，不超过50个字符
- 使用中文描述
- 以动词开头，如"添加"、"修复"、"重构"等
- 描述具体的更改内容

### 正文 (Body) [可选]

- 详细描述更改的原因和影响
- 可以多行
- 每行不超过72个字符

### 页脚 (Footer) [可选]

- 关联的issue或任务
- 破坏性更改的警告
- 其他重要信息

## 示例

### 新功能

```
feat(api): 添加用户登录API接口

添加了基于JWT的用户登录接口，支持密码登录和手机验证码登录

关联任务: #123
```

### 修复bug

```
fix(routes): 解决用户路由冲突问题

修复了/users/:id和/users/:user_id/roles之间的路由冲突

关联issue: #456
```

### 代码重构

```
refactor(service): 重构订单服务逻辑

将订单服务的业务逻辑拆分为多个子函数，提高代码可读性和可维护性
```

### 文档更新

```
docs: 更新项目文档结构

更新了项目的README文件，添加了RBAC相关API和功能说明
```

## 最佳实践

1. **保持一致性**: 团队成员应使用相同的commit message格式
2. **清晰明了**: commit message应该能够快速传达更改的内容和原因
3. **具体详细**: 描述应该具体，避免模糊的表述
4. **关联issue**: 对于有issue的任务，应在commit message中关联
5. **避免大提交**: 每个commit应该只包含一个逻辑更改
6. **使用中文**: 项目使用中文进行commit message描述

## 常见问题

### 1. 提交信息过于简单

**问题**: commit message描述过于简单，无法理解更改的具体内容

**解决方案**: 提供更详细的描述，包括更改的原因和影响

### 2. 提交信息格式不正确

**问题**: commit message格式不符合规范

**解决方案**: 按照`<type>(<scope>): <description>`的格式编写

### 3. 一次提交包含多个不相关的更改

**问题**: 一个commit包含多个不相关的更改，难以 review 和回滚

**解决方案**: 将不同的更改分成多个commit，每个commit只包含一个逻辑更改

### 4. 提交信息使用英文

**问题**: commit message使用英文，与项目规范不符

**解决方案**: 使用中文编写commit message

## 工具集成

### 1. Git Hook

可以使用git hook来检查commit message是否符合规范：

```bash
# .git/hooks/commit-msg
#!/bin/sh

commit_msg="$1"

# 检查commit message格式
if ! grep -qE '^[a-z]+\([^)]+\): .+' "$commit_msg"; then
  echo "Error: commit message format is incorrect"
  echo "Please use format: <type>(<scope>): <description>"
  exit 1
fi
```

### 2. Commitizen

可以使用Commitizen工具来帮助生成规范的commit message：

```bash
# 安装
npm install -g commitizen

# 初始化
commitizen init cz-conventional-changelog --save-dev --save-exact

# 使用
git cz
```

## 总结

本技能旨在帮助开发团队生成符合项目规范的commit message，提高代码管理的质量和效率。通过遵循这些规范和最佳实践，团队可以更清晰地了解代码更改的历史，便于代码审查和问题定位。

## 使用指南

1. **确定更改类型**: 根据更改的性质选择合适的类型（feat、fix、refactor等）
2. **确定更改范围**: 选择描述更改范围的简短标识符
3. **编写描述**: 简洁明了地描述更改的内容
4. **添加正文**: 详细描述更改的原因和影响（可选）
5. **添加页脚**: 关联issue或任务（可选）
6. **检查格式**: 确保commit message符合规范

通过使用本技能，您可以生成符合项目规范的commit message，提高代码管理的质量和效率。