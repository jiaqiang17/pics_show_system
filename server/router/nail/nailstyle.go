package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type NailStyleRouter struct{}

// InitNailStyleRouter 初始化 美甲款式 路由信息
func (s *NailStyleRouter) InitNailStyleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	nailStyleRouter := Router.Group("nailStyle").Use(middleware.OperationRecord())
	nailStyleRouterWithoutRecord := Router.Group("nailStyle")
	nailStyleRouterWithoutAuth := PublicRouter.Group("nailStyle")
	{
		nailStyleRouter.POST("createNailStyle", nailStyleApi.CreateNailStyle)             // 新建美甲款式
		nailStyleRouter.DELETE("deleteNailStyle", nailStyleApi.DeleteNailStyle)           // 删除美甲款式
		nailStyleRouter.DELETE("deleteNailStyleByIds", nailStyleApi.DeleteNailStyleByIds) // 批量删除美甲款式
		nailStyleRouter.PUT("updateNailStyle", nailStyleApi.UpdateNailStyle)              // 更新美甲款式
	}
	{
		nailStyleRouterWithoutRecord.GET("findNailStyle", nailStyleApi.FindNailStyle) // 根据ID获取美甲款式（需要登录）
	}
	{
		nailStyleRouterWithoutAuth.GET("getNailStyleList", nailStyleApi.GetNailStyleList)             // 获取美甲款式列表（公开访问）
		nailStyleRouterWithoutAuth.GET("getNailStyleDataSource", nailStyleApi.GetNailStyleDataSource) // 获取美甲款式数据源（公开访问）
		nailStyleRouterWithoutAuth.GET("getNailStylePublic", nailStyleApi.GetNailStylePublic)         // 美甲款式开放接口
	}
}
