package goblin

import (
    "net/http"
    "log"
)

var (
    GoblinApp *App
)
type App struct {
    Router Router
    Controller Controller
}

func init() {
    GoblinApp = New()
}

func New() *App {
    router   := NewRouter()
    app      := &App{Router: router}
    return app
}

func (app *App) Run(port string) {
    log.Println("=> Goblin server start at port: " + port)
    http.ListenAndServe(port, &app.Router)
}

func Get(pattern string, c ControllerInterface, actionName string) {
    GoblinApp.Router.Get(pattern, c, actionName)
    // return GoblinApp
}

func Post(pattern string, c ControllerInterface, actionName string) {
    GoblinApp.Router.Post(pattern, c, actionName)
    // return GoblinApp
}

func Put(pattern string, c ControllerInterface, actionName string) {
    GoblinApp.Router.Put(pattern, c, actionName)
    // return GoblinApp
}

func Delete(pattern string, c ControllerInterface, actionName string) {
    GoblinApp.Router.Delete(pattern, c, actionName)
    // return GoblinApp
}
