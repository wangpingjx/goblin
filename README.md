# goblin
micro web framework for go


## To do list
- [x] router
- []  controller
- []  view

## sample
```
package main

import(
    "goblin"
    "log"
)

type MainController struct{
    goblin.Controller
}

func (this *MainController) Login() {
    this.Data["Username"] = "astaxie"
    this.Data["Email"]    = "astaxie@gmail.com"
    log.Println("=> in MainController#Login")
}

func main() {
    app := goblin.New()
    app.Router.Get("/login", &MainController{}, "Login")
    app.Run(":9090")
}
```
