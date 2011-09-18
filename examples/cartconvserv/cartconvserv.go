// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	"json"
	"http"
	"fmt"
	"path"
	"github.com/the42/cartconvert/bmn"
	"github.com/the42/cartconvert"
)

const (
	BMNHandler       = "/bmn/"
	JSONFormatSpec   = ".json"
	XMLFormatSpec    = ".xml"
	OutputFormatSpec = "outputformat"
	OFUTM            = "utm"
)

const httpc = "<html><head></head><body>%s</body></html>"

func UTMToJson(utm *cartconvert.UTMCoord) []byte {
	json, _ := json.Marshal(*utm)
	return json
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Cartography transformation")
}

func bundesmeldenetzHandler(w http.ResponseWriter, req *http.Request) {
	// OSGB36 Datum transformation
	// gc := cartconvert.DirectTransverseMercator(&cartconvert.PolarCoord{Latitude: flat, Longitude: flong, El: cartconvert.Airy1830Ellipsoid}, 49, -2, 0.9996012717, 400000, -100000)

	bmnstrval := req.URL.Path[len(BMNHandler):]
	serialformat := path.Ext(bmnstrval)
	bmnstrval = bmnstrval[:len(bmnstrval)-len(serialformat)]
	bmnval, err := bmn.ABMNToStruct(bmnstrval)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	latlong, err := bmn.BMNToWGS84LatLong(bmnval)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	oformat := req.URL.Query().Get(OutputFormatSpec)

	switch oformat {
	case OFUTM:
		utm := cartconvert.LatLongToUTM(latlong)
		if serialformat == JSONFormatSpec {
			json := UTMToJson(utm)
			fmt.Fprintf(w, httpc, json)
		}
	}
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
