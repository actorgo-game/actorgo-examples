# ActorGo HTTP MethodID 示例

该示例使用 ActorGo 的 `httpactor` 组件，将 HTTP 请求通过全局 MethodID 表路由到
`HTTPActor`。无需为每个业务方法手工注册 Gin 路由，所有调用统一使用：

```text
POST /actor/{methodID}
Content-Type: application/json
```

默认监听 `127.0.0.1:18080`，已注册的方法如下：

| MethodID | 请求 | 响应 | 说明 |
|---:|---|---|---|
| `5001` | `{}` | `HealthResponse` | 健康检查 |
| `5002` | `HelloRequest` | `HelloResponse` | Hello 调用 |
| `5003` | `EchoRequest` | `EchoResponse` | JSON 回显与参数校验 |
| `5004` | `ChildHelloRequest` | `ChildHelloResponse` | 父 Actor 调用动态子 Actor |

MethodID 和请求响应结构在 `methodid.go` 中声明。父 Actor 位于 `http_actor.go`，
子 Actor 位于 `http_child_actor.go`；两者分别在 `OnInit()` 中通过
`Methods().Register()` 注册方法。HTTP 入口会根据 `Content-Type` 解码请求，再通过
MethodID 自动定位父 Actor。

`5004` 会根据请求中的 `child_id` 查找子 Actor；子 Actor 不存在时，父 Actor 的
`OnFindChild()` 会动态创建它，然后使用内部 MethodID `5101` 调用子 Actor。
内部 MethodID 不会发布到 HTTP 全局路由，直接请求 `/actor/5101` 会返回 `404`。

## 编译与启动

```bash
cd /data1/home/user00/tools/actor-game/actorgo-examples/test_http
go build -o test_http .
./test_http
```

另开终端执行：

```bash
curl -X POST http://127.0.0.1:18080/actor/5001 \
  -H 'Content-Type: application/json' \
  -d '{}'

curl -X POST http://127.0.0.1:18080/actor/5002 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Codex"}'

curl -X POST http://127.0.0.1:18080/actor/5003 \
  -H 'Content-Type: application/json' \
  -H 'X-ActorGo-Timeout-Ms: 2000' \
  -d '{"message":"HTTP MethodID 测试成功"}'

curl -X POST http://127.0.0.1:18080/actor/5004 \
  -H 'Content-Type: application/json' \
  -d '{"child_id":"user-1001","message":"hello child"}'
```

成功响应带有 `X-ActorGo-Request-ID` 响应头。业务参数错误会映射为 HTTP `400`，
不存在的 MethodID 会返回 HTTP `404`。

可以用 `--path` 和 `--node` 覆盖配置文件与 NodeID。配置文件位于
`../config/test-http.json`，使用 NodeType `3` 和 NodeID `0.0.3.1`。
