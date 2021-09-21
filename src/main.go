package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"crypto/rand"
	"errors"
	"main/base64"
	"main/redis"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const SAVE_DIR string = "/cheese/images/"
const RESULT_IMAGE string = "/cheese/result.jpg"

type Request struct {
	Image string `json:"image"`
}

func MakeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"*",
		},
		MaxAge: 24 * time.Hour,
	}))
	router.POST("/upload", func(c *gin.Context) {
		// フォームデータからファイルを読み込む
		data := Request{}
		err := c.BindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("get form err: %s", err.Error()),
			})
			return
		}
		image := base64.Decode(data.Image) //[]byte

		// 保存用ディレクトリがあるかどうか判定する
		f, err := os.Stat(SAVE_DIR)
		if os.IsNotExist(err) || !f.IsDir() {
			// なければディレクトリを作る
			err := os.MkdirAll(SAVE_DIR, 0777)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("mkdir save dir err: %s", err.Error()),
				})
				return
			}
		}

		// 保存パスの生成
		random, _ := MakeRandomStr(16)
		filepath := SAVE_DIR + random + ".jpg"
		// 保存パスにファイルを保存する
		// err = c.SaveUploadedFile(file, filepath);
		file, err := os.Create(filepath)
		defer file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("upload file err: %s", err.Error()),
			})
			return
		}
		qt := jpeg.Options{
			Quality: 80,
		}
		err = jpeg.Encode(file, image, &qt)
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
		output, err := exec.Command("bash", "-c", "/DisplayImage").CombinedOutput()
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
	router.GET("/download", func(c *gin.Context) {
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
