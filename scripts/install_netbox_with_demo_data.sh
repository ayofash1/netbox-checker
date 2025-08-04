#!/bin/bash

set -euo pipefail

# --- Configuration Variables ---
NETBOX_NAMESPACE="netbox"
NETBOX_HELM_RELEASE_NAME="my-netbox"
NETBOX_CHART_REPO="https://charts.netbox.oss.netboxlabs.com/"
NETBOX_CHART_NAME="netbox"
NETBOX_DEMO_DATA_VERSION="4.3"

SQL_DUMP_FILE="netbox-demo-v${NETBOX_DEMO_DATA_VERSION}.sql"
TEMP_SQL_PATH="/tmp/${SQL_DUMP_FILE}"
PG_USER="netbox"
PG_DB_NAME="netbox"
PG_HOST="localhost"

# --- Function to check for errors and exit ---
check_error() {
    if [ $? -ne 0 ]; then
        echo "ERROR: $1"
        exit 1
    fi
}

echo "--- Starting NetBox Helm Install + Demo Data Script ---"
echo "NetBox Namespace: ${NETBOX_NAMESPACE}"
echo "Helm Release: ${NETBOX_HELM_RELEASE_NAME}"
echo "Demo Data Version: ${NETBOX_DEMO_DATA_VERSION}"
echo "--------------------------------------------------------"

# --- 0. Create namespace if not exists ---
if ! kubectl get ns "${NETBOX_NAMESPACE}" >/dev/null 2>&1; then
  echo "0. Creating namespace: ${NETBOX_NAMESPACE}..."
  kubectl create ns "${NETBOX_NAMESPACE}"
fi

# --- 1. Add Helm repo (if not already added) ---
if ! helm repo list | grep -q "${NETBOX_CHART_REPO}"; then
  echo "1. Adding NetBox Helm repo..."
  helm repo add netbox "${NETBOX_CHART_REPO}"
  helm repo update
fi

# --- 2. Install NetBox via Helm if not already installed ---
if ! helm status "${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" >/dev/null 2>&1; then
  echo "2. Installing NetBox via Helm..."
  helm install "${NETBOX_HELM_RELEASE_NAME}" netbox/${NETBOX_CHART_NAME} \
    -n "${NETBOX_NAMESPACE}"
else
  echo "2. NetBox Helm release already installed. Skipping install."
fi

# --- 3. Wait for NetBox pods to be ready ---
echo "3. Waiting for NetBox pods to be ready..."
kubectl rollout status deployment "${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" --timeout=300s
echo "   NetBox pods ready."

# --- 4. Download demo SQL ---
echo "4. Downloading NetBox demo SQL file..."
DEMO_DATA_URL="https://raw.githubusercontent.com/netbox-community/netbox-demo-data/refs/heads/master/sql/${SQL_DUMP_FILE}"
wget -O "${SQL_DUMP_FILE}" "${DEMO_DATA_URL}"
check_error "Failed to download ${SQL_DUMP_FILE}"
echo "   Downloaded: ${SQL_DUMP_FILE}"

# --- 5. Discover PostgreSQL pod and credentials ---
echo "5. Discovering PostgreSQL pod and password..."
PG_POD_NAME=$(kubectl get pods -n "${NETBOX_NAMESPACE}" -l app.kubernetes.io/name=postgresql,app.kubernetes.io/instance="${NETBOX_HELM_RELEASE_NAME}" -o jsonpath='{.items[0].metadata.name}')
check_error "Failed to find PostgreSQL pod."

PG_PASSWORD=$(kubectl get secret "${NETBOX_HELM_RELEASE_NAME}-postgresql" -n "${NETBOX_NAMESPACE}" -o jsonpath='{.data.password}' | base64 --decode)
check_error "Failed to get PostgreSQL password."
echo "   PostgreSQL password retrieved."

# --- 6. Scale down NetBox deployment ---
echo "6. Scaling down NetBox to allow DB reset..."
kubectl scale deployment "${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" --replicas=0
kubectl wait --for=delete pod -l app.kubernetes.io/name=netbox,app.kubernetes.io/instance="${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" --timeout=120s || true

# --- 7. Copy SQL to DB pod ---
echo "7. Copying SQL to PostgreSQL pod..."
kubectl cp "${SQL_DUMP_FILE}" "${NETBOX_NAMESPACE}/${PG_POD_NAME}:${TEMP_SQL_PATH}"

# --- 8. Reset DB ---
echo "8. Loading SQL into PostgreSQL..."
kubectl exec -n "${NETBOX_NAMESPACE}" "${PG_POD_NAME}" -- bash -c "
  set -e
  export PGPASSWORD='${PG_PASSWORD}'

  echo 'Dropping existing DB...'
  psql -U ${PG_USER} -h ${PG_HOST} -d postgres -c \"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${PG_DB_NAME}' AND pid <> pg_backend_pid();\"
  psql -U ${PG_USER} -h ${PG_HOST} -d postgres -c \"DROP DATABASE IF EXISTS ${PG_DB_NAME};\"
  psql -U ${PG_USER} -h ${PG_HOST} -d postgres -c \"CREATE DATABASE ${PG_DB_NAME};\"

  echo 'Importing SQL dump...'
  psql -U ${PG_USER} -h ${PG_HOST} ${PG_DB_NAME} < ${TEMP_SQL_PATH}
  rm -f ${TEMP_SQL_PATH}
"

# --- 9. Scale NetBox back up ---
echo "9. Scaling NetBox back up..."
kubectl scale deployment "${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" --replicas=1
kubectl rollout status deployment "${NETBOX_HELM_RELEASE_NAME}" -n "${NETBOX_NAMESPACE}" --timeout=300s

# --- 10. Cleanup ---
echo "10. Cleaning up local SQL file..."
rm -f "${SQL_DUMP_FILE}"

# --- 11. Done ---
echo "--- ✅ NetBox Demo Environment Ready! ---"
echo "➡ Access NetBox via port-forward:"
echo "kubectl port-forward svc/${NETBOX_HELM_RELEASE_NAME} 8080:8080 -n ${NETBOX_NAMESPACE}"
echo "Then open: http://localhost:8080"
echo "Login: admin / admin"

