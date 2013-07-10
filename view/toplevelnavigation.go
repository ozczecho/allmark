// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

type ToplevelNavigation struct {
	Entries []*ToplevelEntry
}

func (navigation *ToplevelNavigation) IsAvailable() bool {
	return len(navigation.Entries) > 0
}

type ToplevelEntry struct {
	Title string
	Path  string
}