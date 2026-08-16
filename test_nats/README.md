# NATS 示例

四个示例默认连接 `nats://127.0.0.1:4222`，并使用与 ActorGo 示例配置一致的认证信息：

```text
user: actorgo
password: actorgo2026
subject: actorgo.examples.game.10001
```

可以通过 `--url`、`--user`、`--password`、`--subject` 参数覆盖，也可以设置
`NATS_URL`、`NATS_USER`、`NATS_PASSWORD`、`NATS_SUBJECT` 环境变量。

发布/订阅测试需要先运行订阅端：

```bash
./subscribe/subscribe
./publish/publish --count=10
```

请求/响应测试需要先运行响应端：

```bash
./reply/reply
./request/request --count=10
```

`subscribe` 和 `reply` 使用 `Ctrl+C` 正常退出；`publish` 和 `request` 完成指定次数后自动退出。
