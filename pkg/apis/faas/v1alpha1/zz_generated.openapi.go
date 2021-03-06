// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunction":            schema_pkg_apis_faas_v1alpha1_JSFunction(ref),
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuild":       schema_pkg_apis_faas_v1alpha1_JSFunctionBuild(ref),
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildSpec":   schema_pkg_apis_faas_v1alpha1_JSFunctionBuildSpec(ref),
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildStatus": schema_pkg_apis_faas_v1alpha1_JSFunctionBuildStatus(ref),
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionSpec":        schema_pkg_apis_faas_v1alpha1_JSFunctionSpec(ref),
		"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionStatus":      schema_pkg_apis_faas_v1alpha1_JSFunctionStatus(ref),
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunction(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunction is the Schema for the jsfunctions API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionSpec", "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunctionBuild(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunctionBuild is the Schema for the jsfunctionbuilds API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildSpec", "github.com/openshift-cloud-functions/js-function-operator/pkg/apis/faas/v1alpha1.JSFunctionBuildStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunctionBuildSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunctionBuildSpec defines the desired state of JSFunctionBuild",
				Properties: map[string]spec.Schema{
					"revision": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"image": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"revision", "image"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunctionBuildStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunctionBuildStatus defines the observed state of JSFunctionBuild",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunctionSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunctionSpec defines the desired state of JSFunction",
				Properties: map[string]spec.Schema{
					"func": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"package": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"events": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
				},
				Required: []string{"func"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_faas_v1alpha1_JSFunctionStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JSFunctionStatus defines the observed state of JSFunction",
				Properties: map[string]spec.Schema{
					"nodes": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
				},
				Required: []string{"nodes"},
			},
		},
		Dependencies: []string{},
	}
}
