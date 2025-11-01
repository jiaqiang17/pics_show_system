package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type NailTagRouter struct{}

// InitNailTagRouter 初始化 美甲标签 路由信息
func (s *NailTagRouter) InitNailTagRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	nailTagRouter := Router.Group("nailTag").Use(middleware.OperationRecord())
	nailTagRouterWithoutRecord := Router.Group("nailTag")
	nailTagRouterWithoutAuth := PublicRouter.Group("nailTag")
	{
		nailTagRouter.POST("createNailTag", nailTagApi.CreateNailTag)             // 新建美甲标签
		nailTagRouter.DELETE("deleteNailTag", nailTagApi.DeleteNailTag)           // 删除美甲标签
		nailTagRouter.DELETE("deleteNailTagByIds", nailTagApi.DeleteNailTagByIds) // 批量删除美甲标签
		nailTagRouter.PUT("updateNailTag", nailTagApi.UpdateNailTag)              // 更新美甲标签
	}
	{
		nailTagRouterWithoutRecord.GET("findNailTag", nailTagApi.FindNailTag) // 根据ID获取美甲标签（需要登录）
	}
	{
		nailTagRouterWithoutAuth.GET("getNailTagList", nailTagApi.GetNailTagList)     // 获取美甲标签列表（公开访问）
		nailTagRouterWithoutAuth.GET("getNailTagPublic", nailTagApi.GetNailTagPublic) // 美甲标签开放接口
	}
}
