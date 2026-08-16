# GORM 组件示例

`demo_gorm` 演示在 Actor 中通过 ActorGo GORM 组件访问 MySQL，并执行单条查询和分页查询。

## 环境要求

- 使用相邻目录 `../actorgo` 中的本地 ActorGo 框架。
- MySQL 中已创建 `dev_center` 数据库。
- Go 版本与仓库根目录 `go.mod` 的要求一致。

## 初始化数据库

将 `demo_gorm/db.sql` 导入 MySQL：

```bash
mysql -u root -p < demo_gorm/db.sql
```

然后修改 `config/demo-gorm.json` 中 `center_db_1` 的连接参数：

- `host`
- `user_name`
- `password`
- `db_name`

也可以直接填写 `dsn`。`dsn` 为空时，组件根据上述字段生成连接字符串。

## NodeID

示例沿用 game 节点的 NodeType `5`，默认 NodeID 为：

```text
0.0.5.1
```

配置中同一 NodeType 只有一项，因此可以省略 `node_id`。

## 启动

从 `demo_gorm` 目录启动：

```bash
cd demo_gorm
go run .
```

也可以显式指定配置和节点：

```bash
go run . --path=../config/demo-gorm.json --node=0.0.5.1
```

程序启动后会在一秒后执行分页查询，并每五秒查询一次 `user_bind` 表。日志配置来自 `config/include/logger.json`，默认写入 `demo_gorm/logs/`。
