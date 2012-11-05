<!DOCTYPE HTML>
<html>
<head>
  <title>Cartconvert - API page</title>
  <link rel="shortcut icon" href="../static/images/favicon.png" type="image/png"/> 
  <link rel="icon" href="../static/images/favicon.png" type="image/png"/>
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
    {{range .APIRefs}}<p>
      <a href="{{with $.DOCRoot}}{{.}}{{end}}{{.URL}}">{{.Documentation}}</a>{{end}}
    </p>
  </nav>
</body>
</html>