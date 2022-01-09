package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "ca"
	app.Version = "1.0"
	app.Usage = "a pseudo certificate authority for testing purposes"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "create CA certs and keys if they do not exist",
			Action: initCA,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "local directory where certificates and keys are stored",
					Value:  "fixtures/certs",
					EnvVar: "CA_CERT_DIRECTORY",
				},
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "overwrite keys even if they already exist",
				},
				cli.StringFlag{
					Name:  "o, organization",
					Usage: "name of organization to issue certificates for",
				},
				cli.StringFlag{
					Name:  "C, country",
					Usage: "country of the organization",
				},
				cli.StringFlag{
					Name:  "p, province",
					Usage: "province or state of the organization",
				},
				cli.StringFlag{
					Name:  "l, locality",
					Usage: "locality or city of the organization",
				},
				cli.StringFlag{
					Name:  "a, address",
					Usage: "streed address of the organization",
				},
				cli.StringFlag{
					Name:  "P, postcode",
					Usage: "postal code of the organization",
				},
			},
		},
		{
			Name:   "issue",
			Usage:  "issue a certificate signed by the CA",
			Action: issue,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "local directory where certificates and keys are stored",
					Value:  "fixtures/certs",
					EnvVar: "CA_CERT_DIRECTORY",
				},
				cli.StringFlag{
					Name:  "o, organization",
					Usage: "name of organization to issue certificates for",
				},
				cli.StringFlag{
					Name:  "C, country",
					Usage: "country of the organization",
				},
				cli.StringFlag{
					Name:  "p, province",
					Usage: "province or state of the organization",
				},
				cli.StringFlag{
					Name:  "l, locality",
					Usage: "locality or city of the organization",
				},
				cli.StringFlag{
					Name:  "a, address",
					Usage: "streed address of the organization",
				},
				cli.StringFlag{
					Name:  "P, postcode",
					Usage: "postal code of the organization",
				},
			},
		},
	}

	app.Run(os.Args)
}

func initCA(c *cli.Context) (err error) {
	force := c.Bool("force")
	certPath := filepath.Join(c.String("certs"), "ca.crt")
	keyPath := filepath.Join(c.String("certs"), "ca.key")

	if !force {
		if _, err = os.Stat(certPath); err == nil {
			return cli.NewExitError("certificate file already exists", 1)
		}
		if _, err = os.Stat(keyPath); err == nil {
			return cli.NewExitError("private key file already exists", 1)
		}
	}

	// Create a certificate
	// TODO: create a method to issue the serial number
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1942),
		Subject: pkix.Name{
			Organization:  []string{c.String("organization")},
			Country:       []string{c.String("country")},
			Province:      []string{c.String("province")},
			Locality:      []string{c.String("locality")},
			StreetAddress: []string{c.String("address")},
			PostalCode:    []string{c.String("postcode")},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create private key
	priv, _ := rsa.GenerateKey(rand.Reader, 4096)
	pub := &priv.PublicKey
	signed, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("create ca failed: %s", err), 1)
	}

	// Save the key to a file
	var cf, kf *os.File
	if cf, err = os.Create(certPath); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer cf.Close()
	if err = pem.Encode(cf, &pem.Block{Type: "CERTIFICATE", Bytes: signed}); err != nil {
		return cli.NewExitError(err, 1)
	}

	if kf, err = os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer kf.Close()
	if err = pem.Encode(kf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func issue(c *cli.Context) (err error) {

	if c.String("organization") == "" {
		return cli.NewExitError("specify the name of the organization", 1)
	}

	// Load the CA key pairs
	caPath := filepath.Join(c.String("certs"), "ca.crt")
	keyPath := filepath.Join(c.String("certs"), "ca.key")

	var (
		catls tls.Certificate
		ca    *x509.Certificate
	)

	if catls, err = tls.LoadX509KeyPair(caPath, keyPath); err != nil {
		return cli.NewExitError(err, 1)
	}

	if ca, err = x509.ParseCertificate(catls.Certificate[0]); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Prepare the certificate
	// TODO: how to handle serial numbers?
	// TODO: how to handle subject key ID?
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1945),
		Subject: pkix.Name{
			Organization:  []string{c.String("organization")},
			Country:       []string{c.String("country")},
			Province:      []string{c.String("province")},
			Locality:      []string{c.String("locality")},
			StreetAddress: []string{c.String("address")},
			PostalCode:    []string{c.String("postcode")},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 7),
		SubjectKeyId: []byte{1, 2, 3, 4, 5, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 4096)
	pub := &priv.PublicKey

	// Sign the certificate
	var signed []byte
	if signed, err = x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Write out the certificate to disk
	name := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(c.String("organization")), " ", "_"))

	var cf, kf *os.File
	if cf, err = os.Create(filepath.Join(c.String("certs"), name+".crt")); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer cf.Close()
	if err = pem.Encode(cf, &pem.Block{Type: "CERTIFICATE", Bytes: signed}); err != nil {
		return cli.NewExitError(err, 1)
	}

	if kf, err = os.OpenFile(filepath.Join(c.String("certs"), name+".key"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer kf.Close()
	if err = pem.Encode(kf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
