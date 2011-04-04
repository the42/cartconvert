// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This command reads coordinates in Bundesmeldenetz from stdin, performs a conversion,
// and writes to stdout. Errors are written to stderr.
//
// The target reference ellipsoid is always the WGS84Ellipsoid
//
// Usage of ./bmn2:
//  -of="deg": specify output format. Possible values are:  dms  geohash  utc  deg 
//
package main

import (
	"github.com/the42/cartconvert"
	"github.com/the42/cartconvert/bmn"
	"encoding/line"
	"fmt"
	"flag"
	"os"
	"strings"
)

type displayformat int

const (
	fmtunknown displayformat = iota
	deg
	dms
	utm
	geohash
)

var ofOptions = map[string]displayformat{"deg": deg, "dms": dms, "utm": utm, "geohash": geohash}


func main() {

	var ofcmdlinespec string
	var of displayformat
	var lines uint
	var instring, outstring, paramvalues string

	for key, _ := range ofOptions {
		paramvalues += fmt.Sprintf(" %s ", key)
	}

	flag.StringVar(&ofcmdlinespec, "of", "deg", "specify output format. Possible values are: "+paramvalues)
	flag.Parse()

	of = ofOptions[strings.ToLower(ofcmdlinespec)]

	liner := line.NewReader(os.Stdin, 100)
	longline := false

	for data, prefix, err := liner.ReadLine(); err != os.EOF; data, prefix, err = liner.ReadLine() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "bmn2 %d: %s\n", lines, err)
			continue
		}

		if prefix {
			longline = true
			continue
		}

		if longline {
			longline = false
			continue
		}

		lines++

		instring = strings.TrimSpace(string(data))

		if len(instring) == 0 {
			continue
		}

		bcoord, err := bmn.ABMNToStruct(instring)

		if err != nil {
			fmt.Fprintf(os.Stderr, "bmn2: error on line %d: %s\n", lines, err)
			continue
		}

		pc, err := bmn.BMNToWGS84LatLong(bcoord)

		if err != nil {
			fmt.Fprintf(os.Stderr, "bmn2: error on line %d: %s (BMN does not return a lat/long bearing)\n", lines, err)
			continue
		}

		switch of {
		case deg:
			outstring = cartconvert.LatLongToString(pc, cartconvert.LLFdeg)
		case dms:
			outstring = cartconvert.LatLongToString(pc, cartconvert.LLFdms)
		case utm:
			outstring = cartconvert.LatLongToUTM(pc).String()
		case geohash:
			outstring = cartconvert.LatLongToGeoHash(pc)
		default:
			fmt.Fprintln(os.Stderr, "Unrecognized output specifier")
			flag.Usage()
			fmt.Fprintf(os.Stderr, "possible values are: [%s]\n", paramvalues)
			fmt.Fprintln(os.Stderr, "]")
			os.Exit(2)
		}
		fmt.Fprintf(os.Stdout, "%s\n", outstring)
	}
}
