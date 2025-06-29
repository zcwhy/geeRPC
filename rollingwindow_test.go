package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRollingWindow(t *testing.T) {
	window := NewRollingWindow(5, 50)

	window.Add(1, 1)
	curTotal, curSucc := window.Reduce()
	assert.Equal(t, 1, curTotal)
	assert.Equal(t, 1, curSucc)

	window.Add(1, 0)
	curTotal, curSucc = window.Reduce()
	assert.Equal(t, 2, curTotal)
	assert.Equal(t, 1, curSucc)

	time.Sleep(time.Millisecond * 50)
	window.Add(1, 1)
	curTotal, curSucc = window.Reduce()
	assert.Equal(t, 3, curTotal)
	assert.Equal(t, 2, curSucc)

	time.Sleep(time.Millisecond * 300)
	window.Add(1, 1)
	curTotal, curSucc = window.Reduce()
	assert.Equal(t, 1, curTotal)
	assert.Equal(t, 1, curSucc)
}

func TestNSpan(t *testing.T) {
	window := NewRollingWindow(5, 50)

	window.Add(1, 1)
	assert.Equal(t, 0, window.span())
	assert.Equal(t, 0, window.offset)

	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, 1, window.span())

	window.Add(1, 1)
	assert.Equal(t, 1, window.offset)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, 2, window.span())

	window.Add(1, 1)
	assert.Equal(t, 3, window.offset)

	time.Sleep(time.Millisecond * 300)
	assert.Equal(t, 5, window.span())

	window.Add(1, 1)
	assert.Equal(t, 3, window.offset)
}
