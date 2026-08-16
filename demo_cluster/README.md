# 多节点示例

`demo_cluster` 演示由 master、center、web、gate 和 game 组成的多节点游戏服务。示例没有使用持久化数据库，进程重启后运行数据会还原。

示例包含两种客户端：

- `robot_client`：Go 实现的 AGP/1 TCP/Protobuf 压测客户端。
- `nodes/web/view/`：H5 实现的 AGP/1 WebSocket/Protobuf 客户端。

服务端节点、H5 客户端和压测机器人均已迁移到新框架的 AGP/1 协议。H5 客户端使用 `agp.v1` WebSocket 子协议，机器人使用 AGP/1 TCP 长度帧；二者都使用数字 MethodID。

## 环境要求

- Go 1.26.1 或更高版本。
- NATS Server 2.0 或更高版本，默认监听 `4222` 端口。

`go.mod` 通过 `replace` 固定使用相邻目录中的本地 ActorGo。目录结构应为：

```text
actor-game/
├── actorgo/
└── actorgo-examples/
```

在 `actorgo-examples` 根目录编译节点：

```bash
go build -o demo_cluster/nodes/nodes ./demo_cluster/nodes
```

## 配置文件

- 节点配置：`config/demo-cluster.json`
- 业务数据配置：`config/data/`
- 区服与 game 节点映射：`config/data/areaServerConfig.json`

### NodeID 规则

新框架的启动参数使用以下四段整数格式：

```text
BigWorldID.WorldID.NodeType.NodeInst
```

本示例的 NodeType 分配如下：

| 节点 | NodeType | 启动参数 NodeID | 配置中的内部 NodeID |
|---|---:|---|---:|
| master | 1 | `0.0.1.1` | `262145` |
| center | 2 | `0.0.2.1` | `524289` |
| web | 3 | `0.0.3.1` | `786433` |
| gate | 4 | `0.0.4.1` | `1048577` |
| game | 5 | `1.1.5.1` | `18014467230269441` |

命令行 `--node` 使用四段形式，框架内部会将其编码为 `uint64`。`master_node_id` 可使用四段形式或编码后的十进制字符串，本示例使用 `262145`。

同一 NodeType 只有一条配置时可以省略 `node_id`；框架根据 `--node` 中的 NodeType 选择该配置。同一 NodeType 存在多条配置时必须为每条配置提供 `node_id`，否则框架会报告配置歧义。

业务 `serverId` 与 ActorGo `nodeId` 相互独立，本示例的映射为：

| serverId | game NodeID | 内部 NodeID |
|---:|---|---:|
| 10001 | `1.1.5.1` | `18014467230269441` |
| 10002 | `1.2.5.1` | `18014535949746177` |
| 10003 | `1.3.5.1` | `18014604669222913` |
| 10004 | `1.4.5.1` | `18014673388699649` |

新增 game 区服时，需要更新 `areaServerConfig.json` 中业务 `serverId` 到内部 NodeID 的映射。如果同一个 NodeType 需要使用不同的 address 或 settings，还需要在 `demo-cluster.json` 中拆分配置并增加 `node_id`。

## 启动

### 1. 启动 NATS

Linux 或 macOS：

```bash
nats-server
```

Windows 可在 `actorgo-examples` 根目录执行：

```bat
3rd\nats-server\run_nats.bat
```

日志出现 `Listening for client connections on 0.0.0.0:4222` 表示 NATS 已启动。

### 2. 使用脚本管理五个节点（推荐）

NATS 启动后，在 `actorgo-examples` 根目录执行：

```bash
./demo_cluster/cluster.sh start
./demo_cluster/cluster.sh status
./demo_cluster/cluster.sh stop
```

脚本还支持重启全部节点，或只管理指定节点：

```bash
./demo_cluster/cluster.sh restart
./demo_cluster/cluster.sh restart game
./demo_cluster/cluster.sh start gate
./demo_cluster/cluster.sh stop gate
```

`start` 会先确认 NATS 的 `127.0.0.1:4222` 端口可连接，然后使用相邻目录 `../actorgo` 中的本地框架重新编译节点，并等待每个节点输出就绪日志。进程日志保存在 `demo_cluster/logs/process/`，PID 文件保存在 `demo_cluster/logs/run/`。

脚本使用绝对路径定位配置和工作目录，因此也可以从其他目录调用。

脚本只会停止自己启动并记录了 PID 的进程。如果节点已经通过手工命令启动，`status` 会显示 `EXTERNAL`，`start` 会拒绝重复启动，`stop` 也不会结束该外部进程。

### 3. 手工启动五个节点

所有节点都从 `demo_cluster/nodes/main.go` 启动。为确保 web 节点可以找到 `./static` 和 `./view`，下面所有命令都从 `demo_cluster/nodes/web` 目录执行。

依次打开五个终端，每个终端先执行：

```bash
cd actorgo-examples/demo_cluster/nodes/web
```

然后按 master、center、web、gate、game 的顺序分别运行：

```bash
go run .. 1 --path=../../../config/demo-cluster.json --node=0.0.1.1
```

```bash
go run .. 2 --path=../../../config/demo-cluster.json --node=0.0.2.1
```

```bash
go run .. 3 --path=../../../config/demo-cluster.json --node=0.0.3.1
```

```bash
go run .. 4 --path=../../../config/demo-cluster.json --node=0.0.4.1
```

```bash
go run .. 5 --path=../../../config/demo-cluster.json --node=1.1.5.1
```

默认端口：

| 服务 | 协议 | 端口 |
|---|---|---:|
| NATS | TCP | 4222 |
| web | HTTP | 8081 |
| gate | WebSocket | 10010 |
| gate | TCP | 10011 |

Web 节点默认监听 `:8081`（所有本机网卡）。因此机器人可以使用
`http://127.0.0.1:8081`；从其他机器访问时，将 `127.0.0.1` 替换为服务端 IP。
修改监听地址后需要重启 web 节点才能生效。

## 测试

### H5 客户端

五个节点启动成功后，可访问 [http://127.0.0.1:8081/hello](http://127.0.0.1:8081/hello) 验证 web 到 center 的调用链；访问根路径即可使用 AGP/1 H5 客户端完成注册、登录、选服、创建/选择角色和进入游戏流程。

### 压测机器人

机器人默认执行一次注册、HTTP 登录、AGP 握手、网关登录、查询/创建角色和进入游戏：

```bash
go run ./demo_cluster/robot_client \
  --count=1 \
  --web=http://127.0.0.1:8081 \
  --gate=127.0.0.1:10011 \
  --server=10001 \
  --account-prefix=robot
```

使用 `--count` 调整并发数，`--verbose` 输出详细日志，`--keepalive` 在流程成功后保持连接和心跳。

## 目录说明

- `internal/code`：业务状态码。
- `internal/component`：启动前检查 center 节点等组件。
- `internal/data`：`config/data` 下业务配置的结构与加载逻辑。
- `internal/event`：游戏事件。
- `internal/guid`：全局 ID 生成。
- `internal/pb`：生成后的 Protobuf 代码。
- `internal/protocol`：Protobuf 定义。
- `internal/rpc`：跨节点 RPC 封装。
- `nodes/master`：基于 NATS 的节点发现服务。
- `nodes/center`：帐号与全局业务。
- `nodes/web`：HTTP 接口和 H5 客户端。
- `nodes/gate`：连接管理、消息路由与转发。
- `nodes/game`：分服游戏逻辑。
- `robot_client`：TCP/Protobuf 压测客户端。

## 运行截图

![screenshot](screenshot.png)
