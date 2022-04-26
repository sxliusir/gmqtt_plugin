package monitor

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"

	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/server"
)


var _ server.Plugin = (*Monitor)(nil)

const Name = "monitor"

func init() {
	server.RegisterPlugin(Name, New)
	config.RegisterDefaultPluginConfig(Name, &DefaultConfig)
}
// New 是本插件的构造函数
func New(config config.Config) (server.Plugin, error) {
	cfg := config.Plugins[Name].(*Config)
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db, cfg.Charset)
	open, err := sql.Open("mysql", dbDSN)
	if err != nil {
		panic("数据源配置不正确: " + err.Error())
	}
	return &Monitor{
		mysqlDb: open,
		mysqlDbErr: err,
	}, nil
}

var log *zap.Logger

// Monitor 实现Plugin接口的结构体。
type Monitor struct {
	mysqlDb *sql.DB
	mysqlDbErr error
}

// Load 由Gmqtt按插件的导入顺序，依次执行。
// Load主要的作用就是把server.Server接口传递给插件。
func (m *Monitor) Load(service server.Server) error {
	log = server.LoggerWithField(zap.String("plugin", Name))
	return nil
}

// Unload 当broker退出时调用，可以做一些清理操作。
func (m *Monitor) Unload() error {
	err := m.mysqlDb.Close()
	return err
}

func (m *Monitor) Name() string {
	return Name
}