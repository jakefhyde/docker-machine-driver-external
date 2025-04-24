package main

import (
	"fmt"
	"github.com/rancher/machine/drivers/amazonec2"
	"github.com/rancher/machine/drivers/azure"
	"github.com/rancher/machine/drivers/digitalocean"
	"github.com/rancher/machine/drivers/vmwarevsphere"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/drivers/plugin"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/mcnflag"
	"os"
	"path/filepath"
	"strings"
)

type driverWrapper struct {
	drivers.Driver
}

type genFunc = func(string, string) *driverWrapper

func genWrapper[T drivers.Driver](f func(string, string) T) genFunc {
	return func(a1 string, a2 string) *driverWrapper {
		d := f(a1, a2)
		return &driverWrapper{d}
	}
}

var driverMap = map[string]genFunc{
	"amazonec2":     genWrapper(amazonec2.NewDriver),
	"azure":         genWrapper(azure.NewDriver),
	"digitalocean":  genWrapper(digitalocean.NewDriver),
	"vmwarevsphere": genWrapper(vmwarevsphere.NewDriver),
}

func (d *driverWrapper) DriverName() string {
	return fmt.Sprintf("external%s", d.Driver.DriverName())
}

func (d *driverWrapper) GetSSHKeyPath() string {
	p := d.Driver.GetSSHKeyPath()
	log.Infof("Creating ssh key path: %s", p)
	if err := os.MkdirAll(filepath.Dir(p), 0750); err != nil {
		panic(fmt.Errorf("cannot create the folder to store the SSH private key. %s", err))
	}
	return p
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
			f = fl
		} else if fl, ok := f.(mcnflag.StringSliceFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
			f = fl
		} else if fl, ok := f.(mcnflag.IntFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
			f = fl
		} else if fl, ok := f.(mcnflag.BoolFlag); ok {
			fl.Name = "external" + fl.Name
			fl.EnvVar = "EXTERNAL" + fl.EnvVar
			f = fl
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
		d := driver("", "")
		plugin.RegisterDriver(d)
	} else {
		panic("no driver found for " + name + ".")
	}
}
