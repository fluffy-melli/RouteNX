package logs

import (
	"encoding/csv"
	"os"
	"time"
)

func CsvCreate(filename string, headers []string) error {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	if _, err := os.Stat(filename); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return err
	}

	return nil
}

func CsvAppend(record []string) error {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	filename := "logs/" + time.Now().Format("2006-01-02") + ".csv"

	CsvCreate(filename, []string{
		"time", "forward-ip", "origin-ip", "host", "url",
	})

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}
