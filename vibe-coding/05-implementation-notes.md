# 实现要点与边界处理

- Basic Auth 解析：
  - 上传端：用户名即 username，password 即 machine-id（与 asciinema CLI 保持一致）。
  - 管理端：完整校验用户名/密码匹配 ADMIN_BASIC_*。

- 上传大小限制：
  - Gin 基于 net/http，可通过 c.Request.Body 限制与 MaxMultipartMemory（如 50MB），从 env 调整。

- multipart 字段：
  - 上传表单字段名固定为 `asciicast`，Content-Type 可为 application/octet-stream。
  - Gin 可使用 c.FormFile("asciicast")/c.SaveUploadedFile 或手动读取流以实现原子写。

- 路径与安全：
  - 严禁使用客户端提供的文件名，统一服务端生成 uuid.cast。
  - 路径只由服务端基于 username 与日期构造，防止目录穿越。

- 日期格式：
  - 使用服务器本地时区 YYYYMMDD；如需 UTC，可在配置中切换。

- 幂等性：
  - 同一内容重复上传视为不同条目（不同 uuid）；如需去重，后续可增 content-hash。

- 元数据：
  - 不额外持久化元数据，必要信息（size/mtime）通过文件系统获取。

- 权限：
  - 非管理员只能访问与自己绑定 username 的资源；通过 Basic 中 machine-id 在文件系统读取 STORAGE_ROOT/<username>/machine-id 做等值校验。

- 错误处理：
  - 统一返回 { code, message }；日志记录包含 requestId。
  - Gin 中通过中间件注入 requestId（如 X-Request-Id），并在响应中回显。

- 测试：
  - 单元：auth、path 构造、fs 读写原子性；
  - 集成：创建用户（生成并写 machine-id）、上传（写 .cast）、列表、播放文件 200。

- 部署：
  - 首期单容器即可：server 进程托管 web 静态文件；
  - 卷挂载 STORAGE_ROOT；
  - 备份文件树（含 .by-id 索引与 machine-ids）。

- 文件系统写入原子性：
  - 采用临时文件 + rename 的写法保证原子性；先写 <uuid>.cast.tmp，再 rename 为 .cast。
  - 失败处理：清理残留 .tmp。
