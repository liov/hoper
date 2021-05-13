package errorsi

import "fmt"

// MultiMapError stores multiple decoding errors.
//
// Borrowed from the App Engine SDK.
type MultiMapError map[string]error

func (e MultiMapError) Error() string {
	s := ""
	for _, err := range e {
		s = err.Error()
		break
	}
	switch len(e) {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, len(e)-1)
}

func (e MultiMapError) Merge(errors MultiMapError) {
	for key, err := range errors {
		if e[key] == nil {
			e[key] = err
		}
	}
}

type MultiSliceError []error

func (e MultiSliceError) Error() string {
	s := ""
	for _, err := range e {
		s = err.Error()
		break
	}
	switch len(e) {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, len(e)-1)
}
