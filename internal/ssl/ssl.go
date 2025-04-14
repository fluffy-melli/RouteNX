package ssl

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
)

type SSL struct {
	Email        string
	Domain       []string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *SSL) GetEmail() string {
	return u.Email
}

func (u *SSL) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *SSL) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func NewSSL(domain []string, email string, testing bool) (*SSL, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	ssl := &SSL{
		Domain: domain,
		Email:  email,
		key:    privateKey,
	}

	config := lego.NewConfig(ssl)
	if testing {
		config.CADirURL = lego.LEDirectoryStaging
	} else {
		config.CADirURL = lego.LEDirectoryProduction
	}
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	ssl.Registration = reg

	request := certificate.ObtainRequest{
		Domains: domain,
		Bundle:  true,
	}

	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	_ = os.WriteFile("cert.pem", cert.Certificate, 0600)
	_ = os.WriteFile("cert.key", cert.PrivateKey, 0600)
	_ = os.WriteFile("cert.bundle.pem", cert.IssuerCertificate, 0600)
	return ssl, nil
}

func (ssl *SSL) Renew() error {
	client, err := lego.NewClient(lego.NewConfig(ssl))
	if err != nil {
		return err
	}

	request := certificate.ObtainRequest{
		Domains: ssl.Domain,
		Bundle:  true,
	}

	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	err = os.WriteFile("cert.pem", cert.Certificate, 0600)
	if err != nil {
		return err
	}

	err = os.WriteFile("cert.key", cert.PrivateKey, 0600)
	if err != nil {
		return err
	}

	err = os.WriteFile("cert.bundle.pem", cert.IssuerCertificate, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (ssl *SSL) ApplyToGin(router *gin.Engine, addr string) error {
	return router.RunTLS(addr, "cert.pem", "cert.key")
}
