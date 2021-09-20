package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"main/base64"
	"main/redis"

	"github.com/gin-gonic/gin"
)

const SAVE_DIR string = "/cheese/images/"
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

		// 保存用ディレクトリがあるかどうか判定する
		f, err := os.Stat(SAVE_DIR);
		if os.IsNotExist(err) || !f.IsDir() {
			// なければディレクトリを作る
			err := os.MkdirAll(SAVE_DIR, 0777);
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("mkdir save dir err: %s", err.Error()),
				})
				return
			}
		}

		// 保存パスの生成
		filepath := SAVE_DIR + filepath.Base(file.Filename)

		// 保存パスにファイルを保存する
		err = c.SaveUploadedFile(file, filepath);
		if err != nil {
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

		// opencv製画像処理を実行
		output, err := exec.Command("ls").CombinedOutput()
    log.Printf("opencv output:\n%s :Error:\n%v\n", output, err)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("exec opencv err: %s", err.Error()),
			})
			return
		}

		// ファイルをbase64に変換してその結果をjsonとして返却
		base64Encoding := base64.Encode(bytes)
		c.JSON(http.StatusOK, gin.H{
			"base64": fmt.Sprintf("%s", base64Encoding),
			"output": fmt.Sprintf("%s", output),
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