#!/bin/sh

# setup a kinder cluster with out kubeproxy, run the kubeproxy-operator locally to install

set -eux
cd "$(dirname "$0")" 

addons_root="${PWD}/../"

# --image kindest/node:v1.16.1

cat << EOF | kinder create cluster --config=/dev/stdin
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
- role: control-plane
  extraMounts:
    - containerPath: /addons
      hostPath: ${addons_root}
- role: worker
EOF

kinder do kubeadm-config
kinder do loadbalancer


# TODO(jrjohnson): This is a partial reimplementation of kinder do kubeadm-init
docker exec -it kind-control-plane  /kind/bin/kubeadm init --skip-phases="addon/kube-proxy"  --ignore-preflight-errors="FileContent--proc-sys-net-bridge-bridge-nf-call-iptables,Swap,SystemVerification" --config /kind/kubeadm.conf
kinder exec @all -- sysctl -w net.ipv4.conf.all.rp_filter=1

kinder cp @cp1:/etc/kubernetes/admin.conf $(kinder get kubeconfig-path)
export KUBECONFIG=$(kinder get kubeconfig-path)

kubectl apply -f=https://docs.projectcalico.org/v3.8/manifests/calico.yaml

make install
kubectl apply -f config/samples/

# TODO(jrjohnson): Should the operator be able to detect this on its own? 
ip_port=$(kubectl config view --minify | grep server | cut -f 2- -d ":" | tr -d " " | cut -d/ -f3)
export KUBERNETES_SERVICE_HOST=$(echo "${ip_port}" | cut -d: -f1)
export KUBERNETES_SERVICE_PORT=$(echo "${ip_port}" | cut -d: -f2)

make run
echo "KUBECONFIG=$(kinder get kubeconfig-path)"

