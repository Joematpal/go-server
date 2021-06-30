package grpc

import (
	"crypto/tls"
	"encoding/base64"

	"google.golang.org/grpc/credentials"
)

func Credentials(certs tls.Certificate) credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{certs}})
}
func ParseCertificates(pubCert, privCert string) (tls.Certificate, error) {
	var certs tls.Certificate
	var err error
	certs, err = tls.LoadX509KeyPair(pubCert, privCert)
	if err != nil {
		pubCertByte, err := base64.StdEncoding.DecodeString(pubCert)
		if err != nil {
			return certs, err
		}
		privCertByte, err := base64.StdEncoding.DecodeString(privCert)
		if err != nil {
			return certs, err
		}
		return tls.X509KeyPair(pubCertByte, privCertByte)
	}
	return certs, nil
}
