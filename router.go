package goblin

import (
    "net/http"
    "log"
)

// 路由器, 负责存储路由规则
type Router struct {
    routes []Route
}

// 路由协议, 负责将 URL 映射到正确的 Handle
type Route struct {
    method  string
    pattern string
    handler Handler
}

// 与路由匹配的目标 Handler
type Handler func(w http.ResponseWriter, r *http.Request)

// 路由分发方法
func (m *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path   := r.URL.Path
    method := r.Method

    log.Println("=> " + method + path)

    for _, route := range m.routes {
        isMatch := route.route(method, path)

        // 匹配成功, 调起 Handler
        if isMatch {
            route.handler(w, r)
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

func (p *Router) AddRoute(method string, pattern string, handler Handler) {
    p.routes = append(p.routes, Route{method, pattern, handler})
}

func (p *Router) Get(pattern string, handler Handler) {
    p.AddRoute("GET", pattern, handler)
    p.AddRoute("HEAD", pattern, handler)
}

func (p *Router) Post(pattern string, handler Handler) {
    p.AddRoute("POST", pattern, handler)
}

func (p *Router) Put(pattern string, handler Handler) {
    p.AddRoute("PUT", pattern, handler)
}

func (p *Router) Delete(pattern string, handler Handler) {
    p.AddRoute("DELETE", pattern, handler)
}

func NewRouter() Router {
    router := Router{routes: make([]Route, 5, 5)}
    return router
}
