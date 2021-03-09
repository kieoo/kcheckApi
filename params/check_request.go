package params

import "mime/multipart"

type CRequest struct {
	MultipartFile []*multipart.FileHeader
	CheckBody     CheckBody
}

type CheckBody struct {
	OriYamlString string `json:"ori_yaml"`
	OriYaml       []byte `json:"ori_yaml_byte" `
	CheckYaml     []byte `json:"check_yaml"`
	RuleConfig    string `json:"rule_config"`
	RuleName      string `json:"rule_name"`
}
