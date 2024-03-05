package errors

import "strings"

type Errors []error

func (es Errors) Empty() bool {
	return len(es) == 0
}

func (es *Errors) AddError(e error) *Errors {
	*es = append(*es, e)
	return es
}

func (es Errors) Error() string {
	buf := strings.Builder{}
	for i := range es {
		buf.WriteString(es[i].Error())
		if i < len(es)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}
