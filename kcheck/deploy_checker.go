package kcheck

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

type Corrector interface {
	// Correct is to correct the config
	Correct(org []byte) ([]byte, error)
}

func loadDataFromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Checker is to check the K8S deployment/configuration with the different rules
type Checker interface {
	// Check is to check the file and return the suggestions
	Check(data []byte) (HintsMap, error)
}

type checkBody struct {
	OriYaml    string `json:"ori_yaml"`
	RuleConfig string `json:"rule_config"`
	RuleName   string `json:"rule_name"`
}

// Rule is composed of a set of checkers
type Rule struct {
	// Name is the name of the rule
	Name string
	// Checkers is the checkers relating to the rule's check items
	Checkers []Checker
}

func isStringParamValid(param *string) bool {
	if param == nil || *param == "" {
		return false
	}
	return true
}

func paramChecker(c *gin.Context, param *string, prompt string) {

	isOk := isStringParamValid(param)

	if !isOk {
		c.String(http.StatusBadRequest, *param)
	}
}

func NormalCheck(c *gin.Context) {
	p := checkBody{}
	if err := c.BindJSON(&p); err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusBadRequest, "checkBody error")
	}

	srcYaml := p.OriYaml
	ruleConfig := p.RuleConfig
	ruleName := p.RuleName
	paramChecker(c, &srcYaml, "Set the yaml needing")
	paramChecker(c, &ruleConfig, "Set the rule file for checking")
	paramChecker(c, &ruleName, "Set the rule name for checking")

	ruleSet, _, err := ParserRuleSetConfig("conf/" + ruleConfig)

	if err != nil {
		c.String(http.StatusBadRequest,
			"Failed to load the file '%s'.", ruleConfig)
	}

	var rule *Rule
	for _, r := range ruleSet {
		if r.Name == ruleName {
			rule = r
		}
	}

	if rule == nil {
		c.String(http.StatusBadRequest,
			"Could not find the checking rule '%s'.", ruleName)
	}

	var resultMap []HintsMap

	for _, check := range rule.Checkers {
		hintMap, err := check.Check([]byte(srcYaml))
		if err != nil {
			c.String(http.StatusBadRequest,
				"Checking error")
		}
		if len(hintMap.Hints) > 0 {
			resultMap = append(resultMap, hintMap)
		}
	}

	// jsonResultMap, err := json.Marshal(resultMap)

	if err != nil {
		c.String(http.StatusBadRequest, "Checking error")
	}
	// 想要struct的字段能被marshal, 首字母必须大写+++
	c.JSON(http.StatusOK, resultMap)

}
