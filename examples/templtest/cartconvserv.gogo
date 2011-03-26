package test

import (
	"github.com/the42/cartconvert"
	"github.com/garyburd/twister/server"
	"github.com/garyburd/twister/web"
	"github.com/ziutek/kview"
	"json"
	"io"
)

var layout kview.View

func mainView(req *web.Request) {
	layout.Exec(req.Respond(web.StatusOK))
}

func jsontest(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html")
	x, err := json.Marshal(cartconvert.LatLongToGeoHash(&cartconvert.PolarCoord{Latitude: 49.3, Longitude: 20.0}))
	if err == nil {
		io.WriteString(w, string(x))
	} else {
		io.WriteString(w, "Error: "+err.String())
	}

}

func viewInit() {
	// Load layout template
	kview.TemplatesDir = "../template"
	layout = kview.New("cartconvserverMainView.kt")
}

func main() {
	viewInit()

	router := web.NewRouter().
		Register("/", "GET", mainView).
		Register("/convert", "GET", jsontest)

	h := web.ProcessForm(1000, false, router)
	server.Run("localhost:1111", h)
}
