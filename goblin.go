package goblin

import (
    "net/http"
    "log"
    "fmt"
)
type App struct {
}

// 路由器, 负责存储路由规则。需要实现的方法「Add」「ServeHTTP」
type MyMux struct {
    routes []Route
}

// 路由协议，负责将 URL 映射到正确的 Handle。
type Route struct {
    method string
    pattern string
    // params  map[int]string  // 匹配后的结果，controller、action、params...
    handler Handler
}

type Handler func(w http.ResponseWriter, r *http.Request)

// 路由分发方法
func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path   := r.URL.Path
    method := r.Method

    log.Println("=> " + method + " path")

    for _, route := range m.routes {
        isMatch := route.route(method, path)
        // 匹配成功
        if isMatch {
            sayhello(w, r)
            // m.handle(route.handler) // 调起 Handler
            return
        }
    }
    return
}

func (r *Route) route(method string, path string) (isMatch bool) {
    if (method == r.method && path  == r.pattern) {
        isMatch = true
    } else {
        isMatch = false
    }
    return
}

// func (p *MyMux) handle(routHandler Handler) (err error){
//     sayhello(w http.ResponseWriter, r *http.Request)
//     // err := routHandler(w http.ResponseWriter, r *http.Request)
//     // return
// }

func sayhello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello wangping!")
}

func (p *MyMux) AddRoute(method string, pattern string, handler Handler) {
    p.routes = append(p.routes, Route{method, pattern, handler})
}

func (p *MyMux) Get(pattern string, handler Handler) {
    p.AddRoute("GET", pattern, handler)
    p.AddRoute("HEAD", pattern, handler)
}

func (p *MyMux) Post(pattern string, handler Handler) {
    p.AddRoute("POST", pattern, handler)
}

func (p *MyMux) Put(pattern string, handler Handler) {
    p.AddRoute("PUT", pattern, handler)
}

func (p *MyMux) Delete(pattern string, handler Handler) {
    p.AddRoute("DELETE", pattern, handler)
}

func New() *App {
    app := &App { }
    return app
}

func (app *App) Run(port string) {
    mux := &MyMux{}
    mux.Get("/", sayhello)
    log.Println("=> Goblin server start at port: " + port)
    http.ListenAndServe(port, mux)
}
