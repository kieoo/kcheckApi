package pipeline

import (
	"github.com/gin-gonic/gin"
	h "kcheckApi/helm"
	p "kcheckApi/params"
	"net/http"
)

func HelmCheck(c *gin.Context) {

	in := &p.CRequest{}
	out := &p.CResponse{}
	form, _ := c.MultipartForm()
	files := form.File["files"]

	if len(files) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "should upload chart file"})
		return
	}
	in.MultipartFile = files

	out.FileName = files[0].Filename

	statusCode := http.StatusOK
	// helm check
	if err := h.HelmChange(in, out); err != nil {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, gin.H{"file_name": out.FileName, "result": out.Result, "chart": string(in.CheckBody.OriYaml), "msg": out.Message})
}
