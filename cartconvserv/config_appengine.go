// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.

// +build appengine
package main

func apiroot() (root string) {
	root = "/api/"
	return
}

func binding() (bind string) {
	bind = ":1111"
	return
}
