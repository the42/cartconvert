// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/the42/cartconvert"
	"github.com/the42/cartconvert/bmn"
	"io"
	"fmt"
	"net/http"
	"path"
)

const (
	BMNHandler     = "/bmn/"
	GeoHashHandler = "/geohash/"
	LatLongHandler = "/latlong/"

	JSONFormatSpec = ".json"
	XMLFormatSpec  = ".xml"

	OutputFormatSpec = "outputformat"

	OFUTM          = "utm"
	OFgeohash      = "geohash"
	OFlatlongdeg   = "latlongdeg"
	OFlatlongcomma = "latlongcomma"
)

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

func UTMToSerial(w io.Writer, utm *cartconvert.UTMCoord, serialformat string) (err error) {
	switch serialformat {
	case JSONFormatSpec:
		err = json.NewEncoder(w).Encode(&marshalUTMCoord{UTMCoord: utm})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		err = xml.Marshal(w, &marshalUTMCoord{UTMCoord: utm})
	default:
	    err = fmt.Errorf("Unsupported serialisation format: '%s'", serialformat)
	}
	return
}

func GeoHashToSerial(w io.Writer, geohash string, serialformat string) (err error) {
	switch serialformat {
	case JSONFormatSpec:
		err = json.NewEncoder(w).Encode(&marshalGeoHash{GeoHash: geohash})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		err = xml.Marshal(w, &marshalGeoHash{GeoHash: geohash})
	default:
	  err = fmt.Errorf("Unsupported serialisation format: '%s'", serialformat)
	}
	return
}

func LatLongToSerial(w io.Writer, latlong *cartconvert.PolarCoord, serialformat string, repformat cartconvert.LatLongFormat) (err error) {
  
	lat, long := cartconvert.LatLongToString(latlong, repformat)

	switch serialformat {
	case JSONFormatSpec:
		err = json.NewEncoder(w).Encode(&marshalLatLong{Lat: lat, Long: long, Fmt: repformat.String()})
	case XMLFormatSpec:
		io.WriteString(w, xml.Header)
		err = xml.Marshal(w, &marshalLatLong{Lat: lat, Long: long, Fmt: repformat.String()})
	default:
	  	  err = fmt.Errorf("Unsupported serialisation format: '%s'", serialformat)
	}
	return
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Cartography transformation")
}

type restHandler func(http.ResponseWriter, *http.Request, string, string, string) error

func (fn restHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  
	val := path.Base(req.URL.Path)
	serialformat := path.Ext(val)
	val = val[:len(val)-len(serialformat)]
	
	oformat := req.URL.Query().Get(OutputFormatSpec)
	
	switch serialformat {
	case JSONFormatSpec:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case XMLFormatSpec:
		w.Header().Set("Content-Type", "text/xml")
	}
	
	if err := fn(w, req, val, serialformat, oformat); err != nil {
	  http.Error(w, err.Error(), 500)
	}
}

func geohashHandler(w http.ResponseWriter, req *http.Request, geohashstrval, serialformat, oformat string) (err error) {

	latlong, err := cartconvert.GeoHashToLatLong(geohashstrval, nil)
	if err != nil {
		return
	}

	switch oformat {
	case OFlatlongdeg:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
	default:
	  err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	return
}

func bmnHandler(w http.ResponseWriter, req *http.Request, bmnstrval, serialformat, oformat string) (err error) {

	bmnval, err := bmn.ABMNToStruct(bmnstrval)
	if err != nil {
		return
	}

	latlong, err := bmn.BMNToWGS84LatLong(bmnval)
	if err != nil {
		return
	}

	switch oformat {
	case OFUTM:
		err = UTMToSerial(w, cartconvert.LatLongToUTM(latlong), serialformat)
	case OFgeohash:
		err = GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong), serialformat)
	case OFlatlongdeg:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
	default:
	  err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	
	return
}

func latlongHandler(w http.ResponseWriter, req *http.Request, latlongstrval, serialformat, oformat string) (err error) {

  if len(latlongstrval) > 0 {
    return fmt.Errorf("Latlong doesn't support an input value. Use parameters instead")
  }
  
  var slat, slong string
  slat = req.URL.Query().Get("lat")
  slong = req.URL.Query().Get("long")
  _,_ = slat, slong
  
  bmnval, err := bmn.ABMNToStruct(latlongstrval)
	if err != nil {
		return
	}

	latlong, err := bmn.BMNToWGS84LatLong(bmnval)
	if err != nil {
		return
	}

	switch oformat {
	case OFUTM:
		err = UTMToSerial(w, cartconvert.LatLongToUTM(latlong), serialformat)
	case OFgeohash:
		err = GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong), serialformat)
	case OFlatlongdeg:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
	default:
	  err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	
	return
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.Handle(BMNHandler, restHandler(bmnHandler))
	http.Handle(GeoHashHandler, restHandler(geohashHandler))
	http.Handle(LatLongHandler, restHandler(latlongHandler))
	// TODO: Read from config file
	http.ListenAndServe(":1111", nil)
}
