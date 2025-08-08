# 分阶段开发路线（DoD 驱动）

- 初始化后端（Go）与前端（Vite+React+TS）骨架
- ENV 配置读取（STORAGE_ROOT、ADMIN_BASIC_*、HTTP_PORT）
- 创建 STORAGE_ROOT 目录
- DoD：本地能启动后端 /healthz 返回 ok；首次启动即创建存储根目录

- 管理端 Basic 鉴权中间件（MVP：除上传接口外，其余端点统一使用管理员 Basic）
- POST /api/users（创建用户或返回现有 machine-id，写 machine-id 文件）
- 上传端 Basic（密码=machine-id）校验（读取 user/machine-id）
- POST /api/asciicasts 接收 multipart 文件（字段 asciicast），保存到 STORAGE_ROOT/用户名/日期/uuid.cast
- GET /api/users（扫描 STORAGE_ROOT，列出含 machine-id 的用户）
- 简单列表 GET /api/users/:username/casts（扫描目录收集 *.cast，附带 size/mtime）
- DoD：通过 curl 完成“创建用户 -> 上传 -> 列表”闭环；文件落盘路径正确

## Phase 2：播放与前端 SPA
- GET /api/casts/file?path=<relPath> 提供文件流
- 前端：列表页（按 relPath 展示）、播放页（传 relPath 给播放器），后端托管静态资源
- DoD：浏览器可列表并点击播放

## Phase 3：安全与体验
- 文件大小限制、MIME 校验、错误码统一
- 访问控制：用户仅能看自己的资源（管理员可看所有）
- 简单搜索/分页/排序
- DoD：基本安全措施就绪，常见错误可读

## Phase 4：可部署与运维
- Dockerfile 与容器运行文档
- 简单日志轮转、基础监控指标（可选）
- 备份与恢复说明（文件树）
- DoD：容器一键起，最小备份/恢复路径可演练

里程碑节奏：每个 Phase 以 0.5~1.5 天为单位（视复杂度与验收范围可调）。
