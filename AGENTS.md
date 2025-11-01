<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Repository Guidelines

## 项目结构与模块划分
- `server/` 为 Go 后端，按 `api`→`service`→`router`→`model` 分层；新包保持与现有命名一致，并在 `model/request` 与 `model/response` 同步定义 DTO。测试文件与源码同目录，命名为 `*_test.go`，同时检查 `initialize/` 是否需要注册新资源。
- `web/` 存放 Vue 3 + Vite 前端，核心逻辑位于 `src`（`api/` 请求封装、`view/` 业务页面、`pinia/` 状态）；静态资源归档到 `src/assets/`，复用组件置于 `src/components/`。若新增路由，请同步配置 `src/router/index.ts` 与 `src/permission.ts` 的守卫。
- `docs/`、`SECURITY.md` 记录设计与安全策略；`deploy/` 提供 Docker 与打包脚本；自定义脚本请放入 `deploy/scripts` 以便复用。生成文件应放入 `build/`（由 Makefile 创建），勿手工提交至版本库。
- 功能描述以及必要性描述文件在`.claude/rules/project_rules.md`中有详细说明，请参考。

## 构建、测试与开发命令
```bash
make build              # 容器内编译前后端产物
make build-server-local # 本地编译 Go 后端，输出 server/server
make build-web-local    # 本地构建前端 dist 目录
cd server && go run main.go          # 启动后端开发实例
cd web && npm install && npm run dev # 前端热加载开发（默认 5173 端口）
docker compose -f deploy/docker/docker-compose.yml up -d # 一键启动生产近似环境
```

## 代码风格与命名约定
- Go 代码提交前执行 `go fmt ./...` 与 `goimports`，结构体/方法使用 PascalCase，数据库字段沿用 snake_case；REST 路径保持动词-资源组合（如 `/api/v1/tags`）。新增配置项时补齐 `config.yaml`、`config/config.go` 与 `core/server.go` 的解析。
- 前端遵循 `eslint.config.mjs` 规则，统一使用 2 空格缩进与单引号；组件文件采用 PascalCase，路由与文件夹使用 kebab-case；样式首选 SCSS 模块，公共变量置于 `src/style`。跨模块通用逻辑应抽离到 `src/utils` 并配齐类型声明。

## 测试指引
- 后端使用 Go 标准库测试：`cd server && go test ./...`，新增业务须覆盖核心 service 与 handler，目标函数级覆盖率 ≥70%。涉及数据库的测试可借助 `source/` 中的初始化数据或 mock repository。
- 前端当前无内置单测框架，新增特性需至少提供一次 `npm run lint` 结果；若引入 Vitest，请将用例放在 `src/__tests__/` 并与页面同名，并在 PR 中附带 `npm run test -- --runInBand` 的输出。

## 提交与 PR 规范
- Git 历史采用简洁中文动宾短语（示例：`添加美甲标签和款式的相关功能`），请保持一条提交聚焦单一改动，并在正文列出主要变更点。补丁中若涉及自动生成代码，请注明使用的 GVA 工具及参数。
- PR 描述需包括：变更摘要、关联 Issue/需求链接、测试或截图证据及潜在风险。涉及数据库/配置改动，请附迁移或回滚说明，并在 `docs/CHANGELOG.md`（若新增）同步记录要点。

## 配置与安全提示
- 默认配置位于 `server/config.yaml`，容器环境使用 `config.docker.yaml`；请通过环境变量覆盖敏感信息，避免提交真实密钥或证书。若需要新增配置段落，务必更新示例文件并在 README 中提示。
- 上传资源走 `server/resource` 目录；若接入外部存储，需在 PR 中列明所需凭据与最小权限策略。前后端接入第三方 SDK 时，确保遵循 `SECURITY.md` 列出的依赖审查流程，并补充最小权限策略。

## 环境与自动化提示
- 推荐使用 `golangci-lint`（可在本地安装）对 `server/` 执行静态检查，CI 中如需启用请提交配置文件于仓库根目录。
- 前端自动化构建依赖 Node 20 与 pnpm/yarn 镜像；国内网络可沿用 Makefile 中腾讯源配置。若扩展流水线，参照 `deploy/` 中的 Dockerfile 以保证一致的打包基础镜像。
