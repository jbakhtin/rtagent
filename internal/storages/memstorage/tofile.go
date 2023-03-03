package memstorage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
)

type Writer interface {
	Write(metric *models.Metric)
	Close() error
}

type toFile struct {
	file   *os.File // файл для записи
	writer *bufio.Writer
}

func NewWriter(cfg config.Config) (*toFile, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE, 0777)
	if os.IsNotExist(err) {
		// TODO: файл создается только если присутсвует указанная директория, можно ли как то по дургому?
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
		}
	}

	return &toFile{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (tf *toFile) WriteList(event *map[string]models.Metric) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := tf.writer.Write(data); err != nil {
		return err
	}

	// Удаляем содержимое файла перед перезаписью
	// TODO: Не удалось решить данную проблему дргим способом, например ...|os.O_TRUNC
	err = tf.file.Truncate(0)
	if err != nil {
		return err
	}

	if _, err = tf.file.Seek(0, 0); err != nil {
		return err
	}

	if err = tf.writer.WriteByte('\n'); err != nil {
		return err
	}

	return tf.writer.Flush()
}

func (tf *toFile) Close() error {
	return tf.file.Close()
}
