import service from '@/utils/request'
// @Tags NailTag
// @Summary 创建美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailTag true "创建美甲标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /nailTag/createNailTag [post]
export const createNailTag = (data) => {
  return service({
    url: '/nailTag/createNailTag',
    method: 'post',
    data
  })
}

// @Tags NailTag
// @Summary 删除美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailTag true "删除美甲标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /nailTag/deleteNailTag [delete]
export const deleteNailTag = (params) => {
  return service({
    url: '/nailTag/deleteNailTag',
    method: 'delete',
    params
  })
}

// @Tags NailTag
// @Summary 批量删除美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除美甲标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /nailTag/deleteNailTag [delete]
export const deleteNailTagByIds = (params) => {
  return service({
    url: '/nailTag/deleteNailTagByIds',
    method: 'delete',
    params
  })
}

// @Tags NailTag
// @Summary 更新美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailTag true "更新美甲标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /nailTag/updateNailTag [put]
export const updateNailTag = (data) => {
  return service({
    url: '/nailTag/updateNailTag',
    method: 'put',
    data
  })
}

// @Tags NailTag
// @Summary 用id查询美甲标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.NailTag true "用id查询美甲标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /nailTag/findNailTag [get]
export const findNailTag = (params) => {
  return service({
    url: '/nailTag/findNailTag',
    method: 'get',
    params
  })
}

// @Tags NailTag
// @Summary 分页获取美甲标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取美甲标签列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /nailTag/getNailTagList [get]
export const getNailTagList = (params) => {
  return service({
    url: '/nailTag/getNailTagList',
    method: 'get',
    params
  })
}

// @Tags NailTag
// @Summary 不需要鉴权的美甲标签接口
// @Accept application/json
// @Produce application/json
// @Param data query nailReq.NailTagSearch true "分页获取美甲标签列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /nailTag/getNailTagPublic [get]
export const getNailTagPublic = () => {
  return service({
    url: '/nailTag/getNailTagPublic',
    method: 'get',
  })
}
