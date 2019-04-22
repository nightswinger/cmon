package cmon

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
)

type CSV struct {
	Options *Options

	file   *os.File
	reader *csv.Reader
	writer *csv.Writer
}

type Options struct {
	Headers bool
}

type Option func(*Options)

func Headers(b bool) Option {
	return func(o *Options) { o.Headers = b }
}

func OpenCSVFile(name string, options ...Option) (*CSV, error) {
	opt := &Options{}
	for _, o := range options {
		o(opt)
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	w := csv.NewWriter(f)

	csv := &CSV{
		Options: opt,
		file:    f,
		reader:  r,
		writer:  w,
	}

	return csv, err
}

func (c *CSV) Read() ([][]string, []string, error) {
	data, err := c.reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if c.Options.Headers {
		header := data[0]
		content := data[1:]

		return content, header, err
	}

	return data, nil, err
}

func (c *CSV) Write(record []string) error {
	err := c.writer.Write(record)
	if err != nil {
		return err
	}
	c.writer.Flush()
	if err = c.writer.Error(); err != nil {
		return err
	}

	return err
}

func ToCSV(record []string) (string, error) {
	buf := new(bytes.Buffer)

	w := csv.NewWriter(buf)
	err := w.Write(record)
	if err != nil {
		return "", err
	}

	w.Flush()

	return buf.String(), err
}

func ParseCSV(record string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(record))
	csvArray, err := r.Read()
	if err != nil {
		return nil, err
	}
	return csvArray, err
}
