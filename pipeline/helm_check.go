package pipeline

import (
	"github.com/gin-gonic/gin"
	h "kcheckApi/helm"
	p "kcheckApi/params"
	"net/http"
	"path/filepath"
	"strings"
)

func HelmCheck(c *gin.Context) {

	in := &p.CRequest{}
	out := &p.CResponse{}
	form, _ := c.MultipartForm()

	if form == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "MultipartForm key files is need"})
		return
	}

	if _, ok := form.File["files"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "MultipartForm key files is need"})
		return
	}

	files := form.File["files"]

	if len(files) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "should upload chart file"})
		return
	}
	in.MultipartFile = files

	FileNameList := strings.Split(filepath.ToSlash(files[0].Filename), "/")

	out.FileName = FileNameList[0]

	statusCode := http.StatusOK
	// helm check
	if err := h.HelmChange(in, out); err != nil {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, gin.H{"file_name": out.FileName, "result": out.Result, "chart": string(in.CheckBody.OriYaml), "msg": out.Message})
}
