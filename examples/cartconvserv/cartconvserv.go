// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	"json"
	"io"
	"xml"
	"http"
	"fmt"
	"path"
	"github.com/the42/cartconvert/bmn"
	"github.com/the42/cartconvert"
)

const (
	BMNHandler     = "/bmn/"
	GeoHashHandler = "/geohash/"

	JSONFormatSpec = ".json"
	XMLFormatSpec  = ".xml"

	OutputFormatSpec = "outputformat"

	OFUTM          = "utm"
	OFgeohash      = "geohash"
	OFlatlongdeg   = "latlongdeg"
	OFlatlongcomma = "latlongcomma"
)

// const httpc = "<html><head></head><body>%s</body></html>"

type marshalUTMCoord struct {
	XMLName  xml.Name `json:"-" xml:"payLoad"`
	UTMCoord *cartconvert.UTMCoord
}

type marshalGeoHash struct {
	XMLName xml.Name `json:"-" xml:"payLoad"`
	GeoHash string
}

type marshalLatLong struct {
	XMLName        xml.Name `json:"-" xml:"payLoad"`
	Lat, Long, Fmt string
}

func UTMToSerial(w io.Writer, utm *cartconvert.UTMCoord, serialformat string) {
	// Maybe UTMCoord shoul implement an interface for serialisation
	switch serialformat {
	case JSONFormatSpec:
		json.NewEncoder(w).Encode(&marshalUTMCoord{UTMCoord: utm})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		xml.Marshal(w, &marshalUTMCoord{UTMCoord: utm})
	}
}

func GeoHashToSerial(w io.Writer, geohash string, serialformat string) {
	switch serialformat {
	case JSONFormatSpec:
		json.NewEncoder(w).Encode(&marshalGeoHash{GeoHash: geohash})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		xml.Marshal(w, &marshalGeoHash{GeoHash: geohash})
	}
}

func LatLongToSerial(w io.Writer, latlong *cartconvert.PolarCoord, serialformat string, repformat cartconvert.LatLongFormat) {
	lat, long := cartconvert.LatLongToString(latlong, repformat)
	switch serialformat {

	case JSONFormatSpec:
		json.NewEncoder(w).Encode(&marshalLatLong{Lat: lat, Long: long, Fmt: string(repformat)})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		xml.Marshal(w, &marshalLatLong{Lat: lat, Long: long, Fmt: string(repformat)})
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Cartography transformation")
}

func geohashHandler(w http.ResponseWriter, req *http.Request) {
	// OSGB36 Datum transformation
	// gc := cartconvert.DirectTransverseMercator(&cartconvert.PolarCoord{Latitude: flat, Longitude: flong, El: cartconvert.Airy1830Ellipsoid}, 49, -2, 0.9996012717, 400000, -100000)

	// ONLY placeholders, REWRITE!
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

	switch serialformat {
	case JSONFormatSpec:
		w.Header().Set("Content-Type", "application/json")
	case XMLFormatSpec:
		w.Header().Set("Content-Type", "text/xml")
	}

	switch oformat {
	case OFUTM:
		UTMToSerial(w, cartconvert.LatLongToUTM(latlong), serialformat)
	case OFgeohash:
		GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong), serialformat)
	case OFlatlongdeg:
		LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
	}
}

func bmnHandler(w http.ResponseWriter, req *http.Request) {

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

	switch serialformat {
	case JSONFormatSpec:
		w.Header().Set("Content-Type", "application/json")
	case XMLFormatSpec:
		w.Header().Set("Content-Type", "text/xml")
	}

	switch oformat {
	case OFUTM:
		UTMToSerial(w, cartconvert.LatLongToUTM(latlong), serialformat)
	case OFgeohash:
		GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong), serialformat)
	case OFlatlongdeg:
		LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
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
	http.HandleFunc(BMNHandler, bmnHandler)
	http.HandleFunc(GeoHashHandler, geohashHandler)
	// TODO: Read from config file
	http.ListenAndServe(":1111", nil)
}
