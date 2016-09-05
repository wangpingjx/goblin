package goblin

import (
    "net/http"
)

type Context struct {
    Request        *http.Request
    params         map[string]string
    http.ResponseWriter
}
