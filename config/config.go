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

package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config provides an object model for configuration.
type Config struct {
	Server struct {
		URL         string
		MimeType    string
		Encoding    string
		Language    string
		PrettyPrint bool
		Limit       int
	}
	Logging struct {
		Level   string
		Logfile string
	}
	Metadata struct {
		Identification struct {
			Title             string
			Abstract          string
			Keywords          []string
			KeywordsType      string
			Fees              string
			AccessConstraints string
		}
		Provider struct {
			Name string
			URL  string
		}
		Contact struct {
			Name            string
			Position        string
			Address         string
			City            string
			StateOrProvince string
			PostalCode      string
			Country         string
			Phone           string
			Fax             string
			Email           string
			URL             string
			Hours           string
			Instructions    string
			Role            string
		}
	}
	Repository struct {
		Type     string
		URL      string
		Mappings map[string]string
	}
}

// read environment variables into configuration
func LoadFromEnv() Config {
	var cfg Config
	cfg.Repository.Mappings = make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "GEOCATALOGO_SERVER_URL" {
			cfg.Server.URL = pair[1]
		} else if pair[0] == "GEOCATALOGO_SERVER_MIMETYPE" {
			cfg.Server.MimeType = pair[1]
		} else if pair[0] == "GEOCATALOGO_SERVER_ENCODING" {
			cfg.Server.Encoding = pair[1]
		} else if pair[0] == "GEOCATALOGO_SERVER_LANGUAGE" {
			cfg.Server.Language = pair[1]
		} else if pair[0] == "GEOCATALOGO_SERVER_PRETTY_PRINT" {
			cfg.Server.PrettyPrint, _ = strconv.ParseBool(pair[1])
		} else if pair[0] == "GEOCATALOGO_SERVER_LIMIT" {
			cfg.Server.Limit, _ = strconv.Atoi(pair[1])
		} else if pair[0] == "GEOCATALOGO_LOGGING_LEVEL" {
			cfg.Logging.Level = pair[1]
		} else if pair[0] == "GEOCATALOGO_LOGGING_LOGFILE" {
			cfg.Logging.Logfile = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_TITLE" {
			cfg.Metadata.Identification.Title = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_ABSTRACT" {
			cfg.Metadata.Identification.Abstract = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_KEYWORDS" {
			cfg.Metadata.Identification.Keywords = strings.Split(pair[1], ",")
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_KEYWORDS_TYPE" {
			cfg.Metadata.Identification.KeywordsType = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_FEES" {
			cfg.Metadata.Identification.Fees = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_IDENTIFICATION_ACCESSCONSTRAINTS" {
			cfg.Metadata.Identification.AccessConstraints = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_PROVIDER_NAME" {
			cfg.Metadata.Provider.Name = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_PROVIDER_URL" {
			cfg.Metadata.Provider.URL = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_NAME" {
			cfg.Metadata.Contact.Name = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_POSITION" {
			cfg.Metadata.Contact.Position = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_ADDRESS" {
			cfg.Metadata.Contact.Address = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_CITY" {
			cfg.Metadata.Contact.City = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_STATEORPROVINCE" {
			cfg.Metadata.Contact.StateOrProvince = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_POSTALCODE" {
			cfg.Metadata.Contact.PostalCode = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_COUNTRY" {
			cfg.Metadata.Contact.Country = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_PHONE" {
			cfg.Metadata.Contact.Phone = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_FAX" {
			cfg.Metadata.Contact.Fax = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_EMAIL" {
			cfg.Metadata.Contact.Email = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_URL" {
			cfg.Metadata.Contact.URL = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_HOURS_OF_SERVICE" {
			cfg.Metadata.Contact.Hours = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_CONTACT_INSTRUCTIONS" {
			cfg.Metadata.Contact.Instructions = pair[1]
		} else if pair[0] == "GEOCATALOGO_METADATA_ROLE" {
			cfg.Metadata.Contact.Role = pair[1]
		} else if pair[0] == "GEOCATALOGO_REPOSITORY_TYPE" {
			cfg.Repository.Type = pair[1]
		} else if pair[0] == "GEOCATALOGO_REPOSITORY_URL" {
			cfg.Repository.URL = pair[1]
		} else if strings.HasPrefix(pair[0], "GEOCATALOGO_REPOSITORY_MAPPINGS") {
			tokens := strings.Split(pair[0], "GEOCATALOGO_REPOSITORY_MAPPINGS_")
			key := strings.ToLower(tokens[1])
			cfg.Repository.Mappings[key] = pair[1]
		}
	}
	return cfg
}

// read YAML into configuration
func LoadFromFile(filename string) Config {
	var cfg Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
