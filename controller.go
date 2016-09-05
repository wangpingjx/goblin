package goblin

import (
    "log"
)

type Controller struct {
    Ctx     Context   // 请求上下文
    Data    map[string]interface{}  // 输出参数
    TplName string    // 模板
    TplExt  string    // 扩展名称, Eg.  "tpl"、"xml"

}

type ControllerInterface interface {
	Init(ct Context)
    Before()
    After()
}

// 初始化，设置默认值
func (c *Controller) Init(ctx Context) {
    c.Ctx            = ctx
    c.Data           = make(map[string]interface{})
    c.TplName        = ""
    c.TplExt         = "gtpl"
}

func (c *Controller) Before() {
    log.Println("=> in Controller#Before")
}

func (c *Controller) After() {
    log.Println("=> in Controller#After")
}
