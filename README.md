# kubectl-chartgen 🚀

A powerful kubectl plugin that automatically generates Helm `values.yaml` files from existing Kubernetes resources. Perfect for migrating existing deployments to Helm charts or creating Helm templates from running applications.

## ✨ Features

- 🔍 **Automatic Resource Detection**: Automatically discovers Deployments, Services, and Ingresses
- 🎯 **Helm Values Generation**: Converts Kubernetes resources to Helm-compatible `values.yaml` format
- 📁 **Namespace Support**: Target specific namespaces or use current context
- 💾 **Flexible Output**: Save to file or output to stdout
- 🛡️ **Error Handling**: Comprehensive error handling and user feedback
- ⚡ **Fast & Lightweight**: Built with Go for optimal performance

## 🚀 Installation

### Prerequisites

- Go 1.23 or higher
- kubectl configured with access to a Kubernetes cluster
- Access to a Kubernetes cluster with resources to analyze

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/kubectl-chartgen.git
cd kubectl-chartgen

# Build the binary
go build -o chartgen main.go

# Make it executable
chmod +x chartgen

# Move to a directory in your PATH (optional)
sudo mv chartgen /usr/local/bin/
```

### Install as kubectl Plugin

```bash
# Build the plugin
go build -o kubectl-chartgen main.go

# Make it executable
chmod +x kubectl-chartgen

# Move to kubectl plugins directory
mkdir -p ~/.kube/plugins/chartgen
mv kubectl-chartgen /usr/local/bin/
```

## 📖 Usage

### Basic Usage

```bash
# Generate values.yaml from current namespace
go run main.go generate

# Or if installed as binary
kubectl chartgen generate
```

### Advanced Usage

```bash
# Generate from specific namespace
kubectl chartgen generate -n my-app-namespace

# Output to specific file
kubectl chartgen generate -o my-helm-values.yaml

# Output to stdout
kubectl chartgen generate -o -

# Combine options
kubectl chartgen generate -n production -o prod-values.yaml
```

### Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Target namespace (default: current) | current context |
| `--output` | `-o` | Output file path (use `-` for stdout) | `values.yaml` |
| `--help` | `-h` | Show help information | - |

## 📋 Examples

### Example 1: Basic Generation

```bash
# Generate values.yaml from current namespace
kubectl chartgen generate
```

**Output:**
```
Generating Helm values.yaml from Kubernetes resources...
Fetching Kubernetes resources...
Found 3 resources
Converting to Helm values structure...
Generated values for 2 services
Generating YAML output...
Helm values written to: values.yaml
```

### Example 2: Generate from Specific Namespace

```bash
# Generate from 'production' namespace
kubectl chartgen generate -n production -o prod-values.yaml
```

### Example 3: Preview Output

```bash
# Preview the generated values without saving
kubectl chartgen generate -o -
```

## 🎯 Example Output Format

The tool generates Helm-compatible `values.yaml` files with the following structure:

```yaml
# Frontend Application
frontend:
  replicas: 2
  image:
    repository: myapp/frontend
    tag: "1.2.3"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    ports:
      - 80
      - 443
  environment:
    NODE_ENV: production
    API_URL: https://api.example.com
    CDN_URL: https://cdn.example.com
    GA_TRACKING_ID: GA-123456789
  ingress:
    enabled: true
    host: app.example.com
    hosts:
      - app.example.com
      - www.example.com
    targetPort: 80
  resources:
    limits:
      cpu: "500m"
      memory: "512Mi"
    requests:
      cpu: "250m"
      memory: "256Mi"
  volumes:
    - name: config-volume
      type: configMap
  volumeMounts:
    - name: config-volume
      mountPath: /app/config

---

# Backend API Service
backend:
  replicas: 3
  image:
    repository: myapp/backend
    tag: "2.1.0"
    pullPolicy: IfNotPresent
  service:
    type: ClusterIP
    ports:
      - 8080
  environment:
    DATABASE_URL: postgresql://postgres:5432/myapp
    REDIS_URL: redis://redis:6379
    JWT_SECRET: your-super-secret-jwt-key

---

# Worker Service (no service/ingress)
worker:
  replicas: 2
  image:
    repository: myapp/worker
    tag: "2.1.0"
    pullPolicy: IfNotPresent
  environment:
    DATABASE_URL: postgresql://postgres:5432/myapp
    REDIS_URL: redis://redis:6379
    QUEUE_NAME: email-queue
```

### 📊 Supported Fields

| Field | Description | When Included |
|-------|-------------|---------------|
| `replicas` | Number of pod replicas | Always |
| `image.repository` | Container image repository | Always |
| `image.tag` | Container image tag | Always |
| `image.pullPolicy` | Image pull policy | Always |
| `service.type` | Kubernetes service type | Only if service exists |
| `service.ports` | Array of service ports | Only if service exists |
| `environment` | Environment variables | Only if environment variables exist |
| `ingress.enabled` | Ingress enabled status | Only if ingress exists |
| `ingress.host` | Primary ingress host | Only if ingress exists |
| `ingress.hosts` | Array of ingress hosts | Only if ingress exists |
| `ingress.targetPort` | Target port for ingress | Only if ingress exists |
| `resources.limits` | CPU/Memory limits | Only if resource limits exist |
| `resources.requests` | CPU/Memory requests | Only if resource requests exist |
| `volumes` | Volume configurations | Only if volumes exist |
| `volumeMounts` | Volume mount points | Only if volume mounts exist |

### 🔧 Smart Field Detection

The tool intelligently includes only relevant fields:

- **Service fields** are only included if a Kubernetes Service exists
- **Environment variables** are only included if they are defined in the deployment
- **Ingress configuration** is only included if an Ingress resource exists
- **Resource limits/requests** are only included if they are defined
- **Volumes and volume mounts** are only included if they are configured

For a complete example with all features, see [`example-values.yaml`](example-values.yaml).

## 🏗️ Project Structure

```
chartgen/
├── main.go               # CLI entry point
├── cmd/
│   └── generate.go       # 'chartgen generate' command
├── internal/
│   └── parser.go         # Kubernetes to Helm conversion logic
├── go.mod               # Go module file
├── go.sum               # Dependency checksums
├── example-values.yaml  # Example output format
└── README.md            # This file
```

## 🔧 How It Works

1. **Resource Discovery**: Uses `kubectl get` to fetch Deployments, Services, and Ingresses
2. **Data Extraction**: Parses JSON output to extract relevant configuration
3. **Structure Mapping**: Maps Kubernetes resource fields to Helm values structure
4. **YAML Generation**: Converts the structured data to YAML format
5. **Output**: Writes to file or displays on stdout

### Supported Resource Types

- **Deployments**: Extracts image repository, tag, replica count, environment variables, resources, volumes
- **Services**: Extracts service type and port configuration
- **Ingresses**: Extracts host names and enables ingress configuration

## 🛠️ Development

### Prerequisites

- Go 1.23+
- kubectl
- Access to Kubernetes cluster

### Setup Development Environment

```bash
# Clone and setup
git clone https://github.com/yourusername/kubectl-chartgen.git
cd kubectl-chartgen

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o chartgen main.go
```

### Adding New Resource Types

To support additional Kubernetes resource types:

1. Add the resource type to `GetK8sResources()` in `internal/parser.go`
2. Create a new parsing function (e.g., `parseConfigMap()`)
3. Add the parsing logic to `ParseToHelmValues()`
4. Update the `HelmValues` struct if needed

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [YAML v3](https://gopkg.in/yaml.v3) for YAML processing
- Inspired by the need to migrate existing Kubernetes deployments to Helm

## 📞 Support

If you encounter any issues or have questions:

- Open an issue on GitHub
- Check the existing issues for solutions
- Review the documentation above

---

**Happy Helm-ing! 🎉** 