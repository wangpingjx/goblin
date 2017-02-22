package goblin

import (
    "net/http"
    "log"
    "fmt"
    "reflect"
    "regexp"
    // "io/ioutil"
)

// 路由器, 负责存储路由规则
type Router struct {
    routes []*Route
}

// 路由协议, 负责将 URL 映射到正确的 Handle
type Route struct {
    method          string
    pattern         string
    controllerType  reflect.Type
    actionName      string
    regex           *regexp.Regexp
}

// 路由分发方法
func (p *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path   := r.URL.Path
    method := r.Method

    // 逐个匹配
    for _, route := range p.routes {
        match, vals := route.route(method, path)
        // 如果匹配成功
        if match {
            log.Println("=> mached route is " + string(route.pattern))

            rc := reflect.New(route.controllerType)
            controller, _ := rc.Interface().(ControllerInterface)   // TODO 没看懂

            // 把请求上下文放入 Context
            ctx := &Context{ Params: vals, Request: r, ResponseWriter: w }

            // 获取请求参数
            r.ParseForm()
            r.ParseMultipartForm(2 * 1024 * 1024)

            for k, v := range r.Form {
                ctx.Params[k] = v[0]
            }

            // 初始化 controller
            controller.Init(ctx, route.controllerType.Name(), route.actionName)

            // Before()
            controller.Before()

            // invoke handler
            method := rc.MethodByName(route.actionName)
            method.Call([]reflect.Value{})

            // After()
            controller.After()

            return
        }
    }
    log.Println("=> 404 not found ")
    p.NotFound(w,r)
}

// 借鉴martini，优化route匹配方法
func (r *Route) route(method string, path string) (bool, map[string]string) {
    if method != r.method {
        return false, nil
    }
    matches := r.regex.FindStringSubmatch(path)

    if len(matches) > 0 && matches[0] == path {
        params := make(map[string]string)
        for i, name := range r.regex.SubexpNames() {
            if len(name) > 0 {
                params[name] = matches[i]    // SubexpNames()返回结果第0个元素永远是空字符串
            }
        }
        return true, params
    }
    return false, nil
}

// 生成regexp对象
// Eg:   pattern = /books/:id/users/:user_id
// Then: route.regex = /books/(?P<id>[^/#?]+)/users/(?P<user_id>[^/#?]+)\/?
var routeReg1 = regexp.MustCompile(`:[^/#?()\.\\]+`)

func newRoute(method string, pattern string, t reflect.Type, action string) *Route {
    route  := Route{method, pattern, t, action, nil}
    pattern = routeReg1.ReplaceAllStringFunc(pattern, func(m string) string {
        return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
    })
    pattern += `\/?`
    route.regex = regexp.MustCompile(pattern)
    return &route
}

func (p *Router) AddRoute(method string, pattern string, t reflect.Type, action string) {
    p.routes = append(p.routes, newRoute(method, pattern, t, action))
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
    router := Router{routes: make([]*Route, 0, 0)}
    return router
}

func (p *Router) NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
    fmt.Fprint(w, "404 Not Found.")
}
