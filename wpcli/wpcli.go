package wpcli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Dependency is the key used in the build plan by this buildpack
const Dependency = "wp-cli"

// Contributor is responsibile for deciding what this buildpack will contribute during build
type Contributor struct {
	layer layers.DependencyLayer
}

// NewContributor will create a new Contributor object
func NewContributor(context build.Build) (c Contributor, willContribute bool, err error) {
	_, wantLayer := context.BuildPlan[Dependency]
	if !wantLayer {
		return Contributor{}, false, nil
	}

	deps, err := context.Buildpack.Dependencies()
	if err != nil {
		return Contributor{}, false, err
	}

	version, err := context.Buildpack.DefaultVersion(Dependency)
	if err != nil {
		return Contributor{}, false, err
	}

	dep, err := deps.Best(Dependency, version, context.Stack)
	if err != nil {
		return Contributor{}, false, err
	}

	contributor := Contributor{
		layer: context.Layers.DependencyLayer(dep),
	}

	return contributor, true, nil
}

// Contribute will install wp-cli
func (c Contributor) Contribute() error {
	return c.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Installing to %s", layer.Root)
		if err := helper.CopyFile(artifact, filepath.Join(layer.Root, artifact)); err != nil {
			return err
		}

		if err := writeWrapperScript(layer, "wp"); err != nil {
			return err
		}
		return nil
	}, c.flags()...)
}

func (c Contributor) flags() []layers.Flag {
	return []layers.Flag{layers.Cache, layers.Launch}
}

func writeWrapperScript(layer layers.DependencyLayer, file string, args ...interface{}) error {
	layer.Touch()
	layer.Logger.SubsequentLine("Writing wrapper script bin/%s", file)

	bpBinPath := filepath.Join(layer.Root, "bin")
	script := wrapperScript(bpBinPath)

	if err := os.MkdirAll(bpBinPath, 0755); err != nil {
		return err
	}

	if err := layer.AppendPathSharedEnv("PATH", bpBinPath); err != nil {
		return err
	}

	f := filepath.Join(bpBinPath, file)

	return ioutil.WriteFile(f, []byte(fmt.Sprintf(script, args...)), 0755)
}

func wrapperScript(bpBinPath string) string {
	return fmt.Sprintf(`#!/bin/bash

DEPDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd $DEPDIR

export PATH=$PATH:/layers/org.cloudfoundry.php/php-binary/bin
export PATH=$PATH:%s

php wp-cli-*.phar --path=$HOME/htdocs "$@"
`, bpBinPath)
}
