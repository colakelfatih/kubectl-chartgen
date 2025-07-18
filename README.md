# kubectl-chartgen 🚀

A powerful kubectl plugin that automatically generates Helm `values.yaml` files from existing Kubernetes resources. Perfect for migrating existing deployments to Helm charts or creating Helm templates from running applications.

## ✨ Features

- 🔍 **Automatic Resource Detection**: Automatically discovers Deployments, Services, and Ingresses
- 🎯 **Helm Values Generation**: Converts Kubernetes resources to Helm-compatible `values.yaml` format
- 📁 **Namespace Support**: Target specific namespaces or use current context
- 💾 **Flexible Output**: Save to file or output to stdout
- 🛡️ **Error Handling**: Comprehensive error handling and user feedback
- ⚡ **Fast & Lightweight**: Built with Go for optimal performance
- 🔐 **Secure Connections**: Support for insecure TLS connections with `--insecure-skip-tls-verify`
- 🌐 **Remote Cluster Support**: Connect to remote Kubernetes clusters with custom kubeconfig files

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
mv kubectl-chartgen /usr/local/bin/

# Verify installation
kubectl chartgen --help
```

**Note**: The plugin will be available as `kubectl chartgen` after installation.

## 📖 Usage

### Basic Usage

```bash
# Generate values.yaml from current namespace
go run main.go generate

# Or if installed as binary
kubectl chartgen generate
```

### Using Different Kubeconfig Files

#### Method 1: Using KUBECONFIG Environment Variable

```bash
# Set kubeconfig for a specific cluster
export KUBECONFIG=/path/to/your/kubeconfig.yaml
kubectl chartgen generate

# Or use inline
KUBECONFIG=/path/to/production-kubeconfig.yaml kubectl chartgen generate
```

#### Method 2: Using --kubeconfig Flag

```bash
# Specify kubeconfig file directly
kubectl --kubeconfig=/path/to/your/kubeconfig.yaml chartgen generate

# Generate from specific namespace in different cluster
kubectl --kubeconfig=/path/to/production-kubeconfig.yaml chartgen generate -n production
```

#### Method 3: Using Different Contexts

```bash
# List available contexts
kubectl config get-contexts

# Switch to a specific context
kubectl config use-context production-cluster

# Generate using current context
kubectl chartgen generate

# Or use context inline
kubectl --context=production-cluster chartgen generate
```

#### Method 4: Multiple Kubeconfig Files

```bash
# Combine multiple kubeconfig files
export KUBECONFIG=/path/to/cluster1.yaml:/path/to/cluster2.yaml

# Generate from specific cluster
kubectl --context=cluster1 chartgen generate -n default
kubectl --context=cluster2 chartgen generate -n production
```

### Advanced Usage

```bash
# Generate from specific namespace
kubectl chartgen generate -n my-app-namespace

# Output to specific file
kubectl chartgen generate -o my-helm-values.yaml

# Output to stdout
kubectl chartgen generate -o -

# Combine options with different kubeconfig
kubectl --kubeconfig=/path/to/prod-cluster.yaml chartgen generate -n production -o prod-values.yaml
```

### Remote Cluster Connections

When connecting to remote Kubernetes clusters, you might encounter TLS certificate verification issues. Use the `--insecure-skip-tls-verify` flag to bypass certificate verification:

```bash
# Connect to remote cluster with insecure TLS
kubectl chartgen generate --kubeconfig /path/to/remote-kubeconfig.yaml --insecure-skip-tls-verify

# With specific namespace
kubectl chartgen generate --kubeconfig /path/to/remote-kubeconfig.yaml --insecure-skip-tls-verify -n my-namespace

# Output to file
kubectl chartgen generate --kubeconfig /path/to/remote-kubeconfig.yaml --insecure-skip-tls-verify -o remote-values.yaml
```

**⚠️ Security Note**: Use `--insecure-skip-tls-verify` only when you trust the remote cluster and understand the security implications.

### Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Target namespace (default: current) | current context |
| `--output` | `-o` | Output file path (use `-` for stdout) | `values.yaml` |
| `--kubeconfig` | `-k` | Path to kubeconfig file | default kubeconfig |
| `--insecure-skip-tls-verify` | - | Skip TLS certificate verification | false |
| `--help` | `-h` | Show help information | - |

**Note**: You can also use standard kubectl flags like `--context`, `--cluster`, and `--user` to target different Kubernetes clusters.

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

### Example 4: Remote Cluster with Insecure TLS

```bash
# Connect to remote cluster with TLS verification disabled
kubectl chartgen generate --kubeconfig /path/to/remote-cluster.yaml --insecure-skip-tls-verify -n my-app
```

**Output:**
```
Generating Helm values.yaml from Kubernetes resources...
Target namespace: my-app
Using kubeconfig: /path/to/remote-cluster.yaml
Fetching Kubernetes resources...
Found 5 resources
Converting to Helm values structure...
Generated values for 3 services
Generating YAML output...
Helm values written to: values.yaml
```

## 🔧 Troubleshooting

### Common Issues

#### 1. TLS Certificate Verification Errors

**Error**: `tls: failed to verify certificate: x509: certificate signed by unknown authority`

**Solution**: Use the `--insecure-skip-tls-verify` flag:

```bash
kubectl chartgen generate --kubeconfig /path/to/kubeconfig.yaml --insecure-skip-tls-verify
```

#### 2. Permission Denied Errors

**Error**: `Error from server (Forbidden): deployments.apps is forbidden`

**Solution**: Check your user permissions and target a namespace you have access to:

```bash
# Check your permissions
kubectl auth can-i list deployments

# Use a specific namespace you have access to
kubectl chartgen generate -n your-accessible-namespace
```

#### 3. Cluster Connection Issues

**Error**: `Unable to connect to the server: dial tcp: no route to host`

**Solution**: 
- Ensure your cluster is running (e.g., `minikube start` for local clusters)
- Check your kubeconfig file is correct
- Verify network connectivity to the cluster

#### 4. Plugin Not Found

**Error**: `kubectl chartgen: command not found`

**Solution**: Ensure the plugin is properly installed:

```bash
# Build and install
go build -o kubectl-chartgen main.go
chmod +x kubectl-chartgen
sudo mv kubectl-chartgen /usr/local/bin/

# Verify installation
kubectl chartgen --help
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
git clone https://github.com/colakelfatih/kubectl-chartgen.git
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