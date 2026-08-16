# ActorGo 示例仓库

本仓库通过 `go.mod` 中的 `replace` 固定使用相邻目录 `../actorgo` 的本地框架代码：

```text
actor-game/
├── actorgo/
└── actorgo-examples/
```

## 完整示例

- [单节点 AGP/1 聊天室](demo_chat)
- [多节点游戏集群](demo_cluster)
- [GORM 组件](demo_gorm)

## 独立功能示例

| 目录 | 内容 |
|---|---|
| `test_actor` | MethodID、父子 Actor 和类型化调用 |
| `test_data_config` | JSON 数据配置加载与查询 |
| `test_discovery` | NATS master 与两个 game 节点发现 |
| `test_func_call` | Go 反射函数调用 |
| `test_gin` | Gin Controller 和中间件 |
| `test_gob` | Go GOB 编解码 |
| `test_goroutine` | Goroutine 行为测试 |
| `test_http` | ActorGo HTTP MethodID 路由及动态子 Actor 调用 |
| `test_logger` | ActorGo 日志接口 |
| `test_nats` | NATS publish/subscribe/request/reply |
| `test_protobuf` | Protobuf 基础类型 |
| `test_redis` | Redis 读写和订阅 |
| `test_zap` | Zap 日志切割 |

## 编译与测试

```bash
cd actorgo-examples
go test ./...
```

需要 NATS、MySQL 或 Redis 的运行示例，应先启动对应外部服务并确认 `config/` 中的连接参数。
