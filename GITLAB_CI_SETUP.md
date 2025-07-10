# GitLab CI/CD Pipeline Setup Guide

This guide explains how to set up and configure the GitLab CI/CD pipeline for the WebEnable CMS project.

## Overview

The pipeline includes the following stages:
1. **Validate** - YAML validation and Kustomize manifest validation
2. **Test** - Unit tests for backend, frontend, and admin panel
3. **Build** - Docker image building and pushing to GitLab registry
4. **Security** - Vulnerability scanning with Trivy
5. **Deploy** - Kubernetes deployment using Kustomize

## Prerequisites

### 1. GitLab Project Setup

Ensure your GitLab project has:
- Container Registry enabled
- Kubernetes integration configured (optional, for direct deployment)
- Appropriate permissions for CI/CD

### 2. Required Environment Variables

Set these variables in your GitLab project's CI/CD settings:

#### GitLab Registry (Auto-configured)
- `CI_REGISTRY` - Automatically set by GitLab
- `CI_REGISTRY_USER` - Automatically set by GitLab
- `CI_REGISTRY_PASSWORD` - Automatically set by GitLab

#### Kubernetes Configuration
```bash
# Development Environment
KUBE_CONFIG_DEV=<base64-encoded-kubeconfig>

# Staging Environment  
KUBE_CONFIG_STAGING=<base64-encoded-kubeconfig>

# Production Environment
KUBE_CONFIG_PROD=<base64-encoded-kubeconfig>
```

#### Application Configuration
```bash
# Frontend Environment Variables
NEXT_PUBLIC_API_URL=https://api.webenable.asia
BACKEND_URL=https://api.webenable.asia

# Admin Panel Environment Variables
ADMIN_NEXT_PUBLIC_API_URL=https://api.webenable.asia
ADMIN_BACKEND_URL=https://api.webenable.asia
```

### 3. Kubernetes Cluster Setup

#### Create Namespaces
```bash
# Development
kubectl create namespace webenable-cms-dev

# Staging
kubectl create namespace webenable-cms-staging

# Production
kubectl create namespace webenable-cms-prod
```

#### Generate Registry Secret
Use the provided script to generate the GitLab registry secret:

```bash
# Using environment variables
export CI_REGISTRY_USER="your-gitlab-username"
export CI_REGISTRY_PASSWORD="your-gitlab-token"
./k8s/scripts/generate-registry-secret.sh

# Or with explicit parameters
./k8s/scripts/generate-registry-secret.sh username password namespace
```

Apply the generated secret to your Kubernetes cluster.

## Pipeline Configuration

### Branch Strategy

The pipeline uses the following branch strategy:

- **Main Branch** (`main`/`master`): 
  - Runs all stages
  - Auto-deploys to development
  - Manual deployment to staging/production

- **Feature Branches** (`feature/*`):
  - Runs validation and testing
  - Builds images for testing
  - Auto-deploys to development

- **Release Branches** (`release/*`):
  - Runs all stages
  - Manual deployment to staging/production

- **Merge Requests**:
  - Runs validation and testing only
  - No deployment

### Environment Configuration

#### Development Environment
- **URL**: https://dev.webenable.asia
- **Auto-deploy**: Yes (on main branch and feature branches)
- **Manual approval**: No

#### Staging Environment
- **URL**: https://staging.webenable.asia
- **Auto-deploy**: No (manual approval required)
- **Trigger**: Release branches

#### Production Environment
- **URL**: https://webenable.asia
- **Auto-deploy**: No (manual approval required)
- **Trigger**: Release branches

## Usage

### 1. Initial Setup

1. **Configure Environment Variables**:
   - Go to your GitLab project → Settings → CI/CD → Variables
   - Add all required variables listed above

2. **Set up Kubernetes Clusters**:
   - Create namespaces for each environment
   - Generate and apply registry secrets
   - Configure ingress controllers

3. **Configure DNS**:
   - Point your domains to your Kubernetes ingress
   - Set up SSL certificates

### 2. Development Workflow

1. **Create Feature Branch**:
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make Changes and Push**:
   ```bash
   git add .
   git commit -m "Add new feature"
   git push origin feature/new-feature
   ```

3. **Create Merge Request**:
   - Pipeline will run validation and tests
   - Review and merge when ready

4. **Deploy to Development**:
   - Merging to main branch triggers development deployment
   - Check the deployment in GitLab CI/CD → Environments

### 3. Release Process

1. **Create Release Branch**:
   ```bash
   git checkout -b release/v1.0.0
   git push origin release/v1.0.0
   ```

2. **Deploy to Staging**:
   - Go to GitLab CI/CD → Pipelines
   - Find the release pipeline
   - Click "deploy:staging" job
   - Click "Play" to start deployment

3. **Test in Staging**:
   - Verify all functionality works
   - Run integration tests

4. **Deploy to Production**:
   - If staging tests pass, deploy to production
   - Click "deploy:production" job
   - Click "Play" to start deployment

## Monitoring and Troubleshooting

### 1. Pipeline Monitoring

- **GitLab CI/CD Dashboard**: Monitor pipeline status and logs
- **Environment Status**: Check deployment status in Environments tab
- **Job Logs**: Detailed logs for each pipeline stage

### 2. Kubernetes Monitoring

```bash
# Check deployment status
kubectl get deployments -n webenable-cms-dev
kubectl get deployments -n webenable-cms-staging
kubectl get deployments -n webenable-cms-prod

# Check pod status
kubectl get pods -n webenable-cms-dev
kubectl get pods -n webenable-cms-staging
kubectl get pods -n webenable-cms-prod

# Check logs
kubectl logs -f deployment/backend -n webenable-cms-dev
kubectl logs -f deployment/frontend -n webenable-cms-dev
kubectl logs -f deployment/admin-panel -n webenable-cms-dev
```

### 3. Common Issues

#### Build Failures
- Check Dockerfile syntax
- Verify build arguments
- Check registry permissions

#### Deployment Failures
- Verify Kubernetes cluster access
- Check namespace existence
- Verify image pull secrets

#### Health Check Failures
- Check application logs
- Verify environment variables
- Check service connectivity

## Security Considerations

### 1. Secrets Management
- Use GitLab CI/CD variables for sensitive data
- Never commit secrets to the repository
- Rotate secrets regularly

### 2. Registry Security
- Use GitLab's built-in container registry
- Enable vulnerability scanning
- Regularly update base images

### 3. Kubernetes Security
- Use RBAC for service accounts
- Enable network policies
- Regular security audits

## Performance Optimization

### 1. Pipeline Optimization
- Use caching for dependencies
- Parallel job execution
- Optimize Docker layers

### 2. Deployment Optimization
- Use rolling updates
- Configure resource limits
- Enable horizontal pod autoscaling

## Advanced Configuration

### 1. Custom Build Arguments

You can customize build arguments by modifying the pipeline:

```yaml
build:backend:
  script:
    - docker build --build-arg GO_ENV=production --build-arg VERSION=$CI_COMMIT_SHORT_SHA -t $DOCKER_IMAGE_PREFIX/backend:$BACKEND_VERSION .
```

### 2. Multi-Environment Variables

For different environments, you can use conditional variables:

```yaml
variables:
  BACKEND_URL: $BACKEND_URL_DEV
  rules:
    - if: $CI_ENVIRONMENT_NAME == "development"
      variables:
        BACKEND_URL: $BACKEND_URL_DEV
    - if: $CI_ENVIRONMENT_NAME == "production"
      variables:
        BACKEND_URL: $BACKEND_URL_PROD
```

### 3. Custom Deployment Scripts

You can add custom deployment scripts:

```yaml
deploy:custom:
  script:
    - ./scripts/custom-deploy.sh
    - kubectl apply -f custom-manifests/
```

## Support

For issues and questions:
1. Check the pipeline logs in GitLab CI/CD
2. Review the Kubernetes manifests in the `k8s/` directory
3. Consult the main project README
4. Check the deployment troubleshooting guide 