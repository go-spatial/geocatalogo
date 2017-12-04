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

    "github.com/sirupsen/logrus"
    "github.com/tomkralidis/geocatalogo/config"
)

var LogLevels = map[string]logrus.Level {
    "DEBUG": logrus.DebugLevel,
    "INFO": logrus.InfoLevel,
    "WARN": logrus.WarnLevel,
    "ERROR": logrus.ErrorLevel,
    "FATAL": logrus.FatalLevel,
    "NONE": logrus.PanicLevel,
}

func InitLog(cfg *config.Config) logrus.Logger {
    var log = *logrus.New()

    // set defaults
    log.Level = logrus.PanicLevel
    log.Out =  os.Stderr
    log.Formatter = new(logrus.TextFormatter)
    log.Hooks = make(logrus.LevelHooks)

    // set to optionally write to logfile
    if cfg.Logging.Logfile != "" {
        f, err := os.OpenFile(cfg.Logging.Logfile, os.O_WRONLY | os.O_CREATE, 0644)
        if err != nil {
            panic(err)
        }
        log.Out = f
    }

    // set debug level
    log.SetLevel(LogLevels[cfg.Logging.Level])

    return log
}
