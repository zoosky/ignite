package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

type CreateOptions struct {
	Image  *imgmd.ImageMetadata
	Kernel *kernmd.KernelMetadata
	vm     *vmmd.VMMetadata
	Name   string
	CPUs   int64
	Memory int64
}

func Create(co *CreateOptions) error {
	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	// Create a new name for the VM if none is given
	util.NewName(&co.Name)

	// Create new metadata for the VM and add to createOptions for further processing
	// This enables the generated VM metadata to pass straight to start and attach via run
	co.vm = vmmd.NewVMMetadata(vmID, co.Name, vmmd.NewVMObjectData(co.Image.ID, co.Kernel.ID, co.CPUs, co.Memory))

	// Save the metadata
	if err := co.vm.Save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := co.vm.CopyImage(); err != nil {
		return err
	}

	// Print the ID of the created VM
	fmt.Println(co.vm.ID)

	return nil
}