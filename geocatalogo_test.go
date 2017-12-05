package geocatalogo_test

import (
	"testing"
	"github.com/tomkralidis/geocatalogo"
)

func TestSmokeTest(t *testing.T) {
	mycatalogo := geocatalogo.New()
	if mycatalogo.Config.Server.URL != "http://localhost:8001/" {
		t.Error("Incorrect value")
	}
}
