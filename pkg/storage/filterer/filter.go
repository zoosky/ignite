package filterer

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type ErrAmbiguous error
type ErrNonexistent error

// BaseFilter provides shared functionality for filter types
type BaseFilter interface {
	// ErrAmbiguous specifies what to error if
	// a single request returned multiple matches
	ErrAmbiguous() ErrAmbiguous
	// ErrNonexistent specifies what to error if
	// a single request returned no matches
	ErrNonexistent() ErrNonexistent
}

// ObjectFilter implementations filter fully loaded meta.Objects
type ObjectFilter interface {
	BaseFilter
	// Every Object to be filtered is passed though Filter, which should
	// return the Object on match, or nil if it doesn't match
	Filter(meta.Object) (meta.Object, error)
}

// MetaFilter implementations operate on meta.APIType objects,
// which are more light weight, but provide only name/UID matching.
type MetaFilter interface {
	BaseFilter
	// Every Object to be filtered is passed though FilterMeta, which should
	// return the Object on match, or nil if it doesn't match. The Objects
	// given to FilterMeta are of type meta.APIType, stripped of other contents.
	FilterMeta(meta.Object) (meta.Object, error)
}
