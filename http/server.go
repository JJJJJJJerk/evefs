package http

import (
	"github.com/gin-gonic/gin"
)

func ServerStart() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", uploadFile)
	router.GET("/read", readFile)

	router.Run(":8080")
}
