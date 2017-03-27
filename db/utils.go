package db

import (
    "strings"
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

func Quote(str string) string {
    return "`" + str + "`"
}
