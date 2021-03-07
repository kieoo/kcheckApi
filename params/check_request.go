package params

import "mime/multipart"

type CRequest struct {
	MultipartFile []*multipart.FileHeader
	CheckBody     CheckBody
}

type CheckBody struct {
	OriYaml    []byte `json:"ori_yaml" `
	CheckYaml  []byte `json:"check_yaml" `
	RuleConfig string `json:"rule_config"`
	RuleName   string `json:"rule_name"`
}
