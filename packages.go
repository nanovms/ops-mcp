package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/nanovms/ops/cmd"
	api "github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/provider/onprem"
	"github.com/spf13/cobra"
)

func pkgLoad(pkg string) {
	ccmd := &cobra.Command{}
	flags := ccmd.PersistentFlags()

	configFlags := cmd.NewConfigCommandFlags(flags)
	globalFlags := cmd.NewGlobalCommandFlags(flags)
	nightlyFlags := cmd.NewNightlyCommandFlags(flags)
	nanosVersionFlags := cmd.NewNanosVersionCommandFlags(flags)
	buildImageFlags := cmd.NewBuildImageCommandFlags(flags)
	runLocalInstanceFlags := cmd.NewRunLocalInstanceCommandFlags(flags)
	pkgFlags := cmd.NewPkgCommandFlags(flags)

	pkgFlags.Package = pkg

	c := api.NewConfig()

	mergeContainer := cmd.NewMergeConfigContainer(configFlags, globalFlags, nightlyFlags, nanosVersionFlags, buildImageFlags, runLocalInstanceFlags, pkgFlags)
	err := mergeContainer.Merge(c)
	if err != nil {
		log.Println(err.Error())
	}

	packageFolder := filepath.Base(pkgFlags.PackagePath())
	executableName := c.Program
	if strings.Contains(executableName, packageFolder) {
		executableName = filepath.Base(executableName)
	} else {
		executableName = filepath.Join(api.PackageSysRootFolderName, executableName)
	}

	api.ValidateELF(filepath.Join(pkgFlags.PackagePath(), executableName))

	if c.Mounts != nil {
		err = onprem.AddVirtfsShares(c)
		if err != nil {
		}
	}

	if err = api.BuildImageFromPackage(pkgFlags.PackagePath(), *c); err != nil {
		log.Println(err.Error())
	}

	err = cmd.RunLocalInstance(c)
	if err != nil {
		log.Println(err.Error())
	}
}

func loadPackage(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
	pkgLoad(*arguments.Content.Description)
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("loaded %s", *arguments.Content.Description))), nil
}
