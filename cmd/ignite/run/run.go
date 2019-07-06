package run

import (
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(args []string) (*runOptions, error) {
	// parse the args and the config file
	err := rf.CreateFlags.parseArgsAndConfig(args)
	if err != nil {
		return nil, err
	}

	imageName := rf.VM.Spec.Image.Ref

	// Logic to import the image if it doesn't exist
	if _, err := client.Images().Find(filter.NewIDNameFilter(imageName)); err != nil { // TODO: Use this match in create?
		switch err.(type) {
		case filterer.ErrNonexistent:
			io, err := NewImportOptions(imageName)
			if err != nil {
				return nil, err
			}

			if err := Import(io); err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	co, err := rf.NewCreateOptions(args)
	if err != nil {
		return nil, err
	}

	so := &startOptions{
		StartFlags: rf.StartFlags,
		attachOptions: &attachOptions{
			checkRunning: false,
		},
	}

	return &runOptions{co, so}, nil
}

func Run(ro *runOptions) error {
	if err := Create(ro.createOptions); err != nil {
		return err
	}

	// Copy the pointer over for Start
	ro.vm = ro.newVM

	if err := Start(ro.startOptions); err != nil {
		return err
	}

	return nil
}
