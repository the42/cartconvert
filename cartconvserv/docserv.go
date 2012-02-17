// Copyright 2011,2012 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations - documentation part

package main

import (
	"fmt"
	"html/template"
	"net/http"
	//"path"
)

var docmainTemplate = docroot() + "index.tpl"

type PageLayout struct {
  Navigation []string
}

const errorParsingTemplate = `
<html>
  <head>
  </head>
  <body>
    An error occured: %s
  </body>
</html>`

func docHandler(w http.ResponseWriter, req *http.Request) {
	/*
	 * Idee: Zuerst das allgemeine template laden, falls parameter angegeben wurden, das spezielle template nachladen
	 * 
	 * 
	 */
	
	val := req.URL.Path
	tpl, err := template.ParseFiles(docmainTemplate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, errorParsingTemplate, err)
	} else {
		tpl.Execute(w, nil)
	}
}

func init() {
	http.HandleFunc("/"+docroot(), docHandler)
}
