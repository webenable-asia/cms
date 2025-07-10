#!/bin/bash

# Script to generate GitLab registry secret for Kubernetes
# Usage: ./generate-registry-secret.sh [username] [password] [namespace]

set -e

USERNAME=${1:-$CI_REGISTRY_USER}
PASSWORD=${2:-$CI_REGISTRY_PASSWORD}
NAMESPACE=${3:-webenable-cms}

if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo "Error: Username and password are required"
    echo "Usage: $0 [username] [password] [namespace]"
    echo "Or set CI_REGISTRY_USER and CI_REGISTRY_PASSWORD environment variables"
    exit 1
fi

# Create Docker config JSON
DOCKER_CONFIG=$(cat <<EOF
{
  "auths": {
    "registry.gitlab.com": {
      "auth": "$(echo -n "$USERNAME:$PASSWORD" | base64)"
    }
  }
}
EOF
)

# Encode the config
ENCODED_CONFIG=$(echo "$DOCKER_CONFIG" | base64 -w 0)

# Create the secret YAML
cat <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: gitlab-registry-secret
  namespace: $NAMESPACE
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: $ENCODED_CONFIG
EOF

echo "✅ GitLab registry secret generated for namespace: $NAMESPACE"
echo "💡 Apply this secret to your Kubernetes cluster:"
echo "   kubectl apply -f - <<< '$(cat <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: gitlab-registry-secret
  namespace: $NAMESPACE
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: $ENCODED_CONFIG
EOF
)'" 