package monitor

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt/server"
	"time"
)


func (m *Monitor) HookWrapper() server.HookWrapper {
	return server.HookWrapper{
		OnConnectedWrapper: m.OnConnectedWrapper,
		OnClosedWrapper:    m.OnClosedWrapper,
	}
}

func (m *Monitor) OnConnectedWrapper(pre server.OnConnected) server.OnConnected {
	return func(ctx context.Context, client server.Client) {
		fmt.Println("网关连接成功")
		var id int
		pre(ctx, client)
		//网关名称
		clientName := client.ClientOptions().ClientID
		//网关地址
		ip := client.Connection().RemoteAddr().String()
		opsTime := time.Now().Format("2006-01-02 15:04:05")
		//数据查询
		_ = m.mysqlDb.QueryRow("SELECT id FROM client WHERE name = ?", clientName).Scan(&id)
		if id <= 0 {
			fmt.Println("新增网关")
			//网关不存在，写入一条数据
			_, _ = m.mysqlDb.Exec("INSERT INTO client(name,ip,created_at,status) values(?,?,?,?)",
				clientName, ip, opsTime, 2)
		} else {
			//网关存在，更新数据
			fmt.Println("更新网关状态")
			_, _ = m.mysqlDb.Exec("UPDATE client set updated_at = ? , status = 2 where name = ?",
				opsTime, clientName)
		}
	}
}

func (m *Monitor) OnClosedWrapper(pre server.OnClosed) server.OnClosed {
	return func(ctx context.Context, client server.Client, err error) {
		pre(ctx, client, err)
		fmt.Println("网关已断开，更新网关状态")
		opsTime := time.Now().Format("2006-01-02 15:04:05")
		_, err = m.mysqlDb.Exec("UPDATE client set updated_at = ? , status = 1 where name = ?",
			opsTime, client.ClientOptions().ClientID)
		if err != nil {
			fmt.Println("网关已断开，更新网关失败")
		}
	}
}