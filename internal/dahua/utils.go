package dahua

import (
	"mime"
	"path/filepath"
	"strings"
)

func parseFileExtension(fileName, contentType string) string {
	originalExt := strings.ToLower(filepath.Ext(fileName))

	ext, err := mime.ExtensionsByType(contentType)
	if err != nil || len(ext) == 0 {
		if originalExt == "" {
			return ""
		}

		return originalExt
	}

	for _, e := range ext {
		if e == originalExt {
			return originalExt
		}
	}

	return ext[0]
}
