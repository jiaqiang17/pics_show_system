package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/nail"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(nail.NailStyle{}, nail.NailTag{})
	if err != nil {
		return err
	}
	return nil
}
