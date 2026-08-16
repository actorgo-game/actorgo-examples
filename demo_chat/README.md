# 单节点聊天室示例

`demo_chat` 使用本地 ActorGo 构建一个单节点多人聊天室：

- 浏览器通过 AGP/1 WebSocket 连接 `34590` 端口。
- AGP Packet 使用 Protobuf，业务 Body 使用 JSON。
- HTTP `8081` 端口用于部署 H5 静态页面。
- 示例是 Standalone 模式，不依赖 NATS。

## NodeID 与 MethodID

聊天室使用 NodeType `1`，默认 NodeID 为 `0.0.1.1`。配置中该类型只有一项，因此省略 `node_id`。

| 功能 | MethodID | 类型 |
|---|---:|---|
| 登录 | `1001` | Request |
| 发送消息 | `1002` | Notify |
| 连接退出 | `1003` | 服务端内部 Notify |
| 新用户广播 | `1101` | 服务端 Notify |
| 聊天消息广播 | `1102` | 服务端 Notify |
| 余额通知 | `1103` | 服务端 Notify |

## 启动

```bash
cd actorgo-examples/demo_chat/room
go run .
```

看到以下日志表示服务已启动：

```text
Websocket connector listening at Address :34590
http run. http://:8081
```

在浏览器打开两个 [http://127.0.0.1:8081](http://127.0.0.1:8081) 页面，在任一页面发送消息，两个页面都会收到广播。页面会自动使用当前网页主机连接 WebSocket，无需手工修改公网 IP。

配置文件位于 `config/demo-chat.json`。

## 代码结构

- `room/main.go`：注册 AGP server、room Actor 和 HTTP 静态服务。
- `room/actor_room.go`：连接绑定、用户状态、广播和断线清理。
- `room/protocol.go`：JSON Body 结构和 MethodID。
- `static/agp-client.js`：浏览器 AGP/1 JSON 客户端，包含握手、请求、通知和心跳。
- `static/index.html`：聊天室页面。
