# GitLab CI/CD Quick Reference

## Pipeline Stages

```
validate → test → build → security → deploy
```

## Environment Variables

### Required in GitLab CI/CD Settings

```bash
# Kubernetes Configs (base64 encoded)
KUBE_CONFIG_DEV=<base64-kubeconfig>
KUBE_CONFIG_STAGING=<base64-kubeconfig>
KUBE_CONFIG_PROD=<base64-kubeconfig>

# Application URLs
NEXT_PUBLIC_API_URL=https://api.webenable.asia
BACKEND_URL=https://api.webenable.asia
ADMIN_NEXT_PUBLIC_API_URL=https://api.webenable.asia
ADMIN_BACKEND_URL=https://api.webenable.asia
```

### Auto-configured by GitLab
- `CI_REGISTRY` - GitLab registry URL
- `CI_REGISTRY_USER` - Registry username
- `CI_REGISTRY_PASSWORD` - Registry password

## Branch Strategy

| Branch Type | Validation | Test | Build | Security | Deploy |
|-------------|------------|------|-------|----------|---------|
| `main` | ✅ | ✅ | ✅ | ✅ | dev (auto), staging/prod (manual) |
| `feature/*` | ✅ | ✅ | ✅ | ❌ | dev (auto) |
| `release/*` | ✅ | ✅ | ✅ | ✅ | staging/prod (manual) |
| MR | ✅ | ✅ | ❌ | ❌ | ❌ |

## Quick Commands

### Generate Registry Secret
```bash
./k8s/scripts/generate-registry-secret.sh username password namespace
```

### Deploy to Environment
```bash
# Development (auto on main/feature branches)
git push origin main

# Staging (manual)
# Go to GitLab CI/CD → Pipelines → deploy:staging → Play

# Production (manual)
# Go to GitLab CI/CD → Pipelines → deploy:production → Play
```

### Check Deployment Status
```bash
# Development
kubectl get pods -n webenable-cms-dev

# Staging
kubectl get pods -n webenable-cms-staging

# Production
kubectl get pods -n webenable-cms-prod
```

## Image Tags

| Component | Registry | Tags |
|-----------|----------|------|
| Backend | `registry.gitlab.com/webenable/cms/backend` | `latest`, `$CI_COMMIT_SHORT_SHA` |
| Frontend | `registry.gitlab.com/webenable/cms/frontend` | `latest`, `$CI_COMMIT_SHORT_SHA` |
| Admin Panel | `registry.gitlab.com/webenable/cms/admin-panel` | `latest`, `$CI_COMMIT_SHORT_SHA` |

## Environment URLs

| Environment | Frontend | Admin Panel | API |
|-------------|----------|-------------|-----|
| Development | https://dev.webenable.asia | https://dev.webenable.asia/admin | https://dev.webenable.asia/api |
| Staging | https://staging.webenable.asia | https://staging.webenable.asia/admin | https://staging.webenable.asia/api |
| Production | https://webenable.asia | https://webenable.asia/admin | https://webenable.asia/api |

## Troubleshooting

### Build Failures
```bash
# Check Docker build logs
# GitLab CI/CD → Jobs → build:* → View Logs

# Verify registry access
docker login registry.gitlab.com
```

### Deployment Failures
```bash
# Check Kubernetes cluster access
kubectl cluster-info

# Check namespace exists
kubectl get namespace webenable-cms-dev

# Check image pull secret
kubectl get secret gitlab-registry-secret -n webenable-cms-dev
```

### Health Check Failures
```bash
# Check pod logs
kubectl logs -f deployment/backend -n webenable-cms-dev

# Check pod status
kubectl describe pod -l app.kubernetes.io/component=backend -n webenable-cms-dev
```

## Security Scanning

The pipeline uses Trivy to scan Docker images for vulnerabilities:
- Scans all built images
- Fails on HIGH/CRITICAL vulnerabilities
- Generates security reports

## Monitoring

### GitLab CI/CD
- Pipeline status and logs
- Environment deployment status
- Job artifacts and reports

### Kubernetes
- Pod status and logs
- Service endpoints
- Ingress configuration

## Rollback

To rollback a deployment:
```bash
# Rollback to previous version
kubectl rollout undo deployment/backend -n webenable-cms-dev

# Rollback to specific version
kubectl rollout undo deployment/backend -n webenable-cms-dev --to-revision=2
```

## Performance Tips

1. **Use caching**: Dependencies are cached between pipeline runs
2. **Parallel jobs**: Tests run in parallel for faster execution
3. **Optimized images**: Multi-stage Docker builds reduce image size
4. **Resource limits**: Kubernetes resource limits prevent resource exhaustion 