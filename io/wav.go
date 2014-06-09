package io

import (
	"github.com/r9y9/go-dsp/wav"
	"os"
)

func ReadWav(filename string) (*wav.Wav, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	w, werr := wav.ReadWav(file)
	if werr != nil {
		return nil, werr
	}

	return w, nil
}
