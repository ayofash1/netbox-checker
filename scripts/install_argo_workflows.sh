#!/bin/bash

set -euo pipefail

NAMESPACE="argo"
RELEASE_NAME="argo-workflows"
ARGO_VERSION="0.45.20"  # Helm chart version
ARGO_PASSWORD="admin"
ARGO_USERNAME="admin"
ARGO_SERVER_PORT=2746

# Install Helm if not installed
if ! command -v helm &> /dev/null; then
  echo "Helm not found. Please install Helm before running this script."
  exit 1
fi

# Add Argo Helm repo if not already added
if ! helm repo list | grep -q 'argo'; then
  helm repo add argo https://argoproj.github.io/argo-helm
fi
helm repo update

# Create namespace if not exists
kubectl get namespace "$NAMESPACE" >/dev/null 2>&1 || kubectl create namespace "$NAMESPACE"

# Create values.yaml dynamically
cat <<EOF > argo-values.yaml
server:
  enabled: true
  extraArgs:
    - --auth-mode
    - server
  secure: false
  serviceType: ClusterIP
  ingress:
    enabled: false
  auth:
    enabled: true
    basicAuth:
      enabled: true
      users:
        - username: $ARGO_USERNAME
          password: $(htpasswd -nbBC 10 "$ARGO_USERNAME" "$ARGO_PASSWORD" | sed 's/^.*://')

controller:
  workflowNamespaces:
    - "$NAMESPACE"
  workflowDefaults:
    spec:
      serviceAccountName: argo-workflow-sa
EOF

# Install or upgrade Argo Workflows via Helm
helm upgrade --install "$RELEASE_NAME" argo/argo-workflows \
  --namespace "$NAMESPACE" \
  --version "$ARGO_VERSION" \
  -f argo-values.yaml

# Create service account (if not exists)
kubectl get sa argo-workflow-sa -n "$NAMESPACE" >/dev/null 2>&1 || \
kubectl create sa argo-workflow-sa -n "$NAMESPACE"

# Create RBAC for the service account
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argo-workflow-sa-binding
subjects:
- kind: ServiceAccount
  name: argo-workflow-sa
  namespace: $NAMESPACE
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
EOF

# Port-forward Argo UI
echo "✅ Argo Workflows installed."
echo "🔐 Login with username: $ARGO_USERNAME and password: $ARGO_PASSWORD"
echo "🌐 Access the Argo UI at: http://localhost:$ARGO_SERVER_PORT"

kubectl -n "$NAMESPACE" port-forward svc/$RELEASE_NAME-server "$ARGO_SERVER_PORT":2746
