package goblin

import (
    "net/http"
)

type Context struct {
    Request        *http.Request
    Params         map[string]string
    http.ResponseWriter
}
