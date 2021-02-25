package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"kcheckApi/kcheck"
	"net/http"
	"os"
	"regexp"
)

type Pb struct {
	Code string `json:"code"`
}

func main() {
	logfile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)

	r := gin.New()
	r.Use(CorsMiddleware())
	r.GET("version", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("return", func(c *gin.Context) {
		body := Pb{}
		err := c.BindJSON(&body)
		if err != nil {
			fmt.Println(err.Error())
			c.String(http.StatusBadRequest, "bad")
		} else {
			c.String(http.StatusOK, string(body.Code))
		}
	})

	r.POST("/kcheck", kcheck.NormalCheck)

	r.POST("/upload")

	r.Run(":7002")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var filterHost = [...]string{"10.192.176.36:8080"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = false
		for _, v := range filterHost {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
