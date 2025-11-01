package initialize

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/gva/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	if err := global.GVA_DB.WithContext(ctx).AutoMigrate(new(model.User)); err != nil {
		zap.L().Error(errors.Wrap(err, "register gva plugin tables failed").Error())
	}
}
