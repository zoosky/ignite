package v1alpha1

import (
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ImageSource(obj *ImageSource) {
	obj.Type = ImageSourceTypeDocker
}

func SetDefaults_PoolSpec(obj *PoolSpec) {
	// TODO: These might be nil instead of ignitemeta.EmptySize
	if obj.AllocationSize == ignitemeta.EmptySize {
		obj.AllocationSize = ignitemeta.NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS)
	}

	if len(obj.MetadataPath) == 0 {
		obj.MetadataPath = constants.SNAPSHOTTER_METADATA_PATH
	}

	if len(obj.DataPath) == 0 {
		obj.DataPath = constants.SNAPSHOTTER_DATA_PATH
	}
}

func SetDefaults_VMSpec(obj *VMSpec) {
	if obj.CPUs == 0 {
		obj.CPUs = constants.VM_DEFAULT_CPUS
	}

	// TODO: These might be nil instead of ignitemeta.EmptySize
	if obj.Memory == ignitemeta.EmptySize {
		obj.Memory = ignitemeta.NewSizeFromBytes(constants.VM_DEFAULT_MEMORY)
	}

	if obj.Size == ignitemeta.EmptySize {
		obj.Size = ignitemeta.NewSizeFromBytes(constants.VM_DEFAULT_SIZE)
	}
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.State == "" {
		obj.State = VMStateCreated
	}
}
