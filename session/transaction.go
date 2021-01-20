package session

import "calliphlox/log"

// 事务通过调用 db.Begin() 得到 *sql.Tx 对象，使用 tx.Exec() 执行一系列操作，
// 如果发生错误，通过 tx.Rollback() 回滚，如果没有发生错误，则通过 tx.Commit() 提交
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}