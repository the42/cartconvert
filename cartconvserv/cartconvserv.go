// Copyright 2011,2012 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// RESTFul interface for coordinate transformations.
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/the42/cartconvert"
	"github.com/the42/cartconvert/bmn"
	"github.com/the42/cartconvert/osgb36"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// supported serialization formats
const (
	JSONFormatSpec = ".json"
	XMLFormatSpec  = ".xml"
)

// supported representation/transformation formats
const (
	OutputFormatSpec = "outputformat"

	OFlatlongdeg   = "latlongdeg"
	OFlatlongcomma = "latlongcomma"
	OFgeohash      = "geohash"
	OFUTM          = "utm"
	OFBMN          = "bmn"
	OFOSGB         = "osgb"
)

// Interface type for transparent XML / JSON Encoding
type Encoder interface {
	Encode(v interface{}) error
}

// --------------------------------------------------------------------
// Serialization struct definitions
type (
	// Errors will be put back to the caller in the requested encoding, encapsulated in an error struct
	Error struct {
		Error string
	}

	LatLong struct {
		Lat, Long, Fmt string
		LatLongString  string
	}

	GeoHash struct {
		GeoHash string
	}

	UTMCoord struct {
		UTMCoord  *cartconvert.UTMCoord // MIND: UTMCoord is named, because XML and JSON serialization behave differently. An unnamed struct element will NOT be serialized by the XML encoder
		UTMString string
	}

	BMN struct {
		BMNCoord  *bmn.BMNCoord // MIND: BMNCoord is named, because XML and JSON serialization behave differently. An unnamed struct element will NOT be serialized by the XML encoder
		BMNString string
	}

	OSGB36 struct {
		OSGB36Coord  *osgb36.OSGB36Coord // MIND: OSGB36Coord is named, because XML and JSON serialization behave differently. An unnamed struct element will NOT be serialized by the XML encoder
		OSGB36String string
	}
)

// --------------------------------------------------------------------
// Serialization helper functions. The coordinates will be serialized according to the encoder in enc
//

func latlongToSerial(enc Encoder, latlong *cartconvert.PolarCoord, repformat cartconvert.LatLongFormat) (err error) {

	lat, long := cartconvert.LatLongToString(latlong, repformat)
	return enc.Encode(&LatLong{Lat: lat, Long: long, Fmt: repformat.String(), LatLongString: latlong.String()})
}

func geoHashToSerial(enc Encoder, geohash string) error {
	return enc.Encode(&GeoHash{GeoHash: geohash})
}

func utmToSerial(enc Encoder, utm *cartconvert.UTMCoord) error {
	return enc.Encode(&UTMCoord{UTMCoord: utm, UTMString: utm.String()})
}

func bmnToSerial(enc Encoder, bmn *bmn.BMNCoord) error {
	return enc.Encode(&BMN{BMNCoord: bmn, BMNString: bmn.String()})
}

func osgb36ToSerial(enc Encoder, osgb36 *osgb36.OSGB36Coord) error {
	return enc.Encode(&OSGB36{OSGB36Coord: osgb36, OSGB36String: osgb36.String()})
}

// serialize gets called by the respective handler methods to perform the serialization in the requested output representation
func serialize(enc Encoder, latlong *cartconvert.PolarCoord, oformat string) (err error) {
	switch oformat {
	case OFlatlongdeg:
		err = latlongToSerial(enc, latlong, cartconvert.LLFdms)
	case OFlatlongcomma:
		err = latlongToSerial(enc, latlong, cartconvert.LLFdeg)
	case OFgeohash:
		err = geoHashToSerial(enc, cartconvert.LatLongToGeoHash(latlong))
	case OFUTM:
		err = utmToSerial(enc, cartconvert.LatLongToUTM(latlong))
	case OFBMN:
		var bmnval *bmn.BMNCoord
		bmnval, err = bmn.WGS84LatLongToBMN(latlong, bmn.BMNZoneDet)
		if err == nil {
			err = bmnToSerial(enc, bmnval)
		}
	case OFOSGB:
		var osgb36val *osgb36.OSGB36Coord
		osgb36val, err = osgb36.WGS84LatLongToOSGB36(latlong)
		if err == nil {
			err = osgb36ToSerial(enc, osgb36val)
		}
	default:
		err = fmt.Errorf("Unsupported output format: '%s'", oformat)
	}
	return
}

// --------------------------------------------------------------------
// http handler methods corresponding to the restful methods
//

func latlongHandler(enc Encoder, req *http.Request, latlongstrval, oformat string) (err error) {

	if len(latlongstrval) > 0 {
		return fmt.Errorf("Latlong doesn't accept an input value. Use the parameters 'lat' and 'long' instead")
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
	return serialize(enc, latlong, oformat)
}

func geohashHandler(enc Encoder, req *http.Request, geohashstrval, oformat string) (err error) {
	var latlong *cartconvert.PolarCoord
	if latlong, err = cartconvert.GeoHashToLatLong(geohashstrval, nil); err != nil {
		return
	}
	return serialize(enc, latlong, oformat)
}

func utmHandler(enc Encoder, req *http.Request, utmstrval, oformat string) (err error) {
	var utmval *cartconvert.UTMCoord
	if utmval, err = cartconvert.AUTMToStruct(utmstrval, nil); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = cartconvert.UTMToLatLong(utmval); err != nil {
		return
	}
	return serialize(enc, latlong, oformat)
}

func bmnHandler(enc Encoder, req *http.Request, bmnstrval, oformat string) (err error) {
	var bmnval *bmn.BMNCoord
	if bmnval, err = bmn.ABMNToStruct(bmnstrval); err != nil {
		return
	}

	var latlong *cartconvert.PolarCoord
	if latlong, err = bmn.BMNToWGS84LatLong(bmnval); err != nil {
		return
	}
	return serialize(enc, latlong, oformat)
}

func osgbHandler(enc Encoder, req *http.Request, osgb36strval, oformat string) (err error) {
	var osgb36val *osgb36.OSGB36Coord
	if osgb36val, err = osgb36.AOSGB36ToStruct(osgb36strval, osgb36.OSGB36Leave); err != nil {
		return
	}
	return serialize(enc, osgb36.OSGB36ToWGS84LatLong(osgb36val), oformat)
}

// closure of the restful methods
//    enc: requested encoding scheme
//    req: calling context
//    value: coordinate value to be transformed
//    oformat: requested transformation representation, eg. utm, geohash
type restHandler func(enc Encoder, req *http.Request, value, oformat string) error

func (fn restHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// API error handler 
	// Recover from panic by setting http error 500 and letting the user know the reason
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	// enc keeps the requested encoding scheme as requested by content negotiation
	var enc Encoder
	// allocate buffer to which the http stream is written, until it gets responded. By doing so we keep the chance to trap errors and respond them to the caller
	buf := new(bytes.Buffer)

	// val: coordinate value
	// serialformat: serialization format
	// oformat: requested output format
	val := path.Base(req.URL.Path)
	serialformat := path.Ext(val)
	val = val[:len(val)-len(serialformat)]
	oformat := req.URL.Query().Get(OutputFormatSpec)

	switch serialformat {
	case JSONFormatSpec, "":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc = json.NewEncoder(buf)
	case XMLFormatSpec:
		w.Header().Set("Content-Type", "text/xml")
		enc = xml.NewEncoder(buf)
	default:
		panic(fmt.Sprintf("Unsupported serialization format: '%s'", serialformat))
	}

	if err := fn(enc, req, val, oformat); err != nil {
		// might as well panic(err) but we add some more info
		// we  serialize the error here in the chosen encoding
		enc.Encode(&Error{Error: fmt.Sprint(err)})
		w.WriteHeader(http.StatusInternalServerError)
	}
	buf.WriteTo(w)
}

// --------------------------------------------------------------------
// http part of the restful service: Templates, cache-variables, ...
//
type Link struct {
	*url.URL
	Documentation string
}

// docrootLinks gets initialized in docserv.go. It is used here to display a welcome screen
// with reference to documentation (if there is any)
var docrootLink *Link
var pageCache []*template.Template

type apidocpageLayout struct {
	APIRoot string
	APIRefs []Link
	DOCRoot *Link
}

const rootPage = `<!DOCTYPE HTML>
<html>
<head>
  <title>Cartconvert - Online cartography transformation</title>
</head>
<body>
  <h1>Cartconvert - Online cartography transformation</h1>
  <heading>
    This service provides a RESTFul API to perform cartography transformations.
  </heading>
  <nav>
    <p>
      <a href="{{.APIRoot}}">The API</a>
    </p>
    <p>
      <!-- TODO: set DOC root by template only when documentation is part of the prog -->
      <a href="{{.DOCRoot.URL}}">{{.DOCRoot.Documentation}}</a>
    </p>
  </nav>
</body>
</html>`

func rootHandler(w http.ResponseWriter, req *http.Request) {
	tpl := template.Must(template.New("root").Parse(rootPage))
	apiroot := apiroot()
	rootpage := &apidocpageLayout{APIRoot: apiroot, DOCRoot: docrootLink}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, rootpage); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	buf.WriteTo(w)
}

const apiPage = `<!DOCTYPE HTML>
<html>
<head>
  <title>Cartconvert - API page</title>
</head>
<body>
  <h1>Cartconvert - API page</h1>
  <heading>
    Root of API services
  </heading>
  <nav>
    <p>
      <a href="/">Back to main page</a>
    </p>
    <p>
      <!-- 
      <!-- TODO: iterate over APIs in a sorted manner a href="{{.APIRoot}}">The API</a> -->
      <!-- TODO: set DOC root by template only when documentation is part of the prog -->
      <a href="{{.DOCRoot.URL}}">{{.DOCRoot.Documentation}}</a>
    </p>
  </nav>
</body>
</html>`

func apiHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, apiPage)
}

// Definition of restful methods: combine API URI with handler method.
// For every API URI,there may be a corresponding documentation URI
type httphandlerfunc struct {
	restHandler
	docstring string
}

var httphandlerfuncs = map[string]httphandlerfunc{
	"latlong/": httphandlerfunc{latlongHandler, "Latitude, Longitude"},
	"geohash/": httphandlerfunc{geohashHandler, "Geohash"},
	"utm/":     httphandlerfunc{utmHandler, "UTM"},
	"bmn/":     httphandlerfunc{bmnHandler, "AT:Bundesmeldenetz"},
	"osgb/":    httphandlerfunc{osgbHandler, "UK:OSGB36"},
}

func init() {

	apiroot := apiroot()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc(apiroot, apiHandler)

	for function, handle := range httphandlerfuncs {
		http.Handle(apiroot+function, handle.restHandler)
	}
}
