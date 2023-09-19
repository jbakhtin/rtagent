package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"net/http"
	"sync"
)

type ReportFunction func() string

type sender struct {
	sync.RWMutex
	cfg Configer
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

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

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