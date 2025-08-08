# 架构与技术选型（Go + 文件系统，无数据库）

建议后端：Go 1.22+ + Gin（github.com/gin-gonic/gin）
- 原因：部署简单、单可执行文件、性能好、生态完善；Gin 提供路由、中间件与便捷的请求处理。
- Web/路由：Gin（Logger/Recovery 中间件，统一注册路由与静态资源）。
- 鉴权：Basic Auth（上传端以 machine-id 作为密码；管理端单独管理用户名/密码）。
- 存储：仅使用本地文件系统 storage/
  - 路径：storage/<username>/<YYYYMMDD>/<uuid>.cast
  - 用户绑定：storage/<username>/machine-id（文本，单行单值）
  - 不额外保存 JSON 元数据（最简）。
- 静态资源：后端同时服务前端构建产物（web/dist）。

建议前端：Vite + React + TypeScript
- 播放器：asciinema-player（使用 npm 包或 CDN 资源）
- 路由：react-router-dom
- UI：Tailwind CSS 为主；无需重型组件库。需要基础交互可按需使用 Headless UI 或 Radix Primitives；图标使用 lucide-react。

目录规划（建议）：

- server/
  - cmd/server/main.go（入口）
  - internal/
    - handler/
  - users.go（创建/更新用户，Gin handler）
  - uploads.go（上传、列表、取文件，Gin handler）
      - health.go
    - auth/
  - basic.go（Basic 校验中间件：admin 与 machine-id）
    - storage/
      - fs.go（路径构造、原子写入、扫描）
    - util/
      - time.go、uuid.go、resp.go
  - storage/（运行时生成）
  - go.mod、go.sum

- web/
  - src/
    - main.tsx、App.tsx
    - pages/
      - Home.tsx（概览/欢迎）
  - UserCasts.tsx（某用户的 cast 列表）
  - Player.tsx（嵌入 asciinema-player）
    - components/
      - AsciinemaPlayer.tsx（React 包装）
    - api/
      - client.ts（fetch 封装）
  - index.html、vite.config.ts、package.json

- vibe-coding/（规划与文档）

运行形态：
- 单进程后端服务（端口如 8080），前端构建后由后端托管。
- 配置经由环境变量：
  - STORAGE_ROOT（默认 ./server/storage）
  - ADMIN_BASIC_USER / ADMIN_BASIC_PASS（管理接口鉴权）
  - HTTP_PORT（默认 8080）

日志与可观测：
- 使用 Gin Logger/Recovery 中间件，记录请求与异常；可选集成 slog/zap 作为后端日志实现。
- 简单健康检查 /healthz。
