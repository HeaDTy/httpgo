package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestServer(t *testing.T) {
  err := initServer()
  if err != nil {
    t.Fatal(err)
  }
	assert.Equal(t, true, err)
}
