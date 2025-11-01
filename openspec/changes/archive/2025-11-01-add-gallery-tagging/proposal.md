## Why
当前后台仅支持对单条作品记录维护标签，批量导入后的素材需要逐条修改，效率低且容易遗漏，前端展示页也无法按多个标签组合检索。

## What Changes
- 新增后端批量打标签接口及服务逻辑，支持一次性为多张作品添加或移除标签。
- 扩展作品列表检索条件，允许前端/后台根据标签、款式组合筛选，并返回命中计数。
- 在前端作品管理页面提供批量选择与标签编辑对话框，同时同步更新展示页的标签过滤请求。
- 为操作记录新增审计日志字段，保留批量操作人和标签调整详情。

## Impact
- Affected specs: gallery-management
- Affected code: `server/api/v1/picture`, `server/service/picture`, `server/model/request/picture`, `server/router/picture`, `web/src/view/gallery`, `web/src/api/gallery.ts`, `web/src/pinia/modules/gallery.ts`
