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
	Replicas     int                    `yaml:"replicas"`
	Image        ImageConfig            `yaml:"image"`
	Service      *ServiceConfig         `yaml:"service,omitempty"`
	Environment  map[string]string      `yaml:"environment,omitempty"`
	Ingress      *IngressConfig         `yaml:"ingress,omitempty"`
	Resources    ResourceConfig         `yaml:"resources,omitempty"`
	Volumes      []VolumeConfig         `yaml:"volumes,omitempty"`
	VolumeMounts []VolumeMountConfig    `yaml:"volumeMounts,omitempty"`
}

// ImageConfig represents image configuration
type ImageConfig struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Type   string   `yaml:"type"`
	Ports  []int    `yaml:"ports"`
}

// IngressConfig represents ingress configuration
type IngressConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Host       string   `yaml:"host"`
	Hosts      []string `yaml:"hosts,omitempty"`
	TargetPort int      `yaml:"targetPort"`
}

// ResourceConfig represents resource limits and requests
type ResourceConfig struct {
	Limits   ResourceLimits   `yaml:"limits,omitempty"`
	Requests ResourceRequests `yaml:"requests,omitempty"`
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// ResourceRequests represents resource requests
type ResourceRequests struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// VolumeConfig represents volume configuration
type VolumeConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path,omitempty"`
}

// VolumeMountConfig represents volume mount configuration
type VolumeMountConfig struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

// ServiceValues represents values for a specific service
type ServiceValues struct {
	Name        string     `yaml:"name"`
	HelmValues  *HelmValues `yaml:"values"`
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
	kubeconfig string
}

// NewParser creates a new parser instance
func NewParser(namespace string) *Parser {
	return &Parser{
		namespace: namespace,
	}
}

// NewParserWithKubeconfig creates a new parser instance with kubeconfig
func NewParserWithKubeconfig(namespace, kubeconfig string) *Parser {
	return &Parser{
		namespace: namespace,
		kubeconfig: kubeconfig,
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
	args := []string{"get", resourceType, "-o", "json"}
	
	// Add kubeconfig if specified
	if p.kubeconfig != "" {
		args = append([]string{"--kubeconfig", p.kubeconfig}, args...)
	}
	
	// Add namespace if specified
	if p.namespace != "" {
		args = append(args, "-n", p.namespace)
	}
	
	cmd = exec.Command("kubectl", args...)
	
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

// ParseToMultipleHelmValues converts Kubernetes resources to multiple Helm values (one per service)
func (p *Parser) ParseToMultipleHelmValues(resources []K8sResource) ([]ServiceValues, error) {
	var serviceValues []ServiceValues
	serviceGroups := p.groupResourcesByService(resources)
	for serviceName, group := range serviceGroups {
		values := &HelmValues{
			Environment: make(map[string]string),
		}
		if deployment, exists := group["Deployment"]; exists {
			p.parseDeployment(deployment, values)
		}
		if service, exists := group["Service"]; exists {
			p.parseService(service, values)
		}
		if ingress, exists := group["Ingress"]; exists {
			p.parseIngress(ingress, values)
		}
		// Only set environment if not empty
		if len(values.Environment) == 0 {
			values.Environment = nil
		}
		serviceValues = append(serviceValues, ServiceValues{
			Name:       serviceName,
			HelmValues: values,
		})
	}
	return serviceValues, nil
}

// groupResourcesByService groups resources by their service name
func (p *Parser) groupResourcesByService(resources []K8sResource) map[string]map[string]K8sResource {
	groups := make(map[string]map[string]K8sResource)
	
	for _, resource := range resources {
		// Get resource name from metadata
		var resourceName string
		if metadata, ok := resource.Metadata["name"].(string); ok {
			resourceName = metadata
		} else {
			continue
		}
		
		// Initialize group if it doesn't exist
		if groups[resourceName] == nil {
			groups[resourceName] = make(map[string]K8sResource)
		}
		
		// Add resource to its group
		groups[resourceName][resource.Kind] = resource
	}
	
	return groups
}

// ParseToHelmValues converts Kubernetes resources to Helm values structure (legacy method)
func (p *Parser) ParseToHelmValues(resources []K8sResource) (*HelmValues, error) {
	values := &HelmValues{
		Environment: make(map[string]string),
	}
	
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
	// Set default values
	values.Replicas = 1
	values.Image.PullPolicy = "IfNotPresent"
	
	// Parse replicas
	if replicas, ok := deployment.Spec["replicas"].(float64); ok {
		values.Replicas = int(replicas)
	}
	
	if spec, ok := deployment.Spec["template"].(map[string]interface{}); ok {
		if podSpec, ok := spec["spec"].(map[string]interface{}); ok {
			// Parse containers
			if containers, ok := podSpec["containers"].([]interface{}); ok && len(containers) > 0 {
				if container, ok := containers[0].(map[string]interface{}); ok {
					// Parse image
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
					
					// Parse imagePullPolicy
					if pullPolicy, ok := container["imagePullPolicy"].(string); ok {
						values.Image.PullPolicy = pullPolicy
					}
					
					// Parse environment variables
					if env, ok := container["env"].([]interface{}); ok {
						for _, envVar := range env {
							if envMap, ok := envVar.(map[string]interface{}); ok {
								if name, ok := envMap["name"].(string); ok {
									if value, ok := envMap["value"].(string); ok {
										values.Environment[name] = value
									}
								}
							}
						}
					}
					
					// Parse resources
					if resources, ok := container["resources"].(map[string]interface{}); ok {
						values.Resources = p.parseResources(resources)
					}
					
					// Parse volume mounts
					if volumeMounts, ok := container["volumeMounts"].([]interface{}); ok {
						values.VolumeMounts = p.parseVolumeMounts(volumeMounts)
					}
				}
			}
			
			// Parse volumes
			if volumes, ok := podSpec["volumes"].([]interface{}); ok {
				values.Volumes = p.parseVolumes(volumes)
			}
		}
	}
}

// parseService extracts service information
func (p *Parser) parseService(service K8sResource, values *HelmValues) {
	cfg := &ServiceConfig{}
	if serviceType, ok := service.Spec["type"].(string); ok {
		cfg.Type = serviceType
	}
	if ports, ok := service.Spec["ports"].([]interface{}); ok {
		for _, port := range ports {
			if portMap, ok := port.(map[string]interface{}); ok {
				if portNum, ok := portMap["port"].(float64); ok {
					cfg.Ports = append(cfg.Ports, int(portNum))
				}
			}
		}
	}
	if cfg.Type != "" && len(cfg.Ports) > 0 {
		values.Service = cfg
	}
}

// parseIngress extracts ingress information
func (p *Parser) parseIngress(ingress K8sResource, values *HelmValues) {
	cfg := &IngressConfig{}
	cfg.Enabled = true
	if rules, ok := ingress.Spec["rules"].([]interface{}); ok {
		for _, rule := range rules {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				if host, ok := ruleMap["host"].(string); ok {
					cfg.Host = host
					cfg.Hosts = append(cfg.Hosts, host)
				}
				if http, ok := ruleMap["http"].(map[string]interface{}); ok {
					if paths, ok := http["paths"].([]interface{}); ok && len(paths) > 0 {
						if path, ok := paths[0].(map[string]interface{}); ok {
							if backend, ok := path["backend"].(map[string]interface{}); ok {
								if service, ok := backend["service"].(map[string]interface{}); ok {
									if port, ok := service["port"].(map[string]interface{}); ok {
										if targetPort, ok := port["number"].(float64); ok {
											cfg.TargetPort = int(targetPort)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if cfg.Host != "" || len(cfg.Hosts) > 0 || cfg.TargetPort != 0 {
		values.Ingress = cfg
	}
}

// parseResources extracts resource configuration
func (p *Parser) parseResources(resources map[string]interface{}) ResourceConfig {
	config := ResourceConfig{}
	
	if limits, ok := resources["limits"].(map[string]interface{}); ok {
		if cpu, ok := limits["cpu"].(string); ok {
			config.Limits.CPU = cpu
		}
		if memory, ok := limits["memory"].(string); ok {
			config.Limits.Memory = memory
		}
	}
	
	if requests, ok := resources["requests"].(map[string]interface{}); ok {
		if cpu, ok := requests["cpu"].(string); ok {
			config.Requests.CPU = cpu
		}
		if memory, ok := requests["memory"].(string); ok {
			config.Requests.Memory = memory
		}
	}
	
	return config
}

// parseVolumes extracts volume configuration
func (p *Parser) parseVolumes(volumes []interface{}) []VolumeConfig {
	var configs []VolumeConfig
	
	for _, volume := range volumes {
		if volumeMap, ok := volume.(map[string]interface{}); ok {
			config := VolumeConfig{}
			
			if name, ok := volumeMap["name"].(string); ok {
				config.Name = name
			}
			
			// Parse different volume types
			if _, ok := volumeMap["configMap"].(map[string]interface{}); ok {
				config.Type = "configMap"
			} else if _, ok := volumeMap["secret"].(map[string]interface{}); ok {
				config.Type = "secret"
			} else if _, ok := volumeMap["emptyDir"].(map[string]interface{}); ok {
				config.Type = "emptyDir"
			} else if _, ok := volumeMap["persistentVolumeClaim"].(map[string]interface{}); ok {
				config.Type = "persistentVolumeClaim"
			}
			
			configs = append(configs, config)
		}
	}
	
	return configs
}

// parseVolumeMounts extracts volume mount configuration
func (p *Parser) parseVolumeMounts(volumeMounts []interface{}) []VolumeMountConfig {
	var configs []VolumeMountConfig
	
	for _, mount := range volumeMounts {
		if mountMap, ok := mount.(map[string]interface{}); ok {
			config := VolumeMountConfig{}
			
			if name, ok := mountMap["name"].(string); ok {
				config.Name = name
			}
			
			if mountPath, ok := mountMap["mountPath"].(string); ok {
				config.MountPath = mountPath
			}
			
			configs = append(configs, config)
		}
	}
	
	return configs
}

// GenerateMultipleYAML converts multiple Helm values to YAML format
func (p *Parser) GenerateMultipleYAML(serviceValues []ServiceValues) (string, error) {
	var yamlOutput strings.Builder
	
	for i, sv := range serviceValues {
		if i > 0 {
			yamlOutput.WriteString("\n---\n\n")
		}
		
		// Create a map with service name as key and values as value
		serviceMap := map[string]*HelmValues{
			sv.Name: sv.HelmValues,
		}
		
		yamlData, err := yaml.Marshal(serviceMap)
		if err != nil {
			return "", fmt.Errorf("failed to marshal service %s to YAML: %w", sv.Name, err)
		}
		
		yamlOutput.Write(yamlData)
	}
	
	return yamlOutput.String(), nil
}

// GenerateYAML converts Helm values to YAML format
func (p *Parser) GenerateYAML(values *HelmValues) (string, error) {
	yamlData, err := yaml.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	
	return string(yamlData), nil
} 