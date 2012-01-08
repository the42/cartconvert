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
	BMNHandler     = "/bmn/"
	GeoHashHandler = "/geohash/"
	LatLongHandler = "/latlong/"
	UTMHandler     = "/utm/"

	JSONFormatSpec = ".json"
	XMLFormatSpec  = ".xml"

	OutputFormatSpec = "outputformat"

	OFUTM          = "utm"
	OFgeohash      = "geohash"
	OFlatlongdeg   = "latlongdeg"
	OFlatlongcomma = "latlongcomma"
	OFBMN          = "bmn"
)

type marshalUTMCoord struct {
	XMLName   xml.Name `json:"-" xml:"payLoad"`
	UTMCoord  *cartconvert.UTMCoord
	UTMString string
}

type marshalBMN struct {
	XMLName   xml.Name `json:"-" xml:"payLoad"`
	BMNCoord  *bmn.BMNCoord
	BMNString string
}

type marshalGeoHash struct {
	XMLName xml.Name `json:"-" xml:"payLoad"`
	GeoHash string
}

type marshalLatLong struct {
	XMLName        xml.Name `json:"-" xml:"payLoad"`
	Lat, Long, Fmt string
	LatLongString  string
}

func UTMToSerial(w io.Writer, utm *cartconvert.UTMCoord, serialformat string) (err error) {
	switch serialformat {
	case JSONFormatSpec:
		err = json.NewEncoder(w).Encode(&marshalUTMCoord{UTMCoord: utm, UTMString: utm.String()})
	case XMLFormatSpec:
		err = xml.Marshal(w, &marshalUTMCoord{UTMCoord: utm, UTMString: utm.String()})
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
		err = json.NewEncoder(w).Encode(&marshalLatLong{Lat: lat, Long: long, Fmt: repformat.String(), LatLongString: latlong.String()})
	case XMLFormatSpec:
		err = xml.Marshal(w, &marshalLatLong{Lat: lat, Long: long, Fmt: repformat.String(), LatLongString: latlong.String()})
	default:
		err = fmt.Errorf("Unsupported serialisation format: '%s'", serialformat)
	}
	return
}

func BMNToSerial(w io.Writer, bmn *bmn.BMNCoord, serialformat string) (err error) {

	switch serialformat {
	case JSONFormatSpec:
		err = json.NewEncoder(w).Encode(&marshalBMN{BMNCoord: bmn, BMNString: bmn.String()})
	case XMLFormatSpec:
		err = xml.Marshal(w, &marshalBMN{BMNCoord: bmn, BMNString: bmn.String()})
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
		io.WriteString(w, xml.Header)
	}

	// Recover from panic by setting http error 500 and letting the user know the reason
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), 500)
		}
	}()

	if err := fn(w, req, val, serialformat, oformat); err != nil {
		// might as well panic(err) but maybe treat different
		http.Error(w, fmt.Sprint(err), 500)
	}
}

func serialiser(w http.ResponseWriter, latlong *cartconvert.PolarCoord, oformat, serialformat string) (err error) {
	switch oformat {
	case OFlatlongdeg:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = LatLongToSerial(w, latlong, serialformat, cartconvert.LLFdeg)
	case OFUTM:
		err = UTMToSerial(w, cartconvert.LatLongToUTM(latlong), serialformat)
	case OFgeohash:
		err = GeoHashToSerial(w, cartconvert.LatLongToGeoHash(latlong), serialformat)
	case OFBMN:
		var bmnval *bmn.BMNCoord
		bmnval, err = bmn.WGS84LatLongToBMN(latlong, bmn.BMNZoneDet)
		if err == nil {
			err = BMNToSerial(w, bmnval, serialformat)
		}
	default:
		err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	return
}

func geohashHandler(w http.ResponseWriter, req *http.Request, geohashstrval, serialformat, oformat string) (err error) {

	var latlong *cartconvert.PolarCoord

	if latlong, err = cartconvert.GeoHashToLatLong(geohashstrval, nil); err != nil {
		return
	}
	return serialiser(w, latlong, oformat, serialformat)
}

func bmnHandler(w http.ResponseWriter, req *http.Request, bmnstrval, serialformat, oformat string) (err error) {
	var bmnval *bmn.BMNCoord
	if bmnval, err = bmn.ABMNToStruct(bmnstrval); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = bmn.BMNToWGS84LatLong(bmnval); err != nil {
		return
	}
	return serialiser(w, latlong, oformat, serialformat)
}

func latlongHandler(w http.ResponseWriter, req *http.Request, latlongstrval, serialformat, oformat string) (err error) {

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

	return serialiser(w, latlong, oformat, serialformat)
}

func utmHandler(w http.ResponseWriter, req *http.Request, utmstrval, serialformat, oformat string) (err error) {
	var utmval *cartconvert.UTMCoord
	if utmval, err = cartconvert.AUTMToStruct(utmstrval, nil); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = cartconvert.UTMToLatLong(utmval); err != nil {
		return
	}
	return serialiser(w, latlong, oformat, serialformat)
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.Handle(BMNHandler, restHandler(bmnHandler))
	http.Handle(GeoHashHandler, restHandler(geohashHandler))
	http.Handle(LatLongHandler, restHandler(latlongHandler))
	http.Handle(UTMHandler, restHandler(utmHandler))

	// TODO: Read from config file
	http.ListenAndServe(":1111", nil)
}
