package proxy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ehazlett/tlsutils"
)

var (
	certOrg  = "unknown"
	certBits = 2048
)

func getOrGenerateCA(certPath string) ([]byte, []byte, error) {
	caCertPath := filepath.Join(certPath, "ca.pem")
	caKeyPath := filepath.Join(certPath, "ca-key.pem")

	if contents, err := loadAllFiles(caCertPath, caKeyPath); err == nil {
		return contents[0], contents[1], nil
	} else if !os.IsNotExist(err) {
		return nil, nil, err
	}

	caCert, caKey, err := tlsutils.GenerateCACertificate(certOrg, certBits)
	if err != nil {
		return nil, nil, err
	}
	err = writeAllFiles(map[string][]byte{
		caCertPath: caCert,
		caKeyPath:  caKey,
	})
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Generated CA certs %v", []string{caCertPath, caKeyPath})
	return caCert, caKey, err
}

func loadAllFiles(paths ...string) ([][]byte, error) {
	bs := make([][]byte, len(paths))
	for i, p := range paths {
		var err error
		bs[i], err = ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}

func writeAllFiles(paths map[string][]byte) error {
	for path, content := range paths {
		if err := ioutil.WriteFile(path, content, 0600); err != nil {
			return err
		}
	}
	return nil
}
