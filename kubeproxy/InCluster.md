## This Readme documents how to run the KubeProxy operator in a kinder cluster

# 1. Create a kinder cluster
Ensure kinder is installed. [Installation docs](https://github.com/kubernetes/kubeadm/blob/master/kinder/README.md)

```bash
kinder create cluster --image=kindest/node:v1.18.0

kinder do kubeadm-config
kinder do loadbalancer

docker exec -it kind-control-plane-1  /kind/bin/kubeadm init --skip-phases="addon/kube-proxy"  --ignore-preflight-errors="FileContent--proc-sys-net-bridge-bridge-nf-call-iptables,Swap,SystemVerification" --config /kind/kubeadm.conf
kinder exec @all -- sysctl -w net.ipv4.conf.all.rp_filter=1

kinder cp @cp1:/etc/kubernetes/admin.conf $(kinder get kubeconfig-path)
export KUBECONFIG=$(kinder get kubeconfig-path)
```

You might have set the server ip in the KUBECONFIG to use localhost to reach the cluster, `insecure-skip-tls-verify` to true, and delete the ca certificate. To find the port, run `docker ps | grep kind` and check the port

> insecure-skip-tls-verify: true
> server: https://127.0.0.1:<port>

2. Set the Kubernetes Service host and port in manager.yaml
ssh into the node and get the host and port.
The command below should give the host.
```bash
docker inspect kind-control-plane-1 | grep IPAddress
```

Replace it in the `manager.yaml`

>- name: KUBERNETES_SERVICE_HOST
>  value: "172.17.0.2"
>- name: KUBERNETES_SERVICE_PORT
>  value: "6443"


3. Build and deploy Docker image
```bash
make docker-build

make deploy
```

4. Install CRD

```bash
make install
kubectl apply -f config/samples/
```

5. KubeProxy should be up and running
