// 自动生成模板NailTag
package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 美甲标签 结构体  NailTag
type NailTag struct {
	global.GVA_MODEL
	TagName     *string `json:"tagName" form:"tagName" gorm:"index;comment:标签名称;column:tag_name;size:100;" binding:"required"`          //标签名称
	TagIcon     *string `json:"tagIcon" form:"tagIcon" gorm:"comment:标签图标URL;column:tag_icon;size:255;"`                                //标签图标
	TagColor    *string `json:"tagColor" form:"tagColor" gorm:"default:#FF6B9D;comment:标签颜色（十六进制）;column:tag_color;size:50;"`           //标签颜色
	Description *string `json:"description" form:"description" gorm:"comment:标签描述;column:description;size:500;"`                        //标签描述
	Sort        *int64  `json:"sort" form:"sort" gorm:"default:1;comment:排序号，数字越小越靠前;column:sort;"`                                     //排序号
	Status      *string `json:"status" form:"status" gorm:"default:enabled;comment:状态：启用/禁用;column:status;size:20;" binding:"required"` //状态
}

// TableName 美甲标签 NailTag自定义表名 nail_tags
func (NailTag) TableName() string {
	return "nail_tags"
}
