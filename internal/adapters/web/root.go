package web

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	router.GET("/", RootHandler)
}

func RootHandler(c *gin.Context) {

}
