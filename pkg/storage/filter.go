package storage

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type ErrAmbiguous error
type ErrNonexistent error

// TODO: The filter interface and its filters
type Filter interface {
	// Meta specifies if the Filter can operate on meta.APIType objects
	// returned by a storage's ListMeta. Light weight filters matching
	// names/UIDs should use this as it's much faster.
	Meta() bool
	// Every Object to be filtered is passed though Filter, which should
	// return the Object on match, or nil if it doesn't match
	Filter(meta.Object) (meta.Object, error)
	// ErrAmbiguous specifies what to error if
	// a single request returned multiple matches
	ErrAmbiguous() ErrAmbiguous
	// ErrNonexistent specifies what to error if
	// a single request returned no matches
	ErrNonexistent() ErrNonexistent
}

func getSingle(f Filter, o []meta.Object) (meta.Object, error) {
	var result meta.Object

	for _, object := range o {
		if match, err := f.Filter(object); err != nil { // Filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			if result != nil {
				return nil, f.ErrAmbiguous()
			} else {
				result = match
			}
		}
	}

	if result == nil {
		return nil, f.ErrNonexistent()
	}

	return result, nil
}

func getMultiple(f Filter, o []meta.Object) ([]meta.Object, error) {
	var results []meta.Object

	for _, object := range o {
		if match, err := f.Filter(object); err != nil { // Filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			results = append(results, match)
		}
	}

	return results, nil
}
