# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是 **gin-vue-admin** (GVA)，一个基于 Go + Gin（后端）和 Vue 3 + Vite（前端）的全栈管理系统框架。该项目采用严格的分层架构，前后端目录分离。

**核心技术栈：**
- 后端：Go 1.22+、Gin 1.9.1、GORM 1.25.2、Casbin (RBAC)、Viper (配置)、Zap (日志)
- 前端：Vue 3.5.7、Vite 6.2.3、Pinia 2.2.2、Element Plus 2.10.2、UnoCSS
- 数据库：MySQL（主要）、支持 PostgreSQL、SQLite、SQL Server、MongoDB
- 缓存：Redis
- 认证：JWT，支持多点登录控制

## 开发命令

### 后端 (server/)

**重要：** 始终在 IDE 中直接打开 `server/` 目录，而不是根目录。

```bash
# 安装依赖并生成代码
cd server
go generate

# 运行服务器（默认端口 8888）
go run .

# 生成 Swagger 文档
swag init
# 访问地址：http://localhost:8888/swagger/index.html

# 生产环境构建
go build -o server
```

### 前端 (web/)

```bash
cd web

# 安装依赖
npm install

# 开发服务器（默认端口可在 vite.config.js 中查看）
npm run dev
# 或
npm run serve

# 生产环境构建
npm run build

# 预览生产构建
npm run preview
```

### 使用 Makefile（从根目录）

```bash
# 本地构建前后端
make build-local

# 仅构建前端
make build-web-local

# 仅构建后端
make build-server-local

# 生成 Swagger 文档
make doc

# 打包插件（默认：email）
make plugin PLUGIN=your-plugin-name
```

## 架构概览

### 后端分层架构

后端遵循**严格的四层架构**。依赖关系只能单向流动：

```
Router → API → Service → Model
```

**关键规则：**
1. **绝不**在 API 层直接调用数据库操作 - 必须使用 Service 层
2. **绝不**在 Service 层处理 `gin.Context` - Service 层与 HTTP 无关
3. **绝不**跳过层级或创建循环依赖
4. 使用全局组变量（ServiceGroupApp、ApiGroupApp、RouterGroupApp）避免循环引用

### 各层职责

#### 1. Model 层 (`server/model/`)
- **数据模型**：定义映射到数据库表的 GORM 结构体
  - 继承 `global.GVA_MODEL` 获取标准字段（ID、CreatedAt、UpdatedAt、DeletedAt）
  - 使用清晰的 `json` 和 `gorm` 标签
- **请求模型** (`model/request/`)：用于接收前端数据的 DTO
  - 包含 `json` 和 `form` 标签以便 Gin 绑定
  - 搜索请求应嵌入 `request.PageInfo` 实现分页
- **响应模型** (`model/response/`)：用于返回数据给前端的 DTO

**数据类型一致性：** 确保同一字段在所有模型文件（数据模型、请求、响应）中使用相同的类型。数据模型中的指针类型需要在 Service 层进行仔细的转换。

#### 2. Service 层 (`server/service/`)
- 包含所有业务逻辑和数据库 CRUD 操作
- 函数接收业务参数（模型/请求），返回结果 + error
- **禁止**导入或使用 `gin.Context`
- 每个模块有 `xxx_service.go` + 在 `service/enter.go` 中注册

Service 组结构示例 (`service/enter.go`)：
```go
type ServiceGroup struct {
    SystemServiceGroup
    ExampleServiceGroup
}
var ServiceGroupApp = new(ServiceGroup)
```

#### 3. API 层 (`server/api/`)
- HTTP 请求入口：参数校验、调用服务、返回 JSON 响应
- 每个处理器**必须**包含完整的 Swagger 注解
- 使用 `service.ServiceGroupApp` 调用 Service 层方法
- 使用统一的 `response` 包（OkWithDetailed、FailWithMessage 等）
- 每个模块有 `xxx_api.go` + 在 `api/enter.go` 中注册

#### 4. Router 层 (`server/router/`)
- 将 HTTP 路径映射到 API 处理器
- 配置中间件（认证、日志等）
- 使用 `api.ApiGroupApp` 引用 API 处理器
- 每个模块有 `xxx_router.go` + 在 `router/enter.go` 中注册

### 插件架构

插件是 `server/plugin/[name]/` 和 `web/src/plugin/[name]/` 中的自包含模块。参考实现：`server/plugin/announcement`

**后端插件结构：**
```
server/plugin/[name]/
├── api/           # API 处理器，包含 enter.go
├── config/        # 插件特定配置
├── initialize/    # gorm.go、router.go、menu.go、viper.go、api.go
├── model/         # 数据模型和请求模型
├── router/        # 路由，包含 enter.go
├── service/       # 业务逻辑，包含 enter.go
└── plugin.go      # 插件入口，实现 system.Plugin 接口
```

**前端插件结构：**
```
web/src/plugin/[name]/
├── api/           # API 调用
├── components/    # 可复用组件（可选）
├── view/          # 页面
└── form/          # 表单（可选）
```

### 配置

- **后端配置**：`server/config.yaml`
  - 数据库连接（mysql、pgsql、sqlite 等）
  - Redis 设置
  - JWT 签名密钥和过期时间
  - OSS/云存储设置
  - 系统设置（端口、多点登录等）
- **前端配置**：`web/.env.development`、`web/.env.production`
- **自动代码生成**：在 config.yaml 的 `autocode` 部分设置路径

## 前端架构

### 目录结构

- `src/api/`：API 服务层 - 所有后端调用都通过这里
- `src/components/`：全局可复用组件
- `src/view/`：按功能组织的页面组件
- `src/pinia/modules/`：状态管理 stores（替代 Vuex）
- `src/router/`：Vue Router 配置，支持动态路由
- `src/utils/`：工具函数（request.js 是 HTTP 客户端等）
- `src/permission.js`：路由中间件，用于认证/权限
- `src/style/`：全局样式（可以在这里覆盖 Element Plus）

### 关键模式

**API 调用：** 始终使用 `src/utils/request.js` 封装 axios
```javascript
import service from '@/utils/request'

export const getList = (data) => {
  return service({
    url: '/api/list',
    method: 'post',
    data
  })
}
```

**状态管理：** 使用 Pinia 配合 Composition API
```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useMyStore = defineStore('my', () => {
  const data = ref(null)
  const isLoaded = computed(() => !!data.value)

  const loadData = async () => {
    // 异步 action
  }

  return { data, isLoaded, loadData }
})
```

**组件风格：** 使用 Composition API 配合 `<script setup>`

## 开发工作流

### 创建新功能/插件

**重要：** 开始之前，检查是否有 "GVA Helper" MCP 工具。如果有，请先咨询它以获取针对此框架的最佳实践和代码模板。

1. **设计模型**（后端 `model/` 和 `model/request/`）
2. **实现 Service 层**（业务逻辑、数据库操作）
3. **实现 API 层**（处理器及 Swagger 文档）
4. **实现 Router 层**（路由注册）
5. **创建初始化**（`initialize/gorm.go`、`initialize/router.go` 等）
6. **前端 API**（`web/src/api/`）
7. **前端页面**（`web/src/view/`）
8. **测试集成**

### Swagger 文档

每个 API 处理器**必须**包含 Swagger 注解：

```go
// CreateXxx 创建XXX
// @Tags     XxxModule
// @Summary  创建一个新的XXX
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateXxxRequest true "XXX的名称和描述"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /xxx/createXxx [post]
func (a *XxxApi) CreateXxx(c *gin.Context) {
    // 实现代码
}
```

更新 API 注解后，在 `server/` 目录运行 `swag init`。

## 数据库与初始化

- **自动迁移**：在 `initialize/register_init.go` 中注册的模型会在启动时自动迁移（除非在配置中禁用）
- **初始数据**：种子数据函数在 `server/source/`
- **多数据库**：在 config.yaml 的 `db-list` 部分配置

## 重要说明

### 数据类型一致性
- **关键：** 同一字段必须在后端模型、前端 TypeScript/JavaScript 和 API 文档中使用相同的数据类型
- 特别注意：状态字段、ID、枚举、时间戳
- 后端指针类型需要在 Service 层进行显式的 nil 检查和转换

### 认证与权限
- JWT token 存储在 Redis 中（如果 `use-redis: true`）
- 多点登录限制（如果配置中 `use-multipoint: true`）
- Casbin 处理 RBAC - 按 API 路径定义权限
- 前端使用 `v-auth` 指令实现按钮级权限

### 代码生成
- 内置代码生成器创建 CRUD 脚手架
- 模板在 `server/resource/template/`
- 在 config.yaml 的 `autocode` 部分配置输出路径

### 测试
运行服务器并访问内置测试界面，或使用 Swagger UI 测试 API。

## 常用模式

### 错误处理（API 层）
```go
err := xxxService.DoSomething(params)
if err != nil {
    response.FailWithMessage(err.Error(), c)
    return
}
response.OkWithMessage("操作成功", c)
```

### 分页（Service 层）
```go
var total int64
db := global.GVA_DB.Model(&model.Xxx{})
db.Count(&total)
err := db.Limit(limit).Offset(offset).Find(&list).Error
```

### 全局变量访问
- 数据库：`global.GVA_DB`
- 日志：`global.GVA_LOG`
- 配置：`global.GVA_CONFIG`
- Redis：`global.GVA_REDIS`

## 参考文件

最佳实践和实现示例：
- **后端插件**：`server/plugin/announcement/`
- **前端插件**：`web/src/plugin/announcement/`
- **系统 CRUD**：`server/api/v1/system/`、`server/service/system/`
- **示例 CRUD**：`server/api/v1/example/`、`server/service/example/`
