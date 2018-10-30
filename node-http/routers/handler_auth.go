package routers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func login(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		logrus.WithError(err).Error("form file is empty")
	}
	needle, err := store.PutFile(file)
	if err != nil {
		logrus.WithError(err).Error("write file to haystack failed")
	}
	c.JSON(http.StatusOK, gin.H{"data": needle, "error": err})
}

func refreshToken(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	needle, err := store.GetFile([]byte(id))
	if err != nil {
		logrus.Error(err)
	}
	
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, needle.Name),
	}
	fileReader := bytes.NewBuffer(needle.FileBytes)
	
	c.DataFromReader(http.StatusOK, int64(needle.Size), needle.Mime, fileReader, extraHeaders)
	
}
