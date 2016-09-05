package goblin

import (
    "net/http"
    "log"
)

type App struct {
    Router Router
    Controller Controller
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
