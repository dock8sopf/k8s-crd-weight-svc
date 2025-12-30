# ServiceWeight CRD Project

## Project Overview

ServiceWeight is a Kubernetes Custom Resource Definition (CRD) that extends the functionality of Kubernetes native Service by adding weighted traffic forwarding capability.

## Features

- Fully compatible with all features of Kubernetes native Service
- Support for configuring multiple backend services with weights
- Automatic creation and management of corresponding Kubernetes Service resources
- Status tracking and condition monitoring

## Technology Stack

- Go 1.22.2
- Kubernetes 1.30.14
- Controller Runtime 0.18.4
- Kubebuilder toolchain

## Installation Steps

### 1. Clone the project

```bash
git clone <repository-url>
cd k8s-crd-service-weight
```

### 2. Install CRD

```bash
make install
```

### 3. Deploy controller

```bash
make deploy
```

## Usage Examples

### Basic Usage (same as native Service)

```yaml
apiVersion: example.com/v1alpha1
kind: ServiceWeight
metadata:
  name: my-service
  namespace: default
spec:
  selector:
    app: my-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

### Weighted Traffic Forwarding

```yaml
apiVersion: example.com/v1alpha1
kind: ServiceWeight
metadata:
  name: my-weighted-service
  namespace: default
spec:
  type: ClusterIP
  ports:
    - protocol: TCP
      port: 80
  weightedBackends:
    - name: backend-service-1
      weight: 70
      port: 8080
    - name: backend-service-2
      weight: 30
      port: 8080
```

## Development Guide

### Build the project

```bash
make build
```

### Run tests

```bash
make test
```

### Run controller locally

```bash
make run
```

### Generate API code

```bash
make generate
```

### Generate CRD manifests

```bash
make manifests
```

## Project Structure

```
k8s-crd-service-weight/
├── api/
│   └── v1alpha1/
│       ├── groupversion_info.go    # API group version definition
│       └── serviceweight_types.go   # ServiceWeight CRD type definition
├── controller/
│   └── serviceweight_controller.go  # Controller implementation
├── main.go                          # Entry point
├── Makefile                         # Build scripts
├── go.mod                           # Go module dependencies
└── go.sum                           # Dependency checksum file
```

## API Definition

### ServiceWeightSpec

```go
type ServiceWeightSpec struct {
    corev1.ServiceSpec `json:",inline"`                // Inherit native ServiceSpec
    WeightedBackends   []WeightedBackend `json:"weightedBackends,omitempty"`  // Weighted backends list
}
```

### WeightedBackend

```go
type WeightedBackend struct {
    Name   string `json:"name"`   // Backend service name
    Weight int32  `json:"weight"` // Weight value
    Port   int32  `json:"port"`   // Forwarding port
}
```

## Notes

1. Ensure Kubernetes cluster version >= 1.24
2. Cluster-admin privileges required for CRD installation and deployment
3. Weight values don't need to sum to 100, can be any positive integers
4. When weightedBackends is not specified, ServiceWeight behaves exactly like native Service

## License

MIT License