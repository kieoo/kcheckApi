package kcheck

import (
	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	"kcheckApi/constant"
	p "kcheckApi/params"
)

type SWithGracefulTermination struct {
}

func (c *SWithGracefulTermination) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{CheckName: "WithGracefulTermination"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
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

	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].Lifecycle == nil ||
				stateful.Spec.Template.Spec.Containers[i].Lifecycle.PreStop == nil {
				hintsCount += 1
				hints = hints + " - " + stateful.Spec.Template.Spec.Containers[i].Name + ". \n"

			}
		}
	}
	if hints == "" || hintsCount < len(stateful.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}

	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.WARNING
	return resultMap, nil
}

type SWithLivenessCheck struct {
}

func (c *SWithLivenessCheck) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.Deployment{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{CheckName: "WithHealthCheck"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := " \n"
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

	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].LivenessProbe == nil {
				hintsCount += 1
				hints = hints + " - " + stateful.Spec.Template.Spec.Containers[i].Name + ".\n"

			}

		}
	}
	if hints == "" || hintsCount < len(stateful.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}

	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type SWithReadiness struct {
}

func (c *SWithReadiness) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{CheckName: "WithReadiness"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
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

	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].ReadinessProbe == nil {
				hintsCount += 1
				hints = hints + " - " +
					stateful.Spec.Template.Spec.Containers[i].Name + ".\n"
			}
		}
	}
	if hints == "" || hintsCount < len(stateful.Spec.Template.Spec.Containers) {
		return resultMap, nil
	}

	hintsAll := hintsTitle + hints + hintsContent

	resultMap.Hints = hintsAll
	resultMap.Level = constant.ERROR
	return resultMap, nil
}

type SWithLivenessReadinessDelayCheck struct {
}

func (c *SWithLivenessReadinessDelayCheck) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{CheckName: "WithLivenessReadnessDelayCheck"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
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

	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].ReadinessProbe == nil ||
				stateful.Spec.Template.Spec.Containers[i].LivenessProbe == nil {
				continue
			}
			if stateful.Spec.Template.Spec.Containers[i].ReadinessProbe.InitialDelaySeconds >=
				stateful.Spec.Template.Spec.Containers[i].LivenessProbe.InitialDelaySeconds {
				hints = hints + " - " +
					stateful.Spec.Template.Spec.Containers[i].Name + ".\n"
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

type SWithResourceRequestAndLimit struct {
}

func (c *SWithResourceRequestAndLimit) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)

	resultMap := p.HintsMap{CheckName: "WithResourceRequestAndLimit"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
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

	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].Resources.Requests == nil ||
				stateful.Spec.Template.Spec.Containers[i].Resources.Limits == nil {
				hints = hints + " - " +
					stateful.Spec.Template.Spec.Containers[i].Name + ".\n"
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

type SWithEmptyDirSizeLimit struct {
}

func (c *SWithEmptyDirSizeLimit) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)

	resultMap := p.HintsMap{CheckName: "WithEmptyDirSizeLimit"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "如果需要emptyDir Volumes, 应该设置'Volumes emptyDir limits': \n"
	hints := ""
	hintsContent :=
		`
emptyDir:
  sizeLimit: 4Gi` + "\n"

	if stateful.Spec.Template.Spec.Volumes != nil &&
		len(stateful.Spec.Template.Spec.Volumes) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Volumes); i++ {
			if stateful.Spec.Template.Spec.Volumes[i].EmptyDir != nil &&
				stateful.Spec.Template.Spec.Volumes[i].EmptyDir.SizeLimit == nil {

				hints = hints + " - " + stateful.Spec.Template.Spec.Volumes[i].Name + ".\n"

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

type SWithTerminationGrace struct{}

func (c *SWithTerminationGrace) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)

	resultMap := p.HintsMap{CheckName: "WithTerminationGrace"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
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

	if stateful.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
		hints = hintsTitle + hintsContent
	}

	if hints == "" {
		return resultMap, nil
	}

	resultMap.Hints = hints
	resultMap.Level = constant.WARNING
	return resultMap, nil
}
