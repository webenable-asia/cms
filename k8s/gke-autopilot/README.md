# GKE Autopilot Deployment with Cloudflare Free Tier

This guide covers deploying WebEnable CMS on Google Kubernetes Engine (GKE) Autopilot with Cloudflare free tier managing DNS and SSL certificates.

## Architecture Overview

```
Internet → Cloudflare Free Tier → GKE Autopilot Load Balancer → Ingress → Services → Pods
                ↓
        SSL Termination + DNS Management (Free)
```

## Why GKE Autopilot + Cloudflare Free Tier?

### GKE Autopilot Benefits
- **Fully Managed**: No node management required
- **Cost Effective**: Pay only for running pods
- **Security**: Built-in security features
- **Auto-scaling**: Automatic pod scaling
- **Compliance**: SOC, PCI, HIPAA compliance

### Cloudflare Free Tier Benefits
- **Free SSL Certificates**: Automatic HTTPS
- **Global CDN**: 200+ data centers
- **DDoS Protection**: Basic protection included
- **DNS Management**: Unlimited DNS records
- **Security Features**: Basic WAF included

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

### 2. Cloudflare Free Account Setup

1. **Create Free Account**:
   - Go to https://cloudflare.com
   - Sign up for a free account
   - Add your domain (e.g., `webenable-cms.com`)

2. **Update Nameservers**:
   - Copy the provided nameservers
   - Update at your domain registrar
   - Wait for DNS propagation (up to 24 hours)

3. **Get Zone ID and API Token**:
   - Zone ID: Found in Cloudflare dashboard → Overview
   - API Token: Profile → API Tokens → Create Custom Token

### 3. Required Tools

```bash
# Install kubectl
gcloud components install kubectl

# Install helm
curl https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz | tar xz
sudo mv linux-amd64/helm /usr/local/bin/
```

## GKE Autopilot Cluster Setup

### 1. Create GKE Autopilot Cluster

```bash
# Create production cluster
gcloud container clusters create-auto webenable-cms-autopilot \
  --region us-central1 \
  --project YOUR_PROJECT_ID \
  --release-channel regular \
  --enable-private-nodes \
  --enable-master-authorized-networks \
  --master-authorized-networks 0.0.0.0/0 \
  --enable-autopilot

# Get cluster credentials
gcloud container clusters get-credentials webenable-cms-autopilot --region us-central1
```

### 2. Configure Workload Identity (Optional)

```bash
# Create service account for GitLab CI
gcloud iam service-accounts create gitlab-ci-autopilot \
  --display-name="GitLab CI Service Account for Autopilot"

# Grant necessary permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:gitlab-ci-autopilot@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/container.developer"

# Create key for GitLab CI
gcloud iam service-accounts keys create gitlab-ci-autopilot-key.json \
  --iam-account=gitlab-ci-autopilot@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

## Cloudflare Free Tier Configuration

### 1. DNS Records Setup

Create these DNS records in Cloudflare (all free):

```bash
# A Records (point to GKE load balancer IP)
webenable-cms.com        A    <GKE_LOAD_BALANCER_IP>
www.webenable-cms.com    A    <GKE_LOAD_BALANCER_IP>
api.webenable-cms.com    A    <GKE_LOAD_BALANCER_IP>
admin.webenable-cms.com  A    <GKE_LOAD_BALANCER_IP>

# CNAME Records
*.webenable-cms.com      CNAME webenable-cms.com
```

### 2. Cloudflare Free Tier Settings

Configure these settings in Cloudflare dashboard:

- **SSL/TLS**: Full (strict) - Free
- **Always Use HTTPS**: On - Free
- **HSTS**: Enabled - Free
- **Security Level**: Medium - Free
- **Rate Limiting**: Basic (10,000 requests/day) - Free
- **Bot Fight Mode**: Disabled (Premium feature)
- **Browser Integrity Check**: Enabled - Free

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
kubectl apply -f k8s/gke-autopilot/cluster-issuer.yaml
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
GCP_CLUSTER_NAME=webenable-cms-autopilot
GCP_CLUSTER_REGION=us-central1
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

## Cost Optimization for Autopilot

### 1. Resource Optimization

```yaml
# Optimize resource requests for Autopilot
resources:
  requests:
    cpu: "250m"      # Minimum for Autopilot
    memory: "512Mi"  # Minimum for Autopilot
  limits:
    cpu: "1000m"     # Reasonable limit
    memory: "1Gi"    # Reasonable limit
```

### 2. Pod Scaling

```yaml
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: backend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: backend
  minReplicas: 1
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 3. Storage Optimization

- Use `standard-rwo` storage class (cheaper)
- Implement data retention policies
- Regular cleanup of old data

## Monitoring and Logging

### 1. Google Cloud Monitoring (Free Tier)

```bash
# Enable monitoring (included in Autopilot)
# View metrics in Google Cloud Console
# https://console.cloud.google.com/monitoring
```

### 2. Application Logs

```bash
# View application logs
gcloud logging read "resource.type=k8s_container AND resource.labels.cluster_name=webenable-cms-autopilot" --limit=50

# Or use kubectl
kubectl logs -f deployment/backend -n webenable-cms-prod
kubectl logs -f deployment/frontend -n webenable-cms-prod
kubectl logs -f deployment/admin-panel -n webenable-cms-prod
```

## Security Configuration

### 1. Network Policies

```bash
# Apply network policies
kubectl apply -f k8s/gke-autopilot/network-policies/
```

### 2. Pod Security Standards

```bash
# Enable Pod Security Standards (built into Autopilot)
# Autopilot enforces restricted pod security by default
```

### 3. Cloudflare Security (Free Tier)

Configure these Cloudflare security features:

- **WAF Rules**: Basic rules included
- **Rate Limiting**: 10,000 requests/day (free)
- **Bot Management**: Basic protection
- **Access Control**: IP allowlisting (if needed)

## Backup and Disaster Recovery

### 1. Database Backup

```bash
# Create backup script for Autopilot
cat > backup-couchdb-autopilot.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
kubectl exec -n webenable-cms-prod deployment/couchdb -- \
  curl -X POST http://localhost:5984/_replicate \
  -H "Content-Type: application/json" \
  -d '{"source":"http://localhost:5984/webenable_cms","target":"/backup/webenable_cms_'$DATE'"}'
EOF

chmod +x backup-couchdb-autopilot.sh
```

### 2. Configuration Backup

```bash
# Backup Kubernetes manifests
kubectl get all -n webenable-cms-prod -o yaml > backup-autopilot-$(date +%Y%m%d).yaml

# Backup secrets (encrypted)
kubectl get secrets -n webenable-cms-prod -o yaml > secrets-backup-autopilot-$(date +%Y%m%d).yaml
```

## Troubleshooting

### Common Issues

1. **Autopilot Resource Limits**:
   ```bash
   # Check pod events for resource issues
   kubectl describe pod -n webenable-cms-prod
   ```

2. **SSL Certificate Issues**:
   ```bash
   kubectl describe certificate -n webenable-cms-prod
   kubectl describe order -n webenable-cms-prod
   ```

3. **DNS Resolution Issues**:
   ```bash
   # Check DNS propagation
   dig webenable-cms.com
   nslookup webenable-cms.com
   ```

4. **Cloudflare Free Tier Limits**:
   - Rate limiting: 10,000 requests/day
   - WAF rules: Basic rules only
   - Bot management: Not available

### Support Resources

- [GKE Autopilot Documentation](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview)
- [Cloudflare Free Tier Documentation](https://developers.cloudflare.com/fundamentals/get-started/basic-tasks/free-tier/)
- [cert-manager Documentation](https://cert-manager.io/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

## Maintenance

### Regular Tasks

1. **Update Kubernetes**: Automatic with Autopilot
2. **Update Images**: Regular security patches
3. **Monitor Costs**: Weekly cost review
4. **Backup Verification**: Monthly backup testing
5. **Security Audits**: Quarterly security reviews

### Scaling Considerations

- Monitor resource usage and scale accordingly
- Use horizontal pod autoscaling for traffic spikes
- Consider multi-region deployment for global users
- Implement CDN for static assets (Cloudflare free tier)

## Cost Estimation

### Monthly Costs (Estimated)

- **GKE Autopilot**: $50-150/month (depending on usage)
- **Cloudflare**: $0/month (free tier)
- **Domain**: $10-15/year
- **Total**: $50-150/month

### Cost Optimization Tips

1. **Right-size resources**: Use minimum required resources
2. **Auto-scaling**: Scale down during low traffic
3. **Storage optimization**: Use appropriate storage classes
4. **Monitor usage**: Regular cost monitoring
5. **Cleanup**: Remove unused resources 