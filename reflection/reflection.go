package reflection

import (
	"fmt"
	"reflect"
	"strings"
)

func GetTypeName(t reflect.Type) string {
	buf := strings.Builder{}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		buf.WriteString("*")
	}

	switch t.Kind() {
	case reflect.Slice:
		buf.WriteString(GetSliceName(t))
		break
	case reflect.Map:
		buf.WriteString(GetMapName(t))
		break
	default:
		name := t.PkgPath()
		if name != "" {
			buf.WriteString(strings.Replace(name, "/", ".", -1) + "." + t.Name())
			break
		} else {
			buf.WriteString(t.Name())
			break
		}
	}
	return buf.String()
}

func GetSliceName(t reflect.Type) string {
	elemType := t.Elem()

	name := elemType.PkgPath()
	if name != "" {
		name = strings.Replace(name, "/", ".", -1) + "." + elemType.Name()
		return "[]" + name
	} else {
		return t.String()
	}
}

func GetMapName(t reflect.Type) string {
	keyType := t.Key()
	elemType := t.Elem()

	key := keyType.PkgPath()
	if key != "" {
		key = strings.Replace(key, "/", ".", -1) + "." + keyType.Name()
	}

	name := elemType.PkgPath()
	if name != "" {
		name = strings.Replace(name, "/", ".", -1) + "." + elemType.Name()
		return fmt.Sprintf("map[%s]%s", key, name)
	} else {
		return t.String()
	}
}
