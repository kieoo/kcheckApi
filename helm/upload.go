package helm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Upload struct {
	ChartDir  string `json:"chartDir"`
	checkMark helmCheck
}

type helmCheck struct {
	templates int
	helpers   int
	values    int
	chart     int
}

func (u *Upload) UploadInit() error {
	u.ChartDir = fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Intn(1000))
	u.checkMark = helmCheck{0, 0, 0, 0}
	if err := os.MkdirAll(u.ChartDir+"/templates", os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (u *Upload) UploadClose() {
	_ = os.RemoveAll(u.ChartDir)
}

type UploadSaveError struct {
	S string
}

func (e *UploadSaveError) Error() string {
	return e.S
}

func (u *Upload) UploadSave(c *gin.Context) *UploadSaveError {
	// 多文件上传
	form, _ := c.MultipartForm()
	files := form.File["files"]

	// save Chart file
	for _, file := range files {
		fileName := strings.Split(file.Filename, ".")
		switch fileName[0] {
		case "values":
			if err := c.SaveUploadedFile(file, u.ChartDir+"/"+file.Filename); err != nil {
				return &UploadSaveError{"decode values failure"}
				//return error{ return "decode values failure"}
			}
			u.checkMark.values++
		case "Chart":
			if err := c.SaveUploadedFile(file, u.ChartDir+"/"+file.Filename); err != nil {
				return &UploadSaveError{"decode Chart failure"}
			}
			u.checkMark.chart++
		case "_helpers":
			if err := c.SaveUploadedFile(file, u.ChartDir+"/templates/"+file.Filename); err != nil {
				return &UploadSaveError{"decode Chart failure"}
			}
			u.checkMark.helpers++
		default:
			if err := c.SaveUploadedFile(file, u.ChartDir+"/templates/"+file.Filename); err != nil {
				return &UploadSaveError{"decode template yaml failure"}
			}
			u.checkMark.templates++

		}
	}

	// check Chart file
	if u.checkMark.values*u.checkMark.chart*u.checkMark.helpers*u.checkMark.templates < 1 {
		log.Printf("values: %d, chart: %d, helpers : %d, templates: %d",
			u.checkMark.values,
			u.checkMark.chart,
			u.checkMark.helpers,
			u.checkMark.templates)

		return &UploadSaveError{"values/chart/helpers/template is require"}
	}

	return nil

}
