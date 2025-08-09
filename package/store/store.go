package store

import (
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
	fullFileName := pathName + "/" + fileName

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