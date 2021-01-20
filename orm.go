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

type TxFunc func(*session.Session) (interface{}, error)

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

// 将所有的操作放到一个回调函数中，作为入参传递给 engine.Transaction()，
// 发生任何错误，自动回滚，如果没有错误发生，则提交
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			// 回滚不覆盖 err, 回滚就是因为有业务的报错,
			// 所以不应该被这条语句覆盖掉业务的 err, 业务 err比回滚失败的 err 更重要
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			// commit 失败需要再次回滚
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()

			err = s.Commit()
		}
	}()

	return f(s)
}