package main

import (
	"fmt"
	"github.com/rancher/machine/drivers/amazonec2"
	"github.com/rancher/machine/drivers/azure"
	"github.com/rancher/machine/drivers/digitalocean"
	"github.com/rancher/machine/drivers/vmwarevsphere"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/drivers/plugin"
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

func main() {
	basename, err := os.Executable()
	if err != nil {

	}
	s := strings.Split(filepath.Base(basename), "-")
	name := s[len(s)-1]
	if driver, ok := driverMap[strings.TrimPrefix(name, "external")]; ok {
		plugin.RegisterDriver(&driverWrapper{driver("machine", "")})
	} else {
		panic("no driver found for " + name + ".")
	}
}
