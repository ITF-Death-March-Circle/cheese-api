package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"main/redis"

	"github.com/gin-gonic/gin"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
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

		// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(outputFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	// Print the full base64 representation of the image
	fmt.Println(base64Encoding)

		c.String(http.StatusOK, fmt.Sprintf("Saved file to '%s'. base64: %s", outputFilepath, base64Encoding))
	})
	router.Run(":80")
}