package main

import (
	"github.com/garyburd/twister/server"
	"github.com/garyburd/twister/web"
	"github.com/ziutek/kview"
	"strconv"
	"github.com/the42/cartconvert"
)

type CoordinateTrans struct {
	xcoord, ycoord string
	fromcs, tocs   string
}

// Get an article
func transCoordinate(ct CoordinateTrans) (coord *CoordinateTrans) {
	coord = &CoordinateTrans{ct.xcoord, ct.ycoord, "GSK", "UTM"}
	return
}


type ViewCtx struct {
	edit interface{}
}

// Render edit page
func edit(req *web.Request) {
	xcoord := req.Param.Get("xcoord")
	ycoord := req.Param.Get("ycoord")
	// xcoord, ycoord = ycoord, xcoord

	flat, _ := strconv.Atof64(xcoord)
	flong, _ := strconv.Atof64(ycoord)

	// cart := cartconvert.WGS84Ellipsoid.GeocoordtoCartesian( &cartconvert.GeoCoord{ Latitude: flat, Longitude: flong})
	// cart = helmert.WGS84toMGITransformer.Apply(cart)
	gc := cartconvert.DirectTransverseMercator(&cartconvert.PolarCoord{Latitude: flat, Longitude: flong, El: cartconvert.Airy1830Ellipsoid}, 49, -2, 0.9996012717, 400000, -100000)

	edit_view.Exec(
		req.Respond(web.StatusOK),
		ViewCtx{transCoordinate(CoordinateTrans{xcoord: strconv.Ftoa64(gc.X, 'f', 6), ycoord: strconv.Ftoa64(gc.Y, 'f', 6)})},
	)
}

var edit_view kview.View

func viewInit() {
	// Load layout template
	kview.TemplatesDir = "../../template"
	layout := kview.New("layout.kt")

	// Create edit page
	edit_view = layout.Copy()
	edit_view.Div("Edit", kview.New("edit.kt"))
}


func main() {
	viewInit()

	router := web.NewRouter().
		Register("/", "GET", edit, "POST", edit).
		Register("/style.css", "GET", web.FileHandler("../../static/style.css"))

	h := web.ProcessForm(1000, false, router)
	server.Run("localhost:1111", h)

}
