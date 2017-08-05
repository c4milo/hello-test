package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"

	"os"

	pb "github.com/c4milo/hello-nyt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	os.Setenv("TLS_CERT", `
-----BEGIN CERTIFICATE-----
MIIDUzCCAtmgAwIBAgIJAKTf/aVGhWkYMAkGByqGSM49BAEwgZExCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhOZXcgWW9yazERMA8GA1UEBxMITmV3IFlvcmsxFzAVBgNV
BAoTDkhvb2tsaWZ0LCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmluZzEKMAgGA1UE
AxQBKjEhMB8GCSqGSIb3DQEJARYSY2FtaWxvQGhvb2tsaWZ0LmlvMCAXDTE2MTEx
NDE0NTgwNFoYDzIxMTUwNjA5MTQ1ODA0WjCBkTELMAkGA1UEBhMCVVMxETAPBgNV
BAgTCE5ldyBZb3JrMREwDwYDVQQHEwhOZXcgWW9yazEXMBUGA1UEChMOSG9va2xp
ZnQsIEluYy4xFDASBgNVBAsTC0VuZ2luZWVyaW5nMQowCAYDVQQDFAEqMSEwHwYJ
KoZIhvcNAQkBFhJjYW1pbG9AaG9va2xpZnQuaW8wdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAASH3bmfhqPNDE2YdeBG15Yl13GVWlex0QDCh85koZ3kbKMGdDBqgb5gqgwZ
F1rCCpjff+o3D3JaAMYosACOyHn8lnJOcpryqUkwCklxSQqleLJM4EGSitMm8119
tzYhaCajgfkwgfYwHQYDVR0OBBYEFMNqnVpZOU6jIqWaiHr7AnMXpBwWMIHGBgNV
HSMEgb4wgbuAFMNqnVpZOU6jIqWaiHr7AnMXpBwWoYGXpIGUMIGRMQswCQYDVQQG
EwJVUzERMA8GA1UECBMITmV3IFlvcmsxETAPBgNVBAcTCE5ldyBZb3JrMRcwFQYD
VQQKEw5Ib29rbGlmdCwgSW5jLjEUMBIGA1UECxMLRW5naW5lZXJpbmcxCjAIBgNV
BAMUASoxITAfBgkqhkiG9w0BCQEWEmNhbWlsb0Bob29rbGlmdC5pb4IJAKTf/aVG
hWkYMAwGA1UdEwQFMAMBAf8wCQYHKoZIzj0EAQNpADBmAjEAnvDrqcg7Sl2wK/bH
+98IMGMiYdT1FpSqCT3YyVQeCPELlxmnXbzNesY/R+l8oY9bAjEAhya4BL+ingli
o9FuJqdUS5o9Rgii55nFhNdzQvT/p/ANGHBCfQyUNtAjPp92KvXC
-----END CERTIFICATE-----
`)

	os.Setenv("TLS_KEY", `
-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAD5WLfofxxT8EewIU/VYJ5hXRWAyjxwAhemboTVnnrCmqA5Icxz+oa
kVluFU7LiPWgBwYFK4EEACKhZANiAASH3bmfhqPNDE2YdeBG15Yl13GVWlex0QDC
h85koZ3kbKMGdDBqgb5gqgwZF1rCCpjff+o3D3JaAMYosACOyHn8lnJOcpryqUkw
CklxSQqleLJM4EGSitMm8119tzYhaCY=
-----END EC PRIVATE KEY-----
`)

}

func main() {
	// Set up a connection to the server.
	tlsKeyPair, err := tls.X509KeyPair([]byte(os.Getenv("TLS_CERT")), []byte(os.Getenv("TLS_KEY")))
	if err != nil {
		glog.Fatalf("failed loading TLS certificate and key: %+v", err)
	}

	x509Cert, err := x509.ParseCertificate(tlsKeyPair.Certificate[0])
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(x509Cert)

	creds := credentials.NewClientTLSFromCert(certPool, "")

	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	r, err := c.Counts(context.Background(), new(empty.Empty))
	if err != nil {
		log.Fatalf("💩: %#v", err)
	}
	log.Printf("Counts: %s", r.Counts)
}
