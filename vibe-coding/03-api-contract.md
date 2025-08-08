# API 合约（Go + 文件系统，极简）

所有时间字段均为 ISO8601；错误响应统一为：
- status: HTTP 状态码
- code: 短码（如 INVALID_AUTH、USER_NOT_FOUND、FILE_TOO_LARGE）
- message: 人类可读说明

鉴权（MVP）：
- 非上传接口：统一使用 Basic Auth，需匹配 ADMIN_BASIC_USER/ADMIN_BASIC_PASS。
- 上传接口：Basic Auth，其中 username=用户名、password=machine-id（兼容 asciinema CLI）。

## 列出所有用户（前端 Home 页面使用）
GET /api/users
Headers: Authorization: Basic <admin>
Responses：
- 200 { items: [ "alice", "bob" ], total }
说明：
- 通过遍历 STORAGE_ROOT 下的目录，筛选包含 machine-id 文件的目录名作为合法用户名；按字典序返回。

## 创建/更新用户（管理）
POST /api/users
Headers: Authorization: Basic <admin>
Body (JSON): { "username": "alice" }
行为：
- 若不存在用户目录 STORAGE_ROOT/alice，则创建该目录与 machine-id 文件，并随机生成一个 machine-id 写入。
- 若已存在，则读取并返回该用户已存在的 machine-id。
Responses:
- 201 { username, machineId }（首次创建）
- 200 { username, machineId }（已存在返回现有 machine-id）
- 400 code=VALIDATION_ERROR
- 409 code=CONFLICT（用户名非法或与保留名冲突，如以点开头等）

## 上传 cast 文件（客户端，兼容 asciinema CLI）
POST /api/asciicasts
Headers:
- Authorization: Basic <username:machine-id>
- Content-Type: multipart/form-data; boundary=...
Fields:
- asciicast: .cast 文件（二进制，字段名固定为 asciicast）

Responses:
- 201 { username, relPath: "<username>/<YYYYMMDD>/<uuid>.cast", sizeBytes }
- 401 code=INVALID_AUTH（machine-id 无效或未绑定到该 username）
- 413 code=FILE_TOO_LARGE（可配置）

说明：
- 服务端从 Basic 的用户名获取 username，从密码获取 machine-id 并校验与 username 的绑定关系。
- 保存路径：STORAGE_ROOT/用户名/YYYYMMDD/uuid.cast

## 列出用户的 casts（管理或用户）
GET /api/users/:username/casts
Headers: Authorization: Basic <admin>
Query: page, pageSize（可选；如不需要分页可省略）
Responses：
- 200 { items: [ { relPath, sizeBytes, mtime } ], page?, pageSize?, total? }
- 404 USER_NOT_FOUND

## 获取 cast 原始文件（用于播放器 src）
GET /api/casts/file?path=<username/yyyymmdd/uuid.cast>
Headers: Authorization: Basic <admin>
- 说明：通过 relPath 定位文件；参数需严格校验避免目录穿越。
- 200 Content-Type: application/octet-stream
- 404

## 健康检查
GET /healthz -> 200 { status: "ok" }

错误码列表：
- INVALID_AUTH（鉴权失败）
- USER_NOT_FOUND（用户不存在或未绑定 machine-id）
- CONFLICT（资源冲突）
- FILE_TOO_LARGE（超限）
- CAST_NOT_FOUND（文件或元数据缺失）
- INTERNAL_ERROR（未分类服务端错误）
 - INVALID_USERNAME（用户名不合法）
