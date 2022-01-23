// normal_check.go
package kcheck

import (
	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	"kcheckApi/constant"
	p "kcheckApi/params"
)

type DWithGracefulTermination struct {
}

func (c *DWithGracefulTermination) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{CheckName: "WithGracefulTermination"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}
	hintsTitle := "建议给containers设置'preStop'做退出准备: \n"
	hints := ""
	hintsCount := 0
	hintsContent :=
		`
spec:
  containers:
  - name: lifecycle-demo-container
    image: nginx
    lifecycle:
      preStop:
        exec:
          command: ["/bin/sh","-c","nginx -s quit"]
			` + "\n"
	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].Lifecycle == nil ||
				deploy.Spec.Template.Spec.Containers[i].Lifecycle.PreStop == nil {
				hintsCount += 1
				hints = hints + " - " + deploy.Spec.Template.Spec.Containers[i].Name + ". \n"
			}

		}
	}
	if hints == "" || hintsCount < len(deploy.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}
	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.WARNING
	return resultMap, nil
}

type DWithLivenessCheck struct {
}

func (c *DWithLivenessCheck) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{CheckName: "WithHealthCheck"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "应该给container设置'LivenessProbe': \n"
	hints := ""
	hintsCount := 0
	hintsContent :=
		`
spec:
  containers:
  - name: lifecycle-demo-container
    image: nginx
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5` + "\n"

	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].LivenessProbe == nil {
				hintsCount += 1
				hints = hints + " - " + deploy.Spec.Template.Spec.Containers[i].Name + ".\n"

			}

		}
	}
	if hints == "" || hintsCount < len(deploy.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}
	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type DWithReadiness struct {
}

func (c *DWithReadiness) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{CheckName: "WithReadiness"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "应该给container设置'ReadinessProbe':  \n"
	hints := ""
	hintsCount := 0
	hintsContent :=
		`
spec:
  containers:
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10 ` + "\n"

	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].ReadinessProbe == nil {
				hintsCount += 1
				hints = hints + " - " + deploy.Spec.Template.Spec.Containers[i].Name + ".\n"

			}

		}
	}
	if hints == "" || hintsCount < len(deploy.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}

	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type DWithLivenessReadinessDelayCheck struct {
}

func (c *DWithLivenessReadinessDelayCheck) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{CheckName: "WithLivenessReadnessDelayCheck"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "Liveness.initialDelaySeconds 应该大于 Readiness.initialDelaySeconds:  \n"
	hints := ""
	hintsContent :=
		`
spec:
  containers:
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10 
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 10
      periodSeconds: 5` + "\n"

	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].ReadinessProbe == nil ||
				deploy.Spec.Template.Spec.Containers[i].LivenessProbe == nil {
				continue
			}
			if deploy.Spec.Template.Spec.Containers[i].ReadinessProbe.InitialDelaySeconds >=
				deploy.Spec.Template.Spec.Containers[i].LivenessProbe.InitialDelaySeconds {
				hints = hints + " - " +
					deploy.Spec.Template.Spec.Containers[i].Name + ".\n"
			}
		}
	}
	if hints == "" {
		return resultMap, nil
	}

	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type DWithResourceRequestAndLimit struct {
}

func (c *DWithResourceRequestAndLimit) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{CheckName: "WithResourceRequestAndLimit"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "应该给容器设置'Resource requests & limits': \n"
	hints := ""
	hintsContent :=
		`
spec:
  containers:
	resources:
	  limits:
		cpu: 3
		memory: 2Gi
	  requests:
		cpu: 2
		memory: 1Gi ` + "\n"

	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].Resources.Requests == nil ||
				deploy.Spec.Template.Spec.Containers[i].Resources.Limits == nil {

				hints = hints + " - " + deploy.Spec.Template.Spec.Containers[i].Name + ".\n"
			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}
	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type WithRollingUpdate struct{}

func (c *WithRollingUpdate) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{CheckName: "WithRollingUpdate"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}
	hints := ""
	if deploy.Spec.Strategy.Type == "RollingUpdate" &&
		deploy.Spec.Strategy.RollingUpdate != nil &&
		deploy.Spec.Strategy.RollingUpdate.MaxSurge != nil &&
		deploy.Spec.Strategy.RollingUpdate.MaxSurge.IntVal >= 50 &&
		deploy.Spec.Strategy.RollingUpdate.MaxUnavailable != nil &&
		deploy.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal >= 50 {

		hints = "'MaxSurge & MaxUnavailable' 应该小于 < 50 " +
			`
spec:
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate` + "\n"

	}

	if deploy.Spec.Strategy.Type == "RollingUpdate" &&
		deploy.Spec.Strategy.RollingUpdate != nil &&
		((deploy.Spec.Strategy.RollingUpdate.MaxSurge != nil &&
			deploy.Spec.Strategy.RollingUpdate.MaxSurge.StrVal >= "50%") ||
			(deploy.Spec.Strategy.RollingUpdate.MaxUnavailable != nil &&
				deploy.Spec.Strategy.RollingUpdate.MaxUnavailable.StrVal >= "50%")) {

		hints = "'MaxSurge & MaxUnavailable' 应该小于 < 50 " +
			`
spec:
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate` + "\n"

	}

	if hints == "" {
		return resultMap, nil
	}
	resultMap.Hints = hints
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type DWithEmptyDirSizeLimit struct {
}

func (c *DWithEmptyDirSizeLimit) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{CheckName: "WithEmptyDirSizeLimit"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "如果需要emptyDir Volumes, 应该设置'Volumes emptyDir limits': \n"
	hints := ""
	hintsContent :=
		`
emptyDir:
  sizeLimit: 4Gi` + "\n"

	if deploy.Spec.Template.Spec.Volumes != nil &&
		len(deploy.Spec.Template.Spec.Volumes) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Volumes); i++ {
			if deploy.Spec.Template.Spec.Volumes[i].EmptyDir != nil &&
				deploy.Spec.Template.Spec.Volumes[i].EmptyDir.SizeLimit == nil {

				hints = hints + " - " + deploy.Spec.Template.Spec.Volumes[i].Name + ".\n"

			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}
	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type DWithTerminationGrace struct{}

func (c *DWithTerminationGrace) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{CheckName: "WithTerminationGrace"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "建议设置'TerminationGracePeriodSeconds': \n"
	hints := ""
	hintsContent :=
		`
space:
	terminationGracePeriodSeconds: 30
	# 如果有使用nginx, 这个时间需要大于keep alive时间, 避免缩减pod时由于连接保持而出现502
` + "\n"

	if deploy.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
		hints = hintsTitle + hintsContent
	}

	if hints == "" {
		return resultMap, nil
	}

	resultMap.Hints = hints
	resultMap.Level = constant.WARNING
	return resultMap, nil
}
