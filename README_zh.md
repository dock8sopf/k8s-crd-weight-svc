# Kubernetes ServiceWeight CRD

## 概述

ServiceWeight 是一个 Kubernetes 自定义资源定义 (CRD)，它扩展了 Kubernetes 原生 `Service` 资源的功能，增加了基于权重的流量转发支持。它允许您根据数值权重配置不同后端服务之间的流量分配。

## 功能特性

- **标准 Service 兼容性**：当不指定权重时，行为与 Kubernetes 原生 `Service` 完全一致
- **加权流量转发**：基于配置的权重在后端服务之间分配流量
- **无缝集成**：与现有的 Kubernetes 网络基础设施配合使用
- **可扩展设计**：易于与服务网格解决方案集成，实现高级流量管理

## 架构

ServiceWeight CRD 由以下组件组成：

1. **自定义资源定义**：定义 ServiceWeight 的 API 架构
2. **控制器**：管理 ServiceWeight 资源的生命周期并与后端服务同步
3. **加权后端支持**：允许配置多个后端服务及其流量分配权重

## 安装

### 先决条件

- Kubernetes 集群（建议版本 1.30+）
- Go 1.22.2+ 用于开发
- 已配置好的 kubectl，可以访问您的集群
- 控制器运行时依赖

### 安装步骤

1. **克隆仓库**
   ```bash
   git clone <repository-url>
   cd k8s-crd-weight-svc
   ```

2. **安装 CRD**
   ```bash
   make install
   ```

3. **部署控制器**
   ```bash
   make deploy
   ```

4. **验证安装**
   ```bash
   kubectl get crd serviceweights.example.com
   ```

## 使用方法

### 标准服务（无权重）

当未指定 `WeightedBackends` 时，ServiceWeight 的行为与原生 Kubernetes Service 完全相同：

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

### 加权服务

当指定 `WeightedBackends` 时，流量将根据配置的权重进行分配：

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
      weight: 70  # 70% 的流量
    - name: backend-v2
      port: 8081
      weight: 30  # 30% 的流量
```

## API 参考

### ServiceWeightSpec

| 字段 | 类型 | 描述 |
|------|------|------|
| `ports` | `[]corev1.ServicePort` | 要暴露的端口列表 |
| `selector` | `map[string]string` | 用于查找服务对应的 Pod 的选择器 |
| `type` | `corev1.ServiceType` | 服务类型（ClusterIP、NodePort、LoadBalancer、ExternalName） |
| `clusterIP` | `string` | 服务的集群 IP 地址 |
| `WeightedBackends` | `[]WeightedBackend` | 加权后端服务列表 |

### WeightedBackend

| 字段 | 类型 | 描述 |
|------|------|------|
| `name` | `string` | 后端服务的名称 |
| `port` | `int32` | 转发流量的端口 |
| `weight` | `int32` | 流量分配的权重 |

## 开发

### 构建项目

```bash
make build
```

### 运行测试

```bash
make test
```

### 生成 CRD 清单

```bash
make manifests
```

## 限制

- 当前，加权流量分配需要与服务网格解决方案（如 Istio）或自定义负载均衡器集成
- 控制器会为服务添加加权后端的注解，但实际的流量分配需要由您的网络层实现

## 路线图

- 增加对自动服务网格集成的支持
- 实现高级流量路由策略
- 添加监控和指标支持
- 增强验证和错误处理

## 贡献

欢迎贡献！请参阅贡献指南了解更多详情。

## 许可证

MIT
