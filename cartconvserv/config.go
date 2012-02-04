// Copyright 2011,2012 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.

// +build !appengine
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readConfig(conf *config) {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	if conf == nil {
		conf = &config{}
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		panic("Unable to parse json configuration file")
	}
	return
}

func apiroot() string {
	if conf == nil {
		conf = &config{}
	}
	conf.APIRoot = "/api/"
	readConfig(conf)
	return conf.APIRoot
}

func binding() string {
	if conf == nil {
		conf = &config{}
	}
	conf.Binding = ":1111"
	readConfig(conf)
	return conf.Binding
}
