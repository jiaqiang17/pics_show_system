package nail

import (
	"context"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
	nailReq "github.com/flipped-aurora/gin-vue-admin/server/model/nail/request"
)

type NailTagService struct{}

// CreateNailTag 创建美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) CreateNailTag(ctx context.Context, nailTag *nail.NailTag) (err error) {
	err = global.GVA_DB.Create(nailTag).Error
	return err
}

// DeleteNailTag 删除美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) DeleteNailTag(ctx context.Context, ID string) (err error) {
	// 检查是否有美甲款式绑定了该标签
	var count int64
	err = global.GVA_DB.Model(&nail.NailStyle{}).
		Where("JSON_CONTAINS(tag_ids, CAST(? AS JSON))", ID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(fmt.Sprintf("该标签已被 %d 个美甲款式使用，无法删除。请先解除绑定关系", count))
	}

	// 如果没有款式使用该标签，则可以删除
	err = global.GVA_DB.Delete(&nail.NailTag{}, "id = ?", ID).Error
	return err
}

// DeleteNailTagByIds 批量删除美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) DeleteNailTagByIds(ctx context.Context, IDs []string) (err error) {
	// 检查每个标签是否有美甲款式绑定
	var boundTags []string
	for _, id := range IDs {
		var count int64
		err = global.GVA_DB.Model(&nail.NailStyle{}).
			Where("JSON_CONTAINS(tag_ids, CAST(? AS JSON))", id).
			Count(&count).Error
		if err != nil {
			return err
		}

		if count > 0 {
			// 获取标签名称用于提示
			var tag nail.NailTag
			global.GVA_DB.Where("id = ?", id).First(&tag)
			if tag.TagName != nil {
				boundTags = append(boundTags, fmt.Sprintf("%s(%d个款式)", *tag.TagName, count))
			} else {
				boundTags = append(boundTags, fmt.Sprintf("ID:%s(%d个款式)", id, count))
			}
		}
	}

	if len(boundTags) > 0 {
		return errors.New(fmt.Sprintf("以下标签已被美甲款式使用，无法删除：%v。请先解除绑定关系", boundTags))
	}

	// 如果都没有款式使用这些标签，则可以批量删除
	err = global.GVA_DB.Delete(&[]nail.NailTag{}, "id in ?", IDs).Error
	return err
}

// UpdateNailTag 更新美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) UpdateNailTag(ctx context.Context, nailTag nail.NailTag) (err error) {
	err = global.GVA_DB.Model(&nail.NailTag{}).Where("id = ?", nailTag.ID).Updates(&nailTag).Error
	return err
}

// GetNailTag 根据ID获取美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) GetNailTag(ctx context.Context, ID string) (nailTag nail.NailTag, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&nailTag).Error
	return
}

// GetNailTagInfoList 分页获取美甲标签记录
// Author [yourname](https://github.com/yourname)
func (nailTagService *NailTagService) GetNailTagInfoList(ctx context.Context, info nailReq.NailTagSearch) (list []nail.NailTag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&nail.NailTag{})
	var nailTags []nail.NailTag
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.TagName != nil && *info.TagName != "" {
		db = db.Where("tag_name LIKE ?", "%"+*info.TagName+"%")
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
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

	err = db.Find(&nailTags).Error
	return nailTags, total, err
}
func (nailTagService *NailTagService) GetNailTagPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
