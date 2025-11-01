package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
	nailReq "github.com/flipped-aurora/gin-vue-admin/server/model/nail/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NailStyleApi struct{}

// CreateNailStyle 创建美甲款式
// @Tags NailStyle
// @Summary 创建美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailStyle true "创建美甲款式"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /nailStyle/createNailStyle [post]
func (api *NailStyleApi) CreateNailStyle(c *gin.Context) {
	ctx := c.Request.Context()

	var nailStyle nail.NailStyle
	if err := c.ShouldBindJSON(&nailStyle); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := nailStyleService.CreateNailStyle(ctx, &nailStyle); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteNailStyle 删除美甲款式
// @Tags NailStyle
// @Summary 删除美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailStyle true "删除美甲款式"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /nailStyle/deleteNailStyle [delete]
func (api *NailStyleApi) DeleteNailStyle(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("ID")
	if err := nailStyleService.DeleteNailStyle(ctx, id); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteNailStyleByIds 批量删除美甲款式
// @Tags NailStyle
// @Summary 批量删除美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /nailStyle/deleteNailStyleByIds [delete]
func (api *NailStyleApi) DeleteNailStyleByIds(c *gin.Context) {
	ctx := c.Request.Context()

	ids := c.QueryArray("IDs[]")
	if err := nailStyleService.DeleteNailStyleByIds(ctx, ids); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// BatchUpdateNailStyleTags 批量更新款式标签
// @Tags NailStyle
// @Summary 批量为美甲款式添加或移除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nailReq.NailStyleBatchTagUpdateReq true "批量更新标签请求体"
// @Success 200 {object} response.Response{data=nailReq.NailStyleBatchTagUpdateResult,msg=string} "批量更新标签成功"
// @Router /nailStyle/batchUpdateTags [post]
func (api *NailStyleApi) BatchUpdateNailStyleTags(c *gin.Context) {
	ctx := c.Request.Context()

	var payload nailReq.NailStyleBatchTagUpdateReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := nailStyleService.BatchUpdateNailStyleTags(ctx, payload)
	if err != nil {
		global.GVA_LOG.Error("批量更新标签失败!", zap.Error(err))
		response.FailWithMessage("批量更新标签失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "批量更新标签成功", c)
}

// UpdateNailStyle 更新美甲款式
// @Tags NailStyle
// @Summary 更新美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailStyle true "更新美甲款式"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /nailStyle/updateNailStyle [put]
func (api *NailStyleApi) UpdateNailStyle(c *gin.Context) {
	ctx := c.Request.Context()

	var nailStyle nail.NailStyle
	if err := c.ShouldBindJSON(&nailStyle); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := nailStyleService.UpdateNailStyle(ctx, nailStyle); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindNailStyle 根据ID查询美甲款式
// @Tags NailStyle
// @Summary 用id查询美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询美甲款式"
// @Success 200 {object} response.Response{data=nail.NailStyle,msg=string} "查询成功"
// @Router /nailStyle/findNailStyle [get]
func (api *NailStyleApi) FindNailStyle(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Query("ID")
	nailStyle, err := nailStyleService.GetNailStyle(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(nailStyle, c)
}

// GetNailStyleList 分页获取美甲款式列表
// @Tags NailStyle
// @Summary 分页获取美甲款式列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query nailReq.NailStyleSearch true "分页获取美甲款式列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /nailStyle/getNailStyleList [get]
func (api *NailStyleApi) GetNailStyleList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo nailReq.NailStyleSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := nailStyleService.GetNailStyleInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetNailStyleDataSource 获取数据源
// @Tags NailStyle
// @Summary 获取NailStyle的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /nailStyle/getNailStyleDataSource [get]
func (api *NailStyleApi) GetNailStyleDataSource(c *gin.Context) {
	ctx := c.Request.Context()

	dataSource, err := nailStyleService.GetNailStyleDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetNailStylePublic 公开访问的美甲款式接口
// @Tags NailStyle
// @Summary 不需要鉴权的美甲款式接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /nailStyle/getNailStylePublic [get]
func (api *NailStyleApi) GetNailStylePublic(c *gin.Context) {
	ctx := c.Request.Context()

	nailStyleService.GetNailStylePublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的美甲款式接口信息",
	}, "获取成功", c)
}
