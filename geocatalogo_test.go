package geocatalogo_test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tomkralidis/geocatalogo"
	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/repository"
	"testing"
)

func init() {
	testLog := logrus.New()

	testConfig := config.LoadFromEnv()

	_, err := repository.New(testConfig, testLog)

	if err != nil {
		fmt.Println(err)
	}
}

func TestSmokeTest(t *testing.T) {
	mycatalogo := geocatalogo.New()
	if mycatalogo.Config.Server.URL != "http://localhost:8001/" {
		t.Error("Incorrect value")
	}
}
