// Package scope 用户 Scope Token
package scope

import (
	"github.com/gin-gonic/gin"
)

// RegisterScopeRouter scope token 路由注册
func RegisterScopeRouter(r gin.IRouter) {
	scopeGroup := r.Group("scopes")
	scopeGroup.GET("", scopeList)
	scopeGroup.POST("token", scopeToken)
}
