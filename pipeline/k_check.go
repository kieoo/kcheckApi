package pipeline

import (
	"fmt"
	"github.com/gin-gonic/gin"
	k "kcheckApi/kcheck"
	p "kcheckApi/params"
	"net/http"
)

func KCheck(c *gin.Context) {

	in := &p.CRequest{}
	out := &p.CResponse{}

	if err := c.BindJSON(in.CheckBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("check your post data %s", err)})
		return
	}
	statusCode := http.StatusOK
	// k check
	if err := k.NormalCheck(in, out); err != nil {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, gin.H{"file_name": out.FileName, "result": out.Result, "chart": string(in.CheckBody.OriYaml), "msg": out.Message, "hints": out.Hints})
}
