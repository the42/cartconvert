// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	//"github.com/the42/cartconvert"
	//"json"
	"http"
	"io"
)

const (
	BMNHandler = "/bmn/"
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Cartography transformation")
}

func bundesmeldenetzHandler(w http.ResponseWriter, req *http.Request) {
	// OSGB36 Datum transformation
	// gc := cartconvert.DirectTransverseMercator(&cartconvert.PolarCoord{Latitude: flat, Longitude: flong, El: cartconvert.Airy1830Ellipsoid}, 49, -2, 0.9996012717, 400000, -100000)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, req.URL.Path)
}

/*
func (req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html")
	x, err := json.Marshal(cartconvert.LatLongToGeoHash(&cartconvert.PolarCoord{Latitude: 49.3, Longitude: 20.0}))
	if err == nil {
		io.WriteString(w, string(x))
	} else {
		io.WriteString(w, "Error: "+err.String())
	}

}
*/

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc(BMNHandler, bundesmeldenetzHandler)
	// TODO: Read from config file
	http.ListenAndServe(":1111", nil)
}
