{{define "Back"}}..{{end}}{{define "Payload"}}
  <header>
    <h1><a href=".">Documentation for OSGB36</a></h1>
  </header>
  <h2>Examples</h2>
  <p>
    <a id="osm1" href="#">NN 123 123</a> as <a href="{{.APIRoot}}/osgb/NN 123 123.json?outputformat=osgb">OSGB36, JSON-encoded</a>,
    as <a href="{{.APIRoot}}/osgb/NN 123 123.xml?outputformat=latlongcomma">Lat / Long in fractions, XML-encoded</a>.
  </p>
  <h2>Reference</h2>
  <p>
    <a href="http://en.wikipedia.org/wiki/OSGB36">Wikipedia [EN]</a>, <a href="http://gps.ordnancesurvey.co.uk/etrs89geo_natgrid.asp">UK Ordnance Survey</a>
  </p>
  <h2>Embeded documentation</h2>
  <iframe class="documentation" src="http://markdress.org/raw.github.com/the42/cartconvert/master/cartconvserv/README.md#osgb36conversion">
    <p><a href="https://github.com/the42/cartconvert/blob/master/cartconvserv/README.md#osgb36---conversions-">Documentation on Github</a> (authorative developer source)
    </p>
  </iframe>
  <script>
    document.getElementById("osm1").addEventListener('click', function() {return osmload("{{.APIRoot}}/osgb/NN 123 123.json?outputformat=latlongcomma")});
  </script>
  {{end}}