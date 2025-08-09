package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const ROOT = "CAS"
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
		FileName: hashedString[2:],
		FolderName: ROOT + "/" + hashedString[:2],
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

func (s *Store) Read(key string) ([]byte, error) {
	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	
	if err != nil {
		return nil, err
	}

	fmt.Println(buf)

	return buf.Bytes(), nil
}

// take in a reader stream and write it to the local file system
func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	// create the directory if it does not yet exist
	if err := os.MkdirAll(pathName.FolderName, os.ModePerm); err != nil {
		return err
	}

	// create the file in the local file structure
	fullFileName := pathName.FullPath()

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

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathName := s.PathTransformFunc(key)
	return os.Open(pathName.FullPath())
}

// delete the file by the file name
// panics if the file does not exist
func (s *Store) deleteFile(key string) error {
	pathName := s.PathTransformFunc(key)

	_, err := os.Stat(pathName.FullPath())
	
	if errors.Is(err, os.ErrNotExist) {
		errorMessage := fmt.Sprintf("File %s does not exist", pathName.FullPath())
		panic(errorMessage)
	} else if err != nil {
		return err
	}

	os.Remove(pathName.FullPath())
	fileComponents := strings.Split(pathName.FullPath(), "/")

	return s.removeAllEmptyParentFolders(strings.Join(fileComponents[:len(fileComponents)-1], "/"))
}

// remove all empty directories
// stop only when the CAS folder is reached
func (s *Store) removeAllEmptyParentFolders(currentFolderPath string) error {
	if strings.HasSuffix(currentFolderPath, "CAS") {
		return nil
	}

	fileInfo, err := os.Stat(currentFolderPath)

	if errors.Is(err, os.ErrNotExist) {
		errorMessage := fmt.Sprintf("File %s does not exist", currentFolderPath)
		panic(errorMessage)
	} else if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return nil
	}

	fileCount, err := numFiles(currentFolderPath)
	
	if err != nil || fileCount != 0{
		return err
	}

	os.Remove(currentFolderPath)

	fileNames := strings.Split(currentFolderPath, "/")
	parntFolderPath := strings.Join(fileNames[:len(fileNames) - 1], "/")

	return s.removeAllEmptyParentFolders(parntFolderPath)
}

func numFiles(pathName string) (int, error) {
	dirEntries, err := os.ReadDir(pathName)

	if err != nil {
		return 0, err
	}

	fileCount := 0
	for range dirEntries {
		fileCount++
	}

	return fileCount, nil
}

func (s *Store) Has(key string) bool {
	pathName := CASPathTransformFunc(key)

	_, err := os.Stat(pathName.FullPath())

	return err != nil
}