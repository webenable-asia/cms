# WebEnable CMS - Kubernetes Deployment

This directory contains Kubernetes manifests for deploying WebEnable CMS using Kustomize.

## 🏗️ Architecture Overview

The deployment consists of the following components:

- **CouchDB**: Document database for content storage
- **Valkey**: Redis-compatible cache for sessions and caching
- **Backend**: Go API server (2 replicas)
- **Frontend**: Next.js public site (2 replicas)
- **Admin Panel**: Next.js CMS interface (2 replicas)
- **Ingress**: Nginx ingress controller for external access

## 📁 Directory Structure

```
k8s/
├── base/                    # Base manifests
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── service-accounts.yaml
│   ├── ingress.yaml
│   ├── couchdb/
│   ├── valkey/
│   ├── backend/
│   ├── frontend/
│   ├── admin-panel/
│   └── monitoring/
├── overlays/                # Environment-specific configurations
│   ├── development/
│   ├── staging/
│   └── production/
└── deploy.sh               # Deployment script
```

## 🚀 Quick Start

### Prerequisites

1. **Kubernetes Cluster**: A running Kubernetes cluster (1.20+)
2. **kubectl**: Configured to access your cluster
3. **kustomize**: Installed (v4.0+)
4. **Ingress Controller**: Nginx ingress controller installed
5. **Storage Class**: Default storage class configured

### Installation

1. **Clone and navigate to the project**:
   ```bash
   cd /path/to/webenable-cms
   ```

2. **Make the deployment script executable**:
   ```bash
   chmod +x k8s/deploy.sh
   ```

3. **Deploy to development environment**:
   ```bash
   ./k8s/deploy.sh development apply
   ```

4. **Deploy to production environment**:
   ```bash
   ./k8s/deploy.sh production apply
   ```

## 🔧 Configuration

### Environment Variables

Key configuration is managed through ConfigMaps and Secrets:

#### ConfigMap (`webenable-cms-config`)
- Database URLs
- Environment settings
- CORS origins
- Resource limits

#### Secret (`webenable-cms-secrets`)
- JWT secrets
- Database passwords
- TLS certificates

### Environment-Specific Overlays

#### Development
- Single replica deployments
- Debug mode enabled
- Local development URLs
- Reduced resource limits

#### Staging
- Single replica deployments
- Production-like configuration
- Staging URLs
- Standard resource limits

#### Production
- Multiple replica deployments
- Production configuration
- Production URLs
- Increased resource limits

## 📊 Resource Requirements

### Minimum Requirements
- **CPU**: 4 cores
- **Memory**: 8GB RAM
- **Storage**: 20GB

### Recommended Requirements
- **CPU**: 8 cores
- **Memory**: 16GB RAM
- **Storage**: 50GB

## 🔒 Security Features

### Security Context
- All pods run as non-root users
- Read-only root filesystem where possible
- Security context configured for each component

### Network Security
- Services use ClusterIP (internal access only)
- Ingress provides external access with TLS
- Security headers configured in ingress

### Secrets Management
- Sensitive data stored in Kubernetes secrets
- Environment-specific secret generation
- No hardcoded secrets in manifests

## 📈 Monitoring & Observability

### Health Checks
- Liveness probes for all containers
- Readiness probes for proper traffic routing
- Startup probes for slow-starting containers

### Metrics
- ServiceMonitor for Prometheus integration
- Health endpoint at `/api/health`
- Application metrics exposed

### High Availability
- PodDisruptionBudget configured
- Multiple replicas for stateless services
- Persistent storage for stateful services

## 🛠️ Management Commands

### Deployment Script Usage

```bash
# Deploy to development
./k8s/deploy.sh development apply

# Deploy to staging
./k8s/deploy.sh staging apply

# Deploy to production
./k8s/deploy.sh production apply

# Validate manifests
./k8s/deploy.sh production validate

# Build manifests without applying
./k8s/deploy.sh production build

# Delete deployment
./k8s/deploy.sh production delete
```

### Manual kubectl Commands

```bash
# Check deployment status
kubectl get pods -n webenable-cms-{environment}

# View logs
kubectl logs -f deployment/backend -n webenable-cms-{environment}

# Port forward for local access
kubectl port-forward svc/backend-service 8080:8080 -n webenable-cms-{environment}

# Check ingress
kubectl get ingress -n webenable-cms-{environment}

# View events
kubectl get events -n webenable-cms-{environment} --sort-by='.lastTimestamp'
```

## 🔄 Scaling

### Horizontal Pod Autoscaler (HPA)

To enable automatic scaling, create HPA resources:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: backend-hpa
  namespace: webenable-cms-prod
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: backend
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Manual Scaling

```bash
# Scale backend to 5 replicas
kubectl scale deployment backend --replicas=5 -n webenable-cms-prod

# Scale frontend to 3 replicas
kubectl scale deployment frontend --replicas=3 -n webenable-cms-prod
```

## 🔧 Troubleshooting

### Common Issues

1. **Pods not starting**:
   ```bash
   kubectl describe pod <pod-name> -n webenable-cms-{environment}
   kubectl logs <pod-name> -n webenable-cms-{environment}
   ```

2. **Services not accessible**:
   ```bash
   kubectl get svc -n webenable-cms-{environment}
   kubectl describe svc <service-name> -n webenable-cms-{environment}
   ```

3. **Ingress not working**:
   ```bash
   kubectl get ingress -n webenable-cms-{environment}
   kubectl describe ingress webenable-cms-ingress -n webenable-cms-{environment}
   ```

4. **Storage issues**:
   ```bash
   kubectl get pvc -n webenable-cms-{environment}
   kubectl describe pvc <pvc-name> -n webenable-cms-{environment}
   ```

### Health Checks

```bash
# Check all pods are running
kubectl get pods -n webenable-cms-{environment}

# Check service endpoints
kubectl get endpoints -n webenable-cms-{environment}

# Test API health
kubectl port-forward svc/backend-service 8080:8080 -n webenable-cms-{environment}
curl http://localhost:8080/api/health
```

## 🔄 Updates & Rollbacks

### Rolling Updates

```bash
# Update image
kubectl set image deployment/backend backend=webenable-cms-backend:v2.0.0 -n webenable-cms-{environment}

# Check rollout status
kubectl rollout status deployment/backend -n webenable-cms-{environment}
```

### Rollbacks

```bash
# Rollback to previous version
kubectl rollout undo deployment/backend -n webenable-cms-{environment}

# Rollback to specific revision
kubectl rollout undo deployment/backend --to-revision=2 -n webenable-cms-{environment}
```

## 📝 Customization

### Adding New Environments

1. Create a new overlay directory:
   ```bash
   mkdir -p k8s/overlays/new-environment
   ```

2. Create `kustomization.yaml`:
   ```yaml
   apiVersion: kustomize.config.k8s.io/v1beta1
   kind: Kustomization
   
   namespace: webenable-cms-new-environment
   
   resources:
     - ../../base
   
   configMapGenerator:
     - name: webenable-cms-config
       behavior: merge
       literals:
         - NODE_ENV=production
   ```

3. Add environment-specific patches as needed.

### Custom Resource Limits

Create a patch file in your overlay:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  template:
    spec:
      containers:
      - name: backend
        resources:
          requests:
            cpu: "2000m"
            memory: "2Gi"
          limits:
            cpu: "4000m"
            memory: "4Gi"
```

## 🔐 Security Best Practices

1. **Rotate secrets regularly**
2. **Use RBAC for access control**
3. **Enable network policies**
4. **Scan images for vulnerabilities**
5. **Monitor resource usage**
6. **Backup persistent data**

## 📚 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kustomize Documentation](https://kustomize.io/)
- [Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Prometheus Operator](https://prometheus-operator.dev/) 