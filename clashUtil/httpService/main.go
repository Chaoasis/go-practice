package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	//download static resource
	router.Static("/assets", "./assets")
	//uploadMultipart
	router.MaxMultipartMemory = 8 << 20
	router.POST("/upload", func(c *gin.Context) {
		header := c.GetHeader("appKey")
		if header != "hello world" {
			c.String(http.StatusBadRequest, fmt.Sprintf("%d files uploaded!", 0))
			return
		}
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		for _, file := range files {
			err := c.SaveUploadedFile(file, "./assets/"+file.Filename)
			if err != nil {
				fmt.Println("save file fail:", err)
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":9001")
}
