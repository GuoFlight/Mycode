# skill-maker

用于创建 Claude Code Skill 的 Skill。

## 使用方法

```bash
/skill-maker
```

或直接描述你的需求：

> 帮我创建一个用于生成 Go 单元测试的 skill

## 功能

- 分析 skill 需求和触发条件
- 创建规范的 SKILL.md 文件结构
- 提供最佳实践指导
- 生成示例资源

## 输出

创建完成后会生成：

```
.claude/skills/<skill-name>/
├── SKILL.md          # Skill 定义文件
└── assets/           # 示例资源目录
```

## 示例

### 创建代码生成类 Skill

```
/skill-maker 帮我创建一个 skill，用于为 Go HTTP handler 生成测试用例
```

### 创建文档类 Skill

```
/skill-maker 创建一个 skill，用于生成 API 文档
```

