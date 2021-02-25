package helm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Upload struct {
	ChartDir 	string
	checkMark 	helmCheck
}

type helmCheck struct {
	templates 	int
	helpers 	int
	values 		int
	chart 		int
}

func (u Upload) UloadInit() error {
	u.ChartDir = fmt.Sprintf( "%d-%d", time.Now().Unix(), rand.Intn(1000))
	u.checkMark = helmCheck{0, 0, 0, 0}
	if err := os.MkdirAll(u.ChartDir + "/templates", os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (u Upload) UploadSave( c *gin.Context) {
	// 多文件上传
	form, _:= c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		fileName := strings.Split(file.Filename, ".")
		switch fileName[0] {
		case "values":
			if err := c.SaveUploadedFile(file, u.ChartDir); err != nil {
				c.String(http.StatusBadRequest, "decode values failure")
			}
			u.checkMark.values ++
		case "Chart":
			if err := c.SaveUploadedFile(file, u.ChartDir); err != nil {
				c.String(http.StatusBadRequest, "decode Chart failure")
			}
			u.checkMark.chart ++
		case "_helpers":
			if err := c.SaveUploadedFile(file, u.ChartDir + "/templates"); err != nil {
				c.String(http.StatusBadRequest, "decode Chart failure")
			}
			u.checkMark.helpers ++
		default:
			if err := c.SaveUploadedFile(file, u.ChartDir + "/templates"); err != nil {
				c.String(http.StatusBadRequest, "decode template yaml failure")
			}
			u.checkMark.templates ++

		}
	}

	if u.checkMark.values * u.checkMark.chart * u.checkMark.helpers * u.checkMark.templates < 0 {
		c.
	}
}