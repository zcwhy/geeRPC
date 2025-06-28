package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRollingWindow(t *testing.T) {
	window := NewRollingWindow(5, 1000)

	window.Add(1, 1)
	curTotal, curSucc := window.Reduce()
	assert.Equal(t, curTotal, 1)
	assert.Equal(t, curSucc, 1)

	window.Add(1, 0)
	curTotal, curSucc = window.Reduce()
	assert.Equal(t, curTotal, 2)
	assert.Equal(t, curSucc, 1)

}
