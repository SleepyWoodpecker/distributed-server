package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	// test the creation of a new file
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	store := NewStore(opts)

	storedData := "This is a string"
	fileName := "test"
	data := bytes.NewReader([]byte(storedData))
	
	assert.Nil(t, store.writeStream(fileName, data))

	// test the reading of the new file
	readData, err := store.Read(fileName)
	assert.Nil(t, err)

	assert.Equal(t, storedData, string(readData))

	store.deleteFile(fileName)
}

func TestCASHash(t *testing.T) {
	intialString := "Thisisastring"
	expectedOutput := "CAS/ea/2597d38124fbd43edff2816347b425d8666bd1"

	fileData := CASPathTransformFunc(intialString)
	
	assert.Equal(t, expectedOutput, fileData.FullPath())
}