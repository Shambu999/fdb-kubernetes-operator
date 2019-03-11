/*
Copyright 2019 FoundationDB project authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FoundationDBClusterSpec defines the desired state of FoundationDBCluster
type FoundationDBClusterSpec struct {
	Version          string            `json:"version"`
	ProcessCounts    map[string]int    `json:"processCounts,omitempty"`
	ConnectionString string            `json:"connectionString,omitempty"`
	NextInstanceID   int               `json:"nextInstanceID,omitempty"`
	ReplicationMode  string            `json:"replicationMode,omitempty"`
	StorageEngine    string            `json:"storageEngine,omitempty"`
	Configured       bool              `json:"configured,omitempty"`
	PendingRemovals  map[string]string `json:"pendingRemovals,omitempty"`
}

// FoundationDBClusterStatus defines the observed state of FoundationDBCluster
type FoundationDBClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FoundationDBCluster is the Schema for the foundationdbclusters API
// +k8s:openapi-gen=true
type FoundationDBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FoundationDBClusterSpec   `json:"spec,omitempty"`
	Status FoundationDBClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FoundationDBClusterList contains a list of FoundationDBCluster
type FoundationDBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FoundationDBCluster `json:"items"`
}

// DesiredProcessCount returns the number of processes to configure with a given
// class
func (cluster *FoundationDBCluster) DesiredProcessCount(processClass string) int {
	count := cluster.Spec.ProcessCounts[processClass]
	var minimum int
	if processClass == "storage" {
		switch cluster.Spec.ReplicationMode {
		case "single":
			minimum = 1
		case "double":
			minimum = 3
		default:
			minimum = 1
		}
	}

	if minimum > count {
		return minimum
	}
	return count
}

// DesiredCoordinatorCount returns the number of coordinators to recruit for
// a cluster
func (cluster *FoundationDBCluster) DesiredCoordinatorCount() int {
	switch cluster.Spec.ReplicationMode {
	case "single":
		return 1
	case "double":
		return 3
	default:
		return 1
	}
}

func init() {
	SchemeBuilder.Register(&FoundationDBCluster{}, &FoundationDBClusterList{})
}