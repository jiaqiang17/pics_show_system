import service from '@/utils/request'

// 创建美甲款式
export const createNailStyle = (data) => {
  return service({
    url: '/nailStyle/createNailStyle',
    method: 'post',
    data
  })
}

// 删除美甲款式
export const deleteNailStyle = (params) => {
  return service({
    url: '/nailStyle/deleteNailStyle',
    method: 'delete',
    params
  })
}

// 批量删除美甲款式
export const deleteNailStyleByIds = (params) => {
  return service({
    url: '/nailStyle/deleteNailStyleByIds',
    method: 'delete',
    params
  })
}

// 更新美甲款式
export const updateNailStyle = (data) => {
  return service({
    url: '/nailStyle/updateNailStyle',
    method: 'put',
    data
  })
}

// 查询单个美甲款式
export const findNailStyle = (params) => {
  return service({
    url: '/nailStyle/findNailStyle',
    method: 'get',
    params
  })
}

// 获取美甲款式列表
export const getNailStyleList = (params) => {
  return service({
    url: '/nailStyle/getNailStyleList',
    method: 'get',
    params
  })
}

// 批量更新美甲款式标签
export const batchUpdateNailStyleTags = (data) => {
  return service({
    url: '/nailStyle/batchUpdateTags',
    method: 'post',
    data
  })
}

// 获取美甲款式数据源（标签等）
export const getNailStyleDataSource = () => {
  return service({
    url: '/nailStyle/getNailStyleDataSource',
    method: 'get'
  })
}

// 公开的美甲款式列表（无需鉴权）
export const getNailStylePublic = (params) => {
  return service({
    url: '/nailStyle/getNailStylePublic',
    method: 'get',
    params
  })
}
