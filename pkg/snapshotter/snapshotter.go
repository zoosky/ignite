package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/util"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
)

// A snapshotter Object represents the internal storage format of a device
// that binds the device mapper device to the object metadata
type Object struct {
	device *dm.Device
	object ignitemeta.Object
	parent *Object
}

// The metadata is behind a private field, which enables this getter to load it on-demand
func (o *Object) GetMetaObject() (ignitemeta.Object, error) {
	if o.object == nil && len(o.device.MetadataPath) > 0 {
		if err := scheme.DecodeFileInto(o.device.MetadataPath, o.object); err != nil {
			return nil, err
		}
	}

	return o.object, nil
}

// Snapshotter abstracts the device mapper pool and provides convenience methods
// It's also responsible for (de)serializing the pool
type Snapshotter struct {
	pool    *dm.Pool
	objects []*Object
}

// NewSnapshotter creates a new Snapshotter with a new Pool
// or loads an existing configuration if it exists
// TODO: No support for physical backing devices for now
func NewSnapshotter() (*Snapshotter, error) {
	p := path.Join(constants.SNAPSHOTTER_DIR, constants.METADATA)
	s := &Snapshotter{}

	// If the metadata doesn't exist, return a new Snapshotter
	if !util.FileExists(p) {
		pool := &v1alpha1.Pool{}
		v1alpha1.SetObjectDefaults_Pool(pool)
		s.pool = dm.NewPool(pool)
		return s, nil
	}

	// Load the pool configuration
	if err := scheme.DecodeFileInto(p, s.pool); err != nil {
		return nil, err
	}

	// Create objects from each device
	poolSize := s.pool.Size()
	resolved := make([]int, 0, poolSize)
	s.objects = make([]*Object, poolSize)
	_ = s.pool.ForDevices(func(id ignitemeta.DMID, device *dm.Device) error {
		for _, i := range resolved {
			if i == id.Index() {
				// Already resolved
				return nil
			}
		}

		s.objects = append(s.objects, &Object{
			device: device,
		})

		return nil
	})

	return s, nil
}
