package common

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
)

// Credit: https://github.com/gin-contrib/static/issues/19

type embedFileSystem struct {
	http.FileSystem
}

func (e *embedFileSystem) Exists(prefix string, filepath string) bool {
	p := strings.TrimPrefix(filepath, prefix)
	// 确保 stripped 路径是相对于 embed 子文件系统的
	if p == filepath {
		return false
	}
	_, err := e.Open(p)
	if err != nil {
		return false
	}
	return true
}

func (e *embedFileSystem) Open(name string) (http.File, error) {
	if name == "/" {
		// This will make sure the index page goes to NoRouter handler,
		// which will use the replaced index bytes with analytic codes.
		return nil, os.ErrNotExist
	}
	return e.FileSystem.Open(name)
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return &embedFileSystem{
		FileSystem: http.FS(efs),
	}
}
