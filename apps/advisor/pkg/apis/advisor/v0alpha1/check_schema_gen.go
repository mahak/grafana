//
// Code generated by grafana-app-sdk. DO NOT EDIT.
//

package v0alpha1

import (
	"github.com/grafana/grafana-app-sdk/resource"
)

// schema is unexported to prevent accidental overwrites
var (
	schemaCheck = resource.NewSimpleSchema("advisor.grafana.app", "v0alpha1", &Check{}, &CheckList{}, resource.WithKind("Check"),
		resource.WithPlural("checks"), resource.WithScope(resource.NamespacedScope))
	kindCheck = resource.Kind{
		Schema: schemaCheck,
		Codecs: map[resource.KindEncoding]resource.Codec{
			resource.KindEncodingJSON: &CheckJSONCodec{},
		},
	}
)

// Kind returns a resource.Kind for this Schema with a JSON codec
func CheckKind() resource.Kind {
	return kindCheck
}

// Schema returns a resource.SimpleSchema representation of Check
func CheckSchema() *resource.SimpleSchema {
	return schemaCheck
}

// Interface compliance checks
var _ resource.Schema = kindCheck
