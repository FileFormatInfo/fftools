package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: certinfo url|file\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "          certinfo https://google.com\n")
	fmt.Fprintf(os.Stderr, "          certinfo mycert.pem\n")
	os.Exit(1)
}

func getHttpsCerts(url string) ([]*x509.Certificate, error) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	host := strings.TrimPrefix(url, "https://")
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	if idx := strings.Index(host, ":"); idx == -1 {
		host += ":443" // Default to port 443 for HTTPS
	}

	conn, err := tls.Dial("tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s: %w", host, err)
	}
	defer conn.Close()

	return conn.ConnectionState().PeerCertificates, nil
}

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	var certs []*x509.Certificate

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "https://") {
			httpsCerts, httpsErr := getHttpsCerts(arg)
			if httpsErr != nil {
				log.Fatalf("Error getting HTTPS certificates: %v", httpsErr)
			}
			certs = httpsCerts
		} else {
			fileCert, fileErr := os.ReadFile(arg)
			if fileErr != nil {
				log.Fatalf("Error reading certificate file: %v", fileErr)
			}
			if bytes.HasPrefix(fileCert, []byte("-----BEGIN ")) {
				var block *pem.Block
				rest := fileCert
				certs = make([]*x509.Certificate, 0)
				for {
					block, rest = pem.Decode(rest)
					if block == nil {
						break
					}
					cert, parseErr := x509.ParseCertificate(block.Bytes)
					if parseErr != nil {
						log.Fatalf("Error parsing PEM certificate: %v", parseErr)
					}
					certs = append(certs, cert)
				}
			} else {
				cert, parseErr := x509.ParseCertificate(fileCert)
				if parseErr != nil {
					log.Fatalf("Error parsing certificate: %v", parseErr)
				}
				certs = make([]*x509.Certificate, 1)
				certs[0] = cert
			}
		}

		fmt.Printf("Certificate Information for %s:\n", arg)
		for idx, cert := range certs {
			fmt.Printf("\tCertificate %d:\n", idx+1)
			fmt.Printf("\t\tIssuer      : %s\n", cert.Issuer)
			fmt.Printf("\t\tCommon Name : %s \n", cert.Issuer.CommonName)
			fmt.Printf("\t\tSubject     : %s\n", cert.Subject)
			fmt.Printf("\t\tCommon Name : %s \n", cert.Subject.CommonName)
			fmt.Printf("\t\tStart       : %s \n", cert.NotBefore.Format("2006-01-02"))
			fmt.Printf("\t\tExpiry      : %s \n", cert.NotAfter.Format("2006-01-02"))
			fmt.Printf("\t\tKey Usage   : %v \n", cert.KeyUsage)
		}
	}
}
