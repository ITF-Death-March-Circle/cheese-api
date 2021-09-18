package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"main/redis"

	"github.com/gin-gonic/gin"
)

func toBase64(b []byte) string {
	var base64EncodingPrefix string
	// Determine the content type of the image file
	mimeType := http.DetectContentType(b)
	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64EncodingPrefix = "data:image/jpeg;base64,"
	case "image/jpg":
		base64EncodingPrefix = "data:image/jpeg;base64,"
	case "image/png":
		base64EncodingPrefix = "data:image/png;base64,"
	}
	return base64EncodingPrefix + base64.StdEncoding.EncodeToString(b)
}

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("get form err: %s", err.Error()),
			})
			return
		}

		filename := filepath.Base(file.Filename)
		filepath := "/tmp/images/"+filename
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  fmt.Sprintf("upload file err: %s", err.Error()),
			})
			return
		}

		redis.SetValue("filepath", filepath)
		outputFilepath, err := redis.GetValue("filepath")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":  fmt.Sprintf("get form err: %s", err.Error()),
			})
			return
		}

		// Read the entire file into a byte slice
		bytes, err := ioutil.ReadFile(outputFilepath)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "get form err:" + err.Error(),
			})
			return
		}

		// Append the base64 encoded output
		base64Encoding := toBase64(bytes)
		c.JSON(http.StatusOK, gin.H{
			"base64": base64Encoding,
		})
	})
	router.Run(":80")
}