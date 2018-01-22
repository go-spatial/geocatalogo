package geocatalogo_test

import (
	"fmt"
	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/config"
	"github.com/go-spatial/geocatalogo/repository"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	fmt.Print("INIT")
	testLog := logrus.New()

	fmt.Println("loading from env")
	testConfig := config.LoadFromEnv()
	fmt.Println("creating new repo")

	err := repository.New(testConfig, testLog)

	if err != nil {
		fmt.Println("Repository not created")
	}
}

func TestSmokeTest(t *testing.T) {
	cat, _ := geocatalogo.NewFromEnv()
	if cat.Config.Server.URL != "http://localhost:8001/" {
		t.Error("Incorrect value")
	}
}
