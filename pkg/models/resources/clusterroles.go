/*

 Copyright 2019 The Linkedcare Authors.

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
package resources

import (
	"linkedcare.io/linkedcare/pkg/constants"
	"linkedcare.io/linkedcare/pkg/informers"
	"linkedcare.io/linkedcare/pkg/server/params"
	"linkedcare.io/linkedcare/pkg/utils/k8sutil"
	"linkedcare.io/linkedcare/pkg/utils/sliceutil"
	"sort"
	"strings"

	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type clusterRoleSearcher struct {
}

func (*clusterRoleSearcher) get(namespace, name string) (interface{}, error) {
	return informers.SharedInformerFactory().Rbac().V1().ClusterRoles().Lister().Get(name)
}

// exactly Match
func (*clusterRoleSearcher) match(match map[string]string, item *rbac.ClusterRole) bool {
	for k, v := range match {
		switch k {
		case OwnerKind:
			fallthrough
		case OwnerName:
			kind := match[OwnerKind]
			name := match[OwnerName]
			if !k8sutil.IsControlledBy(item.OwnerReferences, kind, name) {
				return false
			}
		case Name:
			names := strings.Split(v, "|")
			if !sliceutil.HasString(names, item.Name) {
				return false
			}
		case Keyword:
			if !strings.Contains(item.Name, v) && !searchFuzzy(item.Labels, "", v) && !searchFuzzy(item.Annotations, "", v) {
				return false
			}
		case UserFacing:
			if v == "true" {
				if !isUserFacingClusterRole(item) {
					return false
				}
			}
		default:
			// label not exist or value not equal
			if val, ok := item.Labels[k]; !ok || val != v {
				return false
			}
		}
	}
	return true
}

// Fuzzy searchInNamespace
func (*clusterRoleSearcher) fuzzy(fuzzy map[string]string, item *rbac.ClusterRole) bool {
	for k, v := range fuzzy {
		switch k {
		case Name:
			if !strings.Contains(item.Name, v) && !strings.Contains(item.Annotations[constants.DisplayNameAnnotationKey], v) {
				return false
			}
		case Label:
			if !searchFuzzy(item.Labels, "", v) {
				return false
			}
		case annotation:
			if !searchFuzzy(item.Annotations, "", v) {
				return false
			}
			return false
		default:
			if !searchFuzzy(item.Labels, k, v) {
				return false
			}
		}
	}
	return true
}

func (*clusterRoleSearcher) compare(a, b *rbac.ClusterRole, orderBy string) bool {
	switch orderBy {
	case CreateTime:
		return a.CreationTimestamp.Time.Before(b.CreationTimestamp.Time)
	case Name:
		fallthrough
	default:
		return strings.Compare(a.Name, b.Name) <= 0
	}
}

func (s *clusterRoleSearcher) search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	clusterRoles, err := informers.SharedInformerFactory().Rbac().V1().ClusterRoles().Lister().List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*rbac.ClusterRole, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = clusterRoles
	} else {
		for _, item := range clusterRoles {
			if s.match(conditions.Match, item) && s.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			tmp := i
			i = j
			j = tmp
		}
		return s.compare(result[i], result[j], orderBy)
	})

	r := make([]interface{}, 0)
	for _, i := range result {
		r = append(r, i)
	}
	return r, nil
}

// cluster role created by user from linkedcare dashboard
func isUserFacingClusterRole(role *rbac.ClusterRole) bool {
	if role.Annotations[constants.CreatorAnnotationKey] != "" && role.Labels[constants.WorkspaceLabelKey] == "" {
		return true
	}
	return false
}
