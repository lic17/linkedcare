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
	"linkedcare.io/linkedcare/pkg/utils/sliceutil"
	"sort"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type configMapSearcher struct {
}

func (*configMapSearcher) get(namespace, name string) (interface{}, error) {
	return informers.SharedInformerFactory().Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).Get(name)
}

// exactly Match
func (*configMapSearcher) match(match map[string]string, item *v1.ConfigMap) bool {
	for k, v := range match {
		switch k {
		case Name:
			names := strings.Split(v, "|")
			if !sliceutil.HasString(names, item.Name) {
				return false
			}
		case Keyword:
			if !strings.Contains(item.Name, v) && !searchFuzzy(item.Labels, "", v) && !searchFuzzy(item.Annotations, "", v) {
				return false
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
func (*configMapSearcher) fuzzy(fuzzy map[string]string, item *v1.ConfigMap) bool {
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
		case app:
			if !strings.Contains(item.Labels[chart], v) && !strings.Contains(item.Labels[release], v) {
				return false
			}
		default:
			if !searchFuzzy(item.Labels, k, v) {
				return false
			}
		}
	}
	return true
}

func (*configMapSearcher) compare(a, b *v1.ConfigMap, orderBy string) bool {
	switch orderBy {
	case CreateTime:
		return a.CreationTimestamp.Time.Before(b.CreationTimestamp.Time)
	case Name:
		fallthrough
	default:
		return strings.Compare(a.Name, b.Name) <= 0
	}
}

func (s *configMapSearcher) search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	configMaps, err := informers.SharedInformerFactory().Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*v1.ConfigMap, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = configMaps
	} else {
		for _, item := range configMaps {
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
