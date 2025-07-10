# GKE Deployment with Cloudflare DNS

This guide covers deploying WebEnable CMS on Google Kubernetes Engine (GKE) with Cloudflare managing DNS and SSL certificates.

## Architecture Overview

```
Internet → Cloudflare → GKE Load Balancer → Ingress → Services → Pods
                ↓
        SSL Termination + DNS Management
```

## Prerequisites

### 1. Google Cloud Platform Setup

```bash
# Install Google Cloud CLI
# https://cloud.google.com/sdk/docs/install

# Authenticate with GCP
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# Enable required APIs
gcloud services enable container.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable dns.googleapis.com
```

### 2. Cloudflare Setup

1. **Add Domain to Cloudflare**:
   - Log into Cloudflare dashboard
   - Add your domain (e.g., `webenable-cms.com`)
   - Update nameservers at your domain registrar

2. **Create API Token**:
   - Go to Cloudflare Dashboard → Profile → API Tokens
   - Create custom token with:
     - Zone:Zone:Read permissions
     - Zone:DNS:Edit permissions
     - Include specific zone: your domain

### 3. Required Tools

```bash
# Install kubectl
gcloud components install kubectl

# Install helm (optional, for additional tools)
curl https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz | tar xz
sudo mv linux-amd64/helm /usr/local/bin/

# Install cert-manager CLI
curl -L -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-linux-amd64.tar.gz
tar xzf cmctl.tar.gz
sudo mv cmctl /usr/local/bin/
```

## GKE Cluster Setup

### 1. Create GKE Cluster

```bash
# Create production cluster
gcloud container clusters create webenable-cms-prod \
  --zone us-central1-a \
  --num-nodes 3 \
  --min-nodes 1 \
  --max-nodes 10 \
  --machine-type e2-standard-4 \
  --disk-size 100 \
  --disk-type pd-ssd \
  --enable-autoscaling \
  --enable-autorepair \
  --enable-autoupgrade \
  --enable-network-policy \
  --enable-shielded-nodes \
  --workload-pool=YOUR_PROJECT_ID.svc.id.goog \
  --enable-workload-identity

# Get cluster credentials
gcloud container clusters get-credentials webenable-cms-prod --zone us-central1-a
```

### 2. Configure Workload Identity

```bash
# Create service account for GitLab CI
gcloud iam service-accounts create gitlab-ci \
  --display-name="GitLab CI Service Account"

# Grant necessary permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:gitlab-ci@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/container.developer"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:gitlab-ci@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/dns.admin"

# Create key for GitLab CI
gcloud iam service-accounts keys create gitlab-ci-key.json \
  --iam-account=gitlab-ci@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

## Cloudflare Configuration

### 1. DNS Records Setup

Create these DNS records in Cloudflare:

```bash
# A Records (point to GKE load balancer IP)
webenable-cms.com        A    <GKE_LOAD_BALANCER_IP>
www.webenable-cms.com    A    <GKE_LOAD_BALANCER_IP>
api.webenable-cms.com    A    <GKE_LOAD_BALANCER_IP>
admin.webenable-cms.com  A    <GKE_LOAD_BALANCER_IP>

# CNAME Records
*.webenable-cms.com      CNAME webenable-cms.com
```

### 2. Cloudflare Settings

Configure these settings in Cloudflare:

- **SSL/TLS**: Full (strict)
- **Always Use HTTPS**: On
- **HSTS**: Enabled
- **Security Level**: Medium
- **Rate Limiting**: Enabled
- **Bot Fight Mode**: Enabled
- **Browser Integrity Check**: Enabled

## Kubernetes Setup

### 1. Install cert-manager

```bash
# Add cert-manager helm repository
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Install cert-manager
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.13.0 \
  --set installCRDs=true

# Verify installation
kubectl get pods -n cert-manager
```

### 2. Configure Cloudflare DNS01 Challenge

```bash
# Create Cloudflare API secret
kubectl create secret generic cloudflare-api-token \
  --from-literal=api-token=YOUR_CLOUDFLARE_API_TOKEN \
  -n cert-manager

# Apply ClusterIssuer for Let's Encrypt
kubectl apply -f k8s/gke/cluster-issuer.yaml
```

### 3. Deploy Application

```bash
# Create namespaces
kubectl apply -f k8s/base/namespace.yaml

# Apply base configuration
kubectl apply -k k8s/base/

# Deploy to production
kubectl apply -k k8s/overlays/production/
```

## GitLab CI Configuration

### 1. Update Environment Variables

Add these variables to your GitLab CI/CD settings:

```bash
# GCP Configuration
GCP_PROJECT_ID=your-project-id
GCP_CLUSTER_NAME=webenable-cms-prod
GCP_CLUSTER_ZONE=us-central1-a
GCP_SERVICE_ACCOUNT_KEY=<base64-encoded-service-account-key>

# Cloudflare Configuration
CLOUDFLARE_API_TOKEN=your-cloudflare-api-token
CLOUDFLARE_ZONE_ID=your-zone-id

# Application URLs
NEXT_PUBLIC_API_URL=https://api.webenable-cms.com
BACKEND_URL=https://api.webenable-cms.com
ADMIN_NEXT_PUBLIC_API_URL=https://api.webenable-cms.com
ADMIN_BACKEND_URL=https://api.webenable-cms.com
```

### 2. Update GitLab CI Pipeline

The pipeline will automatically:
- Authenticate with GCP
- Build and push Docker images
- Deploy to GKE
- Update DNS records if needed

## Monitoring and Logging

### 1. Google Cloud Monitoring

```bash
# Enable monitoring
gcloud container clusters update webenable-cms-prod \
  --zone us-central1-a \
  --enable-stackdriver-kubernetes

# View metrics in Google Cloud Console
# https://console.cloud.google.com/monitoring
```

### 2. Application Logs

```bash
# View application logs
gcloud logging read "resource.type=k8s_container AND resource.labels.cluster_name=webenable-cms-prod" --limit=50

# Or use kubectl
kubectl logs -f deployment/backend -n webenable-cms-prod
kubectl logs -f deployment/frontend -n webenable-cms-prod
kubectl logs -f deployment/admin-panel -n webenable-cms-prod
```

## Security Configuration

### 1. Network Policies

```bash
# Apply network policies
kubectl apply -f k8s/gke/network-policies/
```

### 2. Pod Security Standards

```bash
# Enable Pod Security Standards
kubectl label namespace webenable-cms-prod \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted
```

### 3. Cloudflare Security

Configure these Cloudflare security features:

- **WAF Rules**: Custom rules for your application
- **Rate Limiting**: Protect against DDoS attacks
- **Bot Management**: Block malicious bots
- **Access Control**: IP allowlisting if needed

## Backup and Disaster Recovery

### 1. Database Backup

```bash
# Create backup script
cat > backup-couchdb.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
kubectl exec -n webenable-cms-prod deployment/couchdb -- \
  curl -X POST http://localhost:5984/_replicate \
  -H "Content-Type: application/json" \
  -d '{"source":"http://localhost:5984/webenable_cms","target":"/backup/webenable_cms_'$DATE'"}'
EOF

chmod +x backup-couchdb.sh
```

### 2. Configuration Backup

```bash
# Backup Kubernetes manifests
kubectl get all -n webenable-cms-prod -o yaml > backup-$(date +%Y%m%d).yaml

# Backup secrets (encrypted)
kubectl get secrets -n webenable-cms-prod -o yaml > secrets-backup-$(date +%Y%m%d).yaml
```

## Cost Optimization

### 1. Node Autoscaling

```bash
# Configure node autoscaling
gcloud container clusters update webenable-cms-prod \
  --zone us-central1-a \
  --enable-autoscaling \
  --min-nodes 1 \
  --max-nodes 5
```

### 2. Resource Optimization

- Use appropriate resource requests/limits
- Enable horizontal pod autoscaling
- Use spot instances for non-critical workloads
- Monitor and optimize storage usage

## Troubleshooting

### Common Issues

1. **SSL Certificate Issues**:
   ```bash
   kubectl describe certificate -n webenable-cms-prod
   kubectl describe order -n webenable-cms-prod
   ```

2. **DNS Resolution Issues**:
   ```bash
   # Check DNS propagation
   dig webenable-cms.com
   nslookup webenable-cms.com
   ```

3. **Load Balancer Issues**:
   ```bash
   # Check load balancer status
   kubectl get service -n webenable-cms-prod
   kubectl describe service -n webenable-cms-prod
   ```

4. **Pod Startup Issues**:
   ```bash
   # Check pod events
   kubectl describe pod -n webenable-cms-prod
   kubectl logs -f pod/POD_NAME -n webenable-cms-prod
   ```

### Support Resources

- [GKE Documentation](https://cloud.google.com/kubernetes-engine/docs)
- [Cloudflare Documentation](https://developers.cloudflare.com/)
- [cert-manager Documentation](https://cert-manager.io/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

## Maintenance

### Regular Tasks

1. **Update Kubernetes**: Monthly cluster upgrades
2. **Update Images**: Regular security patches
3. **Monitor Costs**: Weekly cost review
4. **Backup Verification**: Monthly backup testing
5. **Security Audits**: Quarterly security reviews

### Scaling Considerations

- Monitor resource usage and scale accordingly
- Use horizontal pod autoscaling for traffic spikes
- Consider multi-region deployment for global users
- Implement CDN for static assets 