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
		switch f.(type) {
		case *mcnflag.StringFlag:
			f.(*mcnflag.StringFlag).Name = "external" + f.(*mcnflag.StringFlag).Name
			f.(*mcnflag.StringFlag).EnvVar = "EXTERNAL" + f.(*mcnflag.StringFlag).EnvVar
		case *mcnflag.StringSliceFlag:
			f.(*mcnflag.StringSliceFlag).Name = "external" + f.(*mcnflag.StringSliceFlag).Name
			f.(*mcnflag.StringSliceFlag).EnvVar = "EXTERNAL" + f.(*mcnflag.StringSliceFlag).EnvVar
		case *mcnflag.IntFlag:
			f.(*mcnflag.IntFlag).Name = "external" + f.(*mcnflag.IntFlag).Name
			f.(*mcnflag.IntFlag).EnvVar = "EXTERNAL" + f.(*mcnflag.IntFlag).EnvVar
		case *mcnflag.BoolFlag:
			f.(*mcnflag.BoolFlag).Name = "external" + f.(*mcnflag.BoolFlag).Name
			f.(*mcnflag.BoolFlag).EnvVar = "EXTERNAL" + f.(*mcnflag.BoolFlag).EnvVar
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
