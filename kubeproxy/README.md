# kubeproxy-operator

kubeproxy-operator is a Kubernetes operator for managing kubeproxy.

## Running in a cluster

1. Create a kinder cluster
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

   You might have set the server ip in the KUBECONFIG to use localhost to reach the cluster, you will have to set the server to `localhost`. To find the port, run `docker inspect kind-control-plane-1 -f '{{(index (index .NetworkSettings.Ports "6443/tcp") 0).HostPort}}'` and check the port

   Edit your KUBECONFIG `vi $KUBECONFIG`
   ```yaml
   server: https://localhost:<port>
   ```

2. Set the Kubernetes Service host and port in manager.yaml ssh into the node and get the host and port. The command below should give the host.

   ```bash
   docker inspect kind-control-plane-1 | grep IPAddress
   ```

   Edit `config/manager/patches/apiserver_endpoint.path.yaml`

   ```yaml
   - name: KUBERNETES_SERVICE_HOST
     value: <your-kubernetes-ip>
   - name: KUBERNETES_SERVICE_PORT
     value: <your-kubernetes-port>
   ```


3. Build and deploy Docker image

   ```bash
   make docker-build

   docker image save controller:latest > controller-latest.tar
   kinder cp ./controller-latest.tar @cp1:/kind/
   kinder exec @cp1 -- ctr -n k8s.io image import /kind/controller-latest.tar
   
   IMG=docker.io/library/controller:latest make deploy
   ```

4. Install CRD

   ```bash
   make install
   kubectl apply -f config/samples/
   ```

5. KubeProxy should be up and running

   ```bash
   kubectl get kubeproxy -n kube-system
   kubectl get daemonset -n kube-system kube-proxy
   kubectl get nodes
   ```
