package nail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
	nailReq "github.com/flipped-aurora/gin-vue-admin/server/model/nail/request"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NailStyleService struct{}

// CreateNailStyle 创建美甲款式记录
func (s *NailStyleService) CreateNailStyle(ctx context.Context, nailStyle *nail.NailStyle) error {
	return global.GVA_DB.WithContext(ctx).Create(nailStyle).Error
}

// DeleteNailStyle 删除美甲款式记录
func (s *NailStyleService) DeleteNailStyle(ctx context.Context, id string) error {
	return global.GVA_DB.WithContext(ctx).Delete(&nail.NailStyle{}, "id = ?", id).Error
}

// DeleteNailStyleByIds 批量删除美甲款式记录
func (s *NailStyleService) DeleteNailStyleByIds(ctx context.Context, ids []string) error {
	return global.GVA_DB.WithContext(ctx).Delete(&[]nail.NailStyle{}, "id in ?", ids).Error
}

// UpdateNailStyle 更新美甲款式记录
func (s *NailStyleService) UpdateNailStyle(ctx context.Context, nailStyle nail.NailStyle) error {
	return global.GVA_DB.WithContext(ctx).
		Model(&nail.NailStyle{}).
		Where("id = ?", nailStyle.ID).
		Updates(&nailStyle).Error
}

// GetNailStyle 根据ID获取美甲款式记录
func (s *NailStyleService) GetNailStyle(ctx context.Context, id string) (nail.NailStyle, error) {
	var nailStyle nail.NailStyle
	err := global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&nailStyle).Error
	return nailStyle, err
}

// GetNailStyleInfoList 分页获取美甲款式记录
func (s *NailStyleService) GetNailStyleInfoList(ctx context.Context, info nailReq.NailStyleSearch) (list []nail.NailStyle, total int64, err error) {
	db := global.GVA_DB.WithContext(ctx).Model(&nail.NailStyle{})

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

	switch {
	case len(info.TagIds) > 0:
		tagIds := uniqueUint(info.TagIds)
		if len(tagIds) > 0 {
			if info.MatchAll {
				for _, tagID := range tagIds {
					db = db.Where("JSON_CONTAINS(tag_ids, CAST(? AS JSON))", fmt.Sprintf("%d", tagID))
				}
			} else {
				clauses := make([]string, 0, len(tagIds))
				args := make([]interface{}, 0, len(tagIds))
				for _, tagID := range tagIds {
					clauses = append(clauses, "JSON_CONTAINS(tag_ids, CAST(? AS JSON))")
					args = append(args, fmt.Sprintf("%d", tagID))
				}
				db = db.Where("("+strings.Join(clauses, " OR ")+")", args...)
			}
		}
	case info.TagId != nil:
		db = db.Where("JSON_CONTAINS(tag_ids, CAST(? AS JSON))", fmt.Sprintf("%d", *info.TagId))
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderableFields := map[string]struct{}{
		"id":         {},
		"created_at": {},
		"view_count": {},
		"like_count": {},
		"sort":       {},
	}
	if _, ok := orderableFields[info.Sort]; ok {
		orderExpr := info.Sort
		if info.Order == "descending" {
			orderExpr += " desc"
		}
		db = db.Order(orderExpr)
	}

	if info.PageSize > 0 {
		offset := info.PageSize * (info.Page - 1)
		db = db.Limit(info.PageSize).Offset(offset)
	}

	var nailStyles []nail.NailStyle
	if err = db.Find(&nailStyles).Error; err != nil {
		return nil, 0, err
	}

	filterTags := make([]uint, 0)
	if len(info.TagIds) > 0 {
		filterTags = uniqueUint(info.TagIds)
	} else if info.TagId != nil {
		filterTags = append(filterTags, *info.TagId)
	}

	if len(filterTags) > 0 {
		tagSet := sliceToSet(filterTags)
		for i := range nailStyles {
			var currentTags []uint
			if len(nailStyles[i].TagIds) > 0 {
				if err := json.Unmarshal(nailStyles[i].TagIds, &currentTags); err != nil {
					continue
				}
			}
			matchCount := 0
			for _, tagID := range currentTags {
				if _, ok := tagSet[tagID]; ok {
					matchCount++
				}
			}
			nailStyles[i].TagMatchCount = matchCount
		}
	}

	return nailStyles, total, nil
}

// BatchUpdateNailStyleTags 批量更新款式标签
func (s *NailStyleService) BatchUpdateNailStyleTags(ctx context.Context, req nailReq.NailStyleBatchTagUpdateReq) (nailReq.NailStyleBatchTagUpdateResult, error) {
	result := nailReq.NailStyleBatchTagUpdateResult{
		AppliedAdd:    uniqueUint(req.AddTagIDs),
		AppliedRemove: uniqueUint(req.RemoveTagIDs),
		OperatorNote:  req.OperatorNote,
	}

	styleIDs := uniqueUint(req.StyleIDs)
	result.StyleTotal = len(styleIDs)

	if len(styleIDs) == 0 {
		return result, errors.New("styleIds 不能为空")
	}
	if len(result.AppliedAdd) == 0 && len(result.AppliedRemove) == 0 {
		return result, errors.New("至少指定一个需要新增或移除的标签")
	}

	allTagIDs := append(append([]uint{}, result.AppliedAdd...), result.AppliedRemove...)
	if len(allTagIDs) > 0 {
		var existed []uint
		if err := global.GVA_DB.WithContext(ctx).
			Table("nail_tags").
			Where("id IN ?", allTagIDs).
			Pluck("id", &existed).Error; err != nil {
			return result, err
		}
		existedSet := sliceToSet(existed)

		validAdd := make([]uint, 0, len(result.AppliedAdd))
		for _, id := range result.AppliedAdd {
			if _, ok := existedSet[id]; ok {
				validAdd = append(validAdd, id)
			} else {
				result.InvalidTagIds = append(result.InvalidTagIds, id)
			}
		}

		validRemove := make([]uint, 0, len(result.AppliedRemove))
		for _, id := range result.AppliedRemove {
			if _, ok := existedSet[id]; ok {
				validRemove = append(validRemove, id)
			} else {
				result.InvalidTagIds = append(result.InvalidTagIds, id)
			}
		}

		result.AppliedAdd = uniqueUint(validAdd)
		result.AppliedRemove = uniqueUint(validRemove)
		result.InvalidTagIds = uniqueUint(result.InvalidTagIds)

		if len(result.AppliedAdd) == 0 && len(result.AppliedRemove) == 0 {
			return result, errors.New("没有有效的标签可供更新")
		}
	}

	err := global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var styles []nail.NailStyle
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", styleIDs).
			Find(&styles).Error; err != nil {
			return err
		}

		foundIDSet := make(map[uint]struct{}, len(styles))
		for _, style := range styles {
			foundIDSet[style.ID] = struct{}{}
		}

		for _, id := range styleIDs {
			if _, ok := foundIDSet[id]; !ok {
				result.Skipped = append(result.Skipped, id)
			}
		}

		addSet := sliceToSet(result.AppliedAdd)
		removeSet := sliceToSet(result.AppliedRemove)

		for _, style := range styles {
			var current []uint
			if len(style.TagIds) > 0 {
				if err := json.Unmarshal(style.TagIds, &current); err != nil {
					return err
				}
			}

			beforeSet := sliceToSet(current)
			afterSet := cloneSet(beforeSet)

			for id := range removeSet {
				delete(afterSet, id)
			}
			for id := range addSet {
				afterSet[id] = struct{}{}
			}

			beforeSlice := setToSortedSlice(beforeSet)
			afterSlice := setToSortedSlice(afterSet)
			if slicesEqual(beforeSlice, afterSlice) {
				result.Skipped = append(result.Skipped, style.ID)
				continue
			}

			tagBytes, err := json.Marshal(afterSlice)
			if err != nil {
				return err
			}
			if len(afterSlice) == 0 {
				tagBytes = []byte("[]")
			}

			if err := tx.Model(&nail.NailStyle{}).
				Where("id = ?", style.ID).
				Updates(map[string]any{"tag_ids": datatypes.JSON(tagBytes)}).Error; err != nil {
				return err
			}
			result.Updated++
		}

		result.Skipped = uniqueUint(result.Skipped)
		return nil
	})

	return result, err
}

// GetNailStyleDataSource 获取美甲款式数据源
func (s *NailStyleService) GetNailStyleDataSource(ctx context.Context) (map[string][]map[string]any, error) {
	res := make(map[string][]map[string]any)

	tagIds := make([]map[string]any, 0)
	if err := global.GVA_DB.WithContext(ctx).
		Table("nail_tags").
		Where("deleted_at IS NULL").
		Select("tag_name as label,id as value").
		Scan(&tagIds).Error; err != nil {
		return nil, err
	}
	res["tagIds"] = tagIds
	return res, nil
}

// GetNailStylePublic 留待业务实现
func (s *NailStyleService) GetNailStylePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请根据业务需求自定义实现
}

func uniqueUint(values []uint) []uint {
	if len(values) == 0 {
		return nil
	}
	set := make(map[uint]struct{}, len(values))
	for _, v := range values {
		if v == 0 {
			continue
		}
		set[v] = struct{}{}
	}
	out := make([]uint, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func sliceToSet(values []uint) map[uint]struct{} {
	set := make(map[uint]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

func cloneSet(src map[uint]struct{}) map[uint]struct{} {
	dst := make(map[uint]struct{}, len(src))
	for k := range src {
		dst[k] = struct{}{}
	}
	return dst
}

func setToSortedSlice(set map[uint]struct{}) []uint {
	if len(set) == 0 {
		return nil
	}
	out := make([]uint, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func slicesEqual(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
