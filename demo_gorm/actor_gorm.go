package main

import (
	"time"

	cgorm "github.com/actorgo-game/actorgo/components/gorm"
	clog "github.com/actorgo-game/actorgo/logger"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	"gorm.io/gorm"
)

type ActorDB struct {
	cactor.Base
	centerDB *gorm.DB
}

func (p *ActorDB) AliasID() string {
	return "db"
}

// OnInit Actor初始化前触发该函数
func (p *ActorDB) OnInit() {
	// db配置的注解
	// 打开 demo-gorm.json，找到 NodeType=5 和 db 配置。
	// 当前示例的 NodeID 为 0.0.5.1。
	// db_id_list 配置 center_db_1，表示当前节点需要连接该数据库。
	// main 中注册 GORM 组件后，可以通过组件获取对应的 *gorm.DB。

	component := p.App().Find(cgorm.Name)
	gormComponent, ok := component.(*cgorm.Component)
	if !ok || gormComponent == nil {
		clog.Panic("[component = %s] not found.", cgorm.Name)
	}

	// 获取 db_id = "center_db_1" 的配置
	centerDBID := p.App().Settings().GetConfig("db_id_list").GetString("center_db_id")
	p.centerDB = gormComponent.GetDb(centerDBID)
	if p.centerDB == nil {
		clog.Panic("database %q not found", centerDBID)
	}

	// 每五秒查询一次数据库。
	p.Timer().Add(5*time.Second, p.selectDB)
	// 一秒后进行一次分页查询。
	p.Timer().AddOnce(1*time.Second, p.selectPagination)
}

func (p *ActorDB) selectDB() {
	userBindTable := &UserBindTable{}
	tx := p.centerDB.First(userBindTable)
	if tx.Error != nil {
		clog.Warn(tx.Error.Error())
	}

	clog.Info("%+v", userBindTable)
}

func (p *ActorDB) selectPagination() {
	list, count := p.pagination(1, 10)
	clog.Info("count = %d", count)

	for _, table := range list {
		clog.Info("%+v", table)
	}
}

// pagination 分页查询
func (p *ActorDB) pagination(page, pageSize int) ([]*UserBindTable, int64) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	var list []*UserBindTable
	var count int64

	if err := p.centerDB.Model(&UserBindTable{}).Count(&count).Error; err != nil {
		clog.Warn(err.Error())
		return nil, 0
	}

	if count > 0 {
		list = make([]*UserBindTable, 0, pageSize)
		s := p.centerDB.Limit(pageSize).Offset((page - 1) * pageSize)
		if err := s.Find(&list).Error; err != nil {
			clog.Warn(err.Error())
		}
	}

	return list, count
}
