package storage

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// TODO: The filter interface and its filters
type Filter interface {
	Filter(meta.Object) meta.Object
}
