///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2017 Tom Kralidis
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


// Package geocatalogo
package geocatalogo

import (
    "os"

    "github.com/tomkralidis/geocatalogo/config"
    "github.com/tomkralidis/geocatalogo/repository"
)

// VERSION provides the geocatalogo version installed.
const VERSION string = "0.1.0"

func New() {
    // get configuration
    cfg := config.GetConfig(os.Getenv("GEOCATALOGO_CONFIG"))

    // setup logging
    log := InitLog(&cfg)

    log.Info("geocatalogo Version " + VERSION)
    log.Info("Configuration: " + os.Getenv("GEOCATALOGO_CONFIG"))

    // read backend
    log.Info("Loading repository")
    repo := repository.Open(cfg, &log)

    log.Info("Repository loaded (type) " + repo.Type)
    return
}
