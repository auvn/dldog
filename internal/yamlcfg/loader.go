package yamlcfg

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func LoadFile(file string, dest any) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.WithMessage(err, "os.Open")
	}

	defer f.Close()

	var m map[string]interface{}
	yamldec := yaml.NewDecoder(f)
	if err := yamldec.Decode(&m); err != nil {
		return errors.WithMessage(err, "decode yaml")
	}

	dec, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result:      dest,
			ErrorUnused: true,
			ErrorUnset:  false,
		})
	if err != nil {
		return errors.WithMessage(err, "new decoder")
	}
	return dec.Decode(m)
}
