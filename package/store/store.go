package store

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)
type FullPathname struct {
	FolderName string
	FileName string
}

func (f *FullPathname) FullPath() string {
	return fmt.Sprintf("%s/%s", f.FolderName, f.FileName)
}

type PathTransformFunc func(key string) FullPathname

type StoreOpts struct {
	PathTransformFunc
}

func DefaultPathTransformFunc(key string) FullPathname {
	return FullPathname {
		FolderName: key,
		FileName: key,
	}
}

// create a sha-1 hash and store the files like git would
// first 2 characters are the folder name
// other 2 characters are the file name
// TODO: find some way to hash the filename based on the file type or smth
// Files need to have unique filenames else they are cooked
func CASPathTransformFunc(key string) FullPathname {
	data := []byte(key)

	hasher := sha1.New()
	hasher.Write(data)

	hashedData := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedData)

	return FullPathname{
		FileName: key,
		FolderName: hashedString,
	}
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
	if err := os.MkdirAll(pathName.FolderName, os.ModePerm); err != nil {
		return err
	}

	// create the file in the local file structure
	fullFileName := "CAS/" + pathName.FullPath()

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