import service from '@/utils/request'
// @Tags NailStyle
// @Summary 创建美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailStyle true "创建美甲款式"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /nailStyle/createNailStyle [post]
export const createNailStyle = (data) => {
  return service({
    url: '/nailStyle/createNailStyle',
    method: 'post',
    data
  })
}

// @Tags NailStyle
// @Summary 删除美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailStyle true "删除美甲款式"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /nailStyle/deleteNailStyle [delete]
export const deleteNailStyle = (params) => {
  return service({
    url: '/nailStyle/deleteNailStyle',
    method: 'delete',
    params
  })
}

// @Tags NailStyle
// @Summary 批量删除美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除美甲款式"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /nailStyle/deleteNailStyle [delete]
export const deleteNailStyleByIds = (params) => {
  return service({
    url: '/nailStyle/deleteNailStyleByIds',
    method: 'delete',
    params
  })
}

// @Tags NailStyle
// @Summary 更新美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.NailStyle true "更新美甲款式"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /nailStyle/updateNailStyle [put]
export const updateNailStyle = (data) => {
  return service({
    url: '/nailStyle/updateNailStyle',
    method: 'put',
    data
  })
}

// @Tags NailStyle
// @Summary 用id查询美甲款式
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.NailStyle true "用id查询美甲款式"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /nailStyle/findNailStyle [get]
export const findNailStyle = (params) => {
  return service({
    url: '/nailStyle/findNailStyle',
    method: 'get',
    params
  })
}

// @Tags NailStyle
// @Summary 分页获取美甲款式列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取美甲款式列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /nailStyle/getNailStyleList [get]
export const getNailStyleList = (params) => {
  return service({
    url: '/nailStyle/getNailStyleList',
    method: 'get',
    params
  })
}
// @Tags NailStyle
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /nailStyle/findNailStyleDataSource [get]
export const getNailStyleDataSource = () => {
  return service({
    url: '/nailStyle/getNailStyleDataSource',
    method: 'get',
  })
}

// @Tags NailStyle
// @Summary 不需要鉴权的美甲款式接口
// @Accept application/json
// @Produce application/json
// @Param data query nailReq.NailStyleSearch true "分页获取美甲款式列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /nailStyle/getNailStylePublic [get]
export const getNailStylePublic = () => {
  return service({
    url: '/nailStyle/getNailStylePublic',
    method: 'get',
  })
}
