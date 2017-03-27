package db

import (
    "strings"
    "time"
)

func ToDBName(name string) string {
    return "TODO"
}

func ToColumnName(name string) string {
    return strings.ToLower(name)
}

func ToTableName(name string) string {
    return strings.ToLower(name)
}

func ToVars(value interface{}) interface{} {
    if value, ok := value.(string); ok {
        return "'" + value + "'"
    }
    if value, ok := value.(time.Time); ok {
        return "'" + value.Format("2006-01-02 15:04:05") + "'"
    }
    return value
}

func Quote(str string) string {
    return "`" + str + "`"
}
