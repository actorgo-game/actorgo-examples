package db

import (
	cerror "github.com/actorgo-game/actorgo/error"
	cstring "github.com/actorgo-game/actorgo/extend/string"
	ctime "github.com/actorgo-game/actorgo/extend/time"
	clogger "github.com/actorgo-game/actorgo/logger"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/guid"
)

// DevAccountTable 开发模式的帐号信息表(platform.TypeDevMode)
type DevAccountTable struct {
	AccountId   int64  `gorm:"column:account_id;primary_key;comment:'帐号id'" json:"accountId"`
	AccountName string `gorm:"column:account_name;comment:'帐号名'" json:"accountName"`
	Password    string `gorm:"column:password;comment:'密码'" json:"-"`
	CreateIP    string `gorm:"column:create_ip;comment:'创建ip'" json:"createIP"`
	CreateTime  int64  `gorm:"column:create_time;comment:'创建时间'" json:"createTime"`
}

func (*DevAccountTable) TableName() string {
	return "dev_account"
}

func DevAccountRegister(accountName, password, ip string) int32 {
	devAccount, _ := DevAccountWithName(accountName)
	if devAccount != nil {
		return code.AccountNameIsExist
	}

	devAccountTable := &DevAccountTable{
		AccountId:   guid.Next(),
		AccountName: accountName,
		Password:    password,
		CreateIP:    ip,
		CreateTime:  ctime.Now().Unix(),
	}

	devAccountCache.Put(accountName, devAccountTable)
	// TODO 保存db

	return code.OK
}

func DevAccountWithName(accountName string) (*DevAccountTable, error) {
	val, found := devAccountCache.GetIfPresent(accountName)
	if found == false {
		return nil, cerror.Error("account not found")
	}

	return val.(*DevAccountTable), nil
}

// loadDevAccount 节点启动时预加载帐号数据
func loadDevAccount() {
	// 演示用，直接手工构建几个帐号
	for i := 1; i <= 10; i++ {
		index := cstring.ToString(i)

		devAccount := &DevAccountTable{
			AccountId:   guid.Next(),
			AccountName: "test" + index,
			Password:    "test" + index,
			CreateIP:    "127.0.0.1",
			CreateTime:  ctime.Now().ToMillisecond(),
		}

		devAccountCache.Put(devAccount.AccountName, devAccount)
	}

	clogger.Info("preload DevAccountTable")
}
