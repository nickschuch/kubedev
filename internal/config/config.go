package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func LoadFromFile(path string) (*File, error) {
	var file *File

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return file, errors.Wrap(err, "failed to load file")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return file, errors.Wrap(err, "failed to read file")
	}

	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return file, errors.Wrap(err, "failed to marshal file")
	}

	return file, nil
}
