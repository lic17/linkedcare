package controller

import (
	"fmt"

	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"
)

type ServicePoliciesHandler struct {
}

func newServicePoliciesHandler() *ServicePoliciesHandler {
	return &ServicePoliciesHandler{}
}

func (*ServicePoliciesHandler) Add(obj interface{}) {
	fmt.Printf("Add canary:\n")
	if ca, ok := obj.(servicemeshv1alpha1.ServicePolicy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(ca)
	}
}
func (*ServicePoliciesHandler) Update(oldObj, newObj interface{}) {
	fmt.Printf("Update canary:\n")
	if caOld, ok := oldObj.(servicemeshv1alpha1.ServicePolicy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(caOld)
	}

	if caNew, ok := newObj.(servicemeshv1alpha1.ServicePolicy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(caNew)
	}
}
func (*ServicePoliciesHandler) Delete(obj interface{}) {
	fmt.Printf("Delete canary:\n")
	if ca, ok := obj.(servicemeshv1alpha1.ServicePolicy); !ok {
		fmt.Printf("not canary")
	} else {
		fmt.Println(ca)
	}
}
