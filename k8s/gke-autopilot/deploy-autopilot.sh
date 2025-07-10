#!/bin/bash

# GKE Autopilot Deployment Script for WebEnable CMS with Cloudflare Free Tier
# This script deploys the application to Google Kubernetes Engine Autopilot

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID=${GCP_PROJECT_ID:-"your-project-id"}
CLUSTER_NAME=${GCP_CLUSTER_NAME:-"webenable-cms-autopilot"}
CLUSTER_REGION=${GCP_CLUSTER_REGION:-"us-central1"}
NAMESPACE=${NAMESPACE:-"webenable-cms-prod"}
ENVIRONMENT=${ENVIRONMENT:-"production"}

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists gcloud; then
        print_error "Google Cloud CLI (gcloud) is not installed"
        exit 1
    fi
    
    if ! command_exists kubectl; then
        print_error "kubectl is not installed"
        exit 1
    fi
    
    if ! command_exists helm; then
        print_warning "Helm is not installed. Installing..."
        curl https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz | tar xz
        sudo mv linux-amd64/helm /usr/local/bin/
    fi
    
    print_success "Prerequisites check completed"
}

# Function to authenticate with GCP
authenticate_gcp() {
    print_status "Authenticating with Google Cloud Platform..."
    
    # Check if already authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        gcloud auth login
    fi
    
    # Set project
    gcloud config set project "$PROJECT_ID"
    
    # Get cluster credentials
    gcloud container clusters get-credentials "$CLUSTER_NAME" --region "$CLUSTER_REGION"
    
    print_success "GCP authentication completed"
}

# Function to create Autopilot cluster
create_autopilot_cluster() {
    print_status "Creating GKE Autopilot cluster..."
    
    # Check if cluster already exists
    if ! gcloud container clusters describe "$CLUSTER_NAME" --region="$CLUSTER_REGION" 2>/dev/null; then
        gcloud container clusters create-auto "$CLUSTER_NAME" \
            --region "$CLUSTER_REGION" \
            --project "$PROJECT_ID" \
            --release-channel regular \
            --enable-private-nodes \
            --enable-master-authorized-networks \
            --master-authorized-networks 0.0.0.0/0 \
            --enable-autopilot
        
        print_success "Autopilot cluster created successfully"
    else
        print_warning "Cluster already exists, skipping creation"
    fi
}

# Function to install cert-manager
install_cert_manager() {
    print_status "Installing cert-manager..."
    
    # Add helm repository
    helm repo add jetstack https://charts.jetstack.io
    helm repo update
    
    # Install cert-manager
    helm upgrade --install cert-manager jetstack/cert-manager \
        --namespace cert-manager \
        --create-namespace \
        --version v1.13.0 \
        --set installCRDs=true \
        --wait
    
    print_success "cert-manager installed successfully"
}

# Function to configure Cloudflare DNS
configure_cloudflare() {
    print_status "Configuring Cloudflare DNS..."
    
    # Create Cloudflare API secret
    if [ -n "$CLOUDFLARE_API_TOKEN" ]; then
        kubectl create secret generic cloudflare-api-token \
            --from-literal=api-token="$CLOUDFLARE_API_TOKEN" \
            -n cert-manager \
            --dry-run=client -o yaml | kubectl apply -f -
        
        # Apply ClusterIssuer
        kubectl apply -f k8s/gke-autopilot/cluster-issuer.yaml
        
        print_success "Cloudflare DNS configured"
    else
        print_warning "CLOUDFLARE_API_TOKEN not set. Skipping Cloudflare configuration."
    fi
}

# Function to create static IP
create_static_ip() {
    print_status "Creating static IP address..."
    
    # Check if static IP already exists
    if ! gcloud compute addresses describe webenable-cms-autopilot-ip --region="$CLUSTER_REGION" 2>/dev/null; then
        gcloud compute addresses create webenable-cms-autopilot-ip \
            --region="$CLUSTER_REGION" \
            --description="Static IP for WebEnable CMS Autopilot"
    fi
    
    # Get the IP address
    STATIC_IP=$(gcloud compute addresses describe webenable-cms-autopilot-ip --region="$CLUSTER_REGION" --format="value(address)")
    print_success "Static IP created: $STATIC_IP"
    
    # Update DNS records in Cloudflare
    if [ -n "$CLOUDFLARE_API_TOKEN" ] && [ -n "$CLOUDFLARE_ZONE_ID" ]; then
        print_status "Updating Cloudflare DNS records..."
        
        # Update A records
        for host in "webenable-cms.com" "www.webenable-cms.com" "api.webenable-cms.com" "admin.webenable-cms.com"; do
            curl -X PUT "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/dns_records" \
                -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
                -H "Content-Type: application/json" \
                --data "{
                    \"type\": \"A\",
                    \"name\": \"$host\",
                    \"content\": \"$STATIC_IP\",
                    \"proxied\": true
                }" 2>/dev/null || print_warning "Failed to update DNS record for $host"
        done
        
        print_success "DNS records updated"
    fi
}

# Function to deploy application
deploy_application() {
    print_status "Deploying application to $ENVIRONMENT environment..."
    
    # Create namespace
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply storage classes
    kubectl apply -f k8s/gke-autopilot/storage-class.yaml
    
    # Apply base configuration
    kubectl apply -k k8s/base/
    
    # Apply environment-specific configuration
    kubectl apply -k "k8s/overlays/$ENVIRONMENT/"
    
    # Apply Autopilot-specific ingress
    kubectl apply -f k8s/gke-autopilot/ingress.yaml
    
    # Apply horizontal pod autoscalers
    kubectl apply -f k8s/gke-autopilot/hpa.yaml
    
    print_success "Application deployed successfully"
}

# Function to wait for deployment
wait_for_deployment() {
    print_status "Waiting for deployment to be ready..."
    
    # Wait for all deployments
    kubectl wait --for=condition=available --timeout=300s deployment/backend -n "$NAMESPACE"
    kubectl wait --for=condition=available --timeout=300s deployment/frontend -n "$NAMESPACE"
    kubectl wait --for=condition=available --timeout=300s deployment/admin-panel -n "$NAMESPACE"
    kubectl wait --for=condition=available --timeout=300s deployment/couchdb -n "$NAMESPACE"
    kubectl wait --for=condition=available --timeout=300s deployment/valkey -n "$NAMESPACE"
    
    print_success "All deployments are ready"
}

# Function to verify deployment
verify_deployment() {
    print_status "Verifying deployment..."
    
    # Check pod status
    kubectl get pods -n "$NAMESPACE"
    
    # Check services
    kubectl get services -n "$NAMESPACE"
    
    # Check ingress
    kubectl get ingress -n "$NAMESPACE"
    
    # Check certificates
    kubectl get certificates -n "$NAMESPACE"
    
    # Check HPA
    kubectl get hpa -n "$NAMESPACE"
    
    print_success "Deployment verification completed"
}

# Function to show deployment info
show_deployment_info() {
    print_status "Deployment Information:"
    echo "  Project ID: $PROJECT_ID"
    echo "  Cluster: $CLUSTER_NAME"
    echo "  Region: $CLUSTER_REGION"
    echo "  Namespace: $NAMESPACE"
    echo "  Environment: $ENVIRONMENT"
    
    # Get static IP
    STATIC_IP=$(gcloud compute addresses describe webenable-cms-autopilot-ip --region="$CLUSTER_REGION" --format="value(address)" 2>/dev/null || echo "Not created")
    echo "  Static IP: $STATIC_IP"
    
    echo ""
    print_status "Access URLs:"
    echo "  Frontend: https://webenable-cms.com"
    echo "  Admin Panel: https://admin.webenable-cms.com"
    echo "  API: https://api.webenable-cms.com"
    
    echo ""
    print_status "Cost Optimization Features:"
    echo "  ✅ GKE Autopilot: Pay only for running pods"
    echo "  ✅ Cloudflare Free Tier: $0/month"
    echo "  ✅ Auto-scaling: Automatic pod scaling"
    echo "  ✅ Storage optimization: Standard storage class"
    
    echo ""
    print_status "Useful Commands:"
    echo "  kubectl get pods -n $NAMESPACE"
    echo "  kubectl logs -f deployment/backend -n $NAMESPACE"
    echo "  kubectl describe ingress webenable-cms-ingress -n $NAMESPACE"
    echo "  kubectl get hpa -n $NAMESPACE"
}

# Function to estimate costs
estimate_costs() {
    print_status "Cost Estimation (Monthly):"
    echo "  GKE Autopilot: $50-150/month (depending on usage)"
    echo "  Cloudflare: $0/month (free tier)"
    echo "  Domain: $1-2/month"
    echo "  Total: $51-152/month"
    echo ""
    echo "Cost optimization tips:"
    echo "  - Use minimum resource requests (250m CPU, 512Mi memory)"
    echo "  - Enable auto-scaling to scale down during low traffic"
    echo "  - Use standard storage class for cost savings"
    echo "  - Monitor usage with Google Cloud Console"
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up..."
    
    # Delete application
    kubectl delete -k "k8s/overlays/$ENVIRONMENT/" --ignore-not-found=true
    kubectl delete -k k8s/base/ --ignore-not-found=true
    
    # Delete ingress
    kubectl delete -f k8s/gke-autopilot/ingress.yaml --ignore-not-found=true
    
    # Delete HPA
    kubectl delete -f k8s/gke-autopilot/hpa.yaml --ignore-not-found=true
    
    # Delete namespace
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
    
    print_success "Cleanup completed"
}

# Main function
main() {
    case "${1:-deploy}" in
        "deploy")
            check_prerequisites
            authenticate_gcp
            create_autopilot_cluster
            install_cert_manager
            configure_cloudflare
            create_static_ip
            deploy_application
            wait_for_deployment
            verify_deployment
            show_deployment_info
            estimate_costs
            ;;
        "cleanup")
            authenticate_gcp
            cleanup
            ;;
        "info")
            authenticate_gcp
            show_deployment_info
            ;;
        "verify")
            authenticate_gcp
            verify_deployment
            ;;
        "costs")
            estimate_costs
            ;;
        *)
            echo "Usage: $0 {deploy|cleanup|info|verify|costs}"
            echo ""
            echo "Commands:"
            echo "  deploy   - Deploy the application (default)"
            echo "  cleanup  - Remove the application"
            echo "  info     - Show deployment information"
            echo "  verify   - Verify deployment status"
            echo "  costs    - Show cost estimation"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 