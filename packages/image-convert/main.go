package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
)

// Add terminal colors!

func main() {

	var err error

	// Create man/help page. Specify format integers.

	// format 1 and format 2.
	var from, to *int

	from = flag.Int("from", 0, "format to convert from")
	to = flag.Int("to", 1, "format to convert to")

	// Get target input and output dir from cmd flags, or perhaps env vars.

	var source, target *string

	source = flag.String("source", "", "source directory")
	target = flag.String("target", "", "source directory")

	flag.Parse()

	// Verify that the source format is not the same as the target format.
	if *from == *to {
		log.Fatal(errors.New("source format is the same as target format"))
	}

	if from == nil || to == nil {
		log.Fatal(errors.New("missing format flags"))
	}

	conv := &converter{
		in:  *source,
		out: *target,
	}

	p := PNG{}
	w := WEBP{}

	err = conv.from(p).to(w).convert()
	if err != nil {
		log.Fatal(err)
	}

}

// various extensions
const (
	valJPG int = iota
	valWEBP
)

// ext defines methods that can operate on a file extension. It
// lets a type that implements this interface the ability to be
// used in iterative situations.
type ext interface {
	encode(file *os.File, src image.Image) error
	decode(file *os.File) (image.Image, error)
}

type extension struct {
	// name is the file extension name.
	name string

	// options for specific encodings.
	encOpts *encoder.Options
	decOpts *decoder.Options

	ext
}

// converter takes directories as inputs and outputs.
// It then converts files from one file type to another
// given two different extension struct pointers. These
// are parameters.
type converter struct {

	// These are input and output directories.
	// The in property can also be a single file.
	in, out string

	// These are source and target file extensions.
	source, target *extension
}

// isSet checks if either the source or target has been set.
// This method enforces the fluent pattern.
func (c *converter) isSet() error {
	if c.source == nil && c.target == nil {
		return fmt.Errorf("source and target not set")
	}

	if c.source == nil {
		return fmt.Errorf("source not set")
	}

	return nil
}

// to the extension to convert to.
func (c *converter) to(e ext) *converter {
	c.target.ext = e

	return c
}

// from sets the extension to convert from.
func (c *converter) from(e ext) *converter {
	c.source.ext = e

	return c
}

// convert takes a list of files paths and converts
// them to the provided source and targets extensions.
func (c *converter) convert(paths ...string) error {
	err := c.isSet()
	if err != nil {
		return err
	}

	if len(paths) < 1 {
		return errors.New("no paths provided")
	}

	for _, p := range paths {
		var outputName string

		// fix this.
		outputName += c.out + path.Base(p) + "." + c.target.name

		file, err := os.Open(p)
		if err != nil {
			return err
		}

		img, err := c.source.decode(file)
		if err != nil {
			return err
		}

		output, err := os.Create(path.Join(c.out, outputName))
		if err != nil {
			return err
		}

		defer output.Close()

		err = c.target.encode(file, img)
		if err != nil {
			return err
		}
	}

	return nil
}
