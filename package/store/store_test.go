package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	store := NewStore(opts)

	data := bytes.NewReader([]byte("This is a string"))
	
	assert.Nil(t, store.writeStream("test", data))
}