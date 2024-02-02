// Copyright (c) 2024  The Go-Enjin Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mime

import (
	"sync"
)

type lookup struct {
	m map[string]string
	sync.RWMutex
}

func (l *lookup) unset(k string) {
	l.Lock()
	defer l.Unlock()
	delete(l.m, k)
}

func (l *lookup) set(k, v string) {
	l.Lock()
	defer l.Unlock()
	l.m[k] = v
}

func (l *lookup) get(k string) (v string, ok bool) {
	l.RLock()
	defer l.RUnlock()
	v, ok = l.m[k]
	return
}
