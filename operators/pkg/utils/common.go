/*


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

// Package utils collects all the logic shared between different controllers
package utils

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ParseDockerDirectory returns a valid Docker image directory.
func ParseDockerDirectory(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.ToLower(reg.ReplaceAllString(name, ""))
}

// CheckLabels verifies whether a namespace is characterized by a set of required labels.
func CheckLabels(ns *corev1.Namespace, matchLabels map[string]string) bool {
	for key, value := range matchLabels {
		if v1, ok := ns.Labels[key]; !ok || v1 != value {
			return false
		}
	}
	return true
}

// CheckSelectorLabel checks if the given namespace belongs to the whitelisted namespaces where to perform reconciliation.
func CheckSelectorLabel(ctx context.Context, k8sClient client.Client, namespaceName string, matchLabels map[string]string) (bool, error) {
	ns := corev1.Namespace{}
	namespaceLookupKey := types.NamespacedName{
		Name:      namespaceName,
		Namespace: "",
	}

	// It performs reconciliation only if the InstanceSnapshot belongs to whitelisted namespaces
	// by checking the existence of keys in the namespace of the InstanceSnapshot.
	if err := k8sClient.Get(ctx, namespaceLookupKey, &ns); err == nil {
		if !CheckLabels(&ns, matchLabels) {
			klog.Infof("Namespace %s does not meet the selector labels", namespaceName)
			return false, nil
		}
	} else {
		return false, fmt.Errorf("error when retrieving the InstanceSnapshot namespace -> %w", err)
	}

	klog.Info("Namespace " + namespaceName + " met the selector labels")
	return true, nil
}