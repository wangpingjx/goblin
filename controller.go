package goblin

import (
    "log"
)

type Controller struct {
    Ctx            *Context                // 请求上下文
    controllerName string                  // 控制器名称
    actionName     string                  // Action
}

type ControllerInterface interface {
	Init(ctx *Context, controllerName string, actionName string)
    Before()
    After()
    // Render()
}

// 初始化，设置默认值
func (c *Controller) Init(ctx *Context, controllerName string, actionName string) {
    c.Ctx            = ctx
    c.controllerName = controllerName
    c.actionName     = actionName
}

func (c *Controller) Before() {
    log.Println("=> in Controller#Before")
}

func (c *Controller) After() {
    log.Println("=> in Controller#After")
}

// TODO 放弃 View 模块
func (c *Controller) Render() {
    log.Println("=> should render template")
}

func (c *Controller) RenderJSON(obj interface{}) {
    log.Println("=> should render json")
    c.Ctx.ApplyJSON(obj)
}

func (c *Controller) RenderText(content string) {
    log.Println("=> should render plain text")
    c.Ctx.ApplyString(content)
}
