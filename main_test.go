package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDriver(t *testing.T) {
	driver, ok := driverMap["amazonec2"]
	assert.True(t, ok)

	d := driver("machine", "")

	assert.Equal(t, "amazonec2", d.DriverName())
	assert.Equal(t, "external-amazonec2", (&driverWrapper{d}).DriverName())
}

func TestGetCreateFlags(t *testing.T) {
	driver, ok := driverMap["digitalocean"]
	assert.True(t, ok)

	d := &driverWrapper{driver("machine", "")}
	flags := d.GetCreateFlags()
	for _, f := range flags {
		if strings.Contains(f.String(), "digitalocean") {
			assert.True(t, strings.HasPrefix(f.String(), "externaldigitalocean"))
		} else {
			assert.True(t, !strings.Contains(f.String(), "digitalocean"))
		}
	}
	assert.True(t, len(flags) > 0)
}
