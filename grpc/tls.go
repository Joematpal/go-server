package grpc

import (
	"crypto/tls"
	"encoding/base64"

	"google.golang.org/grpc/credentials"
)

func Credentials(certs tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{certs}})
}

// ParseCertificates will the check (in this order) if the strings coming in are filepath, base64 encoded, or plaintext pem files
// So if you get a x509 key pair error that is not from a file not found it probably means that everything is broken and you are not passing in valid x509 certs
func ParseCertificates(pubCert, privCert string) (tls.Certificate, error) {
	var certs tls.Certificate
	var err error
	certs, err = tls.LoadX509KeyPair(pubCert, privCert)

	if err != nil {
		pubCertByte, err := base64.StdEncoding.DecodeString(pubCert)
		if err != nil {
			// return certs, err
			// bad certs might just be plain text
			return tls.X509KeyPair([]byte(pubCert), []byte(privCert))
		}
		privCertByte, err := base64.StdEncoding.DecodeString(privCert)
		if err != nil {
			return certs, err
		}
		return tls.X509KeyPair(pubCertByte, privCertByte)
	}
	return certs, nil
}
