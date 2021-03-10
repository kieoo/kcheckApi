package pipeline

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	h "kcheckApi/helm"
	k "kcheckApi/kcheck"
	"kcheckApi/model"
	p "kcheckApi/params"
	"kcheckApi/util"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"time"
)

func TotalCheckXML(c *gin.Context) {

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

	out.FileName = filepath.Dir(files[0].Filename)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+out.FileName+".xml")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	testSuite := &model.TestSuite{}
	testSuite.Name = out.FileName

	// helm check >>>>>>>>>>>>>>>>>>>>>>
	HTestCase := &model.TestCase{}
	testSuite.TestCases = append(testSuite.TestCases, HTestCase)
	// 开始检查
	err := h.HelmChange(in, out)

	// result
	HTestCase.ClassN = out.FileName
	HTestCase.Name = "helm check"
	HTestCase.SystemOut.Out = out.Message

	if err != nil {
		HTestCase.Status = model.Fail
		c.Data(http.StatusOK, "text/xml", outXml(testSuite))
		return
	}

	HTestCase.Status = model.PASS

	// kchecker >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	checkYamlList := util.CleanOriYaml(in.CheckBody.OriYaml)

	for fileName, checkYaml := range checkYamlList {
		in.CheckBody.CheckYaml = checkYaml
		out.Hints = []p.HintsMap{}

		if matchState, _ := regexp.Match(`kind:\s*StatefulSet`, checkYaml); matchState {
			in.CheckBody.RuleName = "statefulSet"
		} else {
			in.CheckBody.RuleName = "deployment"
		}

		// 开始检查
		err := k.NormalCheck(in, out)
		// result
		//声明testCase
		if err != nil {
			TestCase := &model.TestCase{}
			TestCase.ClassN = out.FileName
			TestCase.Name = fileName + " normal check"
			TestCase.SystemOut.Out = out.Message
			TestCase.Status = model.Fail
			testSuite.TestCases = append(testSuite.TestCases, HTestCase)
			continue
		}

		if len(out.Hints) > 0 {
			for _, hint := range out.Hints {
				TestCase := &model.TestCase{}
				TestCase.ClassN = out.FileName
				TestCase.Name = fileName + " normal check-" + string(hint.CheckName)
				TestCase.SystemOut.Out = hint.Hints
				TestCase.Status = model.Fail
				testSuite.TestCases = append(testSuite.TestCases, HTestCase)
			}
			continue
		}

		// else
		TestCase := &model.TestCase{}
		TestCase.ClassN = out.FileName
		TestCase.Name = fileName + " normal check"
		TestCase.SystemOut.Out = out.Message
		TestCase.Status = model.PASS
		testSuite.TestCases = append(testSuite.TestCases, HTestCase)
	}

	// finish
	c.Data(http.StatusOK, "text/xml", outXml(testSuite))

	return
}

func outXml(testSuite *model.TestSuite) []byte {
	testSuite.Failures = 0
	testSuite.Tests = int32(len(testSuite.TestCases))

	var tc *model.TestCase
	for _, tc = range testSuite.TestCases {
		if tc.Status != model.PASS {
			testSuite.Failures++
		}
	}
	// time
	t := time.Now()
	str := t.Format("2006-01-02T15:04:05")
	testSuite.Timestamp = str

	// struct to xml
	oXml, err := xml.MarshalIndent(testSuite, "", "  ")
	log.Printf(fmt.Sprintf("%s", err))

	return oXml
}
