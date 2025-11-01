package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/gva/service"

var (
	Api          = new(api)
	serviceGroup = service.Service
)

type api struct{}
