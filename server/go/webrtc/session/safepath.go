package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveUnderRoot 将请求路径限制在 root 下，拒绝符号链接逃逸。
func ResolveUnderRoot(root, req string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	rootAbs = filepath.Clean(rootAbs)
	path := strings.TrimSpace(req)
	if path == "" {
		return rootAbs, nil
	}
	cand := filepath.Clean(path)
	if !filepath.IsAbs(cand) {
		cand = filepath.Join(rootAbs, cand)
	}
	cand, err = filepath.Abs(cand)
	if err != nil {
		return "", err
	}
	if cand != rootAbs && !strings.HasPrefix(cand, rootAbs+string(os.PathSeparator)) {
		return "", fmt.Errorf("session: path outside root")
	}
	fi, err := os.Lstat(cand)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("session: symlink not allowed")
	}
	return cand, nil
}
