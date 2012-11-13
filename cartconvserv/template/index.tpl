<!DOCTYPE HTML>
<html>
<head>
  <title>Cartconvert - Online cartography transformation</title>
  <link rel="shortcut icon" href="./static/images/favicon.png" type="image/png"/> 
  <link rel="icon" href="./static/images/favicon.png" type="image/png"/>
  <link rel="stylesheet" type="text/css" href="./static/css/styles.css"/> 
</head>
<body>
<body>
 <a href="https://github.com/you"><img style="position: absolute; top: 0; right: 0; border: 0;"
   src="https://s3.amazonaws.com/github/ribbons/forkme_right_orange_ff7600.png" alt="Fork me on GitHub"></a>
</body>
  <heading>
    <h1>Cartconvert - Online cartography transformation</h1>
    <p>This service provides a RESTFul API to perform cartography transformations.</p>
  </heading>
  <nav>
  <ul>
    <li><a href="{{.APIRoot}}">The API</a></li>
    {{if .DOCRoot}}<li><a href="{{.DOCRoot.URL}}">{{.DOCRoot.Documentation}}</a></li>{{end}}
  </ul>
  </nav>
  <footer>
<a href="https://twitter.com/myprivate42" class="twitter-follow-button" data-show-count="false" data-size="large">Follow @myprivate42</a>
<script>!function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0];if(!d.getElementById(id)){js=d.createElement(s);js.id=id;js.src="//platform.twitter.com/widgets.js";fjs.parentNode.insertBefore(js,fjs);}}(document,"script","twitter-wjs");</script>
  </footer>
</body>
</html>