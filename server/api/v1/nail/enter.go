package nail

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	NailStyleApi
	NailTagApi
}

var (
	nailStyleService = service.ServiceGroupApp.NailServiceGroup.NailStyleService
	nailTagService   = service.ServiceGroupApp.NailServiceGroup.NailTagService
)
