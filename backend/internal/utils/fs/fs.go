package fs

import (
	"bytes"
	"io"
	"os"
)

func WriteJSONToFile(filePath string, data []byte) error {
    // Create the file, truncating it if it already exists
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the JSON data to the file
    _, err = io.Copy(file, bytes.NewReader(data))
    return err
}

func ReadJSONFromFile(filePath string) ([]byte, error) {
    // Open the file for reading
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Read the JSON data from the file
    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    return data, nil
}

func OpenFile(filePath string) (*os.File, error) {
    return os.Open(filePath)
}

func CreateFile(filePath string) (*os.File, error) {
    return os.Create(filePath)
}
