// 自动生成模板NailStyle
package nail

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// 美甲款式 结构体  NailStyle
type NailStyle struct {
	global.GVA_MODEL
	StyleName     *string        `json:"styleName" form:"styleName" gorm:"comment:款式名称;column:style_name;size:255;" binding:"required"`          //款式名称
	Images        datatypes.JSON `json:"images" form:"images" gorm:"comment:美甲图片集;column:images;" swaggertype:"array,object" binding:"required"` //美甲图片集
	Description   *string        `json:"description" form:"description" gorm:"comment:款式介绍;column:description;type:text;"`                       //款式介绍
	IsRecommended *bool          `json:"isRecommended" form:"isRecommended" gorm:"default:false;comment:是否推荐到首页;column:is_recommended;"`         //是否推荐
	ViewCount     *int64         `json:"viewCount" form:"viewCount" gorm:"default:0;comment:浏览次数;column:view_count;"`                            //浏览次数
	LikeCount     *int64         `json:"likeCount" form:"likeCount" gorm:"default:0;comment:点赞数;column:like_count;"`                             //点赞数
	Sort          *int64         `json:"sort" form:"sort" gorm:"default:1;comment:排序号，数字越小越靠前;column:sort;"`                                     //排序号
	Status        *string        `json:"status" form:"status" gorm:"default:enabled;comment:状态：启用/禁用;column:status;size:20;" binding:"required"` //状态
	TagIds        datatypes.JSON `json:"tagIds" form:"tagIds" gorm:"comment:关联的美甲标签ID列表;column:tag_ids;" swaggertype:"array,object"`             //关联的标签ID列表
}

// TableName 美甲款式 NailStyle自定义表名 nail_styles
func (NailStyle) TableName() string {
	return "nail_styles"
}
