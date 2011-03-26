package main

import (
	"strconv"
	"http"
	"io"
	"template"
	"github.com/the42/cartconvert"
)

// Better read it from an INI-File
const staticWebDir = "../static/"
const templateDir = "../templates/"

// Better read it from an INI-File or by a dirreader
var templateNames = []string{
	"layout.tpl",
	"edit.tpl",
}

var templates = make(map[string]*template.Template)

type CoordinateTrans struct {
	Xcoord, Ycoord string
	Fromcs, Tocs   string
}

func evalTemplate(wr io.Writer, formatter string, data ...interface{}) {
	err := templates[formatter].Execute(wr, data[0])
	if err != nil {
		print(err.String())
	}
}


func initTemplates() {
	fmap := template.FormatterMap{}

	for _, name := range templateNames {
		fmap[name] = evalTemplate
	}

	for _, name := range templateNames {
		templates[name] = template.MustParseFile(templateDir+name, fmap)
	}
}


func transCoordinate(ct CoordinateTrans) (coord CoordinateTrans) {
	coord = CoordinateTrans{ct.Xcoord, ct.Ycoord, "GSK", "UTM"}
	return
}

// Render edit page
func editHandler(w http.ResponseWriter, req *http.Request) {
	xcoord := req.FormValue("xcoord")
	ycoord := req.FormValue("ycoord")
	// xcoord, ycoord = ycoord, xcoord

	flat, _ := strconv.Atof64(xcoord)
	flong, _ := strconv.Atof64(ycoord)

	gc := cartconvert.DirectTransverseMercator(&cartconvert.PolarCoord{Latitude: flat, Longitude: flong, El: cartconvert.Airy1830Ellipsoid}, 49, -2, 0.9996012717, 400000, -100000)

	templates["layout.tpl"].Execute(w,
		transCoordinate(CoordinateTrans{Xcoord: strconv.Ftoa64(gc.X, 'f', 6), Ycoord: strconv.Ftoa64(gc.Y, 'f', 6)}))
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticWebDir+r.URL.Path)
}


func main() {

	initTemplates()

	http.HandleFunc("/", staticFileHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":1111", nil)
}
