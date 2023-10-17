# mysql 脚本处理工具

`mysqlsed` 是一个 sql 文件处理工具，一般用于整库关键字替换。

`mysqlsed` 首先逐行分析每一条语句，只有插入语句才进入后续处理，判断该行包含指定的关键字进行替换，并将插入语句转换为更新语句，然后写入到 `out.sql` 文件。

操作步骤：

1. 导出整个数据库为 sql 文件，导出的数据为一行行的 insert 脚本。

```shell
mysqldump -h {HOST} -P {PORT} -u{USER} -p{PASSWORD} \
--default-character-set=utf8mb4 \
--hex-blob \
--routines \
--skip-extended-insert \
--skip-add-drop-table \
--skip-create-options \
--set-gtid-purged=OFF \
{DB} > {DB}.sql
```

2. 使用 mysqlsed 处理文件。

```shell
mysqlsed {DB}.sql [Keyword>>>TargetWord] [Keyword>>>TargetWord] ...
```

替换后将会生成 `out.sql` 文件。

3. 将 `out.sql` 在原来的数据库上执行。
