package util

import (
	"errors"
	"reflect"
)

func SetOrCopySlice(dest, src reflect.Value, set bool) (int, error) {
	destType := dest.Type()
	if destType.Kind() != reflect.Slice {
		return 0, errors.New("Dest Type is not a slice ptr. " + destType.String())
	}
	destElemType := destType.Elem()
	srcType := src.Type()
	if srcType.Kind() != reflect.Slice {
		return 0, errors.New("Src Type is not a slice. " + srcType.String())
	}
	srcElemType := srcType.Elem()

	if destElemType.Kind() == srcElemType.Kind() {
		if set {
			dest.Set(src)
			return dest.Len(), nil
		} else {
			destTmp := dest
			if destTmp.IsNil() {
				destTmp = reflect.MakeSlice(destType, src.Len(), src.Len())
			}
			_ = reflect.Copy(destTmp, src)
			dest.Set(destTmp)
			return dest.Len(), nil
		}
	} else {
		n := 0
		destTmp := dest
		for i := 0; i < src.Len(); i++ {
			ov := src.Index(i)
			ot := ov.Type()
			// interface
			if ot.AssignableTo(destElemType) {
				destTmp = reflect.Append(destTmp, ov)
				n++
				continue
			}
			if ot.ConvertibleTo(destElemType) {
				destTmp = reflect.Append(destTmp, ov.Convert(destElemType))
				n++
				continue
			}
		}
		dest.Set(destTmp)
		return n, nil
	}
}
