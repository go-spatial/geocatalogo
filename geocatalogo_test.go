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
	fmt.Print("INIT")
	testLog := logrus.New()

	fmt.Println("loading from env")
	testConfig := config.LoadFromEnv()
	fmt.Println("creating new repo")

	status := repository.New(testConfig, testLog)

	if !status {
		fmt.Println("Repository not created")
	}
}

func TestSmokeTest(t *testing.T) {
	mycatalogo := geocatalogo.New()
	if mycatalogo.Config.Server.URL != "http://localhost:8001/" {
		t.Error("Incorrect value")
	}
}
