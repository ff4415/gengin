package svc

import (
    "time"

    "github.com/gin-gonic/gin"
)

// ServiceContext 实现 context.Context 接口
type {{.svcName}}ServiceContext struct {
    GinContext *gin.Context
}

func New{{.svcName}}ServiceContext(c *gin.Context) *{{.svcName}}ServiceContext {
    return &{{.svcName}}ServiceContext{
        GinContext:         c,
    }
}

func (c *{{.svcName}}ServiceContext) Deadline() (deadline time.Time, ok bool) {
    return c.GinContext.Deadline()
}

func (c *{{.svcName}}ServiceContext) Done() <-chan struct{} {
    return c.GinContext.Done()
}

func (c *{{.svcName}}ServiceContext) Err() error {
    return c.GinContext.Err()
}

func (c *{{.svcName}}ServiceContext) Value(key any) any {
    return c.GinContext.Value(key)
}
