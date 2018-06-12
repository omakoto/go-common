package utils

import (
	"strings"
	"github.com/omakoto/go-common/src/common"
)

func HomeExpanded(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = common.MustGetHome() + "/" + path[2:]
	}
	return path
}
