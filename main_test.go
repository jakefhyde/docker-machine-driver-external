package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDriver(t *testing.T) {
	driver, ok := driverMap["amazonec2"]
	assert.True(t, ok)

	d := driver("machine", "")

	assert.Equal(t, "amazonec2", d.DriverName())
	assert.Equal(t, "external-amazonec2", (&driverWrapper{d}).DriverName())
}
