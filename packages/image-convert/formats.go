// formats.go contains format extensions. These extensions are implementations
// of the extension struct defined in main.go. Consequently, these formats
// also implement the "ext" interface.

package main

import (
	"image"
	"image/png"
	"os"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

// WebP

type WEBP extension

func (e WEBP) decode(file *os.File) (image.Image, error) {
	if e.decOpts == nil {
		e.decOpts = &decoder.Options{}
	}

	img, err := webp.Decode(file, e.decOpts)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (e WEBP) encode(file *os.File, src image.Image) error {
	var err error

	// Use default options if no options provided.
	if e.encOpts == nil {
		e.encOpts, err = encoder.NewLossyEncoderOptions(encoder.PresetDefault, 80)
	}
	if err != nil {
		return err
	}

	err = webp.Encode(file, src, e.encOpts)
	if err != nil {
		return err
	}

	return nil
}

// PNG

type PNG extension

func (e PNG) decode(file *os.File) (image.Image, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (e PNG) encode(file *os.File, src image.Image) error { return nil }

// JPG
