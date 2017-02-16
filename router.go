package goblin

import (
    "net/http"
    "log"
    "reflect"
)

// 路由器, 负责存储路由规则
type Router struct {
    routes []Route
}

// 路由协议, 负责将 URL 映射到正确的 Handle
type Route struct {
    method         string
    pattern        string
    controller     reflect.Type
    actionName     string
}

// 路由分发方法
func (m *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path   := r.URL.Path
    method := r.Method

    log.Println("=> " + method + path)

    context := Context {
        Request:        r,
        params:         nil,
        ResponseWriter: w,
    }

    for _, route := range m.routes {
        isMatch := route.route(method, path)

        // 匹配成功
        if isMatch {
            vc := reflect.New(route.controller)
            controller, _ := vc.Interface().(ControllerInterface)

            controller.Init(context, route.controller.Name(), route.actionName)
            controller.Before()

            // invoke handler
            method := vc.MethodByName(route.actionName)
            method.Call([]reflect.Value{})

            controller.After()
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

func (p *Router) AddRoute(method string, pattern string, t reflect.Type, actionName string) {
    p.routes = append(p.routes, Route{method, pattern, t, actionName})
}

func (p *Router) Get(pattern string, c ControllerInterface, actionName string) {
    t := reflect.Indirect(reflect.ValueOf(c)).Type()
    p.AddRoute("GET", pattern, t, actionName)
    p.AddRoute("HEAD", pattern, t, actionName)
}

func (p *Router) Post(pattern string, c ControllerInterface, actionName string) {
    t := reflect.Indirect(reflect.ValueOf(c)).Type()
    p.AddRoute("POST", pattern, t, actionName)
}

func (p *Router) Put(pattern string, c ControllerInterface, actionName string) {
    t := reflect.Indirect(reflect.ValueOf(c)).Type()
    p.AddRoute("PUT", pattern, t, actionName)
}

func (p *Router) Delete(pattern string, c ControllerInterface, actionName string) {
    t := reflect.Indirect(reflect.ValueOf(c)).Type()
    p.AddRoute("DELETE", pattern, t, actionName)
}

func NewRouter() Router {
    router := Router{routes: make([]Route, 5, 5)}
    return router
}
