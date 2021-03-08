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
	hints := ""
	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].Lifecycle == nil ||
				stateful.Spec.Template.Spec.Containers[i].Lifecycle.PreStop == nil {

				hints = "'preStop' should be set for a graceful termination [container: " +
					stateful.Spec.Template.Spec.Containers[i].Name + "]." +

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

			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}
	resultMap.Hints = hints
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
	hints := ""
	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].LivenessProbe == nil {

				hints = "'LivenessProbe' should be set for container: " +
					stateful.Spec.Template.Spec.Containers[i].Name + "." +

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

			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}
	resultMap.Hints = hints
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
	hints := ""
	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].ReadinessProbe == nil {

				hints = "It is nice to have 'ReadinessProbe' setting for container: " +
					stateful.Spec.Template.Spec.Containers[i].Name + "." +

					`
spec:
  containers:
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10 `

			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}

	resultMap.Hints = hints
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
	hints := ""
	if stateful.Spec.Template.Spec.Containers != nil &&
		len(stateful.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(stateful.Spec.Template.Spec.Containers); i++ {
			if stateful.Spec.Template.Spec.Containers[i].Resources.Requests == nil ||
				stateful.Spec.Template.Spec.Containers[i].Resources.Limits == nil {

				hints = "'Resource requests & limits' should be set for container: " +
					stateful.Spec.Template.Spec.Containers[i].Name + "." +

					`
resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"` + "\n"

			}

		}
	}
	if hints == "" {
		return resultMap, nil
	}
	resultMap.Hints = hints
	return resultMap, nil
}
