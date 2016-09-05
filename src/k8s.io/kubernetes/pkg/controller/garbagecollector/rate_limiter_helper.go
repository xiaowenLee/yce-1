/*
Copyright 2016 The Kubernetes Authors.

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

package garbagecollector

import (
	"fmt"
	"strings"
	"sync"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/typed/dynamic"
	"k8s.io/kubernetes/pkg/util/metrics"
)

// RegisteredRateLimiter records the registered RateLimters to avoid
// duplication.
type RegisteredRateLimiter struct {
	rateLimiters map[unversioned.GroupVersion]struct{}
	lock         sync.RWMutex
}

// NewRegisteredRateLimiter returns a new RegisteredRateLimiater.
func NewRegisteredRateLimiter() *RegisteredRateLimiter {
	return &RegisteredRateLimiter{
		rateLimiters: make(map[unversioned.GroupVersion]struct{}),
	}
}

func (r *RegisteredRateLimiter) registerIfNotPresent(gv unversioned.GroupVersion, client *dynamic.Client, prefix string) {
	r.lock.RLock()
	_, ok := r.rateLimiters[gv]
	r.lock.RUnlock()
	if ok {
		return
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.rateLimiters[gv]; !ok {
		if rateLimiter := client.GetRateLimiter(); rateLimiter != nil {
			group := strings.Replace(gv.Group, ".", ":", -1)
			metrics.RegisterMetricAndTrackRateLimiterUsage(fmt.Sprintf("%s_%s_%s", prefix, group, gv.Version), rateLimiter)
		}
		r.rateLimiters[gv] = struct{}{}
	}
}
