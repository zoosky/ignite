package run

import (
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/storage/filterer"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/source"
)

type importOptions struct {
	source   string
	newImage *imgmd.Image
}

func NewImportOptions(source string) (*importOptions, error) {
	return &importOptions{source: source}, nil
}

func Import(bo *importOptions) error {
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(bo.source)
	if err != nil {
		return err
	}

	image := &api.Image{
		Spec: api.ImageSpec{
			Source: *src,
		},
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(bo.source, meta.KindImage)
	if err != nil {
		return err
	}

	// Create new image metadata
	if bo.newImage, err = imgmd.NewImage("", &name, image); err != nil {
		return err
	}
	defer metadata.Cleanup(bo.newImage, false) // TODO: Handle silent

	log.Println("Starting image import...")

	// Create new file to host the filesystem and format it
	if err := bo.newImage.AllocateAndFormat(); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.newImage.AddFiles(dockerSource); err != nil {
		return err
	}

	if err := bo.newImage.Save(); err != nil {
		return err
	}
	log.Printf("Created imported a %s filesystem", image.Spec.Source.Size.HR())

	// If the kernel already exists, don't try to import something with the same name
	if _, err := client.Kernels().Find(filter.NewNameFilter(name)); err == nil {
		return metadata.Success(bo.newImage)
	} else {
		switch err.(type) {
		case filterer.ErrAmbiguous, filterer.ErrNonexistent:
			// With the NameFilter, both of these indicate success
		default:
			return err
		}
	}

	// Import a new kernel from the image if specified
	tmpKernelDir, err := bo.newImage.ExportKernel()
	if err == nil {
		io, err := (&ImportKernelFlags{
			Source: path.Join(tmpKernelDir, constants.KERNEL_FILE),
			Name:   name,
		}).NewImportKernelOptions()
		if err != nil {
			return err
		}

		if err := ImportKernel(io); err != nil {
			return err
		}

		if err := os.RemoveAll(tmpKernelDir); err != nil {
			return err
		}

		//log.Printf("A kernel was imported from the image with name %q and ID %q", name.String(), kernelID)
	} else {
		// Tolerate the kernel to not be found
		if _, ok := err.(*imgmd.KernelNotFoundError); !ok {
			return err
		}
	}

	return metadata.Success(bo.newImage)
}
