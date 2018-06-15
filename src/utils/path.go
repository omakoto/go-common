package utils

import (
	"github.com/omakoto/go-common/src/common"
	"strings"
)

func HomeExpanded(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = common.MustGetHome() + "/" + path[2:]
	}
	return path
}
