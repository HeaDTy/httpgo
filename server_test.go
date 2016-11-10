package main

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestServer(t *testing.T) {
  is, err := initServer(true)
  if err != nil {
    t.Fatal(err)
  }
	assert.Equal(t, true, is)
}
