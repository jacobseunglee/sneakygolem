package protocol

import (
	"io"
	"os"

	"github.com/mr-tron/base58"
)

func OpenFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// func ReadFromFile(file *os.File, length int) (string, error) {
// 	buffer := make([]byte, length)
// 	_, err := file.Read(buffer)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(buffer), nil
// }

func AppendBytesToFile(data []byte, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

// ReadFileBase58 reads the entire file, encodes it in base58, and returns a buffer to read from gradually.
func ReadFileBase58(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	encoded := base58.Encode(data)
	
	return encoded, nil
}
