package rfvclient

import (
	"os"
	"path/filepath"
	"strings"
)

// ThumbCacheDir 环境变量 RB_THUMB_CACHE，默认 .cache/rb_thumbs。
func ThumbCacheDir() string {
	s := strings.TrimSpace(os.Getenv("RB_THUMB_CACHE"))
	if s == "" {
		return ".cache/rb_thumbs"
	}
	return s
}

func thumbCachePath(hash string) string {
	return filepath.Join(ThumbCacheDir(), hash+".webp")
}

// LoadThumbCache 读取本地缩略图缓存。
func LoadThumbCache(hash string) ([]byte, bool) {
	if hash == "" {
		return nil, false
	}
	b, err := os.ReadFile(thumbCachePath(hash))
	if err != nil {
		return nil, false
	}
	return b, true
}

// StoreThumbCache 将缩略图写入本地缓存（已存在则跳过）。
func StoreThumbCache(hash string, data []byte) error {
	if hash == "" || len(data) == 0 {
		return nil
	}
	p := thumbCachePath(hash)
	if _, err := os.Stat(p); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}
