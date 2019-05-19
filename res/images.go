// This file was autogenerated; DO NOT EDIT
package res

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// Images stores decoded Images keyed by their filename as inline compiled resources.
var Images map[string]image.Image

func init() {
	Images = map[string]image.Image{}

	for _, file := range []struct {
		name string
		b    []byte
	}{
		{name: "cursors.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x1, 0x0, 0x51, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0x6e, 0xb0, 0x3d, 0x30, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x1, 0x77, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0x99, 0xb1, 0xe, 0x83, 0x30, 0xc, 0x44, 0x3d, 0xd0, 0x4a, 0x9d, 0x58, 0xc8, 0x62, 0xf5, 0xff, 0xbf, 0xb3, 0x45, 0xa8, 0x81, 0x4, 0x23, 0x44, 0x14, 0xfb, 0xaa, 0xf6, 0x6e, 0xf4, 0xc0, 0x7b, 0x42, 0x80, 0xc2, 0x59, 0x64, 0xc9, 0x33, 0x47, 0x20, 0xd9, 0x60, 0x15, 0xcc, 0x87, 0x18, 0x54, 0xb7, 0x5d, 0xc1, 0xfc, 0x70, 0x3, 0xe3, 0xb1, 0x53, 0x30, 0x3f, 0xd4, 0xe0, 0xe0, 0xb5, 0x8b, 0x33, 0x38, 0x20, 0x8d, 0x82, 0x35, 0x18, 0x23, 0x1f, 0x2, 0x5, 0xf3, 0xd, 0x83, 0x31, 0xfa, 0x43, 0xa0, 0x60, 0x7e, 0x65, 0x0, 0xe0, 0x17, 0x6, 0x10, 0xfe, 0x1b, 0x9b, 0x83, 0xe1, 0xcb, 0x23, 0xa7, 0x9c, 0xdf, 0x73, 0x9c, 0xf9, 0xd3, 0x9a, 0x82, 0x7f, 0x5b, 0x13, 0xc5, 0xdf, 0x1a, 0x6c, 0xf9, 0x9e, 0x6, 0x25, 0x7f, 0x35, 0x28, 0xf9, 0x7e, 0x6, 0x35, 0xff, 0x63, 0x50, 0xf3, 0xbd, 0xc, 0xf6, 0xfc, 0xc5, 0x60, 0xcf, 0xf7, 0x31, 0xb0, 0xf8, 0xb3, 0x81, 0xc5, 0xf7, 0x30, 0xb0, 0xf9, 0xd3, 0x64, 0xf3, 0x1d, 0xc, 0xe, 0x4, 0xd2, 0x81, 0xc0, 0xd0, 0xff, 0x16, 0xd8, 0x7c, 0x91, 0x28, 0xbe, 0x69, 0x90, 0xe6, 0x79, 0x14, 0xdf, 0x30, 0x48, 0xcb, 0x3c, 0x8a, 0xbf, 0x33, 0x48, 0x9f, 0x79, 0x14, 0xbf, 0x32, 0x48, 0xeb, 0x3c, 0x8a, 0x5f, 0x18, 0xa4, 0xed, 0x3c, 0x8a, 0x2f, 0x92, 0x72, 0xca, 0xf9, 0x90, 0x13, 0xfe, 0xc3, 0xc8, 0xc2, 0x80, 0x85, 0x1, 0xb, 0x83, 0x7f, 0x2b, 0xc, 0xce, 0x5f, 0x7b, 0xdf, 0xf, 0xc3, 0xa5, 0x8b, 0x2a, 0x98, 0xef, 0x60, 0x70, 0xf9, 0xa6, 0x2a, 0x98, 0xdf, 0xd9, 0xa0, 0xe9, 0xa1, 0x52, 0x30, 0xbf, 0xa3, 0x41, 0xf3, 0x4b, 0xd5, 0xcb, 0xa0, 0xf9, 0x3a, 0xdd, 0x9a, 0xc, 0x5, 0xf3, 0x1b, 0xd, 0xba, 0x36, 0x39, 0xa, 0xe6, 0x37, 0x18, 0x74, 0x6f, 0xb2, 0x14, 0xcc, 0xbf, 0x68, 0xe0, 0xd2, 0xe4, 0x9d, 0xf7, 0x84, 0x61, 0x4d, 0x22, 0x17, 0x8, 0x5c, 0x20, 0x18, 0x6, 0x5c, 0x20, 0x40, 0x2, 0x5f, 0x20, 0xb0, 0x1f, 0x60, 0x3f, 0xc0, 0x7e, 0x80, 0xfd, 0x0, 0xfb, 0x1, 0xf6, 0x3, 0xec, 0x7, 0xda, 0xc, 0x34, 0x8c, 0xcf, 0x7e, 0x80, 0xfd, 0x0, 0xfb, 0x81, 0xdf, 0xeb, 0x7, 0x14, 0xcc, 0xbf, 0xf2, 0x3d, 0x52, 0x30, 0x9f, 0xe7, 0x1, 0x9e, 0x7, 0xbe, 0xfb, 0x3c, 0x10, 0xc7, 0xe7, 0x79, 0x80, 0xe7, 0x1, 0x9e, 0x7, 0x7e, 0xef, 0x3c, 0xc0, 0x7d, 0xc1, 0x99, 0x1, 0xf7, 0x5, 0x2, 0x35, 0xe0, 0xbe, 0x0, 0x12, 0xf8, 0xbe, 0x80, 0x61, 0x18, 0x86, 0x61, 0x98, 0xbf, 0xc9, 0xb, 0x95, 0x13, 0x35, 0xae, 0x5a, 0xde, 0x34, 0x9b, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "figure.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x51, 0x0, 0x0, 0xe0, 0x5a, 0x40, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xc2, 0x5d, 0x52, 0x38, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x0, 0x86, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0xd3, 0xc1, 0x11, 0x83, 0x30, 0x10, 0x4, 0xc1, 0x7b, 0xeb, 0xd0, 0x5f, 0xf9, 0x47, 0x6a, 0x61, 0xc0, 0x29, 0xac, 0xab, 0xe8, 0x4e, 0x60, 0xa7, 0xe, 0x51, 0x5, 0x0, 0xc0, 0xcf, 0xdc, 0xb2, 0xfb, 0x63, 0xf4, 0x8c, 0x6, 0x1c, 0xe3, 0x98, 0x6f, 0xbe, 0x40, 0xcd, 0xec, 0xfe, 0xd8, 0x7a, 0x4b, 0xee, 0x47, 0xb, 0xf2, 0x1, 0xf7, 0x37, 0xa8, 0x68, 0x40, 0x47, 0x3, 0x2a, 0x1b, 0xf0, 0x9c, 0x20, 0x1b, 0x90, 0xfc, 0xd, 0xcf, 0x80, 0x8e, 0x6, 0xec, 0x84, 0x7d, 0xff, 0xe4, 0xfe, 0x1f, 0x4, 0x54, 0x3e, 0xa0, 0xc2, 0x4f, 0xa0, 0x56, 0xaf, 0x2d, 0x15, 0x70, 0x6e, 0x7f, 0x3, 0x32, 0x5, 0xeb, 0xa, 0x88, 0x9d, 0x60, 0x8d, 0xbe, 0x2, 0xfa, 0xe5, 0x1, 0xb9, 0x27, 0x50, 0xe1, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc8, 0xf9, 0x0, 0x6e, 0x39, 0xa, 0xfd, 0xcc, 0xd4, 0xb1, 0xf6, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "hud.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7b, 0x50, 0x4c, 0x54, 0x45, 0x51, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xd1, 0x6f, 0xc3, 0x5d, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x0, 0x9b, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0xd4, 0x31, 0xe, 0x83, 0x30, 0x10, 0x44, 0xd1, 0xad, 0x56, 0x4a, 0x63, 0xe5, 0x0, 0x56, 0x8a, 0xdc, 0xff, 0x90, 0xc1, 0x40, 0xe0, 0x6, 0x4c, 0xf3, 0x9e, 0xdc, 0xcf, 0x97, 0xb0, 0xa9, 0x3a, 0xcc, 0x59, 0x51, 0xf3, 0x53, 0x33, 0xbc, 0x5f, 0xc9, 0x82, 0x7d, 0x3f, 0x58, 0x70, 0xee, 0xc7, 0xa, 0xae, 0xfd, 0x50, 0xc1, 0xda, 0x7f, 0xef, 0x32, 0x5, 0xe7, 0x7e, 0xf7, 0x51, 0x30, 0x9e, 0xf, 0xa8, 0x15, 0xd0, 0xbd, 0xa, 0x72, 0x1, 0xfd, 0xea, 0xed, 0x24, 0x3, 0x46, 0x6f, 0x27, 0x17, 0xd0, 0xa7, 0x58, 0x40, 0x5d, 0xfb, 0xa9, 0x80, 0xfa, 0xef, 0xc7, 0x2, 0xea, 0xdc, 0x8f, 0x7e, 0x82, 0x5c, 0x40, 0xdf, 0x97, 0x70, 0x7c, 0x9f, 0xff, 0x15, 0x8e, 0x3b, 0xa0, 0x13, 0xfb, 0x5b, 0x41, 0x78, 0xff, 0x2e, 0x48, 0xed, 0xff, 0xb, 0x72, 0xfb, 0x47, 0x41, 0x72, 0x7f, 0x15, 0x64, 0xf7, 0x33, 0xef, 0xf, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x60, 0xf7, 0x3, 0x96, 0x93, 0x9, 0x25, 0x24, 0x98, 0x87, 0xc7, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "terrain.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x51, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xee, 0xa1, 0x5c, 0x1, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x0, 0xb3, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0xd5, 0xc9, 0x11, 0xc3, 0x20, 0xc, 0x0, 0x40, 0x8a, 0x40, 0xea, 0xbf, 0xd4, 0x80, 0x20, 0x49, 0x7, 0xe2, 0xb3, 0x1b, 0xcf, 0x38, 0x3f, 0x5d, 0x58, 0x8c, 0x71, 0x64, 0xc6, 0xfe, 0xad, 0x67, 0x3c, 0xb1, 0x82, 0xe7, 0x3c, 0x39, 0xe4, 0x93, 0xf8, 0x27, 0xf8, 0x7d, 0xbd, 0xa8, 0x7f, 0x7, 0x8f, 0xb9, 0x67, 0x30, 0x1f, 0x4c, 0x61, 0x15, 0x7e, 0xe7, 0x1f, 0xf7, 0x28, 0xf4, 0xc6, 0x3f, 0x81, 0x2b, 0x89, 0x3d, 0x82, 0x1a, 0x44, 0x6b, 0xfd, 0xb7, 0xfa, 0x4a, 0xa1, 0x32, 0x59, 0x7f, 0x66, 0x6b, 0xfd, 0x19, 0x37, 0x89, 0x59, 0x4f, 0x54, 0x22, 0x9d, 0x7, 0x60, 0xde, 0xb3, 0x77, 0xf, 0x40, 0xce, 0xd5, 0x82, 0xec, 0x6d, 0xc1, 0x9, 0x5d, 0x33, 0x38, 0x83, 0x68, 0x3e, 0x4, 0xab, 0x5, 0xd5, 0x87, 0xdd, 0x8d, 0xfa, 0x12, 0xa2, 0x7b, 0xd, 0xfc, 0x3f, 0xc3, 0xa8, 0x14, 0xba, 0x17, 0x41, 0xdc, 0x4d, 0xf0, 0xbd, 0xf, 0xfa, 0x57, 0x61, 0x7c, 0xd7, 0x40, 0x2d, 0x82, 0x47, 0x97, 0x51, 0x9c, 0x65, 0xfc, 0xea, 0x3a, 0xfc, 0xd, 0x20, 0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0xaf, 0xf, 0xcd, 0x67, 0x16, 0xb4, 0xc6, 0x5a, 0xdb, 0x3, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "trees.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x51, 0x0, 0x0, 0xe0, 0x5a, 0x40, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xc2, 0x5d, 0x52, 0x38, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x1, 0x58, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0x96, 0x31, 0x96, 0xc4, 0x30, 0x8, 0x43, 0xdd, 0xc6, 0x66, 0x9d, 0x14, 0x79, 0x3c, 0xdd, 0xff, 0xa4, 0xb, 0xd8, 0x9e, 0x9d, 0xbd, 0x0, 0x34, 0xfa, 0x53, 0x24, 0x9d, 0x64, 0x21, 0x93, 0x69, 0x8d, 0x10, 0xf2, 0x41, 0x55, 0x6b, 0xe5, 0x5f, 0x55, 0xc0, 0x5c, 0x5c, 0xf6, 0x2b, 0xd0, 0x37, 0x79, 0x85, 0xbb, 0xf0, 0x97, 0xec, 0x2c, 0xd4, 0x85, 0x2f, 0xdd, 0x8f, 0x2b, 0xdd, 0x80, 0xbe, 0xaf, 0x86, 0xec, 0xbb, 0x3c, 0x78, 0x12, 0xe9, 0xe9, 0xeb, 0xd1, 0x77, 0x80, 0xe4, 0x1, 0x5c, 0x71, 0xe8, 0x88, 0x1e, 0xa1, 0x7f, 0xa7, 0xf7, 0x2f, 0xea, 0xbf, 0x8c, 0xc0, 0xf4, 0xef, 0xec, 0xa, 0x9a, 0xee, 0xbb, 0xc3, 0xf, 0xfd, 0xec, 0x11, 0x7c, 0xe9, 0x1f, 0x17, 0xa9, 0x6, 0x80, 0xd5, 0x0, 0x7f, 0xae, 0x3d, 0x90, 0xda, 0x82, 0x50, 0xfc, 0x5c, 0x80, 0x5d, 0xc5, 0x54, 0x3, 0x4b, 0xf3, 0x6f, 0xe, 0x11, 0xc5, 0x9d, 0x1a, 0xc0, 0xda, 0x3d, 0x71, 0x7c, 0x9f, 0x86, 0xd7, 0xf0, 0x4e, 0xe, 0x0, 0xff, 0x8f, 0x9f, 0x3b, 0x81, 0x68, 0xfd, 0xf9, 0x18, 0x61, 0x5f, 0xc4, 0xc4, 0x4b, 0x10, 0x8a, 0xe7, 0x63, 0x1c, 0xc7, 0xcf, 0x5e, 0x44, 0xcd, 0x16, 0xcf, 0x5e, 0x2, 0x35, 0xfa, 0x9f, 0x10, 0xd6, 0x16, 0x2e, 0xd0, 0xdf, 0xe, 0xce, 0xf8, 0xb, 0xf4, 0x8d, 0x31, 0xaa, 0xe2, 0x3f, 0xf4, 0x61, 0x26, 0xee, 0x3a, 0x7d, 0x77, 0xd0, 0x87, 0x8c, 0xa7, 0xf0, 0x7f, 0x71, 0xef, 0xcf, 0x23, 0x5a, 0x6a, 0x60, 0x8, 0xa, 0xf5, 0xdb, 0xf3, 0x74, 0x33, 0xa0, 0xa5, 0x25, 0x10, 0x94, 0x8e, 0x60, 0xd8, 0x8, 0x50, 0x5c, 0xc2, 0x62, 0x3, 0x43, 0x26, 0xb4, 0xce, 0x83, 0x1b, 0xb0, 0x16, 0xa2, 0xca, 0x41, 0x1f, 0xb6, 0x6, 0x64, 0x7d, 0x10, 0x4a, 0xf4, 0x45, 0xbc, 0x85, 0x13, 0x55, 0x16, 0xc6, 0x9c, 0x5d, 0xa6, 0x6c, 0x7, 0x15, 0x1, 0x0, 0x5d, 0x7e, 0x96, 0x3, 0xcd, 0x37, 0xa0, 0x98, 0xf0, 0x4d, 0x8c, 0x39, 0xc5, 0x5e, 0xb, 0x12, 0xb0, 0xd8, 0xc7, 0x84, 0x77, 0x70, 0xce, 0xa2, 0x16, 0x7a, 0x1, 0x2b, 0x6f, 0x81, 0xd, 0xa1, 0x29, 0x50, 0x67, 0xe0, 0x6b, 0x16, 0x8d, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x64, 0xf2, 0xb, 0x44, 0x57, 0x2c, 0xc3, 0x96, 0x39, 0x8c, 0xa5, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
	} {
		decoded, err := png.Decode(bytes.NewReader(file.b))
		if err != nil {
			panic(fmt.Sprintf("png.Decode %s: %v", file.name, err))
		}
		Images[file.name] = decoded
	}
}
