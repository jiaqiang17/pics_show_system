package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/gva/api"

var (
	Router   = new(router)
	apiGroup = api.Api
)

type router struct{}
