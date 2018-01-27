package web

import "github.com/gin-gonic/gin"

var Default *gin.Engine

func init() {

	Default = gin.New()

}
