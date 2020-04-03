package goswyftx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joshturge/goswyftw"
)

var (
	client *goswyftx.Client
)

func TestNewClient(t *testing.T) {
	var err error
	client, err = goswyftx.NewClient(os.Getenv("API_KEY"), os.Getenv("TOKEN"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRequest(t *testing.T) {
	TestNewClient(t)

	keys, err := client.Authentication().GetKeys()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(keys) == 0 {
		t.Error("length of keys is 0")
		t.FailNow()
	}

	for _, key := range keys {
		fmt.Println("Key: " + key.ID)
	}
}
