# bolbox

bolbox 是一个 Golang 项目的基础工具库，提供了一系列常用的功能模块，帮助开发者快速构建 Golang 应用。

## 功能模块

### 1. 配置管理 (pkg/configs)
- 支持从环境变量、命令行参数和动态映射中加载配置
- 优先级：默认值 < 环境变量 < 命令行参数 < 动态映射
- 支持配置变更回调
- 类型安全的配置管理

### 2. 错误处理 (pkg/errors)
- 基于 cockroachdb/errors 包装
- 提供丰富的错误处理功能，包括堆栈跟踪、错误包装、错误链等
- 兼容标准库 errors 包的接口

### 3. 日志管理 (pkg/log)
- 支持默认日志和 zap 日志实现
- 提供统一的日志接口
- 支持不同级别的日志输出
- 支持结构化日志

### 4. HTTP 响应包装 (pkg/mix)
- 提供链式调用的 HTTP 响应包装器
- 支持 JSON 和文本响应
- 支持常见 HTTP 错误状态码的快速响应

### 5. 服务管理 (pkg/services)
- 支持模块的生命周期管理
- 支持模块依赖排序和循环依赖检测
- 支持优雅启动和关闭

### 6. 信号处理 (pkg/signals)
- 提供优雅关闭上下文
- 处理系统信号（SIGTERM、SIGINT）

### 7. 类型定义 (pkg/types)
- 提供常用的类型定义和工具函数

## 安装

```bash
go get github.com/wolfbolin/bolbox
```

## 使用示例

### 配置管理

```go
import (
    "github.com/wolfbolin/bolbox/pkg/configs"
)

// 定义配置结构体
type AppConfig struct {
    ServerPort int    `env:"SERVER_PORT" flag:"server-port"`
    Debug      bool   `env:"DEBUG" flag:"debug"`
    DatabaseURL string `env:"DATABASE_URL" flag:"database-url"`
}

// 创建配置管理器
conf, err := configs.NewManager(&AppConfig{
    ServerPort: 8080, // 默认值
})
if err != nil {
    // 处理错误
}

// 获取配置值
config := conf.Vars()
fmt.Println("Server port:", config.ServerPort)

// 监听配置变更
serverPortConf, err := conf.Conf("ServerPort")
if err != nil {
    // 处理错误
}
serverPortConf.OnChange(func(val any) {
    fmt.Println("Server port changed:", val)
})
```

### 日志管理

```go
import (
    "github.com/wolfbolin/bolbox/pkg/log"
    "github.com/wolfbolin/bolbox/pkg/log/zap"
    "go.uber.org/zap"
)

// 使用默认日志
log.Infof("Hello, %s", "world")
log.Errorf("Error: %v", err)

// 使用 zap 日志
zapLogger, _ := zap.NewProduction()
log.SetLogger(zap.NewLogger(zapLogger))
log.Infof("Hello with zap", "key", "value")
```

### 服务管理

```go
import (
    "context"
    "github.com/wolfbolin/bolbox/pkg/services"
)

// 实现 Module 接口
type MyModule struct {
    name string
    status *services.ModuleStatus
}

func (m *MyModule) Name() string {
    return m.name
}

func (m *MyModule) Status() *services.ModuleStatus {
    return m.status
}

func (m *MyModule) Run(ctx context.Context) {
    m.status.Set(services.StatusRunning)
    // 模块逻辑
    <-ctx.Done()
    m.status.Set(services.StatusStopped)
}

func (m *MyModule) Requires() []string {
    return []string{"dependency-module"}
}

// 创建服务管理器
manager := services.NewManager()

// 添加模块
manager.AddModule("my-module", &MyModule{
    name: "my-module",
    status: services.NewModuleStatus(),
})

// 启动服务
ctx, cancel := context.WithCancel(context.Background())
go manager.StartAndServe(ctx)

// 优雅关闭
<-manager.Done(cancel)
```

### 信号处理

```go
import (
    "github.com/wolfbolin/bolbox/pkg/signals"
)

// 创建优雅关闭上下文
ctx, closeChan := signals.GracefulShutdownContext()

// 使用 ctx 控制服务生命周期
go func() {
    <-ctx.Done()
    // 处理关闭逻辑
}()

// 等待强制关闭信号
<-closeChan
```

### HTTP 响应包装

```go
import (
    "net/http"
    "github.com/wolfbolin/bolbox/pkg/mix"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // JSON 响应
    mix.HttpRsp(w).Code(http.StatusOK).Json(map[string]string{
        "message": "Hello, world!",
    })

    // 文本响应
    mix.HttpRsp(w).Code(http.StatusOK).Text("Hello, %s!", "world")

    // 错误响应
    mix.HttpRsp(w).BadRequest(errors.New("Bad request"))
    mix.HttpRsp(w).ServerError(errors.New("Internal server error"))
}
```

## 依赖

- [github.com/agiledragon/gomonkey/v2 v2.14.0](https://github.com/agiledragon/gomonkey) - 用于测试
- [github.com/cockroachdb/errors v1.11.1](https://github.com/cockroachdb/errors) - 错误处理
- [github.com/stretchr/testify v1.9.0](https://github.com/stretchr/testify) - 测试断言
- [go.uber.org/zap v1.27.0](https://github.com/uber-go/zap) - 日志库
- [gopkg.in/natefinch/lumberjack.v2 v2.2.1](https://github.com/natefinch/lumberjack) - 日志轮转

## 许可证

[LICENSE](LICENSE)
