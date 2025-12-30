# Kubernetes ServiceWeight CRD

## Overview

ServiceWeight is a Kubernetes Custom Resource Definition (CRD) that extends the functionality of Kubernetes native `Service` resources by adding support for weighted traffic forwarding. It allows you to configure traffic distribution between different backend services based on numerical weights.

## Features

- **Standard Service Compatibility**: Behaves exactly like Kubernetes native `Service` when no weights are specified
- **Weighted Traffic Forwarding**: Distribute traffic between backend services based on configured weights
- **Seamless Integration**: Works with existing Kubernetes networking infrastructure
- **Extensible Design**: Easy to integrate with service mesh solutions for advanced traffic management

## Architecture

The ServiceWeight CRD consists of:

1. **Custom Resource Definition**: Defines the API schema for ServiceWeight
2. **Controller**: Manages the lifecycle of ServiceWeight resources and syncs with backend services
3. **Weighted Backend Support**: Allows configuring multiple backend services with traffic distribution weights

## Installation

### Prerequisites

- Kubernetes cluster (version 1.30+ recommended)
- Go 1.22.2+ for development
- kubectl configured to access your cluster
- Controller runtime dependencies

### Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd k8s-crd-weight-svc
   ```

2. **Install the CRD**
   ```bash
   make install
   ```

3. **Deploy the controller**
   ```bash
   make deploy
   ```

4. **Verify installation**
   ```bash
   kubectl get crd serviceweights.example.com
   ```

## Usage

### Standard Service (No Weights)

When no `WeightedBackends` are specified, ServiceWeight behaves exactly like a native Kubernetes Service:

```yaml
apiVersion: example.com/v1alpha1
kind: ServiceWeight
metadata:
  name: standard-service
  namespace: default
spec:
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app: my-app
  type: ClusterIP
```

### Weighted Service

When `WeightedBackends` are specified, traffic is distributed based on the configured weights:

```yaml
apiVersion: example.com/v1alpha1
kind: ServiceWeight
metadata:
  name: weighted-service
  namespace: default
spec:
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app: my-app
  type: ClusterIP
  
  WeightedBackends:
    - name: backend-v1
      port: 8080
      weight: 70  # 70% of traffic
    - name: backend-v2
      port: 8081
      weight: 30  # 30% of traffic
```

## API Reference

### ServiceWeightSpec

| Field | Type | Description |
|-------|------|-------------|
| `ports` | `[]corev1.ServicePort` | List of ports to expose |
| `selector` | `map[string]string` | Selector to find pods for the service |
| `type` | `corev1.ServiceType` | Type of service (ClusterIP, NodePort, LoadBalancer, ExternalName) |
| `clusterIP` | `string` | Cluster IP address of the service |
| `WeightedBackends` | `[]WeightedBackend` | List of weighted backend services |

### WeightedBackend

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Name of the backend service |
| `port` | `int32` | Port to forward traffic to |
| `weight` | `int32` | Weight for traffic distribution |

## Development

### Building the project

```bash
make build
```

### Running tests

```bash
make test
```

### Generating CRD manifests

```bash
make manifests
```

## Limitations

- Currently, weighted traffic distribution requires integration with a service mesh solution (like Istio) or custom load balancer
- The controller adds annotations to services for weighted backends, but actual traffic distribution needs to be implemented by your networking layer

## Roadmap

- Add support for automatic service mesh integration
- Implement advanced traffic routing strategies
- Add monitoring and metrics support
- Enhance validation and error handling

## Contributing

Contributions are welcome! Please see the contributing guidelines for more details.

## License

MIT
