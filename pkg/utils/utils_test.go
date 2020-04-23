package utils

import (
	"reflect"
	"testing"

	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestContainsPullSecret(t *testing.T) {
	superset := []corev1.LocalObjectReference{{Name: "foo"}, {Name: "bar"}}
	type args struct {
		pullSecrets []corev1.LocalObjectReference
		ps          corev1.LocalObjectReference
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Contains pull secret",
			args{
				pullSecrets: superset,
				ps:          corev1.LocalObjectReference{Name: "foo"},
			},
			true,
		},
		{
			"Does not contain pull secret",
			args{
				pullSecrets: superset,
				ps:          corev1.LocalObjectReference{Name: "baz"},
			},
			false,
		},
		{
			"Empty list",
			args{
				pullSecrets: []corev1.LocalObjectReference{},
				ps:          corev1.LocalObjectReference{Name: "baz"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsPullSecret(tt.args.pullSecrets, tt.args.ps); got != tt.want {
				t.Errorf("ContainsPullSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsMap(t *testing.T) {
	superset := map[string]string{
		"hello":     "world",
		"goodnight": "moon",
		"yip":       "yip",
	}
	type args struct {
		all      map[string]string
		expected map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Superset",
			args{
				all:      superset,
				expected: map[string]string{"hello": "world", "yip": "yip"},
			},
			true,
		},
		{
			"Partial overlap",
			args{
				all:      superset,
				expected: map[string]string{"hello": "world", "greetings": "traveler"},
			},
			false,
		},
		{
			"Empty superset",
			args{
				all:      map[string]string{},
				expected: map[string]string{"yip": "yip"},
			},
			false,
		},
		{
			"Empty subset",
			args{
				all:      superset,
				expected: map[string]string{},
			},
			true,
		},
		{
			"Same keys, different values",
			args{
				all:      superset,
				expected: map[string]string{"hello": "moon", "yip": "yip"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsMap(tt.args.all, tt.args.expected); got != tt.want {
				t.Errorf("ContainsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMchIsValid(t *testing.T) {
	replicas := int(1)
	validMCH := &operatorsv1beta1.MultiClusterHub{
		TypeMeta:   metav1.TypeMeta{Kind: "MultiClusterHub"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1beta1.MultiClusterHubSpec{
			ImagePullSecret: "test",
			ReplicaCount:    &replicas,
			Mongo: operatorsv1beta1.Mongo{
				Storage:      "mongoStorage",
				StorageClass: "mongoStorageClass",
				ReplicaCount: &replicas,
			},
			Etcd: operatorsv1beta1.Etcd{
				Storage:      "etcdStorage",
				StorageClass: "etcdStorageClass",
			},
		},
		Status: operatorsv1beta1.MultiClusterHubStatus{
			CurrentVersion: "1.0.0",
		},
	}
	noMongoReplicas := validMCH.DeepCopy()
	noMongoReplicas.Spec.Mongo.ReplicaCount = nil

	type args struct {
		m *operatorsv1beta1.MultiClusterHub
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid MCH",
			args{validMCH},
			true,
		},
		{
			"Zero Mongo Replicas",
			args{noMongoReplicas},
			false,
		},
		{
			"Empty object",
			args{&operatorsv1beta1.MultiClusterHub{}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MchIsValid(tt.args.m); got != tt.want {
				t.Errorf("MchIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistributePods(t *testing.T) {
	t.Run("Returns pod affinity", func(t *testing.T) {
		if got := DistributePods("app", "testapp"); reflect.TypeOf(got) != reflect.TypeOf((*corev1.Affinity)(nil)) {
			t.Errorf("DistributePods() did not return an affinity type")
		}
	})
}
