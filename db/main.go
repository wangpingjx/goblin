/**
 * 实现方式借鉴 gorm、xorm
 */
package db

import (
    "database/sql"
    "errors"
    "log"
)

/* 如果某个对象实现了此接口的所有方法，那么这个对象就实现了这个接口 */
type sqlCommon interface {
    Exec(query string, args ...interface{}) (sql.Result, error)
    Prepare(query string) (*sql.Stmt, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}

/* 数据库连接信息（尽量保持长连接、避免频繁open&close）*/
type DB struct {
    db       sqlCommon       // *sql.DB
    qb       *QueryBuilder   // SQL构造器
    parent   *DB
}

/**
 * eg: db.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
 */
func Open(driver string, value string) (*DB, error) {
    var sqlDB *sql.DB

    sqlDB, err := sql.Open(driver, value)

    db       := &DB { db: sqlDB}
    db.qb     = &QueryBuilder{ db: db }
    db.parent = db    // TODO 作用暂时不是非常明确

    return db, err
}

func (s *DB) Close() error {
    if db, ok := s.parent.db.(*sql.DB); ok {
        return db.Close()
    }
    return errors.New("close db failed.")
}

/* 返回 DB 中的 *sql.DB */
func (s *DB) DB() *sql.DB {
	return s.db.(*sql.DB)
}

func (s *DB) Table(name string) *DB {
   return s.qb.Table(name).db
}

/************************
 *     Schema 相关       *
 ************************/
func (s *DB) NewTable(value interface{}) *Table {
    t :=  &Table{ db: s }
    return t.New(value)
}

func (s *DB) Migrate(values ...interface{}) bool {
    for _, value := range values {
        s.NewTable(value).Migrate()
    }
    return true
}

func (s *DB) DropTable(values ...interface{}) *DB {
    for _, value := range values {
        s.NewTable(value).dropTable()
    }
    return s
}

func (s *DB) ModifyColumn(value interface{}, column string, tag string) *DB {
    s.NewTable(value).modifyColumn(column, tag)
    return s
}

func (s *DB) DropColumn(value interface{}, column string) *DB{
    s.NewTable(value).dropColumn(column)
    return s
}

func (s *DB) AddIndex(value interface{}, indexName string, column ...string) *DB {
    s.NewTable(value).addIndex(false, indexName, column...)
    return s
}

func (s *DB) AddUniqueIndex(value interface{}, indexName string , column ...string) *DB {
    s.NewTable(value).addIndex(true, indexName, column...)
    return s
}

func (s *DB) RemoveIndex(value interface{}, indexName string) *DB {
    s.NewTable(value).removeIndex(indexName)
    return s
}

/************************
 *       查询构造器       *
 ************************/
func (s *DB) NewSession(value interface{}) *Session {
    session :=  &Session{ db: s}
    return session.New(value)
}

func (s *DB) Create(value interface{}) *DB {
    return s.NewSession(value).Create().db
}

func (s *DB) First(value interface{})  (*sql.Rows, error) {
    return s.NewSession(value).qb.Limit(1).Query()
}

func (s *DB) Query(query string) (*sql.Rows, error) {
    return s.db.Query(query)
}




func (s *DB) Select(selects string) *DB {
    return s.qb.Select(selects).db
}

func (s *DB) Where(query interface{}, args ...interface{}) *DB {
    return s.qb.Where(query, args...).db
}

func (s *DB) Limit(limit int) *DB {
    return s.qb.Limit(limit).db
}

func (s *DB) Find() (*sql.Rows, error) {
    sql := s.qb.buildSelect()
    log.Println("=> sql: " + sql)
    return s.Query(sql)
}

/************************
 *     对象关系映射       *
 ************************/



// func (s *DB) ScanRows(rows *sql.Rows, fields []string) error {
//     var ignored interface{}
//     var columns, err = rows.Columns()
//
//     for index, column := range columns {
//         values[index] = &ignored
//
//         for findex, field := range fields {
//
//         }
//
//     }
// }
