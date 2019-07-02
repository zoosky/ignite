package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

// This package represents kernel objects, which reside in /var/lib/firecracker/kernel/{id}/metadata.json
type Kernel struct {
	*v1alpha1.Kernel
	device *dm.Device
	resize *Resize
}

// Get the metadata filename for the image
func (k *Kernel) MetadataPath() string {
	// TODO: This
	return ""
}

func (k *Kernel) GetImage() *Image {
	k.ss.pool.GetDevice(k.device.Parent)
}
