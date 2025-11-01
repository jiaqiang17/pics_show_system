package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
	nailReq "github.com/flipped-aurora/gin-vue-admin/server/model/nail/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NailTagApi struct{}

// CreateNailTag 创建美甲标签
// @Tags NailTag
// @Summary 创建美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailTag true "创建美甲标签"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /nailTag/createNailTag [post]
func (nailTagApi *NailTagApi) CreateNailTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var nailTag nail.NailTag
	err := c.ShouldBindJSON(&nailTag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = nailTagService.CreateNailTag(ctx, &nailTag)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteNailTag 删除美甲标签
// @Tags NailTag
// @Summary 删除美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailTag true "删除美甲标签"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /nailTag/deleteNailTag [delete]
func (nailTagApi *NailTagApi) DeleteNailTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := nailTagService.DeleteNailTag(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteNailTagByIds 批量删除美甲标签
// @Tags NailTag
// @Summary 批量删除美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /nailTag/deleteNailTagByIds [delete]
func (nailTagApi *NailTagApi) DeleteNailTagByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := nailTagService.DeleteNailTagByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateNailTag 更新美甲标签
// @Tags NailTag
// @Summary 更新美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body nail.NailTag true "更新美甲标签"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /nailTag/updateNailTag [put]
func (nailTagApi *NailTagApi) UpdateNailTag(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var nailTag nail.NailTag
	err := c.ShouldBindJSON(&nailTag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = nailTagService.UpdateNailTag(ctx, nailTag)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindNailTag 用id查询美甲标签
// @Tags NailTag
// @Summary 用id查询美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询美甲标签"
// @Success 200 {object} response.Response{data=nail.NailTag,msg=string} "查询成功"
// @Router /nailTag/findNailTag [get]
func (nailTagApi *NailTagApi) FindNailTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	renailTag, err := nailTagService.GetNailTag(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(renailTag, c)
}

// GetNailTagList 分页获取美甲标签列表
// @Tags NailTag
// @Summary 分页获取美甲标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query nailReq.NailTagSearch true "分页获取美甲标签列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /nailTag/getNailTagList [get]
func (nailTagApi *NailTagApi) GetNailTagList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo nailReq.NailTagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := nailTagService.GetNailTagInfoList(ctx, pageInfo)
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

// GetNailTagPublic 不需要鉴权的美甲标签接口
// @Tags NailTag
// @Summary 不需要鉴权的美甲标签接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /nailTag/getNailTagPublic [get]
func (nailTagApi *NailTagApi) GetNailTagPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	nailTagService.GetNailTagPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的美甲标签接口信息",
	}, "获取成功", c)
}
