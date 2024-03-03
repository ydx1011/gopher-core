package util

import (
	"errors"
	"fmt"
	"reflect"
)

func SetOrCopyMap(dest, src reflect.Value, set bool) (int, error) {
	destType := dest.Type()
	if destType.Kind() != reflect.Map {
		return 0, errors.New("Dest Type is not a map ptr. " + destType.String())
	}
	destKeyType := destType.Key()
	destElemType := destType.Elem()
	srcType := src.Type()
	if srcType.Kind() != reflect.Map {
		return 0, errors.New("Src Type is not a map. " + srcType.String())
	}
	srcKeyType := srcType.Key()
	srcElemType := srcType.Elem()

	if destKeyType != srcKeyType {
		return 0, fmt.Errorf("Expect Map key type: %s but get %s. ", destKeyType.String(), srcKeyType.String())
	}

	if set && (destElemType.Kind() == srcElemType.Kind()) {
		dest.Set(src)
		return dest.Len(), nil
	} else {
		n := 0
		destTmp := dest
		if destTmp.IsNil() {
			destTmp = reflect.MakeMapWithSize(destType, src.Len())
		}
		keys := src.MapKeys()
		for _, key := range keys {
			ov := src.MapIndex(key)
			ot := ov.Type()
			if ot.AssignableTo(destElemType) {
				destTmp.SetMapIndex(key, ov)
				n++
				continue
			}
			if ot.ConvertibleTo(destElemType) {
				destTmp.SetMapIndex(key, ov.Convert(destElemType))
				n++
				continue
			}

		}
		dest.Set(destTmp)
		return n, nil
	}
}
