package ods

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func Read(filepath string) (ODS, error) {
	// Open the ODS file as a zip archive
	files, err := os.Open(filepath)
	if err != nil {
		return ODS{}, fmt.Errorf("error opening ODS file: %v", err)
	}
	defer func(files *os.File) {
		err := files.Close()
		if err != nil {
			panic(err)
		}
	}(files)

	// Get the file stats
	fileInfo, err := files.Stat()
	if err != nil {
		return ODS{}, fmt.Errorf("error getting file info: %v", err)
	}

	return ReadFrom(files, fileInfo.Size())
}

// ReadFrom Read the data from io reader of provided size.
func ReadFrom(reader io.Reader, size int64) (ODS, error) {
	buf := new(bytes.Buffer)
	data := ODS{}

	// Copy data from the reader to the buffer
	_, err := io.Copy(buf, reader)
	if err != nil {
		return data, fmt.Errorf("error copying data: %v", err)
	}

	files, err := zip.NewReader(bytes.NewReader(buf.Bytes()), size)
	if err != nil {
		return data, fmt.Errorf("error creating zip archive: %v", err)
	}

	// Iterate through the files within the zip archive
	for _, file := range files.File {
		// Open the file for processing
		rc, err := file.Open()
		if err != nil {
			return data, fmt.Errorf("error opening file: %v", err)
		}
		defer func(rc io.ReadCloser) {
			err := rc.Close()
			if err != nil {
				panic(err)
			}
		}(rc)

		// Decode the file contents into the ODS struct
		switch file.Name {
		case "content.xml":
			if err = xml.NewDecoder(rc).Decode(&data.Content); err != nil {
				return ODS{}, fmt.Errorf("error decoding file: %v", err)
			}
		}
	}
	// Return the decoded content, the zip reader (for potential further use), and any error
	return data, nil
}
