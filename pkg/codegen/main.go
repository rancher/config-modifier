package main

import (
	"os"

	v1 "github.com/rancher/config-modifier/pkg/apis/some.api.group/v1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/config-modifier/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"some.api.group": {
				Types: []interface{}{
					v1.Foo{},
				},
				GenerateTypes: true,
			},
		},
	})
}