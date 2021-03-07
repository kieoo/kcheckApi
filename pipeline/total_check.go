package pipeline

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	h "kcheckApi/helm"
	"kcheckApi/model"
	p "kcheckApi/params"
	"kcheckApi/util"
	"log"
	"net/http"
	"time"
)

func TotalCheckXML(c *gin.Context) {

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

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+out.FileName+".xml")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	testSuite := &model.TestSuite{}
	testSuite.Name = out.FileName

	// helm check
	HTestCase := &model.TestCase{}
	testSuite.TCs = append(testSuite.TCs, HTestCase)
	// 开始检查
	err := h.HelmChange(in, out)

	// result
	HTestCase.ClassN = out.FileName
	HTestCase.Name = "helm check"
	HTestCase.SystemOut = model.SystemOUTStart + out.Message + model.SystemOUTEnd

	if err != nil {
		HTestCase.Status = model.Fail
		c.Data(http.StatusOK, "text/xml", outXml(testSuite))
		return
	}

	HTestCase.Status = model.PASS

	util.CleanOriYaml(in.CheckBody.OriYaml)
	// kchecker

}

func outXml(testSuite *model.TestSuite) []byte {
	testSuite.Failures = 0
	testSuite.Tests = int32(len(testSuite.TCs))

	var tc model.TestCase
	for tc = range testSuite.TCs {
		if tc.Status != model.PASS {
			testSuite.Failures++
		}
	}
	// time
	t := time.Now()
	str := t.Format("2006-01-02T15:04:05")
	testSuite.Timestamp = str

	// struct to xml
	oXml, err := xml.MarshalIndent(testSuite, "", "'")
	log.Printf(fmt.Sprintf("%s", err))

	return oXml
}
