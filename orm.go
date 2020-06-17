package orm

import (
	"calliphlox/dialect"
	"calliphlox/log"
	"calliphlox/session"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// 负责与数据库交互前的准备工作（连接、测试数据库），和交互后的收尾工作（关闭连接）
type Engine struct {
	db *sql.DB

	dialect dialect.Dialect
}

// 创建 Engine 实例时，获取 diver 对应的 dialect
func NewEngine(driver, source string) (e *Engine, err error) {
	// 连接数据库
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// ping 检查数据库是否正常连接
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// 判断方言是否存在
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Fount", driver)
		return
	}

	e = &Engine{
		db: db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

// 关闭数据库
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// 创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}