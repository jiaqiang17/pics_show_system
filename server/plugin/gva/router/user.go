package router

import "github.com/gin-gonic/gin"

type user struct{}

func (u *user) Init(public, private *gin.RouterGroup) {}
