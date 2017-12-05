package geocatalogo_test

import (
	"github.com/tomkralidis/geocatalogo"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	mycatalogo := geocatalogo.New()
	if mycatalogo.Config.Server.URL != "http://localhost:8001/" {
		t.Error("Incorrect value")
	}
}
