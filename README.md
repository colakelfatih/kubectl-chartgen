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
mv kubectl-chartgen ~/.kube/plugins/chartgen/
```

## 📖 Usage

### Basic Usage

```bash
# Generate values.yaml from current namespace
go run main.go generate

# Or if installed as binary
chartgen generate
```

### Advanced Usage

```bash
# Generate from specific namespace
chartgen generate -n my-app-namespace

# Output to specific file
chartgen generate -o my-helm-values.yaml

# Output to stdout
chartgen generate -o -

# Combine options
chartgen generate -n production -o prod-values.yaml
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
chartgen generate
```

**Output:**
```
Generating Helm values.yaml from Kubernetes resources...
Fetching Kubernetes resources...
Found 3 resources
Converting to Helm values structure...
Generating YAML output...
Helm values written to: values.yaml
```

### Example 2: Generate from Specific Namespace

```bash
# Generate from 'production' namespace
chartgen generate -n production -o prod-values.yaml
```

### Example 3: Preview Output

```bash
# Preview the generated values without saving
chartgen generate -o -
```

**Sample Output:**
```yaml
image:
  repository: nginx
  tag: "1.21"
service:
  type: ClusterIP
  port: 80
ingress:
  enabled: true
  hosts:
    - myapp.example.com
replicaCount: 3
```

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
└── README.md            # This file
```

## 🔧 How It Works

1. **Resource Discovery**: Uses `kubectl get` to fetch Deployments, Services, and Ingresses
2. **Data Extraction**: Parses JSON output to extract relevant configuration
3. **Structure Mapping**: Maps Kubernetes resource fields to Helm values structure
4. **YAML Generation**: Converts the structured data to YAML format
5. **Output**: Writes to file or displays on stdout

### Supported Resource Types

- **Deployments**: Extracts image repository, tag, and replica count
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