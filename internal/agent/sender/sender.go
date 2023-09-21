package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"github.com/jbakhtin/rtagent/pkg/crypto"
	"net/http"
)

type ReportFunction func() string

type sender struct {
	cfg Configer
}

func New(cfg Configer) *sender {
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

	publicKey, err := crypto.GetPublicKey(r.cfg.GetCryptoKey())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var encryptedKey string
	if publicKey != nil {
		buf, encryptedKey, err = crypto.GetEncryptedMessage(publicKey, buf)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	fmt.Println(string(buf))

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Encrypted-Key", encryptedKey)

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
