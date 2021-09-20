package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"main/base64"
	"main/redis"

	"github.com/gin-gonic/gin"
)

const RESULT_IMAGE string = "/cheese/result.png"

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// フォームデータからファイルを読み込む
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("get form err: %s", err.Error()),
			})
			return
		}

		// ファイルを保存する
		filename := filepath.Base(file.Filename)
		filepath := "/tmp/images/"+filename
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("upload file err: %s", err.Error()),
			})
			return
		}

		// 保存パスをRedisに保存
		err = redis.SetValue("filepath", filepath)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("get redis err: %s", err.Error()),
			})
			return
		}

		// 保存パスをRedisから取得
		outputFilepath, err := redis.GetValue("filepath")
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("get redis err: %s", err.Error()),
			})
			return
		}

		// 保存パスからファイルを読込
		bytes, err := ioutil.ReadFile(outputFilepath)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("get form err: %s", err.Error()),
			})
			return
		}

		// ファイルをbase64に変換してその結果をjsonとして返却
		base64Encoding := base64.Encode(bytes)
		c.JSON(http.StatusOK, gin.H{
			"base64": fmt.Sprintf("%s", base64Encoding),
		})
	})
	router.GET("/download",func(c*gin.Context){
		// OpenCVからの出力画像を取得
		bytes, err := ioutil.ReadFile(RESULT_IMAGE)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("read file err: %s", err.Error()),
			})
			return
		}

		// 画像をbase64に変換してその結果をjsonとして返却
		base64Encoding := base64.Encode(bytes)
		c.JSON(http.StatusOK, gin.H{
			"base64": fmt.Sprintf("%s", base64Encoding),
		})
	})
	router.Run(":80")
}