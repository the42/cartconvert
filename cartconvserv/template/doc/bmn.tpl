{{define "Back"}}../{{end}}{{define "Payload"}}
  <header>
    <h1><a href=".">Documentation for Bundesmeldenetz</a></h1>
  </header>
  <h2>Examples</h2>
  <p><a href="../..{{.APIRoot}}/bmn/M34 703168 374510.json?outputformat=latlongdeg">M34 703168 374510</a> as Latitude / Longitude, result JSON-encoded. View on <a href="#">OpenStreetMap</a></p>
  <p><a href="../..{{.APIRoot}}/bmn/M34 703168 374510.xml?outputformat=utm">M34 703168 374510</a> as UTM, result XML-encoded. View on <a href="#">OpenStreetMap</a></p>
  <iframe class="documentation" src="http://markdress.org/raw.github.com/the42/cartconvert/master/cartconvserv/README.md#bmnconversion">
    <p><a href="https://github.com/the42/cartconvert/tree/master/cartconvserv#bmn---conversions">Documentation on Github</a> (authorative developer source)
    </p>
  </iframe>
  {{end}}