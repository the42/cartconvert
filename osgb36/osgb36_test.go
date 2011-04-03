// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// Automated tests for the cartconvert/osgb36 package
package osgb36

import (
	"fmt"
	"github.com/the42/cartconvert"
	"testing"
)

// ## AOSGB36ToStruct
type oSGB36StringToStructTest struct {
	in  string
	out *OSGB36Coord
}

var oSGB36StringToStructTestssuc = []oSGB36StringToStructTest{
	{
		"NN1660071200",
		&OSGB36Coord{Zone: "NN", Easting: 166, Northing: 712, GridLen: 3},
	},
		{
		"NN",
		&OSGB36Coord{Zone: "NN", Easting: 0, Northing: 0, GridLen: 0},
	},

}

func osgb36equal(osgb1, osgb2 *OSGB36Coord) bool {
	p1 := fmt.Sprintf("%s", osgb1)
	p2 := fmt.Sprintf("%s", osgb2)
	return p1 == p2
}

func TestOSGB36StringToStruct(t *testing.T) {
	for cnt, test := range oSGB36StringToStructTestssuc {
		out, err := AOSGB36ToStruct(test.in, OSGB36Auto)

		if err != nil {
			t.Errorf("TestOSGB36StringToStruct [%d]: Error: %s", cnt, err)
		} else {
			if !osgb36equal(test.out, out) {
				t.Errorf("TestOSGB36StringToStruct [%d]: Expected %s, got %s", cnt, test.out, out)
			}
		}
	}
}


// ## OSGB36ToWGS84LatLong
type oSGB36ToWGS84LatLongTest struct {
	in  *OSGB36Coord
	out *cartconvert.PolarCoord
}

var oSGB36ToWGS84LatLongTests = []oSGB36ToWGS84LatLongTest{
	{
		&OSGB36Coord{Zone: "SE", Easting: 29793, Northing: 33798, GridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 53.79965, Longitude: -1.54915},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 166, Northing: 712, GridLen: 3, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796088, Longitude: -5.0039304},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 16600, Northing: 71200, GridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0047120},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 16650, Northing: 71250, GridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0039304},
	},
}

func latlongequal(pcp1, pcp2 *cartconvert.PolarCoord) bool {
	pp1s := fmt.Sprintf("%.5g %.5g", pcp1.Latitude, pcp1.Longitude)
	pp2s := fmt.Sprintf("%.5g %.5g", pcp2.Latitude, pcp2.Longitude)

	return pp1s == pp2s
}

func TestOSGB36ToWGS84LatLong(t *testing.T) {
	for cnt, test := range oSGB36ToWGS84LatLongTests {

		out, err := OSGB36ToWGS84LatLong(test.in)

		if err != nil {
			t.Errorf("OSGB36ToWGS84LatLong [%d]: Error: %s", cnt, err)
		} else {
			if !latlongequal(test.out, out) {
				t.Errorf("OSGB36ToWGS84LatLong:%d [%s]: Expected %s, got %s", cnt, test.in, test.out, out)
			}
		}
	}
}

// ## WGS84LatLongToOSGB36
type wGS84LatLongToOSGB36Test struct {
	in  *cartconvert.PolarCoord
	out *OSGB36Coord
}

var wGS84LatLongToBMNTests = []wGS84LatLongToOSGB36Test{
	{
		&cartconvert.PolarCoord{Latitude: 53.79965, Longitude: -1.54915},
		&OSGB36Coord{Zone: "SE", Easting: 29793, Northing: 33798, el: cartconvert.Airy1830Ellipsoid},
	},
	{
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0039304},
		&OSGB36Coord{Zone: "NN", Easting: 16650, Northing: 71250, el: cartconvert.Airy1830Ellipsoid},
	},
}

func TestWGS84LatLongToOSGB36(t *testing.T) {
	for cnt, test := range wGS84LatLongToBMNTests {
		out, err := WGS84LatLongToOSGB36(test.in)
		if err != nil {
			t.Errorf("OSGB36ToWGS84LatLong [%d]: Error: %s", cnt, err)
		} else {
			if !osgb36equal(test.out, out) {
				t.Errorf("WGS84LatLongToOSGB36:%d [%s]: Expected %s, got %s", cnt, test.in, test.out, out)
			}
		}
	}
}
