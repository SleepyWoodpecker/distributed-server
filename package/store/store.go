package store

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type PathTransformFunc func(key string)string

type StoreOpts struct {
	PathTransformFunc
}

func DefaultPathTransformFunc(key string) string {
	return key
}

// create a sha-1 hash and store the files like git would
// first 2 characters are the folder name
// other 2 characters are the file name
func CASPathTransformFunc(key string) string {
	data := []byte(key)

	hasher := sha1.New()
	hasher.Write(data)

	hashedData := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedData)

	return hashedString[:2] + "/" + hashedString[2:]
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

// take in a reader stream and write it to the local file system
func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	// create the directory if it does not yet exist
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	// create the file in the local file structure
	fileName := "somepath"
	fullFileName := "CAS/" + pathName + "/" + fileName

	f, err := os.Create(fullFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	
	// write to the file
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("Created a file of %d bytes\n", n)

	return nil
}