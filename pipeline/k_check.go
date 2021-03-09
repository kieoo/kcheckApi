package pipeline

import (
	"fmt"
	"github.com/gin-gonic/gin"
	k "kcheckApi/kcheck"
	p "kcheckApi/params"
	"kcheckApi/util"
	"net/http"
	"regexp"
)

func KCheck(c *gin.Context) {

	in := &p.CRequest{}
	out := &p.CResponse{}

	if c.BindJSON(&in.CheckBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("check your post data %s")})
		return
	}
	statusCode := http.StatusOK
	in.CheckBody.OriYaml = []byte(in.CheckBody.OriYamlString)

	checkYamlList := util.CleanOriYaml(in.CheckBody.OriYaml)

	for _, checkYaml := range checkYamlList {
		in.CheckBody.CheckYaml = checkYaml

		if matchState, _ := regexp.Match(`kind:\s*StatefulSet`, checkYaml); matchState {
			in.CheckBody.RuleName = "statefulSet"
		} else {
			in.CheckBody.RuleName = "deployment"
		}
		// k check
		if err := k.NormalCheck(in, out); err != nil {
			statusCode = http.StatusBadRequest
		}
	}

	c.JSON(statusCode, gin.H{"file_name": out.FileName, "result": out.Result, "chart": string(in.CheckBody.OriYaml), "msg": out.Message, "hints": out.Hints})
}
