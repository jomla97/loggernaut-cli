package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/viper"
)

type ApiErrorResponse struct {
	Error string `json:"error"`
}

// SendAll sends all log files in the outbox folder to the Loggernaut API, returning the number of logs sent
func SendAll(debug bool) (int, error) {
	files, err := walkOutbox()
	if err != nil {
		return 0, err
	}
	var sent int
	for _, file := range files {
		err := Send(file, debug)
		if err != nil {
			return sent, err
		}
		sent++
	}
	return sent, nil
}

// Send sends the log file at the specified path to the Loggernaut API
func Send(logPath string, debug bool) error {
	// Open the source file
	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	// Copy the file data to a buffer
	var log bytes.Buffer
	_, err = io.Copy(&log, file)
	if err != nil {
		return fmt.Errorf("failed to copy data from log file: %w", err)
	}

	// Get the meta data
	meta, err := ReadMetaFile(logPath)
	if err != nil {
		return err
	}

	// Create a new form writer
	reqBody := new(bytes.Buffer)
	form := multipart.NewWriter(reqBody)

	// Create a list of form files
	formFiles := []struct {
		Fieldname string
		Filename  string
		Data      []byte
	}{
		{"log", filepath.Base(logPath), log.Bytes()},
		{"meta", filepath.Base(*meta.MetaPath), meta.Bytes()},
	}

	// Add the form files to the form writer
	for _, f := range formFiles {
		// Create a form part for the file
		part, err := form.CreateFormFile(f.Fieldname, f.Filename)
		if err != nil {
			return fmt.Errorf("failed to create form part: %v", err.Error())
		}

		// Write the file data to the form part
		if _, err := part.Write(f.Data); err != nil {
			return fmt.Errorf("failed to write to form part: %w", err)
		}
	}

	// Close the form writer
	form.Close()

	// Create http request
	req, err := http.NewRequest("POST", viper.GetString("api_url")+"/ingest", reqBody)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Add("User-Agent", "loggernaut-cli/v"+config.Version)
	req.Header.Add("Content-Type", form.FormDataContentType())
	req.Header.Add("Accept", "application/json")

	// Dump the request if debugging is enabled
	if debug {
		if b, err := httputil.DumpRequest(req, true); err != nil {
			fmt.Println("Error dumping request:", err)
		} else {
			fmt.Println(string(b))
		}
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Dump the response if debugging is enabled
	if debug {
		if b, err := httputil.DumpResponse(resp, true); err != nil {
			fmt.Println("Error dumping request:", err)
		} else {
			fmt.Println(string(b))
		}
	}

	// Check the response
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var respBody ApiErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err == nil {
			return fmt.Errorf("unexpected response status: %v, error: %v", resp.Status, respBody.Error)
		}
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	// Remove the log file and meta file
	for _, file := range []string{logPath, *meta.MetaPath} {
		err := os.Remove(file)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove file: %w", err)
		}
	}

	return nil
}

// walkOutbox walks the outbox folder and returns a list of log files
func walkOutbox() ([]string, error) {
	var paths []string
	err := filepath.Walk(config.OutboxPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".log" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
