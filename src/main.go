package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"main/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		filepath := "/tmp/images/"+filename
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		redis.SetValue("filepath", filepath)
		outputFilepath, err := redis.GetValue("filepath")
		if err != nil {
			c.String(http.StatusBadGateway, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		log.Print(outputFilepath)

		c.String(http.StatusOK, fmt.Sprintf("Saved file to '%s'.", outputFilepath))
	})
	router.GET("/download",func(c*gin.Context){
		
	})
	router.Run(":80")
}