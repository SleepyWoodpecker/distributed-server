package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	store := NewStore(opts)

	data := bytes.NewReader([]byte("This is a string"))
	
	assert.Nil(t, store.writeStream("test", data))
}

func TestCASHash(t *testing.T) {
	intialString := "Thisisastring"
	expectedOutput := "f7/2017485fbf6423499baf9b240daa14f5f095a1/Thisisiastring"

	fileData := CASPathTransformFunc(intialString)
	
	assert.Equal(t, fileData.FullPath(), expectedOutput)
}