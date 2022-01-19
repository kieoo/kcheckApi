// normal_check.go
package kcheck

import (
	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	p "kcheckApi/params"
)

type DWithGracefulTermination struct {
}

func (c *DWithGracefulTermination) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{"", "WithGracefulTermination"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}
	hintsTitle := "'preStop' should be set for a graceful termination containers: \n"
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
	return resultMap, nil
}

type DWithLivenessCheck struct {
}

func (c *DWithLivenessCheck) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{"", "WithHealthCheck"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "'LivenessProbe' should be set for container: \n"
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
	return resultMap, nil
}

type DWithReadiness struct {
}

func (c *DWithReadiness) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{"", "WithReadiness"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "It is nice to have 'ReadinessProbe' setting for container:  \n"
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
	return resultMap, nil
}

type DWithLivenessReadinessDelayCheck struct {
}

func (c *DWithLivenessReadinessDelayCheck) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	resultMap := p.HintsMap{"", "WithLivenessReadnessDelayCheck"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "Liveness.initialDelaySeconds should bigger then Readiness.initialDelaySeconds  \n"
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
	return resultMap, nil
}

type DWithResourceRequestAndLimit struct {
}

func (c *DWithResourceRequestAndLimit) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{"", "WithResourceRequestAndLimit"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "'Resource requests & limits' should be set for container: \n"
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
	return resultMap, nil
}

type WithRollingUpdate struct{}

func (c *WithRollingUpdate) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{"", "WithRollingUpdate"}
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

		hints = "'MaxSurge & MaxUnavailable' should be set and < 50 " +
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

		hints = "'MaxSurge & MaxUnavailable' should be set and < 50 " +
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
	return resultMap, nil
}

type DWithEmptyDirSizeLimit struct {
}

func (c *DWithEmptyDirSizeLimit) Check(data []byte) (p.HintsMap, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)

	resultMap := p.HintsMap{"", "WithEmptyDirSizeLimit"}
	if err != nil {
		return resultMap, err
	}
	if deploy.Kind != "Deployment" {
		return resultMap, nil
	}

	hintsTitle := "'Volumes emptyDir limits' should be set for emptyDir: \n"
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
	return resultMap, nil
}
