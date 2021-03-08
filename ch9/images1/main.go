package images1

import (
	"image"
	"sync"
)

var (
	icons map[string]image.Image
	mu    sync.RWMutex
)

func loadIcon(name string) image.Image {
	var x image.Image
	return x
}

func loadIcons() {
	icons = map[string]image.Image{
		"spades.png":   loadIcon("spades.png"),
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

// Concurrency-safe
func Icon(name string) image.Image {
	mu.RLock()
	if icons != nil {
		icon := icons[name]
		mu.RUnlock()
		return icon
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	if icons == nil { // NOTE: we must recheck for nil
		loadIcons()
	}
	return icons[name]
}
