package main

import (
	"fmt"
	"github.com/rancher/machine/drivers/amazonec2"
	"github.com/rancher/machine/drivers/azure"
	"github.com/rancher/machine/drivers/digitalocean"
	"github.com/rancher/machine/drivers/vmwarevsphere"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/drivers/plugin"
	"github.com/rancher/machine/libmachine/mcnflag"
	"os"
	"path/filepath"
	"strings"
)

type driverFunc = func(string, string) drivers.Driver

func generateWrapper[T drivers.Driver](f func(string, string) T) driverFunc {
	return func(a1 string, a2 string) drivers.Driver {
		return f(a1, a2)
	}
}

var driverMap = map[string]driverFunc{
	"amazonec2":     generateWrapper(amazonec2.NewDriver),
	"azure":         generateWrapper(azure.NewDriver),
	"digitalocean":  generateWrapper(digitalocean.NewDriver),
	"vmwarevsphere": generateWrapper(vmwarevsphere.NewDriver),
}

type driverWrapper struct {
	drivers.Driver
}

func (d *driverWrapper) DriverName() string {
	return fmt.Sprintf("external%s", d.Driver.DriverName())
}

type driverOptions struct {
	drivers.DriverOptions
	name string
}

func (d *driverOptions) String(key string) string {
	if strings.HasPrefix(key, d.name) {
		key = "external" + key
	}
	return d.DriverOptions.String(key)
}

func (d *driverOptions) StringSlice(key string) []string {
	if strings.HasPrefix(key, d.name) {
		key = "external" + key
	}
	return d.DriverOptions.StringSlice(key)
}

func (d *driverOptions) Int(key string) int {
	if strings.HasPrefix(key, d.name) {
		key = "external" + key
	}
	return d.DriverOptions.Int(key)
}

func (d *driverOptions) Bool(key string) bool {
	if strings.HasPrefix(key, d.name) {
		key = "external" + key
	}
	return d.DriverOptions.Bool(key)
}

func (d *driverWrapper) SetConfigFromFlags(opts drivers.DriverOptions) error {
	return d.Driver.SetConfigFromFlags(&driverOptions{opts, d.Driver.DriverName()})
}

func (d *driverWrapper) GetCreateFlags() []mcnflag.Flag {
	flags := d.Driver.GetCreateFlags()
	for i, f := range flags {
		// skip flags that are prefixed with the driver name, as they are just generic flags
		if !strings.HasPrefix(f.String(), d.Driver.DriverName()) {
			continue
		}
		if fl, ok := f.(mcnflag.StringFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
		} else if fl, ok := f.(mcnflag.StringSliceFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
		} else if fl, ok := f.(mcnflag.IntFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
		} else if fl, ok := f.(mcnflag.BoolFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
		} else {
			panic("unknown flag type")
		}
		flags[i] = f
	}
	return flags
}

func main() {
	basename, err := os.Executable()
	if err != nil {
		panic(err)
	}
	s := strings.Split(filepath.Base(basename), "-")
	name := s[len(s)-1]
	if driver, ok := driverMap[strings.TrimPrefix(name, "external")]; ok {
		plugin.RegisterDriver(&driverWrapper{driver("machine", "")})
	} else {
		panic("no driver found for " + name + ".")
	}
}
