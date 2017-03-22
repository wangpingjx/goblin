package db

import (
    "reflect"
    "fmt"
    "strings"
    "log"
)

/* 数据表的抽象类，存储Model对应的表名&表字段信息 */
type Table struct {
    db          *DB
    name        string
    PrimaryKey  *Field
    Fields      []*Field
    modelType   reflect.Type
}

type Field struct {
    Name         string
    Tag          reflect.StructTag
    Value        reflect.Value
    IsPrimaryKey bool
}

// 缓存映射信息，提高效率
var tableCache = map[reflect.Type]*Table{}

func (t *Table) New(value interface{}) *Table {
    modelType := reflect.ValueOf(value).Type().Elem()
    if v := tableCache[modelType]; v == nil {
        t.modelType = modelType
        t.name      = strings.ToLower(t.modelType.Name())   // 放到Uitl中

        // 获取所有列
        for i := 0; i < t.modelType.NumField(); i++ {
            f := t.modelType.Field(i)
            field := &Field{ Name: strings.ToLower(f.Name), Tag: f.Tag}
            t.Fields = append(t.Fields, field)
        }
        tableCache[modelType] = t
    }
    return tableCache[modelType]
}

func (t *Table) TableName() string {
    return t.name
}

func (t *Table) ColumnTypeOf(field *Field) string {
    return field.Tag.Get("schema")
}

func (t *Table) ColumnIndexOf(field *Field) string {
    return field.Tag.Get("index")
}

// func (t *Table) Quote(str string) string {
//     return "`" + str + "`"
// }

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

// 列字段是否存在
func (t *Table) HasColumn(tableName string, columnName string) bool {
    var count int
    t.db.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", t.CurrentDatabase(), tableName, columnName).Scan(&count)
    return count > 0
}

// 索引是否存在
func (t *Table) HasIndex(tableName string, indexName string) bool {
    var count int
    t.db.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", t.CurrentDatabase(), tableName, indexName).Scan(&count)
    return count > 0
}


func (t *Table) Migrate() {
    if !t.HasTable(t.name) {
        t.CreateTable()
    } else {
         log.Println("TODO: alter table")
         t.UpdateTable()
    }
}

// 创建表
func (t *Table) CreateTable() {
    var tags    []string
    var indexes []string
    for _, field := range t.Fields {
        tags  = append(tags, Quote(field.Name) + " " + t.ColumnTypeOf(field))
        if index := t.ColumnIndexOf(field); index != "" {
            indexes = append(indexes, index + " (" + Quote(field.Name) + ")")
        }
    }
    if len(indexes) == 0 {
        indexes = append(indexes, "PRIMARY KEY (`id`)")
    }
    tags = append(tags, indexes...)
    additionSQL := "ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"

    sql := fmt.Sprintf("CREATE TABLE %v (%v) %s", Quote(t.TableName()), strings.Join(tags, ","), additionSQL)
    log.Println("=> sql: " + sql)
    t.db.db.Exec(sql)
}

// 更新表
func (t *Table) UpdateTable() {
    for _, field := range t.Fields {
        if !t.HasColumn(t.TableName(), field.Name) {
            sql := fmt.Sprintf("ALTER TABLE %v ADD %v %v", Quote(t.TableName()), Quote(field.Name), t.ColumnTypeOf(field))
            t.db.db.Exec(sql)
        }
    }
}

func (t *Table) dropTable() {
    if t.HasTable(t.TableName()) {
        t.db.db.Exec(fmt.Sprintf("DROP TABLE %v", Quote(t.TableName())))
    }
}

func (t *Table) modifyColumn(column string ,tag string) {
    t.db.db.Exec(fmt.Sprintf("ALTER TABLE %v MODIFY %v %v", Quote(t.TableName()), Quote(column), tag))
}

func (t *Table) dropColumn(column string) {
    t.db.db.Exec(fmt.Sprintf("ALTER TABLE %v DROP COLUMN %v", Quote(t.TableName()), Quote(column)))
}

func (t *Table) addIndex(unique bool, indexName string, column ...string) {
    if t.HasIndex(t.TableName(), indexName) {
        return
    }
    sqlStr := "CREATE INDEX"
    if unique {
        sqlStr = "CREATE UNIQUE INDEX"
    }
    t.db.db.Exec(fmt.Sprintf("%v %v ON %v(%v)", sqlStr, indexName, Quote(t.TableName()), strings.Join(column, ",")))
}

func (t *Table) removeIndex(indexName string) {
    t.db.db.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, Quote(t.TableName())))
}
