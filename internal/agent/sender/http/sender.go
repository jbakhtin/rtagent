package http

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	sender2 "github.com/jbakhtin/rtagent/internal/agent/sender"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"github.com/jbakhtin/rtagent/pkg/crypto"
	"net/http"
)

type ReportFunction func() string

type sender struct {
	cfg sender2.Configer
}

func New(cfg sender2.Configer) *sender {
	return &sender{
		cfg: cfg,
	}
}

func (r *sender) Send(key string, value types.Metricer) error {
	endpoint := fmt.Sprintf("%s/update/", fmt.Sprintf("http://%s", r.cfg.GetServerAddress()))
	model, err := models.ToJSON(r.cfg, key, value)
	if err != nil {
		return err
	}

	model.Hash, err = model.CalcHash([]byte(r.cfg.GetKeyApp()))
	if err != nil {
		return err
	}

	hash, err := model.CalcHash([]byte(r.cfg.GetKeyApp()))
	if err != nil {
		return err
	}
	model.Hash = hash

	buf, err := json.Marshal(model)
	if err != nil {
		return err
	}

	var encryptedKey string
	var publicKey *rsa.PublicKey
	if r.cfg.GetCryptoKey() != "" {
		publicKey, err = crypto.GetPublicKey(r.cfg.GetCryptoKey())
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

	if r.cfg.GetCryptoKey() != "" {
		request.Header.Set("Encrypted-Key", encryptedKey)
	}

	if r.cfg.GetTrustedSubnet() != "" {
		request.Header.Set("X-Real-IP", r.cfg.GetTrustedSubnet())
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
