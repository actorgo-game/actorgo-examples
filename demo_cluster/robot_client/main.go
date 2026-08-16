package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	chttp "github.com/actorgo-game/actorgo/extend/http"
	clog "github.com/actorgo-game/actorgo/logger"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	robotCount := flag.Int("count", 1, "number of robots")
	webURL := flag.String("web", "http://127.0.0.1:8081", "web node URL")
	gateAddress := flag.String("gate", "127.0.0.1:10011", "gate TCP address")
	serverID := flag.Int("server", 10001, "game server ID")
	pid := flag.String("pid", "2126001", "SDK package ID")
	accountPrefix := flag.String("account-prefix", "test", "robot account prefix")
	printLog := flag.Bool("verbose", false, "print detailed robot logs")
	keepAlive := flag.Bool("keepalive", false, "keep successful robot connections alive")
	flag.Parse()

	if *robotCount < 1 {
		clog.Warn("robot count must be positive")
		return
	}

	wg := sync.WaitGroup{}
	var failed atomic.Int32

	accounts := make(map[string]string)
	for i := 1; i <= *robotCount; i++ {
		key := fmt.Sprintf("%s%d", *accountPrefix, i)
		accounts[key] = key
	}

	if err := RegisterDevAccount(*webURL, accounts); err != nil {
		clog.Error("register through web node %s failed: %v; start the web node or set --web", *webURL, err)
		os.Exit(1)
	}

	for userName, password := range accounts {
		time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
		wg.Add(1)
		go func(userName, password string) {
			defer wg.Done()
			robot := RunRobot(*webURL, *pid, userName, password, *gateAddress, int32(*serverID), *printLog)
			if robot == nil {
				failed.Add(1)
				return
			}
			if !*keepAlive {
				robot.Close()
			}
		}(userName, password)
	}

	wg.Wait()
	if failed.Load() > 0 {
		clog.Error("%d robot(s) failed", failed.Load())
		os.Exit(1)
	}
	if *keepAlive {
		select {}
	}
}

func RegisterDevAccount(url string, accounts map[string]string) error {
	requestURL := fmt.Sprintf("%s/register", url)

	for key, val := range accounts {
		params := map[string]string{
			"account":  key,
			"password": val,
		}

		jsonBytes, _, err := chttp.GET(requestURL, params)
		if err != nil {
			return fmt.Errorf("register account %s: %w", key, err)
		}

		rsp := &code.Result{}
		err = jsoniter.Unmarshal(jsonBytes, rsp)
		if err != nil {
			return fmt.Errorf("decode register response for account %s: %w", key, err)
		}

		clog.Debug("register account = %s, result = %+v", key, rsp)
	}
	return nil
}

func RunRobot(url, pid, userName, password, addr string, serverId int32, printLog bool) *Robot {

	// 创建客户端
	cli := New(
		NewAGPClient(10 * time.Second),
	)
	cli.PrintLog = printLog

	// 登录获取token
	if err := cli.GetToken(url, pid, userName, password); err != nil {
		clog.Error(err.Error())
		return nil
	}

	// 根据地址连接网关
	if err := cli.Connect(addr); err != nil {
		clog.Error(err.Error())
		return nil
	}

	if cli.PrintLog {
		clog.Info("tcp connect %s is ok", addr)
	}

	// 随机休眠
	cli.RandSleep()

	// 用户登录到游戏节点
	err := cli.UserLogin(serverId)
	if err != nil {
		clog.Warn(err.Error())
		return nil
	}

	if cli.PrintLog {
		clog.Info("user login is ok. [user = %s, serverId = %d]", userName, serverId)
	}

	//cli.RandSleep()

	// 查看是否有角色
	err = cli.PlayerSelect()
	if err != nil {
		clog.Warn(err.Error())
		return nil
	}

	//cli.RandSleep()

	// 创建角色
	err = cli.ActorCreate()
	if err != nil {
		clog.Warn(err.Error())
		return nil
	}

	//cli.RandSleep()

	// 角色进入游戏
	err = cli.ActorEnter()
	if err != nil {
		clog.Warn(err.Error())
		return nil
	}

	elapsedTime := cli.StartTime.NowDiffMillisecond()
	clog.Debug("[%s] is enter to game. elapsed time:%dms", cli.TagName, elapsedTime)

	// cli.Close()

	return cli
}
