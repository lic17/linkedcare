package k8sutil

import (
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"linkedcare.io/linkedcare/pkg/models"
)

func IsControlledBy(reference []metav1.OwnerReference, kind string, name string) bool {
	for _, ref := range reference {
		if ref.Kind == kind && (name == "" || ref.Name == name) {
			return true
		}
	}
	return false
}

func GetControlledWorkspace(reference []metav1.OwnerReference) string {
	for _, ref := range reference {
		if ref.Kind == "Workspace" {
			return ref.Name
		}
	}
	return ""
}

func ContainsUser(subjects interface{}, username string) bool {
	switch subjects.(type) {
	case []*v1.Subject:
		for _, subject := range subjects.([]*v1.Subject) {
			if subject.Kind == v1.UserKind && subject.Name == username {
				return true
			}
		}
	case []v1.Subject:
		for _, subject := range subjects.([]v1.Subject) {
			if subject.Kind == v1.UserKind && subject.Name == username {
				return true
			}
		}
	case []models.User:
		for _, u := range subjects.([]models.User) {
			if u.Username == username {
				return true
			}
		}

	case []*models.User:
		for _, u := range subjects.([]*models.User) {
			if u.Username == username {
				return true
			}
		}
	}
	return false
}
