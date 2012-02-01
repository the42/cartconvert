// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/the42/cartconvert"
	"github.com/the42/cartconvert/bmn"
	"io"
	"net/http"
	"path"
)

const (
	BMNHandler     = "bmn/"
	GeoHashHandler = "geohash/"
	LatLongHandler = "latlong/"
	UTMHandler     = "utm/"

	JSONFormatSpec = ".json"
	XMLFormatSpec  = ".xml"

	OutputFormatSpec = "outputformat"

	OFUTM          = "utm"
	OFgeohash      = "geohash"
	OFlatlongdeg   = "latlongdeg"
	OFlatlongcomma = "latlongcomma"
	OFBMN          = "bmn"
)

type Encoder interface {
	Encode(v interface{}) error
}

type UTMCoord struct {
	UTMCoord  *cartconvert.UTMCoord
	UTMString string
}

type BMN struct {
	BMNCoord  *bmn.BMNCoord
	BMNString string
}

type LatLong struct {
	Lat, Long, Fmt string
	LatLongString  string
}

func UTMToSerial(w Encoder, utm *cartconvert.UTMCoord) error {
	return w.Encode(&UTMCoord{UTMCoord: utm, UTMString: utm.String()})
}

func GeoHashToSerial(w Encoder, geohash string) error {
	return w.Encode(geohash)
}

func LatLongToSerial(w Encoder, latlong *cartconvert.PolarCoord, repformat cartconvert.LatLongFormat) (err error) {

	lat, long := cartconvert.LatLongToString(latlong, repformat)
	return w.Encode(&LatLong{Lat: lat, Long: long, Fmt: repformat.String(), LatLongString: latlong.String()})
}

func BMNToSerial(w Encoder, bmn *bmn.BMNCoord) error {
	return w.Encode(&BMN{BMNCoord: bmn, BMNString: bmn.String()})
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Cartography transformation")
}

type restHandler func(Encoder, *http.Request, string, string) error

func (fn restHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var enc Encoder
	val := path.Base(req.URL.Path)
	serialformat := path.Ext(val)
	val = val[:len(val)-len(serialformat)]

	oformat := req.URL.Query().Get(OutputFormatSpec)

	switch serialformat {
	case JSONFormatSpec, "":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc = json.NewEncoder(w)
	case XMLFormatSpec:
		w.Header().Set("Content-Type", "text/xml")
		io.WriteString(w, xml.Header)
		enc = xml.NewEncoder(w)
	default:
		http.Error(w, fmt.Sprintf("Unsupported serialisation format: '%s'", serialformat), 500)
		return
	}

	// Recover from panic by setting http error 500 and letting the user know the reason
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), 500)
		}
	}()

	if err := fn(enc, req, val, oformat); err != nil {
		// might as well panic(err) but maybe treat different
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func serialiser(w Encoder, latlong *cartconvert.PolarCoord, oformat string) (err error) {
	switch oformat {
	case OFlatlongdeg:
		err = LatLongToSerial(w, latlong, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = LatLongToSerial(w, latlong, cartconvert.LLFdeg)
	case OFUTM:
		err = UTMToSerial(w, cartconvert.LatLongToUTM(latlong))
	case OFgeohash:
		err = GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong))
	case OFBMN:
		var bmnval *bmn.BMNCoord
		bmnval, err = bmn.WGS84LatLongToBMN(latlong, bmn.BMNZoneDet)
		if err == nil {
			err = BMNToSerial(w, bmnval)
		}
	default:
		err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	return
}

func geohashHandler(w Encoder, req *http.Request, geohashstrval, oformat string) (err error) {

	var latlong *cartconvert.PolarCoord

	if latlong, err = cartconvert.GeoHashToLatLong(geohashstrval, nil); err != nil {
		return
	}
	return serialiser(w, latlong, oformat)
}

func bmnHandler(w Encoder, req *http.Request, bmnstrval, oformat string) (err error) {
	var bmnval *bmn.BMNCoord
	if bmnval, err = bmn.ABMNToStruct(bmnstrval); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = bmn.BMNToWGS84LatLong(bmnval); err != nil {
		return
	}
	return serialiser(w, latlong, oformat)
}

func latlongHandler(w Encoder, req *http.Request, latlongstrval, oformat string) (err error) {

	if len(latlongstrval) > 0 {
		return fmt.Errorf("Latlong doesn't accept an input value. Use parameters instead")
	}

	slat := req.URL.Query().Get("lat")
	slong := req.URL.Query().Get("long")

	var lat, long float64
	lat, err = cartconvert.ADegMMSSToNum(slat)
	if err != nil {
		lat, err = cartconvert.ADegCommaToNum(slat)
		if err != nil {
			return fmt.Errorf("Not a bearing: '%s'", slat)
		}
	}

	long, err = cartconvert.ADegMMSSToNum(slong)
	if err != nil {
		long, err = cartconvert.ADegCommaToNum(slong)
		if err != nil {
			return fmt.Errorf("Not a bearing: '%s'", slong)
		}
	}

	latlong := &cartconvert.PolarCoord{Latitude: lat, Longitude: long, El: cartconvert.DefaultEllipsoid}

	return serialiser(w, latlong, oformat)
}

func utmHandler(w Encoder, req *http.Request, utmstrval, oformat string) (err error) {
	var utmval *cartconvert.UTMCoord
	if utmval, err = cartconvert.AUTMToStruct(utmstrval, nil); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = cartconvert.UTMToLatLong(utmval); err != nil {
		return
	}
	return serialiser(w, latlong, oformat)
}

type config struct {
	APIRoot string
	Binding string
}

var conf *config

func init() {

	apiroot := apiroot()

	http.HandleFunc("/", rootHandler)
	http.Handle(apiroot+BMNHandler, restHandler(bmnHandler))
	http.Handle(apiroot+GeoHashHandler, restHandler(geohashHandler))
	http.Handle(apiroot+LatLongHandler, restHandler(latlongHandler))
	http.Handle(apiroot+UTMHandler, restHandler(utmHandler))
}
