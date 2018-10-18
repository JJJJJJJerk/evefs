package http

import (
	. "github.com/dejavuzhou/spookyfs/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

var store *Store

func init() {
	store = NewBarn("127.0.0.1:1212", "temp01", 8)
}

func uploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	needle, err := store.PutFile(file)
	c.JSON(http.StatusOK, gin.H{"data": needle, "error": err})
}

func readFile(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	needle := store.GetOneWithId([]byte(id))

	c.JSON(http.StatusOK, gin.H{"data": needle})

}
