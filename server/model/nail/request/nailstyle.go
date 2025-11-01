package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// NailStyleSearch 定义美甲款式搜索条件
type NailStyleSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	StyleName      *string     `json:"styleName" form:"styleName"`
	IsRecommended  *bool       `json:"isRecommended" form:"isRecommended"`
	Status         *string     `json:"status" form:"status"`
	TagId          *uint       `json:"tagId" form:"tagId"`       // 兼容旧版的单标签筛选
	TagIds         []uint      `json:"tagIds" form:"tagIds[]"`   // 多标签筛选
	MatchAll       bool        `json:"matchAll" form:"matchAll"` // 是否需要匹配所有标签
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}

// NailStyleBatchTagUpdateReq 用于批量更新款式标签
type NailStyleBatchTagUpdateReq struct {
	StyleIDs     []uint `json:"styleIds" binding:"required,min=1"`       // 需要更新的款式ID
	AddTagIDs    []uint `json:"addTagIds"`                               // 需要新增的标签ID
	RemoveTagIDs []uint `json:"removeTagIds"`                            // 需要移除的标签ID
	OperatorNote string `json:"operatorNote" binding:"max=200" form:"-"` // 操作备注，记录审计信息
}

// NailStyleBatchTagUpdateResult 返回批量标签更新结果
type NailStyleBatchTagUpdateResult struct {
	Updated       int    `json:"updated"`
	Skipped       []uint `json:"skipped"`
	InvalidTagIds []uint `json:"invalidTagIds"`
	AppliedAdd    []uint `json:"appliedAdd"`
	AppliedRemove []uint `json:"appliedRemove"`
	OperatorNote  string `json:"operatorNote"`
	StyleTotal    int    `json:"styleTotal"`
}
