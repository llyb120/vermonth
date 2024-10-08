package vermouth

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"sort"
	"strings"
)

type aopItem struct {
	Expression *regexp.Regexp
	Fn         func(*Context)
	Order      int
	// controllerCaller interface{}
}

type Context struct {
	// 调用方法
	Fn func()
	// 参数表
	Arguments     []interface{}
	ArgumentNames []string
	// 返回值
	Result []interface{}
	// 上下文环境
	GinContext *gin.Context

	// // 控制器信息
	// ControllerInformation *ControllerDefinition
	// // 方法信息
	// MethodInformation *RequestMapping

	ControllerInformation *ControllerInformation

	// 是否自动返回前端
	AutoReturn bool
}

type ControllerInformation struct {
	Path        string
	Transaction bool

	Attributes map[string]string
}

func NewControllerInformation() *ControllerInformation {
	return &ControllerInformation{
		// Attributes: make(map[string]string),
	}
}

func (aopContext *Context) Call() {
	aopContext.Fn()
}

func newAopContext(argumentsLength int) *Context {
	return &Context{
		Arguments:     make([]interface{}, argumentsLength),
		ArgumentNames: make([]string, argumentsLength),
		AutoReturn:    true,
	}
}

var aopItems []*aopItem = make([]*aopItem, 0)

func RegisterAop(exp string, order int, fn func(*Context)) {
	// 替换.为\.
	exp = strings.Replace(exp, "**", "(.+)", -1)
	exp = strings.Replace(exp, "*", "[^/]{0,}", -1)
	exp = "^" + exp + "$"
	reg, err := regexp.Compile(exp)
	if err != nil {
		return
	}
	aopItems = append(aopItems, &aopItem{Expression: reg, Fn: fn, Order: order})
	sort.Slice(aopItems, func(a, b int) bool {
		return aopItems[a].Order > aopItems[b].Order
	})
}

// func main(){
// 	RegisterAop("*.*", func (aopContext *Context)  {
// 		aopContext.Arguments[0] = reflect.ValueOf(1)
// 		aopContext.Fn()
// 	})
// }

// type ControllerContext struct {
// 	caller *Caller
// }

// func test(){
// 	RegisterAop("*.*", func (caller *Caller)  {
// 			caller.call()
// 	})
// }
