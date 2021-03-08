package images2

import (
	"image"
	"sync"
)

var (
	icons         map[string]image.Image
	loadIconsOnce sync.Once
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
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}
