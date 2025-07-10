# Quick Start: GKE Autopilot + Cloudflare Free Tier

This guide will get you up and running with WebEnable CMS on GKE Autopilot using Cloudflare's free tier in under 30 minutes.

## 🚀 Prerequisites (5 minutes)

### 1. Google Cloud Account
- Create a Google Cloud account
- Enable billing
- Install Google Cloud CLI: `gcloud`

### 2. Cloudflare Free Account
- Sign up at https://cloudflare.com (free)
- Add your domain
- Get API token from Profile → API Tokens

### 3. Domain
- Purchase a domain (e.g., `webenable-cms.com`)
- Point nameservers to Cloudflare

## ⚡ Quick Deployment (15 minutes)

### 1. Set Environment Variables

```bash
export GCP_PROJECT_ID="your-project-id"
export GCP_CLUSTER_NAME="webenable-cms-autopilot"
export GCP_CLUSTER_REGION="us-central1"
export CLOUDFLARE_API_TOKEN="your-cloudflare-api-token"
export CLOUDFLARE_ZONE_ID="your-zone-id"
```

### 2. Run Deployment Script

```bash
# Make script executable
chmod +x k8s/gke-autopilot/deploy-autopilot.sh

# Deploy everything
./k8s/gke-autopilot/deploy-autopilot.sh deploy
```

### 3. Verify Deployment

```bash
# Check deployment status
./k8s/gke-autopilot/deploy-autopilot.sh verify

# Get deployment info
./k8s/gke-autopilot/deploy-autopilot.sh info
```

## 🌐 Access Your Application

Once deployed, access your application at:
- **Frontend**: https://webenable-cms.com
- **Admin Panel**: https://admin.webenable-cms.com
- **API**: https://api.webenable-cms.com

## 💰 Cost Breakdown

| Service | Cost |
|---------|------|
| GKE Autopilot | $50-150/month |
| Cloudflare | $0/month (free) |
| Domain | $1-2/month |
| **Total** | **$51-152/month** |

## 🔧 Management Commands

```bash
# Check deployment status
kubectl get pods -n webenable-cms-prod

# View logs
kubectl logs -f deployment/backend -n webenable-cms-prod

# Check auto-scaling
kubectl get hpa -n webenable-cms-prod

# Monitor costs
gcloud billing accounts list
```

## 🛠️ Troubleshooting

### Common Issues

1. **Cluster Creation Fails**
   ```bash
   # Check billing
   gcloud billing accounts list
   
   # Enable APIs
   gcloud services enable container.googleapis.com
   ```

2. **DNS Not Working**
   ```bash
   # Check DNS propagation
   dig webenable-cms.com
   
   # Verify Cloudflare settings
   # Go to Cloudflare Dashboard → DNS
   ```

3. **SSL Certificate Issues**
   ```bash
   # Check certificate status
   kubectl get certificates -n webenable-cms-prod
   
   # Check cert-manager logs
   kubectl logs -n cert-manager deployment/cert-manager
   ```

## 📈 Scaling

### Auto-scaling is enabled by default:
- **CPU**: Scales at 70% utilization
- **Memory**: Scales at 80% utilization
- **Min replicas**: 1
- **Max replicas**: 3-5 (depending on service)

### Manual scaling:
```bash
# Scale backend
kubectl scale deployment backend --replicas=3 -n webenable-cms-prod

# Scale frontend
kubectl scale deployment frontend --replicas=2 -n webenable-cms-prod
```

## 🔒 Security Features

### Included by default:
- ✅ HTTPS with Let's Encrypt certificates
- ✅ Cloudflare DDoS protection
- ✅ GKE Autopilot security
- ✅ Network policies
- ✅ Pod security standards

### Cloudflare Free Tier Security:
- ✅ SSL/TLS encryption
- ✅ Basic WAF rules
- ✅ Rate limiting (10,000 requests/day)
- ✅ Browser integrity check

## 📊 Monitoring

### Google Cloud Console:
- Go to https://console.cloud.google.com
- Navigate to Kubernetes Engine → Clusters
- View metrics and logs

### Application Monitoring:
```bash
# Check application health
curl https://api.webenable-cms.com/health

# View application logs
kubectl logs -f deployment/backend -n webenable-cms-prod
```

## 🧹 Cleanup

To remove everything:
```bash
./k8s/gke-autopilot/deploy-autopilot.sh cleanup
```

## 🆘 Support

### Documentation:
- [GKE Autopilot Guide](k8s/gke-autopilot/README.md)
- [Main Project README](../README.md)
- [GitLab CI Setup](../GITLAB_CI_SETUP.md)

### Useful Links:
- [GKE Autopilot Pricing](https://cloud.google.com/kubernetes-engine/pricing#autopilot)
- [Cloudflare Free Tier](https://www.cloudflare.com/plans/free/)
- [cert-manager Documentation](https://cert-manager.io/docs/)

## 🎯 Next Steps

1. **Customize Configuration**: Edit `k8s/overlays/production/` for your needs
2. **Set up Monitoring**: Configure alerts and dashboards
3. **Backup Strategy**: Implement regular database backups
4. **CI/CD Pipeline**: Set up GitLab CI for automated deployments
5. **Domain Management**: Configure additional subdomains if needed

---

**Need help?** Check the troubleshooting section or refer to the detailed documentation in `k8s/gke-autopilot/README.md`. 