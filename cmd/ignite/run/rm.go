package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	"github.com/weaveworks/ignite/pkg/metadata"

	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmFlags struct {
	Force bool
}

type rmOptions struct {
	*RmFlags
	vms []*vmmd.VM
}

func (rf *RmFlags) NewRmOptions(vmMatches []string) (*rmOptions, error) {
	ro := &rmOptions{RmFlags: rf}

	for _, match := range vmMatches {
		if vm, err := client.VMs().Find(filter.NewVMFilter(match)); err == nil {
			ro.vms = append(ro.vms, &vmmd.VM{vm})
		} else {
			return nil, err
		}
	}

	return ro, nil
}

func Rm(ro *rmOptions) error {
	for _, vm := range ro.vms {
		// Check if the VM is running
		if vm.Running() {
			// If force is set, kill the vm
			if ro.Force {
				if err := Stop(&stopOptions{
					&StopFlags{
						Kill: true,
					},
					[]*vmmd.VM{vm},
					true,
				}); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%v is running", vm)
			}
		}

		if err := metadata.Remove(vm, logs.Quiet); err != nil {
			return err
		}
	}

	return nil
}
