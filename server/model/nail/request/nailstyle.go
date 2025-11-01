package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type NailStyleSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	StyleName      *string     `json:"styleName" form:"styleName"`
	IsRecommended  *bool       `json:"isRecommended" form:"isRecommended"`
	Status         *string     `json:"status" form:"status"`
	TagId          *uint       `json:"tagId" form:"tagId"` // 按标签ID筛选
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
