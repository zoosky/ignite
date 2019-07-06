package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"

	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type ImportKernelFlags struct {
	Source string
	Name   string
}

type importKernelOptions struct {
	*ImportKernelFlags
}

func (i *ImportKernelFlags) NewImportKernelOptions() (*importKernelOptions, error) {
	return &importKernelOptions{ImportKernelFlags: i}, nil
}

func ImportKernel(ao *importKernelOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not a kernel image: %s", ao.Source)
	}

	// TODO: Kernel importing from docker when moving to pool/snapshotter
	kernel := &api.Kernel{
		Spec: api.KernelSpec{
			Version: "unknown",
			Source: api.ImageSource{
				Type: "file",
				ID:   "-",
				Name: "-",
			},
		},
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(ao.Name, meta.KindKernel)
	if err != nil {
		return err
	}

	// Create new kernel metadata
	md, err := kernmd.NewKernel("", &name, kernel)
	if err != nil {
		return err
	}
	defer metadata.Cleanup(md, false) // TODO: Handle silent

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	return metadata.Success(md)
}
