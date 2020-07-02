#!/bin/sh

# setup a kinder cluster with out kubeproxy, run the kubeproxy-operator locally to install

set -eux
cd "$(dirname "$0")"

addons_root="${PWD}/../"

# --image kindest/node:v1.16.1

kinder create cluster --worker-nodes 1 --image=kindest/node:v1.18.0

kinder do kubeadm-config
kinder do loadbalancer


# # TODO(jrjohnson): This is a partial reimplementation of kinder do kubeadm-init
docker exec -it kind-control-plane-1  /kind/bin/kubeadm init --skip-phases="addon/kube-proxy"  --ignore-preflight-errors="FileContent--proc-sys-net-bridge-bridge-nf-call-iptables,Swap,SystemVerification" --config /kind/kubeadm.conf
kinder exec @all -- sysctl -w net.ipv4.conf.all.rp_filter=1

kinder cp @cp1:/etc/kubernetes/admin.conf $(kinder get kubeconfig-path)
export KUBECONFIG=$(kinder get kubeconfig-path)

export HOSTPORT=$(docker inspect kind-control-plane-1 -f '{{(index (index .NetworkSettings.Ports "6443/tcp") 0).HostPort}}')
export IPADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' kind-control-plane-1)
sed -i '' 's/'"${IPADDRESS}"'.*$/localhost:'"${HOSTPORT}"'/' $KUBECONFIG

kubectl apply -f=https://docs.projectcalico.org/v3.8/manifests/calico.yaml

make install
kubectl apply -f config/samples/

# TODO(jrjohnson): Should the operator be able to detect this on its own?
ip_port=$(kubectl config view --minify | grep server | cut -f 2- -d ":" | tr -d " " | cut -d/ -f3)
export KUBERNETES_SERVICE_HOST=${IPADDRESS}
export KUBERNETES_SERVICE_PORT=6443

make run
echo "KUBECONFIG=$(kinder get kubeconfig-path)"
