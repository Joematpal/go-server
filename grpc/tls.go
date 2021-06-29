package grpc

import (
	"crypto/tls"
	"encoding/base64"

	"google.golang.org/grpc/credentials"
)

func ParseCredentials(pubCert, privCert string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewServerTLSFromFile(pubCert, privCert)
	if err != nil {
		pubCert, err := base64.StdEncoding.DecodeString(pubCert)
		if err != nil {
			return nil, err
		}
		privCert, err := base64.StdEncoding.DecodeString(privCert)
		if err != nil {
			return nil, err
		}
		creds, err := tls.X509KeyPair(pubCert, privCert)
		if err != nil {
			return nil, err
		}
		return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{creds}}), nil
	}
	return creds, nil
}
