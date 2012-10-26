// Copyright 2011,2012 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// +build !appengine

// RESTFul interface for coordinate transformations - configuration for stand alone server
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var configFileName = flag.String("config", "config.json", "location of JSON configuration file")

type config struct {
	APIRoot string
	DocRoot string
	Binding string
}

var conf = &config{APIRoot: "/api/", DocRoot: "/doc/", Binding: ":1111"}

func readConfig(filename string, conf *config) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if conf == nil {
		conf = &config{}
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic("Unable to parse json configuration file")
	}
	return
}

func apiroot() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	path := path.Base(conf.APIRoot)
	if path == "/" || path == "." {
		panic("configuration error: APIRoot must be a non-empty path")
	}
	return conf.APIRoot
}

func docroot() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	path := path.Base(conf.DocRoot)
	if path == "/" || path == "." {
		panic("configuration error: DocRoot must be a non-empty path")
	}
	return conf.DocRoot
}

func binding() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	return conf.Binding
}
