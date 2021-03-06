package fileutils

import (
	"github.com/omakoto/go-common/src/common"
	"path/filepath"
)

func MustGetRealPath(path string) string {
	ret, err := filepath.EvalSymlinks(path)
	common.Check(err, "EvalSymlinks() failed")
	return ret
}

func SamePath(path1, path2 string) bool {
	return MustGetRealPath(path1) == MustGetRealPath(path2)
}
