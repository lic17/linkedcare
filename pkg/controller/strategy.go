package controller

import (
	"fmt"

	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"
)

type StrategyHandler struct {
}

func newStrategyHandler() *StrategyHandler {
	return &StrategyHandler{}
}

func (*StrategyHandler) Add(obj interface{}) {
	fmt.Printf("Add canary:\n")
	if ca, ok := obj.(servicemeshv1alpha1.Strategy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(ca)
	}
}
func (*StrategyHandler) Update(oldObj, newObj interface{}) {
	fmt.Printf("Update canary:\n")
	if caOld, ok := oldObj.(servicemeshv1alpha1.Strategy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(caOld)
	}

	if caNew, ok := newObj.(servicemeshv1alpha1.Strategy); !ok {
		fmt.Printf("not canaries\n")
	} else {
		fmt.Println(caNew)
	}
}
func (*StrategyHandler) Delete(obj interface{}) {
	fmt.Printf("Delete canary:\n")
	if ca, ok := obj.(servicemeshv1alpha1.Strategy); !ok {
		fmt.Printf("not canary")
	} else {
		fmt.Println(ca)
	}
}
