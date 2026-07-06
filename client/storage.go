package client

import (
	"encoding/json"
	"os"

	"github.com/gofrs/flock"
)

const logFileName = "logs/apilogs.json"

func LogApiCall(responseBody string, statusCode int, username string) error {
	os.MkdirAll("logs", 0755)

	fileLock := flock.New(logFileName + ".lock")
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Close()

	if err := fileLock.Lock(); err != nil {
		return err
	}
	defer fileLock.Unlock()

	data, err := os.ReadFile(logFileName)
	var logData []interface{}
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, &logData)
	}

	logEntry := map[string]interface{}{
		"username":      username,
		"status_code":   statusCode,
		"response_body": responseBody,
	}
	logData = append(logData, logEntry)
	out, _ := json.Marshal(logData)
	return os.WriteFile(logFileName, out, 0644)
}
