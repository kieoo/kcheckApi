package helm

import (
	"fmt"
	"kcheckApi/model"
	p "kcheckApi/params"
)

func HelmChange(in *p.CRequest, out *p.CResponse) error {

	// save helm chart
	up := Upload{}
	if err := up.UploadInit(); err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("upload error : %s", err)
		return err
	}
	defer up.UploadClose()

	if err := up.UploadSave(in.MultipartFile); err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("upload error : %s", err)
		return err
	}

	// chart template to deploy
	chartDir := up.TmpDir + "/" + string(up.ChartName)
	template, err := Chart2Deploy(chartDir)

	if err != nil {
		out.Result = model.Fail
		out.Message = fmt.Sprintf("chage error : %s", string(template))
		return err
	}

	out.Result = model.PASS
	in.CheckBody.OriYaml = template
	out.Message = ""

	return nil
}
