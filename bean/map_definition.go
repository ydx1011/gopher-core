package bean

import (
	"github.com/ydx1011/gopher-core/reflection"
	"reflect"
)

type mapDefinition struct {
	name string
	o    interface{}
	t    reflect.Type
}

func newMapDefinition(o interface{}) (Definition, error) {
	t := reflect.TypeOf(o)
	return &mapDefinition{
		name: reflection.GetMapName(t),
		o:    o,
		t:    t,
	}, nil
}

func (d *mapDefinition) Type() reflect.Type {
	return d.t
}

func (d *mapDefinition) Name() string {
	return d.name
}

func (d *mapDefinition) Value() reflect.Value {
	return reflect.ValueOf(d.o)
}

func (d *mapDefinition) Interface() interface{} {
	return d.o
}

func (d *mapDefinition) IsObject() bool {
	return false
}

func (d *mapDefinition) AfterSet() error {
	return nil
}

func (d *mapDefinition) Destroy() error {
	return nil
}

func (d *mapDefinition) Classify(classifier Classifier) (bool, error) {
	return false, nil
}
