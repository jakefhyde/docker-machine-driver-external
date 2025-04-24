package main

import (
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDriver(t *testing.T) {
	driver, ok := driverMap["amazonec2"]
	assert.True(t, ok)

	d := driver("", "")

	assert.Equal(t, "amazonec2", d.Driver.DriverName())
	assert.Equal(t, "external-amazonec2", d.DriverName())
}

func TestGetCreateFlags(t *testing.T) {
	driver, ok := driverMap["digitalocean"]
	assert.True(t, ok)

	d := driver("", "")
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

func TestGetSSHKeyPath(t *testing.T) {
	driver, ok := driverMap["digitalocean"]
	assert.True(t, ok)

	d := driver("", "")
	var base drivers.Driver = d
	assert.NotEmpty(t, base.GetSSHKeyPath())

}
