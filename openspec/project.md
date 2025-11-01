# Project Context

## Purpose
pics_show_system 基于 Gin-Vue-Admin 定制，主要目标是为美甲作品与相关素材提供一体化的管理后台。系统需要支持作品图片、标签、款式信息的录入与展示，同时向前端画廊提供经过审核与分类的素材数据。

## Tech Stack
- 后端：Go 1.22、Gin、Gorm、Gin-Vue-Admin 内置中间件
- 前端：Vue 3、Vite、Element Plus、Pinia、Axios
- 数据库/缓存：MySQL（主库）、Redis（可选）、OSS/S3 兼容存储用于资源托管

## Project Conventions

### Code Style
- Go 代码遵循 `go fmt` 与 `goimports`，结构体/方法使用 PascalCase，数据表字段采用 snake_case；路由路径为动词-资源组合（例：`/api/v1/tags`）。
- 前端统一使用 2 空格缩进与单引号，遵守 `eslint.config.mjs` 规则；组件文件以 PascalCase 命名，路由与目录使用 kebab-case；共享样式存放在 `web/src/style`。

### Architecture Patterns
- 后端延续 GVA 分层：`api` 处理请求、`service` 承载业务、`router` 注册路由、`model` 定义实体；初始化逻辑集中在 `initialize/`，配置通过 `config/config.go` 解析。
- 前端以模块化视图组织，`src/api` 封装后端接口，`src/pinia` 管理状态，`src/utils` 存放跨模块工具函数；路由守卫位于 `src/permission.ts` 负责鉴权与动态路由装载。

### Testing Strategy
- 后端统一使用 Go testing：`go test ./...`，新增 service/handler 需提供核心路径用例，推荐覆盖率 ≥70%；数据库相关逻辑可通过初始化数据或 mock 实现。
- 前端暂未启用单测框架，至少保证 `npm run lint` 通过；若添加 Vitest 等测试，文件置于 `web/src/__tests__/` 并与目标页面同名。

### Git Workflow
- 使用 GitHub Flow：功能分支自 `main` 派生，完成后通过 PR 合并；保持小步提交。
- 提交信息使用简洁中文动宾短语（示例：`添加美甲标签和款式的相关功能`），扩展说明放在正文；PR 描述需列出变更摘要、关联需求、测试截图与风险评估。

## Domain Context
- 核心实体包括美甲作品（含多张图片）、标签、款式、展示页配置等，需要支持多维筛选与批量上下架。
- 管理员通过后台维护素材，前端画廊向访客展示精选作品并支持按标签检索。

## Important Constraints
- 必须兼容 Gin-Vue-Admin 的权限与菜单体系，新增模块需同步自动化生成或登记 API、菜单、按钮权限。
- 所有用户上传内容需通过后台审核后才可公开展示，敏感信息与凭据不得提交至仓库，配置通过环境变量覆盖。
- 构建流水线依赖 Docker 镜像（Go、Node），CI 必须通过 `make build` 与 `go test` 才可合并。

## External Dependencies
- 对象存储：阿里云 OSS、AWS S3 或 Cloudflare R2（通过 `server/config.yaml` 配置）
- 第三方库：Element Plus UI 组件、WangEditor 富文本、ECharts 可视化、UnoCSS 样式工具
- 基础设施：MySQL 数据库实例、Redis 缓存服务（可选）
