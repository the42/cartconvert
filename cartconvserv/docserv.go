// Copyright 2011,2012 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations - documentation part

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

// These constants specifiy the directory in which the documentation files are saved
const (
	docfileroot     = "doc/"
	docmainTemplate = "index.tpl" // The main documentation file. Other filenames are created from the requested API documentation
)

// defines the layout of a documentation page and is used by html/template
type docPageLayout struct {
	ConcreteHeading  string // Heading for detailed documentation pages
	APIRoot, DocRoot string // APIRoot is used for inline examples
	Navigation       []Link // Keeps an url and the text for a href
}

// user defined APIRoot and DocRoot are constant throughout program execution
var docPage = docPageLayout{APIRoot: apiroot(), DocRoot: docroot()}

func docHandler(w http.ResponseWriter, req *http.Request) {

	// Error handler for documentation	
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "An error occurred: "+fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	base := path.Base(req.URL.Path)
	var filename string

	// check if the incoming url is the base url for documentation
	if base == path.Base(docPage.DocRoot) {
		// if the incoming url is the base url for documentation, load the generic help template
		filename = docfileroot + docmainTemplate
	} else {
		// else load the specific help template. The filename is constructed from the API function
		filename = docfileroot + base + ".tpl"
		docPage.ConcreteHeading = httphandlerfuncs[base+"/"].docstring
	}

	tpl, err := template.ParseFiles(filename)
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(w, docPage)
	if err != nil {
		panic(err)
	}
}

func init() {

	// parse all REST handlers and create corresponding documentation links
	for function, val := range httphandlerfuncs {
		url, err := url.Parse(function)
		if err != nil {
			panic(fmt.Sprintf("%s: %s is not a valid url", err.Error(), function))
		}
		docitem := Link{URL: url, Documentation: val.docstring}
		docPage.Navigation = append(docPage.Navigation, docitem)
	}

	// if documentation is compiled in, we want it included as a link on the main page
	url, _ := url.Parse(docfileroot)
	docrootLinks = append(docrootLinks, Link{URL: url, Documentation: "API Documentation"})
	http.HandleFunc("/"+docroot(), docHandler)
}
