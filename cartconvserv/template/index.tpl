<!DOCTYPE HTML>
<html>
<head>
  <title>Cartconvert - Online cartography transformation</title>
  <link rel="shortcut icon" href="./static/images/favicon.png" type="image/png"/> 
  <link rel="icon" href="./static/images/favicon.png" type="image/png"/>
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
    {{if .DOCRoot}}<p>
      <a href="{{.DOCRoot.URL}}">{{.DOCRoot.Documentation}}</a>
    </p>{{end}}
  </nav>
</body>
</html>