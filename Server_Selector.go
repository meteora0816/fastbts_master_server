package main

import (
	"fmt"
	"net/http"
	"os"
)
import "github.com/gin-gonic/gin"

var r *gin.Engine
var filesPath string

func init() {
	filesPath = "../files"
}

func main() {
	fmt.Println("start")
	r = gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, `hello, FastBTS!`)
	})
	r.GET("/file/:filesize", func(c *gin.Context) {
		fileSize := c.Param("filesize")
		filePath := filesPath + "/" + fileSize + "file.txt"
		file, err := os.Open(filePath)
		defer file.Close()
		fileName := fileSize + "file.txt"
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			c.Header("Content-Type", "text/plain")
			c.Header("Content-Disposition", "attachment; filename="+fileName)
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Cache-Control", "no-cache")
			c.File(filePath)
		}
	})
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
