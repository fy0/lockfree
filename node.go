// Copyright 2020 dustinxie
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

package lockfree

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type (
	hashNode struct {
		hash uint64
		key  unsafe.Pointer
		val  unsafe.Pointer
		nxt  unsafe.Pointer
	}
)

func newFence() *hashNode {
	return &hashNode{hash: math.MaxUint64}
}

func isFence(n *hashNode) bool {
	return n.key == nil
}

func (n *hashNode) linkTo(next *hashNode) {
	n.nxt = unsafe.Pointer(next)
}

func (n *hashNode) next() *hashNode {
	return (*hashNode)(atomic.LoadPointer(&n.nxt))
}

func (n *hashNode) value() unsafe.Pointer {
	return atomic.LoadPointer(&n.val)
}

func (n *hashNode) casNext(expected, target unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n.nxt, expected, target)
}

func (n *hashNode) casValue(expected, target unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n.val, expected, target)
}