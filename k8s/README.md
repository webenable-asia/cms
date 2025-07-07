# K3s Deployment Manifests

This directory contains Kubernetes manifests for deploying the CMS application on K3s.

## Directory Structure
```
k8s/
├── namespace.yaml           # Application namespace
├── configmap.yaml          # Configuration data
├── secrets.yaml            # Sensitive data (base64 encoded)
├── persistent-volumes.yaml # Storage definitions
├── database/
│   ├── couchdb-deployment.yaml
│   ├── couchdb-service.yaml
│   └── couchdb-pvc.yaml
├── cache/
│   ├── valkey-deployment.yaml
│   ├── valkey-service.yaml
│   └── valkey-pvc.yaml
├── backend/
│   ├── backend-deployment.yaml
│   ├── backend-service.yaml
│   └── backend-hpa.yaml
├── frontend/
│   ├── frontend-deployment.yaml
│   ├── frontend-service.yaml
│   └── frontend-hpa.yaml
├── ingress/
│   ├── ingress.yaml
│   └── cluster-issuer.yaml
└── monitoring/
    ├── prometheus.yaml
    └── grafana.yaml
```

## Quick Deployment

```bash
# Apply all manifests
kubectl apply -k k8s/

# Check deployment status
kubectl get pods -n cms

# Check services
kubectl get svc -n cms

# Check ingress
kubectl get ingress -n cms
```

## Resource Requirements

### Minimum Single Node
- **CPU**: 6 vCPU
- **RAM**: 4GB
- **Storage**: 80GB SSD

### Recommended Multi-Node
- **Master**: 2 vCPU / 2GB RAM
- **Worker 1**: 4 vCPU / 4GB RAM  
- **Worker 2**: 4 vCPU / 4GB RAM

## Features Included

✅ **Auto-scaling**: HPA for frontend/backend  
✅ **SSL/TLS**: Automatic cert-manager integration  
✅ **Health Checks**: Kubernetes probes  
✅ **Rolling Updates**: Zero-downtime deployments  
✅ **Service Discovery**: Internal DNS resolution  
✅ **Resource Limits**: CPU/Memory constraints  
✅ **Persistent Storage**: Database/cache persistence  
✅ **Monitoring**: Prometheus + Grafana ready  
