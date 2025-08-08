# Phase 1 任务清单（Go + 文件系统，极简）

- [ ] server: 初始化 Go + Gin 项目（/cmd/server），/healthz，注册 Logger/Recovery
- [ ] server: 管理端 Basic 中间件（env 校验）
- [ ] server: POST /api/users（创建用户目录与 machine-id，若存在则返回现有 machine-id）
- [ ] server: 上传端 Basic 中间件（用户名=用户名，password=machine-id -> 读取 user/machine-id 校验）
- [ ] server: POST /api/asciicasts（multipart，表单字段 asciicast，原子写入 <uuid>.cast）
- [ ] server: GET /api/users（扫描 STORAGE_ROOT，返回包含 machine-id 的目录名）
- [ ] server: GET /api/users/:username/casts（扫描目录，返回 relPath、size、mtime）
- [ ] server: GET /api/casts/file?path=<relPath>（读取并返回文件）
- [ ] server: 错误码统一与日志
- [ ] test: 集成测试（创建 -> 上传 -> 列表 -> 获取文件）
- [ ] docs: README 用法与 curl 示例
