package nail

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
	nailReq "github.com/flipped-aurora/gin-vue-admin/server/model/nail/request"
)

type NailStyleService struct{}

// CreateNailStyle 创建美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) CreateNailStyle(ctx context.Context, nailStyle *nail.NailStyle) (err error) {
	err = global.GVA_DB.Create(nailStyle).Error
	return err
}

// DeleteNailStyle 删除美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) DeleteNailStyle(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&nail.NailStyle{}, "id = ?", ID).Error
	return err
}

// DeleteNailStyleByIds 批量删除美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) DeleteNailStyleByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]nail.NailStyle{}, "id in ?", IDs).Error
	return err
}

// UpdateNailStyle 更新美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) UpdateNailStyle(ctx context.Context, nailStyle nail.NailStyle) (err error) {
	err = global.GVA_DB.Model(&nail.NailStyle{}).Where("id = ?", nailStyle.ID).Updates(&nailStyle).Error
	return err
}

// GetNailStyle 根据ID获取美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) GetNailStyle(ctx context.Context, ID string) (nailStyle nail.NailStyle, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&nailStyle).Error
	return
}

// GetNailStyleInfoList 分页获取美甲款式记录
// Author [yourname](https://github.com/yourname)
func (nailStyleService *NailStyleService) GetNailStyleInfoList(ctx context.Context, info nailReq.NailStyleSearch) (list []nail.NailStyle, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&nail.NailStyle{})
	var nailStyles []nail.NailStyle
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.StyleName != nil && *info.StyleName != "" {
		db = db.Where("style_name LIKE ?", "%"+*info.StyleName+"%")
	}
	if info.IsRecommended != nil {
		db = db.Where("is_recommended = ?", *info.IsRecommended)
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	if info.TagId != nil {
		// 按标签ID筛选：使用 JSON_CONTAINS 查询包含指定标签ID的款式
		// MySQL JSON_CONTAINS 需要 JSON 格式的参数
		db = db.Where("JSON_CONTAINS(tag_ids, CAST(? AS JSON))", fmt.Sprintf("%d", *info.TagId))
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["view_count"] = true
	orderMap["like_count"] = true
	orderMap["sort"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&nailStyles).Error
	return nailStyles, total, err
}
func (nailStyleService *NailStyleService) GetNailStyleDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	tagIds := make([]map[string]any, 0)

	global.GVA_DB.Table("nail_tags").Where("deleted_at IS NULL").Select("tag_name as label,id as value").Scan(&tagIds)
	res["tagIds"] = tagIds
	return
}
func (nailStyleService *NailStyleService) GetNailStylePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
