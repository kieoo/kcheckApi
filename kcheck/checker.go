package kcheck

import (
	"fmt"
	"io/ioutil"
	"kcheckApi/model"
	p "kcheckApi/params"
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
	Check(data []byte) (p.HintsMap, error)
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

type SaveError struct {
	S string
}

func (e *SaveError) Error() string {
	return e.S
}

func paramChecker(param *string, prompt string) *SaveError {

	isOk := isStringParamValid(param)

	if !isOk {
		return &SaveError{prompt}
	}

	return nil
}

func NormalCheck(in *p.CRequest, out *p.CResponse) error {

	checkBody := p.CheckBody{}

	checkBody = in.CheckBody

	srcYaml := string(checkBody.CheckYaml)
	ruleConfig := checkBody.RuleConfig
	ruleName := checkBody.RuleName

	out.Result = model.PASS

	if len(ruleConfig) <= 0 {
		ruleConfig = "default.yaml"
	}

	if err := paramChecker(&srcYaml, "Set the yaml needing"); err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("%s", err)
		return err
	}

	if err := paramChecker(&ruleConfig, "Set the rule file for checking"); err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("%s", err)
		return err
	}

	if err := paramChecker(&ruleName, "Set the rule name for checking"); err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("%s", err)
		return err
	}

	ruleSet, _, err := ParserRuleSetConfig("conf/" + ruleConfig)

	if err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("Failed to load the file '%s'.", ruleConfig)
		return &SaveError{out.Message}
	}

	var rule *Rule
	for _, r := range ruleSet {
		if r.Name == ruleName {
			rule = r
		}
	}

	if rule == nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("Could not find the checking rule '%s'.", ruleName)
		return &SaveError{out.Message}
	}

	var resultMap []p.HintsMap

	for _, check := range rule.Checkers {
		hintMap, err := check.Check([]byte(srcYaml))
		if err != nil {
			out.Result = model.Fail
			out.Message = fmt.Sprintf("Checking erro %s", err)
			return &SaveError{out.Message}
		}
		if len(hintMap.Hints) > 0 {
			out.Result = model.Fail
			resultMap = append(resultMap, hintMap)
		} else {
			resultMap = append(resultMap, hintMap)
		}
	}

	// jsonResultMap, err := json.Marshal(resultMap)

	// 想要struct的字段能被marshal, 首字母必须大写+++

	out.Hints = resultMap
	return nil

}
