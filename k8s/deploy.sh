#!/bin/bash

# WebEnable CMS Kubernetes Deployment Script
# Usage: ./deploy.sh [environment] [action]
# Environments: development, staging, production
# Actions: apply, delete, build, validate

set -e

ENVIRONMENT=${1:-development}
ACTION=${2:-apply}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
    print_error "Invalid environment: $ENVIRONMENT"
    echo "Valid environments: development, staging, production"
    exit 1
fi

# Validate action
if [[ ! "$ACTION" =~ ^(apply|delete|build|validate)$ ]]; then
    print_error "Invalid action: $ACTION"
    echo "Valid actions: apply, delete, build, validate"
    exit 1
fi

print_status "Deploying WebEnable CMS to $ENVIRONMENT environment..."

# Set namespace based on environment
NAMESPACE="webenable-cms-$ENVIRONMENT"

case $ACTION in
    "build")
        print_status "Building manifests for $ENVIRONMENT environment..."
        kustomize build k8s/overlays/$ENVIRONMENT
        print_success "Manifests built successfully"
        ;;
    
    "validate")
        print_status "Validating manifests for $ENVIRONMENT environment..."
        kustomize build k8s/overlays/$ENVIRONMENT | kubectl apply --dry-run=client -f -
        print_success "Manifests validated successfully"
        ;;
    
    "apply")
        print_status "Applying manifests to $ENVIRONMENT environment..."
        
        # Create namespace if it doesn't exist
        kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
        
        # Apply the manifests
        kustomize build k8s/overlays/$ENVIRONMENT | kubectl apply -f -
        
        print_success "Deployment applied successfully"
        
        # Wait for deployments to be ready
        print_status "Waiting for deployments to be ready..."
        kubectl wait --for=condition=available --timeout=300s deployment/backend -n $NAMESPACE
        kubectl wait --for=condition=available --timeout=300s deployment/frontend -n $NAMESPACE
        kubectl wait --for=condition=available --timeout=300s deployment/admin-panel -n $NAMESPACE
        kubectl wait --for=condition=available --timeout=300s deployment/couchdb -n $NAMESPACE
        kubectl wait --for=condition=available --timeout=300s deployment/valkey -n $NAMESPACE
        
        print_success "All deployments are ready!"
        
        # Show service URLs
        print_status "Service URLs:"
        kubectl get svc -n $NAMESPACE
        
        # Show ingress if available
        print_status "Ingress configuration:"
        kubectl get ingress -n $NAMESPACE
        ;;
    
    "delete")
        print_warning "Deleting WebEnable CMS from $ENVIRONMENT environment..."
        kustomize build k8s/overlays/$ENVIRONMENT | kubectl delete -f -
        print_success "Deployment deleted successfully"
        ;;
esac

print_success "Operation completed successfully!" 