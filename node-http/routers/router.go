package routers

import "github.com/gin-gonic/gin"
import (
	. "github.com/dejavuzhou/evefs/store"
)

var Router *gin.Engine
var store *Store

func init() {
	store = NewStore("127.0.0.1:1212", "temp", 8)
	
	Router = gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	Router.POST("/upload", uploadFile)
	Router.GET("/read", readFile)
}
