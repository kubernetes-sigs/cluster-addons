/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"net"
	"regexp"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/coredns/corefile-migration/migration"
	"github.com/pkg/errors"
)

// getCoreDNSService fetches the CoreDNS Service
func getCoreDNSService(ctx context.Context, c client.Client) (*corev1.Service, error) {
	kubernetesService := &corev1.Service{}
	id := client.ObjectKey{Namespace: metav1.NamespaceDefault, Name: "kubernetes"}

	// Get the CoreDNS Service
	err := c.Get(ctx, id, kubernetesService)

	return kubernetesService, err
}

// getCoreDNSConfigMap fetches the CoreDNS ConfigMap
func getCoreDNSConfigMap(ctx context.Context, c client.Client, configMapName string) (*corev1.ConfigMap, error) {
	coreDNSConfigMap := &corev1.ConfigMap{}
	id := client.ObjectKey{Namespace: metav1.NamespaceSystem, Name: configMapName}

	// Get the CoreDNS ConfigMap
	err := c.Get(ctx, id, coreDNSConfigMap)

	return coreDNSConfigMap, err
}

// getCoreDNSDeployment fetches the CoreDNS Deployment
func getCoreDNSDeployment(ctx context.Context, c client.Client) (*appsv1.Deployment, error) {
	coreDNSDeploy := &appsv1.Deployment{}
	deployId := client.ObjectKey{Namespace: metav1.NamespaceSystem, Name: coreDNSName}

	// Get the CoreDNS Deployment
	err := c.Get(ctx, deployId, coreDNSDeploy)

	return coreDNSDeploy, err
}

// findDNSClusterIP tries to find the Cluster IP to be used by the DNS service
// It is usually the 10th address to the Kubernetes Service Cluster IP
// If the Kubernetes Service Cluster IP is not found, we default it to be "10.96.0.10"
func findDNSClusterIP(ctx context.Context, c client.Client) (string, error) {
	kubernetesService, err := getCoreDNSService(ctx, c)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}

	if apierrors.IsNotFound(err) {
		// If it cannot determine the Cluster IP, we default it to "10.96.0.10"
		return coreDNSIP, nil
	}

	ip := net.ParseIP(kubernetesService.Spec.ClusterIP)
	if ip == nil {
		return "", errors.Errorf("cannot parse kubernetes ClusterIP %q", kubernetesService.Spec.ClusterIP)
	}

	// The kubernetes Service ClusterIP is the 1st IP in the Service Subnet.
	// Increment the right-most byte by 9 to get to the 10th address, canonically used for CoreDNS.
	// This works for both IPV4, IPV6, and 16-byte IPV4 addresses.
	ip[len(ip)-1] += 9

	result := ip.String()
	klog.Infof("determined ClusterIP for CoreDNS should be %q", result)
	return result, nil
}

// getDNSDomain returns Kubernetes DNS cluster domain
// If it cannot determine the domain, we default it to "cluster.local"
// TODO (rajansandeep): find a better way to implement this?
func getDNSDomain() string {
	svc := "kubernetes.default.svc"

	cname, err := net.LookupCNAME(svc)
	if err != nil {
		// If it cannot determine the domain, we default it to "cluster.local"
		klog.Infof("determined DNS Domain for CoreDNS should be %q", coreDNSDomain)
		return coreDNSDomain
	}

	domain := strings.TrimPrefix(cname, svc)
	domain = strings.TrimSuffix(coreDNSDomain, ".")

	klog.Infof("determined DNS Domain for CoreDNS should be %q", domain)

	return domain
}

// getCorefile defines the Corefile to be generated for use in the CoreDNS ConfigMap
// depending on whether this is a fresh install or a change of CoreDNS version. This function
// tries to extract the existing Corefile from the cluster if available.
// If there is none available, we set CoreDNS to configure with the default Corefile, which is defined in each version.
func getCorefile(ctx context.Context, c client.Client) (string, error) {
	var err error

	// Get the CoreDNS Deployment
	coreDNSDeploy, err := getCoreDNSDeployment(ctx, c)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}

	if apierrors.IsNotFound(err) {
		klog.Infof("CoreDNS deployment not found")
		// If the CoreDNS Deployment isn't found, it is assumed that it is a new install
		// of CoreDNS and we proceed to use the default Corefile
		return "", nil
	}

	// Due to Kustomize, the CoreDNS ConfigMap name is hashed.
	// The name of the ConfigMap is extracted from the deployment
	var configMapName string
	coreDNSDeploySpec := coreDNSDeploy.Spec.Template.Spec
	for _, volume := range coreDNSDeploySpec.Volumes {
		if volume.Name == "config-volume" && volume.ConfigMap != nil {
			configMapName = volume.ConfigMap.Name
		}
	}

	if configMapName == "" {
		klog.Infof("CoreDNS deployment did not have config-volume")
		return "", nil
	}

	// Get the CoreDNS ConfigMap
	coreDNSConfigMap, err := getCoreDNSConfigMap(ctx, c, configMapName)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}

	if apierrors.IsNotFound(err) {
		// If the CoreDNS ConfigMap isn't found, use the default Corefile
		klog.Infof("CoreDNS deployment had config-volume %q, but it was not found", configMapName)
		return "", nil
	}

	// Get the Corefile
	corefile, ok := coreDNSConfigMap.Data["Corefile"]
	if !ok {
		klog.Infof("CoreDNS deployment had ConfigMap %q, but it did not have a Corefile entry", configMapName)
		return "", errors.New("unable to find the CoreDNS Corefile data")
	}
	klog.Infof("Found corefile: %q", corefile)

	return corefile, nil
}

// corefileMigration if necessary/possible, migrates the Corefile to reflect the latest changes, when the CoreDNS version is being upgraded.
func corefileMigration(ctx context.Context, c client.Client, coreDNSVersion, corefile string) (string, error) {
	var err error

	// Get the CoreDNS Deployment
	coreDNSDeploy, err := getCoreDNSDeployment(ctx, c)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", err
	}
	// The deployment is used to extract the CoreDNS Image version
	var coreDNSImage string
	coreDNSDeploySpec := coreDNSDeploy.Spec.Template.Spec
	for _, container := range coreDNSDeploySpec.Containers {
		if container.Name == coreDNSName {
			coreDNSImage = container.Image
		}
	}
	if coreDNSImage == "" {
		klog.Warningf("unable to find coredns container (%q) in pod", coreDNSName)
	}

	coreDNSImageParts := strings.Split(coreDNSImage, ":")
	currentCoreDNSVersion := coreDNSImageParts[len(coreDNSImageParts)-1]

	if currentCoreDNSVersion == "" {
		klog.Warningf("cannot extract coredns version from %q", coreDNSImage)
	}

	if currentCoreDNSVersion != coreDNSVersion {
		// Check if Corefile Migration is necessary and get the migrated Corefile
		// If the Corefile from the previous version is untouched, we can proceed to replace it with the
		// Corefile of the current version
		klog.Infof("from: %q  to: %q", currentCoreDNSVersion, coreDNSVersion)
		isDefault := migration.Default("", corefile)
		switch isDefault {
		case true:
			corefile = ""
			klog.Infof("the default Corefile will be applied")
		case false:
			corefile, err = performMigrationOfCorefile(ctx, c, corefile, currentCoreDNSVersion, coreDNSVersion)
			if err != nil {
				return "", errors.Errorf("unable to migrate the CoreDNS Corefile data: %v", err)
			}
			klog.Infof("determined Corefile for CoreDNS should be %q", corefile)
		}
	}
	return corefile, nil
}

// performMigrationOfCorefile checks for a possibility or requirement of Migrating the Corefile
// during a change of the CoreDNS Version. Currently, it is NOT recommended to support CoreDNS downgrades,
// but it has full capability of supporting CoreDNS upgrades.
// It first checks, whether the current Corefile is modified or not. If it the default Corefile, it skips migrating.
// In case the Corefile is modified, it goes through the entire Corefile to check for any configuration that can break
// the functionality of CoreDNS or is deprecated and migrates accordingly.
func performMigrationOfCorefile(ctx context.Context, c client.Client, corefile, fromVersion, toVersion string) (string, error) {
	if fromVersion == "" || fromVersion == toVersion {
		return corefile, nil
	}

	// Checks if the CoreDNS version is officially supported
	isVersionSupported, err := isCoreDNSVersionSupported(ctx, c)
	if !isVersionSupported {
		klog.Warningf("the CoreDNS Configuration will not be migrated due to unsupported version of CoreDNS. " +
			"The existing CoreDNS Corefile configuration and deployment has been retained.")
		return corefile, err
	}

	migratedCorefile, err := migration.Migrate(fromVersion, toVersion, corefile, false)
	if err != nil {
		klog.Warningf("the CoreDNS Configuration was not migrated: %v. The existing CoreDNS Corefile configuration has been retained.", err)
		return corefile, err
	}

	// show the migration changes
	klog.Infof("the CoreDNS configuration has been migrated and applied: %v.", migratedCorefile)

	return migratedCorefile, nil
}

var (
	// imageDigestMatcher is used to match the SHA256 digest from the ImageID of the CoreDNS pods
	imageDigestMatcher = regexp.MustCompile(`^.*(?i:sha256:([[:alnum:]]{64}))$`)
)

func isCoreDNSVersionSupported(ctx context.Context, c client.Client) (bool, error) {
	var err error
	isValidVersion := true

	coreDNSPodList := &corev1.PodList{}
	err = c.List(ctx, coreDNSPodList, client.InNamespace(metav1.NamespaceSystem), client.MatchingLabels(map[string]string{"k8s-app": "kube-dns"}))
	if err != nil && !apierrors.IsNotFound(err) {
		return false, errors.Errorf("error getting CoreDNS Pod list: %v", err)
	}

	for _, pod := range coreDNSPodList.Items {
		imageID := imageDigestMatcher.FindStringSubmatch(pod.Status.ContainerStatuses[0].ImageID)
		klog.Info(imageID)
		if len(imageID) != 2 {
			return false, errors.Errorf("unable to match SHA256 digest ID in %q", pod.Status.ContainerStatuses[0].ImageID)
		}
		// The actual digest should be at imageID[1]
		if !migration.Released(imageID[1]) {
			isValidVersion = false
		}
	}

	return isValidVersion, nil
}

// prepCorefileFormat indents the output of the Corefile and replaces tabs with spaces
// to neatly format the configmap, making it readable and ensure there is no error converting YAML to JSON.
func prepCorefileFormat(s string, indentation int) string {
	var str []string
	if s == "" {
		return ""
	}

	for _, line := range strings.Split(s, "\n") {
		indented := strings.Repeat(" ", indentation) + line
		str = append(str, indented)
	}
	corefile := strings.Join(str, "\n")
	corefile = strings.TrimSpace(corefile)

	return strings.Replace(corefile, "\t", "   ", -1)
}

const (
	coreDNSDomain = "cluster.local"
	coreDNSIP     = "10.96.0.10"
	coreDNSName   = "coredns"
)
