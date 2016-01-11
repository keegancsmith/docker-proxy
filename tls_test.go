package proxy

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func init() {
	// Smaller cert size for faster tests
	certBits = 512
}

func TestGetOrGenerateCA(t *testing.T) {
	certPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { os.RemoveAll(certPath) }()

	caCert, caKey, err := getOrGenerateCA(certPath)
	if err != nil {
		t.Fatal(err)
	}

	caCert1, caKey1, err := getOrGenerateCA(certPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(caCert, caCert1) {
		t.Error("getOrGenerateCA was not idempotent for caCert")
	}
	if !reflect.DeepEqual(caKey, caKey1) {
		t.Error("getOrGenerateCA was not idempotent for caKey")
	}
}

func TestGetOrGenerateServerCert(t *testing.T) {
	certPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { os.RemoveAll(certPath) }()

	c, err := getOrGenerateServerCert(certPath, []string{"127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.TLSConfig()
	if err != nil {
		t.Fatal(err)
	}

	c2, err := getOrGenerateServerCert(certPath, []string{"127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, c2) {
		t.Error("getOrGenerateServerCert was not idempotent")
	}
}
