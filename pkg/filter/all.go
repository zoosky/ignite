package filter

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

// The AllFilter matches anything it's given
type AllFilter struct{}

// It's more efficient for this to be an ObjectFilter, as it loads everything anyways
var _ filterer.ObjectFilter = &AllFilter{}

func NewAllFilter() *AllFilter {
	return &AllFilter{}
}

func (f *AllFilter) Filter(object meta.Object) (meta.Object, error) {
	return object, nil
}

// The AllFilter shouldn't be used to match single Objects
func (f *AllFilter) ErrAmbiguous() filterer.ErrAmbiguous {
	return fmt.Errorf("ambiguous query: AllFilter used to match single Object")
}

func (f *AllFilter) ErrNonexistent() filterer.ErrNonexistent {
	return fmt.Errorf("no results: AllFilter used to match single Object")
}
