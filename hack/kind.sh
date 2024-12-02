#!/bin/bash

# This script should be somewhat idempotent

[ -d hack ] || {
  echo "Run this script from the project root with: ./hack/$(basename $0)" >&2
  exit 1
}

set -x

./hack/keyring.sh

set -e

. .env
CERTMANAGER_VERSION="${CERTMANAGER_VERSION:-1.16.1}"

# create the Kind cluster
C="${KIND_NAME:-tsigan}"
kind create cluster -n "$C" || true
CTX="kind-$C"
K="kubectl --context $CTX"
$K cluster-info

# wait for nodes to be Ready
$K wait --timeout=1h --for=condition=Ready=true node -l node-role.kubernetes.io/control-plane
sleep 2

# deploy the telepresence ambassador
#if [[ "$($K get ns ambassador -o jsonpath='{.metadata.name}')" == "ambassador" ]]
#then
#    telepresence helm upgrade
#else
#    telepresence helm install
#fi

# deploy cert-manager
$K apply -f "https://github.com/cert-manager/cert-manager/releases/download/v$CERTMANAGER_VERSION/cert-manager.yaml"
$K -n cert-manager wait --timeout=1h --for=condition=Available=true deployment/cert-manager
$K -n cert-manager wait --timeout=1h --for=condition=Available=true deployment/cert-manager-webhook
sleep 2

# connect to the telepresence ambassador
#telepresence connect --name tsigan --context $CTX --namespace cert-manager
#telepresence status

# create a secret with the TSIG key
TSIG_SECRET_KEY="$(sed -n 's/PrivateKey: \(.*\)$/\1/p' .keyring/tsigan-ed25519.private)"
$K --dry-run=client -o yaml -n cert-manager create secret generic tsig-tsigan --from-literal=secret-key="$TSIG_SECRET_KEY" | $K apply -f -

# create ClusterIssuers for Let's Encrypt staging server
HOST_ADDRESS="$(docker network inspect kind | jq -r '.[0].IPAM.Config[] | select(.Gateway) | .Gateway')"
$K apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging-tsigan
spec:
  acme:
    email: tsigan-dev@enix.github.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-staging-tsigan-account-key
    solvers:
    - dns01:
        rfc2136:
          nameserver: "${HOST_ADDRESS}:53000"
          tsigKeyName: tsigan
          tsigAlgorithm: HMACSHA512
          tsigSecretSecretRef:
            name: tsig-tsigan
            key: secret-key
EOF
$K -n cert-manager wait --timeout=1h --for=condition=Ready=true clusterissuer.cert-manager.io/letsencrypt-staging-tsigan
$K apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging-tsigan-noauth
spec:
  acme:
    email: tsigan-dev@enix.github.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-staging-tsigan-account-key
    solvers:
    - dns01:
        rfc2136:
          nameserver: "${HOST_ADDRESS}:53000"
EOF
$K -n cert-manager wait --timeout=1h --for=condition=Ready=true clusterissuer.cert-manager.io/letsencrypt-staging-tsigan-noauth
