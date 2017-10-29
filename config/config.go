package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

// Config provides an object model for configuration.
type Config struct {
    Server struct {
        URL string
        MimeType string
        Encoding string
        Language string
        PrettyPrint bool
        Limit int
    }
    Logging struct {
        Level string
        Logfile string
    }
    Metadata struct {
        Identification struct {
            Title string
            Abstract string
            Keywords []string
            KeywordsType []string
            Fees string
            AccessConstraints string
        }
        Provider struct {
            Name string
            URL string
        }
        Contact struct {
            Name string
            Position string
            Address string
            City string
            StateOrProvince string
            PostalCode string
            Country string
            Phone string
            Fax string
            Email string
            URL string
            Hours string
            Instructions string
            Role string
        }
        Repository struct {
            Type string
            URL string
        }
    }
}

// GetConfig reads YAML configuration.
func GetConfig(filename string) Config {
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
