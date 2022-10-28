package fs

import (
	"github.com/actliboy/hoper/server/go/lib/utils/errors/multierr"
	"os"
)

func RangeDir(dir string, callback func(path string, entry os.DirEntry) error) error {
	entities, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	errs := &multierr.MultiError{}
	for _, e := range entities {
		path := dir + PathSeparator + e.Name()
		if e.IsDir() {
			err = RangeDir(path, callback)
			if err != nil {
				errs.Append(err)
			}
		}
		err = callback(path, e)
		if err != nil {
			errs.Append(err)
		}
	}
	return errs
}
