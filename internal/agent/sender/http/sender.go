package http

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"github.com/jbakhtin/rtagent/pkg/crypto"
	"net/http"
)

type Configer interface {
	GetServerAddress() string
	GetKeyApp() string
	GetCryptoKey() string
	GetTrustedSubnet() string
}

type ReportFunction func() string

type HTTPSender struct {
	Cfg Configer
}

func New(cfg Configer) (*HTTPSender, error) {
	return &HTTPSender{
		Cfg: cfg,
	}, nil
}

func (s *HTTPSender) Send(key string, value types.Metricer) error {
	endpoint := fmt.Sprintf("%s/update/", fmt.Sprintf("http://%s", s.Cfg.GetServerAddress()))
	model, err := models.ToJSON(s.Cfg, key, value)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(model)
	if err != nil {
		return err
	}

	var encryptedKey string
	var publicKey *rsa.PublicKey

	if s.Cfg.GetCryptoKey() != "" {
		publicKey, err = crypto.GetPublicKey(s.Cfg.GetCryptoKey())
		if err != nil {
			return err
		}

		if publicKey != nil {
			buf, encryptedKey, err = crypto.GetEncryptedMessage(publicKey, buf)
			if err != nil {
				return err
			}
		}

	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	if s.Cfg.GetCryptoKey() != "" {
		request.Header.Set("Encrypted-Key", encryptedKey)
	}

	if s.Cfg.GetTrustedSubnet() != "" {
		request.Header.Set("X-Real-IP", s.Cfg.GetTrustedSubnet())
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if err = response.Body.Close(); err != nil {
		return err
	}

	return nil
}
