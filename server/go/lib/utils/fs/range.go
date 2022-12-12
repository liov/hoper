package fs

import (
	"github.com/liov/hoper/server/go/lib/utils/errors/multierr"
	"os"
)

func RangeDir(dir string, callback func(dir string, entry os.DirEntry) error) error {
	entities, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	errs := &multierr.MultiError{}
	for _, e := range entities {
		if e.IsDir() {
			err = RangeDir(dir+PathSeparator+e.Name(), callback)
			if err != nil {
				errs.Append(err)
			}
		}
		err = callback(dir, e)
		if err != nil {
			errs.Append(err)
		}
	}
	return errs
}
