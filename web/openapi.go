///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2019 Tom Kralidis
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.
//
///////////////////////////////////////////////////////////////////////////////

// Package web - Simple HTTP Wrapper
package web

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"reflect"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/config"
)

var SwaggerHTML = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI - FIXME TITLE</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.22.2/swagger-ui.css" >
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@3.22.2/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@3.22.2/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }
      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.22.2/swagger-ui-bundle.js"> </script>
    <script src="https://unpkg.com/swagger-ui-dist@3.22.2/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      ui = SwaggerUIBundle({
        url: 'FIXME',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
         SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: 'StandaloneLayout'
      })
      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  <style>.swagger-ui .topbar .download-url-wrapper { display: none } undefined</style>
  </body>
</html>`

// GenerateOpenAPIDocument generates an OpenAPI Document
func GenerateOpenAPIDocument(cfg config.Config) []byte {
	var b []byte
	openAPI3Schema := &openapi3.Swagger{
		OpenAPI: "3.0.0",
		Info: openapi3.Info{
			Title:       cfg.Metadata.Identification.Title,
			Description: cfg.Metadata.Identification.Abstract,
			Version:     "0.0.1",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
	}
	b = geocatalogo.Struct2JSON(&openAPI3Schema, cfg.Server.PrettyPrint)
	fmt.Println(reflect.TypeOf(openAPI3Schema).String())
	fmt.Println(b)
	return b
}
