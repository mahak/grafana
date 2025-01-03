// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v0alpha1 "github.com/grafana/grafana/pkg/apis/service/v0alpha1"
	internal "github.com/grafana/grafana/pkg/generated/applyconfiguration/internal"
	servicev0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/service/v0alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=service.grafana.app, Version=v0alpha1
	case v0alpha1.SchemeGroupVersion.WithKind("ExternalName"):
		return &servicev0alpha1.ExternalNameApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("ExternalNameSpec"):
		return &servicev0alpha1.ExternalNameSpecApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
