package goblin

import (
    "log"
    "strings"
)

type Controller struct {
    Ctx            Context                 // 请求上下文
    Data           map[string]interface{}  // 输出参数
    controllerName string                  // 控制器名称
    actionName     string                  // Action
    TplName        string                  // 模板
    TplExt         string                  // 扩展名称, Eg.  "tpl"、"xml"

}

type ControllerInterface interface {
	Init(ctx Context, controllerName string, actionName string)
    Before()
    After()
}

// 初始化，设置默认值
func (c *Controller) Init(ctx Context, controllerName string, actionName string) {
    c.Ctx            = ctx
    c.Data           = make(map[string]interface{})
    c.TplName        = ""
    c.TplExt         = "gtpl"
    c.controllerName = controllerName
    c.actionName     = actionName
}

func (c *Controller) Before() {
    log.Println("=> in Controller#Before")
}

func (c *Controller) After() {
    log.Println("=> in Controller#After")
}

func (c *Controller) RenderText(text string) {
    c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
    c.Ctx.ResponseWriter.Write([]byte(text))
}

func (c *Controller) Render(tplname string) {
    if c.TplName == "" {
        c.TplName = strings.Split(strings.ToLower(c.controllerName), "controller")[0] + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
    }
    BuildTemplate("/Users/vikki/iwork/Mine/Go/src/sample/views", c.TplName)
    c.RenderText(c.TplName)
}
