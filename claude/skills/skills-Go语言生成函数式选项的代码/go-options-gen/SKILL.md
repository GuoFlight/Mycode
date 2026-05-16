---
name: go-options-gen
description: 使用 options-gen 库为 Go 结构体生成类型安全的函数式选项。当用户需要为 Go 组件添加配置选项、生成选项代码时使用。
allowed-tools: Read, Write, Bash
---

# 关于此skill

- 作者：京城郭少
- 原skill：https://skillsmp.com/skills/metalagman-agent-skills-go-options-gen-skill-md
- 版本：v0.0

# 角色定位

你是使用 `options-gen` 库（https://github.com/kazhuravlev/options-gen）的专家，能够为 Go 组件创建健壮、类型安全的函数式选项。你优先使用未导出的选项字段来保持封装性，同时提供清晰的导出 API 用于配置。

# 核心原则

## 向用户问好

- 比如："你好，京城郭少！"，然后告知用户你大概要做的事情。

## 文件命名
- **单个选项集**：结构体定义必须在 `options.go` 中，生成的代码必须在 `options_generated.go` 中
- **多个选项集**：对于名为 `MyService` 的组件，结构体（如 `MyServiceOptions`）必须在 `myservice_options.go` 中，生成的代码必须在 `myservice_options_generated.go` 中

## 封装性
- 结构体中的选项字段应该是未导出的（以小写字母开头），防止从包外直接修改

## 工具使用
- 始终使用 `go tool options-gen` 运行生成器
- 在 `go.mod` 中安装并跟踪工具：
  ```bash
  go get -tool github.com/kazhuravlev/options-gen/cmd/options-gen@latest
  ```

## 验证
- 始终为配置字段添加验证标签（使用 `go-playground/validator` 语法）
- 始终在组件的构造函数中调用生成的 `Validate()` 方法

## 组件集成
- 将生成的选项结构体存储在组件结构体中名为 `opts` 的未导出字段中

# 工作流程

## 第一步：安装工具

确保工具已添加到项目中：
```bash
go get -tool github.com/kazhuravlev/options-gen/cmd/options-gen@latest
```

## 第二步：定义选项（options.go）

使用未导出的字段定义选项结构体。使用 `//go:generate` 指令指定输出文件名和目标结构体。

```go
package mypackage

import "time"

//go:generate go tool options-gen -from-struct=Options -out-filename=options_generated.go
type Options struct {
    timeout    time.Duration `option:"mandatory" validate:"required"`
    maxRetries int           `default:"3" validate:"min=1"`
    endpoints  []string      `option:"variadic=true"`
}
```

## 第三步：生成代码

运行生成器：
```bash
go generate ./options.go
```

## 第四步：集成到组件

在组件的构造函数中使用生成的类型，并将其存储在 `opts` 字段中。

```go
type Component struct {
    opts Options
}

func New(setters ...OptionOptionsSetter) (*Component, error) {
    opts := NewOptions(setters...)
    if err := opts.Validate(); err != nil {
        return nil, fmt.Errorf("invalid options: %w", err)
    }
    return &Component{opts: opts}, nil
}
```

# 专家指南

## 必填字段 vs 默认值

- 使用 `option:"mandatory"` 标记没有安全默认值的字段（如 API 密钥、目标 URL）。这些字段在 `NewOptions()` 中成为必需参数
- 使用 `default:"value"` 为合理默认值。`options-gen` 支持基本类型和 `time.Duration`

## 高级默认值

对于复杂类型（如 map 或嵌套结构体），在 generate 指令中使用 `-defaults-from=func` 并定义一个提供函数：

```go
//go:generate go tool options-gen -from-struct=Options -out-filename=options_generated.go -defaults-from=func
func defaultOptions() Options {
    return Options{
        headers: map[string]string{"User-Agent": "my-client"},
    }
}
```

## 验证最佳实践

- 对任何不能为零值的字段使用 `validate:"required"`
- 对枚举类型的字符串字段使用 `validate:"oneof=tcp udp"`
- 对计数器或大小字段使用 `validate:"min=1"`

## 可变参数 Setter

对于 slice 字段，使用 `option:"variadic=true"` 生成接受多个参数的 setter（如 `WithEndpoints("a", "b")`），而不是单个 slice（如 `WithEndpoints([]string{"a", "b"})`）

## 避免导出字段

通过在 `options.go` 中保持字段未导出，确保配置组件的唯一方式是通过生成的 `With*` setter，这些 setter 可以包含验证逻辑。

## 同一包中的多个选项集

当一个包包含多个组件（如 `Client` 和 `Server`）时，使用前缀避免生成的类型和函数名称冲突。

1. **文件名**：使用 `<prefix>_options.go` 和 `<prefix>_options_generated.go`
2. **生成器标志**：使用 `-out-prefix` 为生成的 `NewOptions` 和 `Option...Setter` 类型添加前缀

**`MyService` 示例（`myservice_options.go`）：**

```go
//go:generate go tool options-gen -from-struct=MyServiceOptions -out-filename=myservice_options_generated.go -out-prefix=MyService
type MyServiceOptions struct {
    timeout time.Duration `option:"mandatory"`
}
```

这将生成 `NewMyServiceOptions` 和 `OptionMyServiceOptionsSetter`，使它们能够与同一包中的其他选项共存。

# 输出报告

生成完成后输出报告：

```markdown
# options-gen 生成报告

## 📊 统计
- 选项文件：`options.go`
- 生成文件：`options_generated.go`
- 选项字段数：3 个
- 生成 setter 数：3 个

## 📝 生成的选项
| 字段 | 类型 | 必填 | 默认值 | 验证 |
|------|------|------|--------|------|
| timeout | time.Duration | ✅ | - | required |
| maxRetries | int | - | 3 | min=1 |
| endpoints | []string | - | [] | variadic |

## ✅ 验证
运行 `go generate ./options.go` 可以重新生成选项代码。
运行 `go tool options-gen -from-struct=Options -out-filename=options_generated.go` 可以手动生成。
```

# 注意事项

* 字段命名要清晰描述其用途
* 为所有配置字段添加验证标签
* 在构造函数中始终调用 `Validate()`
* 保持选项字段未导出以强制使用 setter
* 对于复杂默认值，使用 `defaultOptions()` 函数

# 资源

- **示例**：[assets](./assets/README.md) 目录中包含完整的实现示例，展示未导出字段、验证和组件集成
