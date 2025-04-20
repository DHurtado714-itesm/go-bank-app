package csvwriter

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVWriter struct {
	headers []string
	rows    [][]string
}

func NewCSVWriter(headers []string) *CSVWriter {
	return &CSVWriter{
		headers: headers,
		rows:    [][]string{},
	}
}

func (w *CSVWriter) AddRow(row []string) error {
	if len(row) != len(w.headers) {
		return fmt.Errorf("row length (%d) does not match header length (%d)", len(row), len(w.headers))
	}
	w.rows = append(w.rows, row)
	return nil
}

func (w *CSVWriter) WriteToFile(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write(w.headers); err != nil {
		return err
	}

	for _, row := range w.rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
