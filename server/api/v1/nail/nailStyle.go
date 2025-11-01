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
func (nailStyleApi *NailStyleApi) CreateNailStyle(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var nailStyle nail.NailStyle
	err := c.ShouldBindJSON(&nailStyle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = nailStyleService.CreateNailStyle(ctx, &nailStyle)
	if err != nil {
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
func (nailStyleApi *NailStyleApi) DeleteNailStyle(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := nailStyleService.DeleteNailStyle(ctx, ID)
	if err != nil {
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
func (nailStyleApi *NailStyleApi) DeleteNailStyleByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := nailStyleService.DeleteNailStyleByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
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
func (nailStyleApi *NailStyleApi) UpdateNailStyle(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var nailStyle nail.NailStyle
	err := c.ShouldBindJSON(&nailStyle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = nailStyleService.UpdateNailStyle(ctx, nailStyle)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindNailStyle 用id查询美甲款式
// @Tags NailStyle
// @Summary 用id查询美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询美甲款式"
// @Success 200 {object} response.Response{data=nail.NailStyle,msg=string} "查询成功"
// @Router /nailStyle/findNailStyle [get]
func (nailStyleApi *NailStyleApi) FindNailStyle(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	renailStyle, err := nailStyleService.GetNailStyle(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(renailStyle, c)
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
func (nailStyleApi *NailStyleApi) GetNailStyleList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo nailReq.NailStyleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
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

// GetNailStyleDataSource 获取NailStyle的数据源
// @Tags NailStyle
// @Summary 获取NailStyle的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /nailStyle/getNailStyleDataSource [get]
func (nailStyleApi *NailStyleApi) GetNailStyleDataSource(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口为获取数据源定义的数据
	dataSource, err := nailStyleService.GetNailStyleDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetNailStylePublic 不需要鉴权的美甲款式接口
// @Tags NailStyle
// @Summary 不需要鉴权的美甲款式接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /nailStyle/getNailStylePublic [get]
func (nailStyleApi *NailStyleApi) GetNailStylePublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	nailStyleService.GetNailStylePublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的美甲款式接口信息",
	}, "获取成功", c)
}
