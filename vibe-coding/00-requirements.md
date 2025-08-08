# 需求与约束清单

目标：实现一个自托管的 asciinema 服务器，满足如下功能：

1) 上传接口：兼容 asciinema 命令行上传 cast 文件。
- 鉴权：使用 Basic Auth（username=用户名，password=machine-id）。
- 存储：文件保存路径为 base_dir/用户名/年月日/uuid.cast。

2) 用户创建接口：
- Basic Auth 保护（建议使用“管理凭据”与普通上传凭据区分）。
- 请求体仅包含 username；若不存在则创建并生成 machine-id；若存在则返回既有 machine-id。

3) 浏览与播放：
- 提供用户列表、用户 cast 列表与文件读取接口（MVP 统一使用管理员 Basic）。
- 提供前端 SPA（React）集成 asciinema-player 进行在线播放。

4) 非功能性要求：
- 简单、可部署（单进程 + 本地文件存储，无数据库）。
- 具备最小化日志、错误码与基础测试。
- 可运行于容器（可选）。

隐含需求与假设：
- 假设 machine-id 全局唯一，且一个 machine-id 只绑定一个用户（每用户仅一个 machine-id）。
- 上传接口需支持 multipart/form-data（asciinema upload 常见形态）；也可兼容 application/octet-stream。
- 日期粒度按服务器收到请求的本地日期（YYYYMMDD），而非 cast 内 header 时间。
- 文件名使用服务端生成的 uuid（v4）。
- 不做公开匿名访问，需登录后端 API（Basic）或仅前端公开只读（可选）。

待澄清问题：
- 管理员凭据的发放方式与轮换策略？
- 是否需要用户自助注册（无需管理员）？
- 是否需要软/硬配额与清理策略？
- 播放页面是否需要分享链接与一次性 token？

交付物边界：
- 首期落地：本地文件存储（无数据库）；REST API；React SPA 播放。
- 未来可选：对象存储（S3 兼容）、OIDC 登录、审计日志、K8s/Helm。
