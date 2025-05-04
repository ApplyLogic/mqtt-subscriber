package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

// Helper function to create TLS configuration
func newTLSConfig(certPath string) (*tls.Config, error) {
	// Load CA certificate
	//certPool := x509.NewCertPool()
	//pemCerts, err := ioutil.ReadFile(certPath)
	//if err != nil {
	//	return nil, err
	//}

	//if !certPool.AppendCertsFromPEM(pemCerts) {
	//	return nil, fmt.Errorf("failed to parse certificate")
	//}

	// Create TLS configuration
	//return &tls.Config{
	//	RootCAs:            certPool,
	//	InsecureSkipVerify: false,
	//	MinVersion:         tls.VersionTLS12,
	//}, nil

	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair("/Users/golanshay/Workspace/tls-certs/server.crt", "/Users/golanshay/Workspace/tls-certs/server.key")
	if err != nil {
		log.Fatalf("failed to load certificates: %v", err)
	}

	caCert, err := os.ReadFile("/Users/golanshay/Workspace/ls-certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to load ca cert: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert, //Require Client Certificate

		ClientCAs:          caCertPool,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}, nil

}
