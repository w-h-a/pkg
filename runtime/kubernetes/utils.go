package kubernetes

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"
)

func DetectNamespace() (string, error) {
	nsPath := path.Join(serviceAccountPath, "namespace")

	// make sure it's a file and that we can read it
	if file, err := os.Stat(nsPath); err != nil {
		return "", err
	} else if file.IsDir() {
		return "", ErrReadNamespace
	}

	// read it
	ns, err := os.ReadFile(nsPath)
	if err != nil {
		return "", err
	}

	return string(ns), nil
}

func CertPoolFromFile(filename string) (*x509.CertPool, error) {
	certs, err := CertificatesFromFile(filename)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()

	for _, cert := range certs {
		pool.AddCert(cert)
	}

	return pool, nil
}

func CertificatesFromFile(file string) ([]*x509.Certificate, error) {
	if len(file) == 0 {
		return nil, errors.New("error reading certificates from an empty filename")
	}

	pemBlock, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	certs, err := CertsFromPEM(pemBlock)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", file, err)
	}

	return certs, nil
}

func CertsFromPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false

	certs := []*x509.Certificate{}

	for len(pemCerts) > 0 {
		var block *pem.Block

		block, pemCerts = pem.Decode(pemCerts)

		if block == nil {
			break
		}

		// only use PEM "CERTIFICATE" blocks without extra headers
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}

		certs = append(certs, cert)

		ok = true
	}

	if !ok {
		return certs, errors.New("could not read any certificates")
	}

	return certs, nil
}
