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
type oSGB36StringToStructParam struct {
	osgb36coord string
	prec        OSGB36prec
}

type oSGB36StringToStructTest struct {
	in  oSGB36StringToStructParam
	out *OSGB36Coord
}

func newNewOSGB36CoordHelper(zone string, easting, northing uint, prec OSGB36prec, relheight float64) *OSGB36Coord {
	coord, _ := NewOSGB36Coord(zone, easting, northing, prec, relheight)
	return coord
}

var oSGB36StringToStructTestssuc = []oSGB36StringToStructTest{
	{
		oSGB36StringToStructParam{"NN1660071200", OSGB36Auto},
		newNewOSGB36CoordHelper("NN", 1660, 7120, OSGB36Auto, 0),
	},
	{
		oSGB36StringToStructParam{"NN", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 0, Northing: 0, gridLen: 0},
	},
	{
		oSGB36StringToStructParam{"NN11", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 1, Northing: 1, gridLen: 1},
	},
	{
		oSGB36StringToStructParam{"NN1212", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 12, Northing: 12, gridLen: 2},
	},
	{
		oSGB36StringToStructParam{"NN123123", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 123, Northing: 123, gridLen: 3},
	},
	{
		oSGB36StringToStructParam{"NN12341234", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 1234, Northing: 1234, gridLen: 4},
	},
	{
		oSGB36StringToStructParam{"NN1234512345", OSGB36Auto},
		&OSGB36Coord{Zone: "NN", Easting: 12345, Northing: 12345, gridLen: 5},
	},
	{
		oSGB36StringToStructParam{"NN1234512345", OSGB36_2},
		&OSGB36Coord{Zone: "NN", Easting: 12, Northing: 12, gridLen: 2},
	},
	{
		oSGB36StringToStructParam{"NN1234512345", OSGB36_2},
		newNewOSGB36CoordHelper("NN", 12, 12, OSGB36Auto, 0),
	},
	{
		oSGB36StringToStructParam{"NN166712", OSGB36_5},
		newNewOSGB36CoordHelper("NN", 1660, 7120, OSGB36_5, 0),
	},
}

func osgb36equal(osgb1, osgb2 *OSGB36Coord) bool {
	p1 := fmt.Sprintf("%s", osgb1)
	p2 := fmt.Sprintf("%s", osgb2)
	return p1 == p2
}

func TestOSGB36StringToStruct(t *testing.T) {
	for cnt, test := range oSGB36StringToStructTestssuc {
		out, err := AOSGB36ToStruct(test.in.osgb36coord, test.in.prec)

		if err != nil {
			t.Errorf("TestOSGB36StringToStruct [%d]: Error: %s", cnt, err)
		} else {
			formattspec := "TestOSGB36StringToStruct [%d]: Expected %s, got %s"
			if !osgb36equal(test.out, out) {
				t.Errorf(formattspec, cnt, test.out, out)
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
		&OSGB36Coord{Zone: "SE", Easting: 29793, Northing: 33798, gridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 53.79965, Longitude: -1.54915},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 166, Northing: 712, gridLen: 3, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796088, Longitude: -5.0039304},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 16600, Northing: 71200, gridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0047120},
	},
	{
		&OSGB36Coord{Zone: "NN", Easting: 16650, Northing: 71250, gridLen: 5, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0039304},
	},
	{
		&OSGB36Coord{Zone: "SV", Easting: 0, Northing: 0, gridLen: 0, el: cartconvert.Airy1830Ellipsoid},
		&cartconvert.PolarCoord{Latitude: 49.766795, Longitude: -7.557349},
	},
}

func latlongequal(pcp1, pcp2 *cartconvert.PolarCoord) bool {
	pp1s := fmt.Sprintf("%.5g %.5g", pcp1.Latitude, pcp1.Longitude)
	pp2s := fmt.Sprintf("%.5g %.5g", pcp2.Latitude, pcp2.Longitude)

	return pp1s == pp2s
}

func TestOSGB36ToWGS84LatLong(t *testing.T) {
	for cnt, test := range oSGB36ToWGS84LatLongTests {

		out := OSGB36ToWGS84LatLong(test.in)

		if !latlongequal(test.out, out) {
			t.Errorf("OSGB36ToWGS84LatLong:%d [%s]: Expected %s, got %s", cnt, test.in, test.out, out)
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
		&OSGB36Coord{Zone: "SE", Easting: 29793, Northing: 33798, gridLen: 5, el: cartconvert.Airy1830Ellipsoid},
	},
	{
		&cartconvert.PolarCoord{Latitude: 56.796557, Longitude: -5.0039304},
		&OSGB36Coord{Zone: "NN", Easting: 16650, Northing: 71250, gridLen: 5, el: cartconvert.Airy1830Ellipsoid},
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
