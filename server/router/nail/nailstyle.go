package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type NailStyleRouter struct{}

// InitNailStyleRouter 初始化美甲款式路由
func (s *NailStyleRouter) InitNailStyleRouter(router *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	protected := router.Group("nailStyle").Use(middleware.OperationRecord())
	protectedNoRecord := router.Group("nailStyle")
	public := publicRouter.Group("nailStyle")

	{
		protected.POST("createNailStyle", nailStyleApi.CreateNailStyle)   // 新建美甲款式
		protected.DELETE("deleteNailStyle", nailStyleApi.DeleteNailStyle) // 删除美甲款式
		protected.DELETE("deleteNailStyleByIds", nailStyleApi.DeleteNailStyleByIds)
		protected.PUT("updateNailStyle", nailStyleApi.UpdateNailStyle)           // 更新美甲款式
		protected.POST("batchUpdateTags", nailStyleApi.BatchUpdateNailStyleTags) // 批量更新标签
	}

	{
		protectedNoRecord.GET("findNailStyle", nailStyleApi.FindNailStyle) // 根据ID获取美甲款式
	}

	{
		public.GET("getNailStyleList", nailStyleApi.GetNailStyleList)             // 获取美甲款式列表
		public.GET("getNailStyleDataSource", nailStyleApi.GetNailStyleDataSource) // 获取数据源
		public.GET("getNailStylePublic", nailStyleApi.GetNailStylePublic)         // 开放接口
	}
}
