# 前端 SPA 规划（React + asciinema-player）

路由：
- /                 欢迎与概览
- /users/:username  用户的 cast 列表（通过 relPath 展示）
- /play/*           播放页（url 中直接包含 relPath，内嵌 asciinema-player）

数据来源（MVP 均需管理 Basic）：
- 用户列表（Home）：GET /api/users -> [username]
- 列表：GET /api/users/:username/casts -> [{ relPath, sizeBytes, mtime }]
- 播放源：/api/casts/file?path=<relPath>
说明：前端需在请求中附带管理员 Basic 认证头（由部署环境配置）。

UI 与播放器：
- 方案 A：npm 包 asciinema-player（需要引入 CSS 与 JS）
- 方案 B：CDN 脚本 + 自定义 React 封装组件
	- 在播放页将 src 指向 /api/casts/file?path=<relPath>
样式：Tailwind CSS；如需对话框/下拉等，按需使用 Headless UI 或 Radix；图标 lucide-react。

UI 要点：
- 简洁卡片式列表（显示日期、大小）
- 播放页支持自动适配宽度/主题切换（可选）
- 错误与空态

构建与部署：
- Vite 构建产物 dist/ 由后端静态托管
- 前端 .env 使用相对路径请求后端（同域）
