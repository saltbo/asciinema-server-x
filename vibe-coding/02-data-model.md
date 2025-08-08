# 数据模型（文件系统，极简）

根目录：STORAGE_ROOT（默认 ./server/storage）

用户目录结构：
- STORAGE_ROOT/<username>/
	- machine-id   （文本文件；仅一行，一个 machineId）
	- YYYYMMDD/    （日期目录）
		- <uuid>.cast  （原始 cast 文件，仅此）

说明：
- 不保存额外元数据文件；所有展示通过“路径扫描 + 文件信息（ctime/size）”实现。
- 若需要获取大小或时间，可用 os.Stat 获取 size 与 mtime/ctime。

查询方式：
- 按用户列出：遍历 STORAGE_ROOT/<username> 下的日期目录，收集 *.cast 路径。
- 播放：通过相对路径回到具体文件。

一致性（写入流程）：
- 先将上传数据写入 <uuid>.cast.tmp，完成后 rename 为 <uuid>.cast（原子）。失败则清理 .tmp。
