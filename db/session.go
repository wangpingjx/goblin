package db

import (
    "reflect"
    "log"
    "strings"
    "fmt"
)

type Session struct {
    db        *DB
    QB        *QueryBuilder    // SQL组装器

    Value     interface{}

    SQL           string
    SQLVars      []interface{}
    RowsAffected int64
    LastInsertId int64
}

type ModelStruct struct {
    Fields    []*Field
    TableName string
}

//TODO Close release the connection from pool
// func (s *Session) Close() {
//
// }

func (session *Session) Init() {
    session.QB      = &QueryBuilder{ db: session.db }
    session.SQL     = ""
    session.SQLVars = make([]interface{}, 0)
}

// 从 Value 中得到目标 Model 的 TableName、Fields 等信息
func (session *Session) GetModelStruct(value interface{}) ModelStruct {
    var modelStruct ModelStruct

    reflectValue := reflect.ValueOf(value)
    reflectType  := reflectValue.Type()

    var isSlice bool
    for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
        if reflectType.Kind() == reflect.Slice {
            isSlice = true
        }
        reflectType = reflectType.Elem()
    }

    for i := 0; i < reflectType.NumField(); i++ {
        f := reflectType.Field(i)
        if isSlice {
            reflectValue = reflect.ValueOf(reflect.New(reflectType).Interface())
        }
        field := &Field{ Name: ToColumnName(f.Name), Tag: f.Tag, Value: reflect.Indirect(reflectValue).FieldByName(f.Name) }
        modelStruct.Fields = append(modelStruct.Fields, field)
    }
    modelStruct.TableName = ToTableName(reflectType.Name())

    return modelStruct
}

func (session *Session) DB() sqlCommon{
    return session.db.db
}

func (session *Session) Query() error {
    var err error
    session.SQL = session.QB.ToSQL()
    log.Printf("session.SQL is: %v", session.SQL)
    if rows, err := session.DB().Query(session.SQL, session.SQLVars...); err == nil {
        session.scan(rows)
    }
    return err
}

func (session *Session) Where(query interface{}, args ...interface{}) *Session {
    session.QB.Where(query, args...)
    return session
}

func (session *Session) Find(value interface{}) error {
    modelStruct   := session.GetModelStruct(value)
    session.Value  = value
    session.QB.Table(modelStruct.TableName)

    return session.Query()
}

func (session *Session) First(value interface{}) error {
    modelStruct   := session.GetModelStruct(value)
    session.Value  = value
    session.QB.Table(modelStruct.TableName)
    session.QB.Limit(1)

    return session.Query()
}

func (session *Session) Last(value interface{}) error {
    modelStruct   := session.GetModelStruct(value)
    session.Value  = value
    session.QB.Table(modelStruct.TableName)
    session.QB.Order("id DESC").Limit(1)

    return session.Query()
}

func (session *Session) Order(order string) *Session {
    session.QB.Order(order)
    return session
}

func (session *Session) Limit(limit int) *Session {
    session.QB.Limit(limit)
    return session
}

func (session *Session) Offset(offset int) *Session {
    session.QB.Offset(offset)
    return session
}

func (session *Session) Join(joinOperator string, tableName string, condition string) *Session {
    session.QB.Join(joinOperator, tableName, condition)
    return session
}

func (session *Session) Group(column string) *Session {
    session.QB.Group(column)
    return session
}

func (session *Session) Having(condition string) *Session {
    session.QB.Having(condition)
    return session
}

func (session *Session) Select(str string) *Session {
    session.QB.Select(str)
    return session
}

func (session *Session) Create() (int64, error) {
    var (
        columns       []string
        placeholders  []string
        modelStruct   ModelStruct
    )
    modelStruct = session.GetModelStruct(session.Value)

    for _, field := range modelStruct.Fields {
        if field.Name == "id" {
            continue
        }
        columns         = append(columns, Quote(field.Name))
        placeholders    = append(placeholders, "?")
        session.SQLVars = append(session.SQLVars, field.Value.Interface())
    }
    session.SQL = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", Quote(modelStruct.TableName), strings.Join(columns, ","), strings.Join(placeholders, ","))
    log.Printf("Insert SQL is %v", session.SQL)

    if result, err := session.DB().Exec(session.SQL, session.SQLVars...); err == nil {
        session.RowsAffected, _ = result.RowsAffected()
        session.LastInsertId, _ = result.LastInsertId()

        return session.LastInsertId, err
    } else {
        return 0, err
    }
}

func (session *Session) Update(column string, attr interface{}) (int64, error) {
    var (
        updateSqls    []string
        modelStruct   ModelStruct
    )
    modelStruct = session.GetModelStruct(session.Value)

    // TODO 待扩展支持多个参数更新
    updateSqls = append(updateSqls, fmt.Sprintf("%v = %v", Quote(column), ToVars(attr)))

    if len(updateSqls) > 0 {
        session.SQL = fmt.Sprintf("UPDATE %v SET %v%v", Quote(modelStruct.TableName), strings.Join(updateSqls, ","), session.QB.buildWhereSQL())
    }
    log.Printf("Update SQL is %v", session.SQL)
    if result, err := session.DB().Exec(session.SQL, session.SQLVars...); err == nil {
        session.RowsAffected, _ = result.RowsAffected()

        return session.RowsAffected, err
    } else {
        return 0, err
    }
}

// Create or Update
func (session *Session) Save() (int64, error) {
    var (
        isNewRecord  bool
        updateSqls   []string
        modelStruct   ModelStruct
    )
    modelStruct = session.GetModelStruct(session.Value)

    for _, field := range modelStruct.Fields {
        if field.Name == "id" {
            if id, ok := field.Value.Interface().(int); ok {
                if id > 0 {
                    isNewRecord = false
                    session.QB.Where("id = ?", id)
                } else {
                    isNewRecord = true
                    break
                }
            }
        }
        updateSqls = append(updateSqls, fmt.Sprintf("%v = %v", Quote(field.Name), ToVars(field.Value.Interface())))
    }

    if isNewRecord {
        return session.Create()
    }
    if len(updateSqls) > 0 {
        session.SQL = fmt.Sprintf("UPDATE %v SET %v%v", Quote(modelStruct.TableName), strings.Join(updateSqls, ","), session.QB.buildWhereSQL())
    }
    log.Printf("Save SQL is %v", session.SQL)
    if result, err := session.DB().Exec(session.SQL, session.SQLVars...); err == nil {
        session.RowsAffected, _ = result.RowsAffected()

        return session.RowsAffected, err
    } else {
        return 0, err
    }
}

func (session *Session) Delete(value interface{}) (int64, error) {
    var (
        modelStruct ModelStruct
    )
    modelStruct = session.GetModelStruct(value)
    for _, field := range modelStruct.Fields {
        if field.Name == "id" {
            if id, ok := field.Value.Interface().(int); ok {
                if id > 0 {
                    session.QB.Where("id = ?", id)
                }
            }
            break
        }
    }
    session.SQL = fmt.Sprintf("DELETE FROM %v%v", Quote(modelStruct.TableName), session.QB.buildWhereSQL())
    log.Printf("Delete SQL is %v", session.SQL)

    if result, err := session.DB().Exec(session.SQL, session.SQLVars...); err == nil {
        session.RowsAffected, _ = result.RowsAffected()
        return session.RowsAffected, err
    } else {
        return 0, err
    }
}
