package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

// ---------- HTTP Server

type HTTPServer struct {
	e               *echo.Echo
	address         string
	cert            *Certificate
	shutdownTimeout time.Duration
}

func NewHTTPServer(
	e *echo.Echo,
	address string,
	cert *Certificate,
) HTTPServer {
	return HTTPServer{
		e:               e,
		address:         address,
		cert:            cert,
		shutdownTimeout: 3 * time.Second,
	}
}

func (s HTTPServer) String() string {
	return fmt.Sprintf("server.HTTP(address=%s)", s.address)
}

func (s HTTPServer) Serve(ctx context.Context) error {
	s.e.HideBanner = true
	s.e.HidePort = true
	log.Info().Str("address", s.address).Msg("Starting HTTP server")

	errC := make(chan error, 1)
	go func() {
		if s.cert == nil {
			errC <- s.e.Start(s.address)
		} else {
			errC <- s.e.StartTLS(s.address, s.cert.CertFile, s.cert.KeyFile)
		}
	}()

	select {
	case err := <-errC:
		return errors.Join(suture.ErrTerminateSupervisorTree, err)
	case <-ctx.Done():
		log.Info().Msg("Gracefully shutting down HTTP server...")

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := s.e.Shutdown(ctx); err != nil {
			log.Err(err).Msg("HTTP Server failed to shutdown gracefully")
			return err
		}

		return nil
	}
}

// ---------- HTTP Router

func NewHTTPRouter() *echo.Echo {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Middleware
	e.Use(echoext.LoggerWithConfig(echoext.LoggerConfig{
		Format: []string{
			"remote_ip",
			"host",
			"method",
			"user_agent",
			"status",
			"error",
			"latency_human",
		},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: echoext.RecoverLogErrorFunc,
	}))

	return e
}

// ---------- HTTP Redirect

func NewHTTPRedirect(httpsPort string) *echo.Echo {
	e := echo.New()

	e.Any("*", func(c echo.Context) error {
		r := c.Request()

		host, _ := core.SplitAddress(r.Host)

		http.Redirect(c.Response(), r, "https://"+host+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
		return nil
	})

	return e
}

// ---------- Certificate

type Certificate struct {
	CertFile string
	KeyFile  string
}

func (c Certificate) ForceGenerate() error {
	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2024),
		Subject: pkix.Name{
			OrganizationalUnit: []string{"IPCManView"},
			CommonName:         "IPCManView",
			Country:            []string{"US"},
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(c.CertFile, certPEM.Bytes(), 0600); err != nil {
		return err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err := os.WriteFile(c.KeyFile, certPrivKeyPEM.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
