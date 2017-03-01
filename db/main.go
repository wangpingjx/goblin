/**
 * 本文件实现方式借鉴 gorm
 */
package db

import (
    "database/sql"
    // "log"
)
// 数据库连接信息（尽量保持长连接、避免频繁open&close）
type DB struct {
    db *sql.DB
}

// eg: db.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
func Open(driver string, value string) (*DB, error) {
    var sqlDB *sql.DB

    sqlDB, err := sql.Open(driver, value)
    db := &DB { db: sqlDB }

    return db, err
}

func (s *DB) Close() error {
    return s.db.Close()
}

// 获取数据库连接
func (s *DB) DB() *sql.DB {
	return s.db
}

func (s *DB) Query(sql string) (rows *sql.Rows, err error) {
    rows, err = s.db.Query(sql)
    return rows, err
}
