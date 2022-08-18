// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "sync/atomic"

type Atomic[T any] struct {
	value atomic.Value
}

func (a *Atomic[T]) Load() T {
	return a.value.Load().(T)
}

func (a *Atomic[T]) Store(val T) {
	a.value.Store(val)
}
