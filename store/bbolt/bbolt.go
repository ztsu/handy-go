package bbolt

import (
	"go.etcd.io/bbolt"
	"os"
)

func Open(path string, mode os.FileMode, options *bbolt.Options) (*bbolt.DB, error) {
	return bbolt.Open(path, mode, options)
}
