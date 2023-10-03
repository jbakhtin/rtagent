package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-faster/errors"
	"os"
)

func (cb *Builder) WithAllFromJSONFile() *Builder {
	var err error
	defer func() {
		if cb.err != nil {
			cb.err = errors.Wrap(cb.err, "-config flag")
		}
	}()

	jsonFile := flag.String("config", _config, _configLabel)
	flag.Parse()

	fileInfo, err := os.Stat(*jsonFile)
	if err != nil {
		cb.err = err
		return cb
	}

	if fileInfo.IsDir() {
		cb.err = errors.New(fmt.Sprintf("%v is dirdirectory", *jsonFile))
		return cb
	}

	buf, err := os.ReadFile(*jsonFile)
	if err != nil {
		cb.err = err
		return cb
	}

	err = json.Unmarshal(buf, &cb.config)
	if err != nil {
		cb.err = err
		return cb
	}

	return cb
}
