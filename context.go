package goblin

import (
    "net/http"
    "encoding/json"
)

type Context struct {
    Params         map[string]string
    Request        *http.Request
    http.ResponseWriter
}

func (c *Context) ApplyJSON(obj interface{}) error{
    c.ResponseWriter.Header().Set("Content-Type", "application/json")

    var b []byte
    b, err := json.Marshal(obj)
    if err != nil {
        return err
    }
    c.ResponseWriter.Write(b)
    return nil
}

func (c *Context) ApplyString(content string) error{
    c.ResponseWriter.Header().Set("Content-Type", "text/html")
    c.ResponseWriter.Write([]byte(content))
    return nil
}
