package kcheck

import (
	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	p "kcheckApi/params"
)

type SWithGracefulTermination struct {
}

func (c *SWithGracefulTermination) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{"", "WithGracefulTermination"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}
	hintsTitle := "'preStop' should be set for a graceful termination containers: \n"
	hints := ""
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

				hints = hints + stateful.Spec.Template.Spec.Containers[i].Name + ". \n"

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

type SWithLivenessCheck struct {
}

func (c *SWithLivenessCheck) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.Deployment{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{"", "WithHealthCheck"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "'LivenessProbe' should be set for container: \n"
	hints := ""
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

				hints = hints + stateful.Spec.Template.Spec.Containers[i].Name + ".\n"

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

type SWithReadiness struct {
}

func (c *SWithReadiness) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)
	resultMap := p.HintsMap{"", "WithReadiness"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "It is nice to have 'ReadinessProbe' setting for container:  \n"
	hints := ""
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

				hints = hints +
					stateful.Spec.Template.Spec.Containers[i].Name + "."
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

type SWithResourceRequestAndLimit struct {
}

func (c *SWithResourceRequestAndLimit) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)

	resultMap := p.HintsMap{"", "WithResourceRequestAndLimit"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "'Resource requests & limits' should be set for container: \n"
	hints := ""
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
			if stateful.Spec.Template.Spec.Containers[i].Resources.Requests == nil ||
				stateful.Spec.Template.Spec.Containers[i].Resources.Limits == nil {
				hints = hints +
					stateful.Spec.Template.Spec.Containers[i].Name + ".\n"
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

type SWithEmptyDirSizeLimit struct {
}

func (c *SWithEmptyDirSizeLimit) Check(data []byte) (p.HintsMap, error) {
	stateful := &v1.StatefulSet{}
	err := yaml.Unmarshal(data, stateful)

	resultMap := p.HintsMap{"", "DWithEmptyDirSizeLimit"}
	if err != nil {
		return resultMap, err
	}
	if stateful.Kind != "StatefulSet" {
		return resultMap, nil
	}

	hintsTitle := "'Volumes emptyDir limits' should be set for emptyDir: \n"
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

				hints = hints + stateful.Spec.Template.Spec.Volumes[i].Name + ".\n"

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
