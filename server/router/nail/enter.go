package nail

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	NailStyleRouter
	NailTagRouter
}

var (
	nailStyleApi = api.ApiGroupApp.NailApiGroup.NailStyleApi
	nailTagApi   = api.ApiGroupApp.NailApiGroup.NailTagApi
)
