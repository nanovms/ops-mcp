package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/cobra"

	"github.com/nanovms/ops/cmd"
	"github.com/nanovms/ops/lepton"
	api "github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/provider"
	"github.com/nanovms/ops/provider/onprem"
	"github.com/nanovms/ops/types"
)

type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}

type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

type InstanceArguments struct {
	ImageName string `json:"longitude" jsonschema:"required,description=The image name of the instance to create"`
}

func getProviderAndContext(c *types.Config, providerName string) (api.Provider, *api.Context, error) {
	p, err := provider.CloudProvider(providerName, &c.CloudConfig)
	if err != nil {
		return nil, nil, err
	}

	ctx := api.NewContext(c)

	return p, ctx, nil
}

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

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	// actually this is foregrounded.
	err := server.RegisterTool("pkg_load", "Load package", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		pkgLoad(*arguments.Content.Description)
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("loaded %s", *arguments.Content.Description))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("instance_create", "Instance create", func(arguments InstanceArguments) (*mcp_golang.ToolResponse, error) {

		c := lepton.NewConfig()
		log.Println(arguments.ImageName)

		log.Println(arguments)

		c.CloudConfig.ImageName = arguments.ImageName

		p, ctx, err := getProviderAndContext(c, "onprem")
		if err != nil {
			log.Println(err)
		}

		c.RunConfig.Kernel = lepton.GetOpsHome() + "/0.1.53-arm/kernel.img"

		err = p.CreateInstance(ctx)
		if err != nil {
			log.Println(err)
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("spun that instance up")), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("list_instances", "List instances", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {

		c := &types.Config{}

		p, ctx, err := getProviderAndContext(c, "onprem")
		if err != nil {
			log.Println(err)
		}
		ctx.Config().RunConfig.JSON = true

		instances, err := p.GetInstances(ctx)
		if err != nil {
			log.Println(err)
		}

		imgs, err := json.Marshal(instances)
		if err != nil {
			log.Println(err)
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(imgs))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("list_images", "List images", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {

		c := &types.Config{}

		p, ctx, err := getProviderAndContext(c, "onprem")
		if err != nil {
			log.Println(err)
		}
		ctx.Config().RunConfig.JSON = true

		images, err := p.GetImages(ctx, "")
		if err != nil {
			log.Println(err)
		}

		imgs, err := json.Marshal(images)
		if err != nil {
			log.Println(err)
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(imgs))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
