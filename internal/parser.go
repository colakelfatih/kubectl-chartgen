package internal

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"gopkg.in/yaml.v3"
)

// HelmValues represents the structure of Helm values.yaml
type HelmValues struct {
	Image struct {
		Repository string `yaml:"repository"`
		Tag        string `yaml:"tag"`
	} `yaml:"image"`
	Service struct {
		Type string `yaml:"type"`
		Port int    `yaml:"port"`
	} `yaml:"service"`
	Ingress struct {
		Enabled bool     `yaml:"enabled"`
		Hosts   []string `yaml:"hosts"`
	} `yaml:"ingress"`
	ReplicaCount int `yaml:"replicaCount"`
}

// K8sResource represents a generic Kubernetes resource
type K8sResource struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       map[string]interface{} `json:"spec"`
}

// Parser handles the conversion from Kubernetes resources to Helm values
type Parser struct {
	namespace string
}

// NewParser creates a new parser instance
func NewParser(namespace string) *Parser {
	return &Parser{
		namespace: namespace,
	}
}

// GetK8sResources fetches Kubernetes resources using kubectl
func (p *Parser) GetK8sResources() ([]K8sResource, error) {
	var resources []K8sResource
	
	// Get deployments
	deployments, err := p.getResources("deployment")
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}
	resources = append(resources, deployments...)
	
	// Get services
	services, err := p.getResources("service")
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	resources = append(resources, services...)
	
	// Get ingresses
	ingresses, err := p.getResources("ingress")
	if err != nil {
		return nil, fmt.Errorf("failed to get ingresses: %w", err)
	}
	resources = append(resources, ingresses...)
	
	return resources, nil
}

// getResources fetches a specific type of Kubernetes resource
func (p *Parser) getResources(resourceType string) ([]K8sResource, error) {
	var cmd *exec.Cmd
	if p.namespace != "" {
		cmd = exec.Command("kubectl", "get", resourceType, "-n", p.namespace, "-o", "json")
	} else {
		cmd = exec.Command("kubectl", "get", resourceType, "-o", "json")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Items []K8sResource `json:"items"`
	}
	
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}
	
	return result.Items, nil
}

// ParseToHelmValues converts Kubernetes resources to Helm values structure
func (p *Parser) ParseToHelmValues(resources []K8sResource) (*HelmValues, error) {
	values := &HelmValues{}
	
	for _, resource := range resources {
		switch resource.Kind {
		case "Deployment":
			p.parseDeployment(resource, values)
		case "Service":
			p.parseService(resource, values)
		case "Ingress":
			p.parseIngress(resource, values)
		}
	}
	
	return values, nil
}

// parseDeployment extracts deployment information
func (p *Parser) parseDeployment(deployment K8sResource, values *HelmValues) {
	if spec, ok := deployment.Spec["template"].(map[string]interface{}); ok {
		if podSpec, ok := spec["spec"].(map[string]interface{}); ok {
			if containers, ok := podSpec["containers"].([]interface{}); ok && len(containers) > 0 {
				if container, ok := containers[0].(map[string]interface{}); ok {
					if image, ok := container["image"].(string); ok {
						parts := strings.Split(image, ":")
						if len(parts) >= 2 {
							values.Image.Repository = parts[0]
							values.Image.Tag = parts[1]
						} else {
							values.Image.Repository = image
							values.Image.Tag = "latest"
						}
					}
				}
			}
		}
	}
	
	if replicas, ok := deployment.Spec["replicas"].(float64); ok {
		values.ReplicaCount = int(replicas)
	} else {
		values.ReplicaCount = 1
	}
}

// parseService extracts service information
func (p *Parser) parseService(service K8sResource, values *HelmValues) {
	if serviceType, ok := service.Spec["type"].(string); ok {
		values.Service.Type = serviceType
	}
	
	if ports, ok := service.Spec["ports"].([]interface{}); ok && len(ports) > 0 {
		if port, ok := ports[0].(map[string]interface{}); ok {
			if portNum, ok := port["port"].(float64); ok {
				values.Service.Port = int(portNum)
			}
		}
	}
}

// parseIngress extracts ingress information
func (p *Parser) parseIngress(ingress K8sResource, values *HelmValues) {
	values.Ingress.Enabled = true
	
	if rules, ok := ingress.Spec["rules"].([]interface{}); ok {
		for _, rule := range rules {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				if host, ok := ruleMap["host"].(string); ok {
					values.Ingress.Hosts = append(values.Ingress.Hosts, host)
				}
			}
		}
	}
}

// GenerateYAML converts Helm values to YAML format
func (p *Parser) GenerateYAML(values *HelmValues) (string, error) {
	yamlData, err := yaml.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	
	return string(yamlData), nil
} 