package bean

import (
	"github.com/ydx1011/gopher-core/reflection"
	"reflect"
)

type sliceDefinition struct {
	name string
	o    interface{}
	t    reflect.Type
}

func newSliceDefinition(o interface{}) (Definition, error) {
	t := reflect.TypeOf(o)
	return &sliceDefinition{
		name: reflection.GetSliceName(t),
		o:    o,
		t:    t,
	}, nil
}

func (d *sliceDefinition) Type() reflect.Type {
	return d.t
}

func (d *sliceDefinition) Name() string {
	return d.name
}

func (d *sliceDefinition) Value() reflect.Value {
	return reflect.ValueOf(d.o)
}

func (d *sliceDefinition) Interface() interface{} {
	return d.o
}

func (d *sliceDefinition) IsObject() bool {
	return false
}
func (d *sliceDefinition) AfterSet() error {
	return nil
}

func (d *sliceDefinition) Destroy() error {
	return nil
}

func (d *sliceDefinition) Classify(classifier Classifier) (bool, error) {
	return false, nil
}
