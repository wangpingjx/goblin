package db

import (
    "reflect"
    "fmt"
    "strings"
    "log"
)

type Table struct {
    db          *DB
    name        string
    PrimaryKey  *Field
    Fields      []*Field
    modelType   reflect.Type
}

type Field struct {
    Name         string
}

// TODO 需要缓存起来
func (t *Table) New(value interface{}) *Table {
    t.modelType = reflect.ValueOf(value).Type().Elem()
    t.name      = t.modelType.Name()

    // 获取所有列
    for i := 0; i < t.modelType.NumField(); i++ {
        field := &Field{ Name: t.modelType.Field(i).Name }
        t.Fields = append(t.Fields, field)
    }
    return t
}

func (t *Table) Migrate() bool {
    if !t.HasTable(t.name) {
        t.CreateTable()
    } else {
         log.Println("TODO: alter table")
    }
    return true
}

func (t *Table) TableName() string {
    return t.name
}

func (t *Table) ColumnTypeOf(field *Field) string {
    return "int(11) NOT NULL"
}

// 当前数据库名称
func (t *Table) CurrentDatabase() (name string){
   t.db.db.QueryRow("SELECT DATABASE()").Scan(&name)
   return
}

func (t *Table) HasTable(name string) bool {
    var count int
    t.db.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", t.CurrentDatabase(), name).Scan(&count)
    return count > 0
}

func (t *Table) CreateTable() bool{
    var tags []string
    for _, field := range t.Fields {
        tags = append(tags, t.Quote(field.Name) + " " + t.ColumnTypeOf(field))
    }
    pkSQL       := ", PRIMARY KEY (`id`)"
    additionSQL := "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"
    sql         := fmt.Sprintf("CREATE TABLE %v (%v %v) %s", t.Quote(t.TableName()), strings.Join(tags, ","), pkSQL, additionSQL)
    log.Println("=> sql: " + sql)
    t.db.db.Exec(sql)

    // TODO 创建索引
    return true
}

func (t *Table) Quote(str string) string {
    return "`" + str + "`"
}
