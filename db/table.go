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
    Name  string
    Tag   reflect.StructTag
}

// TODO 需要缓存起来
func (t *Table) New(value interface{}) *Table {
    t.modelType = reflect.ValueOf(value).Type().Elem()
    t.name      = strings.ToLower(t.modelType.Name())

    // 获取所有列
    for i := 0; i < t.modelType.NumField(); i++ {
        f := t.modelType.Field(i)
        field := &Field{ Name: strings.ToLower(f.Name), Tag: f.Tag}
        t.Fields = append(t.Fields, field)
    }
    return t
}

func (t *Table) TableName() string {
    return t.name
}

func (t *Table) ColumnTypeOf(field *Field) string {
    return field.Tag.Get("db")
}

func (t *Table) ColumnIndexOf(field *Field) string {
    return field.Tag.Get("index")
}

func (t *Table) Quote(str string) string {
    return "`" + str + "`"
}

func (t *Table) Migrate() bool {
    if !t.HasTable(t.name) {
        t.CreateTable()
    } else {
         log.Println("TODO: alter table")
         // 遍历字段、字段是否存在/变化、执行变化
         // 变更索引
    }
    return true
}

// 当前数据库名称
func (t *Table) CurrentDatabase() (name string){
   t.db.db.QueryRow("SELECT DATABASE()").Scan(&name)
   return
}

// 表是否存在
func (t *Table) HasTable(name string) bool {
    var count int
    t.db.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", t.CurrentDatabase(), name).Scan(&count)
    return count > 0
}

// 建表
func (t *Table) CreateTable() bool{
    var tags    []string
    var indexes []string
    for _, field := range t.Fields {
        tags  = append(tags, t.Quote(field.Name) + " " + t.ColumnTypeOf(field))
        if index := t.ColumnIndexOf(field); index != "" {
            indexes = append(indexes, index + " (" + t.Quote(field.Name) + ")")
        }
    }
    if len(indexes) == 0 {
        indexes = append(indexes, "PRIMARY KEY (`id`)")
    }
    tags = append(tags, indexes...)
    additionSQL := "ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"

    sql := fmt.Sprintf("CREATE TABLE %v (%v) %s", t.Quote(t.TableName()), strings.Join(tags, ","), additionSQL)
    log.Println("=> sql: " + sql)
    t.db.db.Exec(sql)

    return true
}
