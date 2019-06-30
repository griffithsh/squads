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
		{name: "font.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x80, 0x8, 0x6, 0x0, 0x0, 0x0, 0xbb, 0x81, 0x6f, 0x6a, 0x0, 0x0, 0x3, 0x8c, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0x5c, 0xd9, 0x92, 0xe3, 0x20, 0xc, 0x34, 0x14, 0xff, 0xff, 0xcb, 0xd9, 0xa7, 0xd9, 0x62, 0x59, 0x49, 0xb4, 0xe, 0xe2, 0x4c, 0xdc, 0xfd, 0xe6, 0xd8, 0x1c, 0x12, 0x48, 0xb4, 0x4, 0xa4, 0x5d, 0x1b, 0xbc, 0x5e, 0xaf, 0x57, 0x6b, 0xad, 0x69, 0xcf, 0xeb, 0xb7, 0xf3, 0xb3, 0xf6, 0x9d, 0x55, 0xce, 0xaa, 0xfb, 0xe7, 0x1d, 0xf2, 0x7d, 0x45, 0xbf, 0xfe, 0x16, 0x5c, 0x2b, 0x58, 0x9f, 0xb5, 0xdf, 0xb5, 0xef, 0xbc, 0xdf, 0x58, 0xfd, 0xc9, 0xb6, 0xd1, 0xaf, 0x1b, 0x61, 0xcd, 0x26, 0x14, 0xad, 0xb5, 0xe6, 0x51, 0xe2, 0x8a, 0xa1, 0x4d, 0x91, 0xb9, 0x73, 0xeb, 0x73, 0x76, 0x1a, 0x7a, 0x84, 0xf6, 0xa, 0xe7, 0x6d, 0x67, 0x58, 0x82, 0xed, 0x34, 0x8f, 0xda, 0x67, 0x74, 0xa4, 0xbd, 0x7d, 0xb3, 0xfc, 0xd5, 0x4f, 0xd9, 0xb5, 0x1f, 0xa3, 0x52, 0xfb, 0x96, 0x62, 0xde, 0x65, 0x1a, 0xda, 0xc, 0xd0, 0xea, 0x1b, 0x1e, 0x61, 0xaa, 0x6d, 0x5f, 0x33, 0x2d, 0xaf, 0x89, 0x20, 0x82, 0xaa, 0x26, 0x10, 0x75, 0x22, 0x9a, 0xa6, 0xa5, 0xdf, 0xd7, 0x36, 0xd6, 0xe, 0xa3, 0xef, 0xa4, 0xf6, 0xd1, 0xba, 0xa2, 0x3e, 0x88, 0x78, 0xc, 0x56, 0x72, 0x31, 0x3, 0x2d, 0x3f, 0x7f, 0x6f, 0x95, 0x43, 0x8, 0x13, 0xda, 0x3e, 0xda, 0xae, 0xd6, 0x86, 0xc8, 0xb0, 0x32, 0x6b, 0x6f, 0x25, 0x63, 0xf4, 0x28, 0x1f, 0x16, 0x76, 0x79, 0xee, 0x21, 0x8e, 0x5c, 0xcc, 0xe6, 0x3c, 0x42, 0x57, 0xac, 0x46, 0x73, 0x9f, 0xbb, 0xb7, 0xa0, 0xc7, 0x2c, 0xa2, 0xd4, 0x36, 0xf3, 0x1e, 0x19, 0xf9, 0xb9, 0x8e, 0xb1, 0x2b, 0xa8, 0xb1, 0x39, 0x74, 0xe4, 0x91, 0x65, 0xc9, 0x5a, 0xda, 0x3c, 0xef, 0x23, 0x91, 0x6c, 0xd7, 0xec, 0xfe, 0x87, 0xa8, 0x64, 0xa7, 0x70, 0x9b, 0xb0, 0x2a, 0x63, 0x6d, 0xa3, 0x72, 0xf4, 0x25, 0x61, 0x25, 0x5e, 0xd1, 0x67, 0xb2, 0x72, 0x47, 0x34, 0x76, 0xd2, 0x94, 0xa4, 0x7e, 0xad, 0xbf, 0x77, 0x6f, 0x85, 0x92, 0xf, 0x98, 0xdf, 0x55, 0x3b, 0x43, 0x6f, 0x0, 0x14, 0x51, 0xc2, 0x11, 0x6f, 0x7d, 0x9a, 0xaf, 0x54, 0x29, 0xbc, 0x65, 0x6c, 0xfb, 0xe, 0xc1, 0xcb, 0xdb, 0x45, 0x48, 0xc3, 0x29, 0x42, 0x94, 0x6d, 0x3f, 0xca, 0x54, 0xff, 0x2b, 0x53, 0xa5, 0x4, 0xa9, 0xf2, 0x5d, 0x27, 0xb5, 0x32, 0x5e, 0xea, 0x1e, 0x65, 0x8b, 0xdd, 0xe3, 0x14, 0x90, 0x6f, 0xe6, 0xba, 0x50, 0xae, 0xe0, 0x75, 0x4a, 0x52, 0xe, 0x20, 0x3a, 0x6b, 0x47, 0xc6, 0x16, 0x35, 0xbb, 0x94, 0x12, 0x1d, 0xd1, 0x98, 0xfc, 0x74, 0xa, 0xdc, 0xad, 0x0, 0x89, 0x19, 0x6a, 0xa4, 0x3, 0x2d, 0xaf, 0x29, 0x2d, 0xca, 0x3e, 0x3d, 0xf9, 0xca, 0x71, 0x3d, 0x8, 0x92, 0xf2, 0xfa, 0x3b, 0x22, 0xbc, 0x6a, 0x56, 0xb9, 0x7e, 0x97, 0xe9, 0xef, 0xd0, 0x4, 0xbf, 0x6b, 0xdd, 0x5f, 0x7d, 0xb, 0x92, 0x2f, 0xac, 0xf0, 0x35, 0x5f, 0x99, 0xd9, 0x42, 0xd0, 0xaf, 0x87, 0xc3, 0xb5, 0x14, 0x9d, 0xda, 0x89, 0x45, 0x77, 0x87, 0x23, 0x61, 0xb2, 0xe6, 0x1f, 0x60, 0xbf, 0xa1, 0x31, 0xae, 0xaa, 0x1d, 0x62, 0xed, 0x7b, 0xab, 0x1e, 0xf, 0xb, 0x4c, 0xe5, 0x1b, 0x77, 0x14, 0xd9, 0x93, 0x7c, 0x8c, 0xd8, 0x6c, 0x36, 0x59, 0x8a, 0xd4, 0xd1, 0xa1, 0x60, 0xe1, 0x50, 0x50, 0x52, 0x91, 0xf, 0x88, 0x2a, 0xfe, 0x9f, 0x65, 0x30, 0xba, 0x43, 0xac, 0xd9, 0x92, 0x95, 0x7, 0xdc, 0xa5, 0xa9, 0xac, 0x2d, 0xf9, 0x13, 0x18, 0x19, 0x41, 0x35, 0x7, 0xb4, 0xa3, 0xb3, 0x92, 0xb0, 0x77, 0x91, 0xb0, 0xee, 0x25, 0x25, 0x99, 0xce, 0x49, 0x89, 0xd0, 0xbb, 0x31, 0x76, 0x14, 0x34, 0xab, 0x84, 0xb9, 0x8c, 0x76, 0x78, 0xa2, 0x82, 0x71, 0xee, 0x4c, 0xc9, 0xb3, 0x74, 0x5e, 0x9f, 0x98, 0xe5, 0xbd, 0xc5, 0x4, 0x1e, 0x8d, 0xe8, 0x1a, 0x7c, 0xf2, 0x7d, 0x26, 0xf5, 0xf5, 0x6b, 0x66, 0x80, 0x66, 0xb3, 0xb7, 0x84, 0xe8, 0x51, 0x72, 0x51, 0xcd, 0xe0, 0xa2, 0xdb, 0xf5, 0x48, 0xb9, 0x5e, 0x25, 0x64, 0xf6, 0x4c, 0x81, 0xb5, 0xdb, 0x64, 0x11, 0xad, 0xb4, 0xe3, 0xde, 0x75, 0x4, 0x3d, 0x2e, 0xbb, 0x2b, 0x67, 0x75, 0xda, 0x9a, 0x31, 0xc8, 0x6c, 0xca, 0x28, 0xa4, 0x47, 0xd2, 0xd8, 0x55, 0x4b, 0xa6, 0xc5, 0x2b, 0xd6, 0xa3, 0x74, 0xd1, 0xf4, 0x59, 0x89, 0xcd, 0x20, 0xa3, 0x17, 0x1d, 0x61, 0xd4, 0x9c, 0xb2, 0xc7, 0x70, 0x52, 0x23, 0xe9, 0x15, 0x0, 0x79, 0x8f, 0x9a, 0x5a, 0xd4, 0xef, 0x7c, 0x3d, 0x99, 0xab, 0x12, 0x90, 0x39, 0xc1, 0xa8, 0x86, 0x77, 0xb7, 0x37, 0x76, 0x39, 0x86, 0x5d, 0x88, 0x6c, 0x39, 0x47, 0x6f, 0x7f, 0x76, 0x65, 0xc2, 0x53, 0xcf, 0xb2, 0x63, 0x74, 0x99, 0xf2, 0xfa, 0x82, 0x6a, 0x1f, 0xe1, 0x3e, 0x26, 0x87, 0x64, 0x6d, 0xd1, 0x7d, 0xc2, 0xe8, 0x28, 0xa0, 0x69, 0x35, 0xe4, 0xfe, 0xd1, 0x40, 0x63, 0xe5, 0x9d, 0xf0, 0x9e, 0x75, 0x39, 0xc3, 0xeb, 0x25, 0x56, 0xe8, 0x39, 0x15, 0x6e, 0xb6, 0x1d, 0x5d, 0x6b, 0xd1, 0xcb, 0x4c, 0xda, 0xa8, 0x49, 0x23, 0x1a, 0x49, 0x7d, 0x67, 0xf8, 0x44, 0xd8, 0xe6, 0xd0, 0x4e, 0x7b, 0x63, 0x8c, 0x97, 0x80, 0xa8, 0x3f, 0xb1, 0xfa, 0xd0, 0x77, 0x1e, 0xb7, 0xfa, 0xe8, 0xaa, 0xc7, 0xab, 0xcf, 0xf9, 0x43, 0xc4, 0x9e, 0x43, 0x26, 0x15, 0xe5, 0xd9, 0xe8, 0xb2, 0x63, 0x2d, 0x45, 0xe8, 0x45, 0x48, 0x8f, 0x9d, 0xef, 0xde, 0xf1, 0xb6, 0x8, 0x21, 0x30, 0x41, 0xcf, 0x45, 0xa3, 0xaf, 0x9b, 0x46, 0xde, 0x4, 0xc8, 0xa9, 0x28, 0xeb, 0x74, 0xfd, 0xc, 0x86, 0x3e, 0x59, 0x1, 0xd2, 0x3d, 0xe5, 0x77, 0xb5, 0x3d, 0xbc, 0x5b, 0x4a, 0xa8, 0x30, 0xbf, 0xc6, 0x57, 0x78, 0x6f, 0x70, 0x55, 0x9f, 0x1d, 0x28, 0x4b, 0x69, 0x55, 0x9b, 0x40, 0x55, 0xc, 0xfd, 0x69, 0xbb, 0xc1, 0xae, 0xf0, 0xd2, 0xa, 0x52, 0x76, 0x31, 0xc0, 0x63, 0x36, 0x58, 0xbd, 0x51, 0xdb, 0x3b, 0xae, 0xde, 0x11, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0x97, 0xe3, 0xaf, 0xb5, 0xbf, 0x1, 0xfc, 0xab, 0x4d, 0xe3, 0xf9, 0x71, 0xa, 0xa0, 0x12, 0x9e, 0xa8, 0x0, 0xa, 0x2f, 0x8, 0xcf, 0x3f, 0x53, 0x23, 0x68, 0xff, 0x4, 0x41, 0x10, 0x4, 0x41, 0x10, 0xcf, 0xc4, 0x1f, 0x1e, 0x7e, 0x45, 0x2d, 0xac, 0x5c, 0x67, 0x45, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "hud.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xd3, 0x6b, 0x47, 0x38, 0x0, 0x0, 0x0, 0x7b, 0x50, 0x4c, 0x54, 0x45, 0x71, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0x58, 0xa1, 0x76, 0xe2, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x7, 0x8e, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0x9a, 0x7d, 0x8f, 0xdb, 0x28, 0x10, 0xc6, 0xad, 0x22, 0x85, 0x9c, 0xc, 0xe8, 0x24, 0x2b, 0x2a, 0x96, 0xd8, 0x55, 0xb5, 0x7f, 0x54, 0xfd, 0xfe, 0x9f, 0xf0, 0x78, 0x67, 0x6, 0xe3, 0x97, 0xac, 0x71, 0x8e, 0xcb, 0x31, 0x69, 0x37, 0x89, 0xe9, 0x76, 0xf3, 0xfc, 0x78, 0x98, 0x19, 0xf0, 0xe, 0x83, 0xb, 0xa5, 0xe2, 0xb, 0x1c, 0xc3, 0xc9, 0x10, 0x17, 0x5, 0x5f, 0x79, 0x7c, 0xfc, 0x29, 0x3f, 0xa4, 0x28, 0x3f, 0x82, 0xec, 0xcf, 0x21, 0x12, 0x0, 0x9f, 0x5e, 0xf1, 0xa1, 0x55, 0x0, 0x6b, 0x4, 0x56, 0xf4, 0x7f, 0xac, 0xe8, 0x97, 0x51, 0x7f, 0x50, 0xae, 0x82, 0x9, 0x2a, 0x11, 0x78, 0x5a, 0xd8, 0x3c, 0xcf, 0x2f, 0x77, 0x80, 0xd5, 0xef, 0x9, 0x28, 0xb0, 0x1c, 0xfe, 0xd, 0x0, 0xf3, 0x3a, 0x0, 0x29, 0x65, 0x15, 0x7, 0xcc, 0xb9, 0x3, 0xbc, 0x7e, 0x47, 0x40, 0xa1, 0x45, 0xf0, 0x7a, 0x0, 0x5a, 0xe5, 0xa, 0x1, 0x89, 0x0, 0x6c, 0x38, 0x80, 0xe8, 0x88, 0xca, 0xcd, 0x1b, 0xec, 0x80, 0x99, 0xcf, 0xd8, 0x1, 0x51, 0xff, 0x92, 0xc0, 0xeb, 0x73, 0xc0, 0x7c, 0x14, 0xc0, 0x86, 0x3, 0x8, 0x24, 0x10, 0x5e, 0x3, 0xfd, 0x4c, 0x7c, 0x41, 0x7, 0x18, 0xfd, 0x7f, 0xdb, 0x58, 0xcc, 0xff, 0x55, 0x0, 0x74, 0x86, 0xd9, 0x2, 0x20, 0xe7, 0x35, 0xfd, 0xc7, 0x1d, 0x80, 0x2, 0x3a, 0x40, 0xeb, 0xe7, 0xe2, 0xb, 0x38, 0xc0, 0xeb, 0xbf, 0xdd, 0x1c, 0x1, 0x6d, 0x9, 0x5e, 0x53, 0xff, 0x50, 0xd6, 0xbf, 0x4a, 0x60, 0x5e, 0x25, 0xe0, 0x46, 0xe, 0xe5, 0x80, 0xa8, 0x1a, 0xbc, 0xa, 0x8a, 0x67, 0xad, 0x1f, 0x39, 0xc0, 0x4c, 0xb8, 0xd1, 0x6f, 0x8, 0x38, 0x0, 0x81, 0x80, 0xaa, 0xa1, 0x7f, 0x15, 0xc0, 0x2a, 0x1, 0x99, 0x9, 0x8d, 0xb9, 0x11, 0x5b, 0x60, 0xd3, 0x1, 0x7e, 0xe5, 0xa7, 0xf9, 0x87, 0xe, 0x30, 0x0, 0xa0, 0x3, 0xc, 0x80, 0xdb, 0x5f, 0x37, 0xfd, 0x27, 0x0, 0x18, 0x38, 0xb7, 0xd3, 0x5f, 0x43, 0xff, 0xb3, 0xe, 0x8, 0x8b, 0xa0, 0x10, 0xf2, 0xa8, 0x3, 0x8, 0x5e, 0x5, 0xa8, 0xa, 0xb8, 0x9f, 0xfc, 0xa5, 0x32, 0x7, 0xdc, 0xc4, 0x4d, 0xff, 0x89, 0x0, 0x8c, 0x9, 0xaa, 0x4c, 0xff, 0x3a, 0x80, 0xad, 0x3a, 0x30, 0x67, 0x0, 0x64, 0x94, 0x7f, 0xd0, 0x1, 0x30, 0xf, 0x7c, 0x64, 0xe, 0x50, 0xc4, 0xfc, 0xf8, 0xcc, 0x1, 0x37, 0x1f, 0x9, 0x40, 0x35, 0xfd, 0x2b, 0x49, 0x50, 0x1c, 0x0, 0x10, 0x8d, 0x20, 0x3, 0x82, 0xf9, 0xb0, 0x3, 0x2, 0x81, 0xf0, 0x2e, 0xe5, 0x0, 0x65, 0x2e, 0xab, 0x2c, 0x7, 0xc, 0x51, 0x7f, 0x13, 0x0, 0x82, 0x60, 0xb, 0x42, 0x4a, 0xbf, 0x26, 0x9e, 0xa8, 0x2, 0x1f, 0xa8, 0x16, 0x92, 0xcc, 0x1, 0x6, 0x40, 0xe6, 0x0, 0x47, 0x60, 0x0, 0x0, 0xd4, 0xc5, 0x4b, 0x60, 0xaf, 0x15, 0x88, 0xae, 0x77, 0xba, 0x2d, 0x11, 0x3, 0x41, 0x3c, 0xef, 0x80, 0xf, 0x54, 0x5, 0xf4, 0xa, 0x28, 0x39, 0xc0, 0x10, 0x18, 0x12, 0x0, 0xdb, 0xb, 0x5f, 0x96, 0x4, 0x95, 0xda, 0x49, 0x2, 0x6e, 0xd6, 0xdd, 0xe4, 0xfb, 0x95, 0xb0, 0x4, 0x70, 0xc8, 0x1, 0xae, 0x16, 0x0, 0x7, 0x98, 0xc, 0xa0, 0x2f, 0xfe, 0x50, 0x99, 0x3, 0x6e, 0xc8, 0x1, 0xaa, 0xe6, 0x32, 0x28, 0xa7, 0xc0, 0x4d, 0x2, 0x73, 0x94, 0xec, 0x31, 0xf8, 0x3c, 0x88, 0xc, 0x70, 0xb8, 0xa, 0x20, 0x7, 0x4, 0x0, 0x89, 0x80, 0x51, 0x7b, 0x4b, 0x49, 0x50, 0x7c, 0xc1, 0xed, 0xd0, 0x15, 0x8d, 0x90, 0x51, 0xbf, 0x97, 0x7, 0x92, 0xee, 0x28, 0x5e, 0xe6, 0x0, 0xb6, 0x76, 0x83, 0xa1, 0xfe, 0xc1, 0x3a, 0x88, 0x1c, 0x0, 0x8, 0x98, 0x39, 0x8f, 0x0, 0x6e, 0x41, 0xff, 0x75, 0xad, 0xb0, 0xd1, 0xad, 0xf6, 0x32, 0x21, 0x48, 0x7c, 0x73, 0x5a, 0x6, 0x87, 0x1d, 0x0, 0x32, 0x60, 0xde, 0x9, 0x16, 0x1c, 0xa0, 0x3f, 0x22, 0xd2, 0x8f, 0x0, 0x7c, 0xd6, 0x7, 0x20, 0x2, 0x80, 0xcd, 0x24, 0x60, 0xf6, 0x43, 0x61, 0x11, 0xb8, 0x37, 0x52, 0x3e, 0xe9, 0x80, 0x3f, 0xb8, 0x17, 0x86, 0x55, 0xd0, 0x0, 0x80, 0xe7, 0x1, 0x62, 0x4d, 0xff, 0xf0, 0x79, 0xc9, 0x76, 0x58, 0x1d, 0x3b, 0x13, 0x50, 0x32, 0x0, 0x90, 0xee, 0xcd, 0x61, 0x7, 0x90, 0x8c, 0x0, 0xaa, 0x2, 0x58, 0xbf, 0x3b, 0x11, 0x12, 0x49, 0xbf, 0x29, 0x80, 0x49, 0xff, 0x25, 0x0, 0x76, 0xaa, 0x60, 0x68, 0xfc, 0x43, 0x2a, 0x8, 0xdb, 0xe3, 0x63, 0xfa, 0x83, 0x3, 0x40, 0x1d, 0xc0, 0xbb, 0x41, 0xa5, 0xe5, 0xff, 0x50, 0xf9, 0x99, 0xa0, 0x88, 0xfa, 0xa3, 0xfc, 0x4f, 0x1b, 0xd7, 0x0, 0xd8, 0xf3, 0x80, 0xf2, 0xc7, 0x42, 0xb1, 0x1d, 0xf4, 0xcd, 0x40, 0x95, 0x33, 0x41, 0xdc, 0x9, 0xcb, 0xf0, 0x31, 0x83, 0xfe, 0xec, 0x58, 0x98, 0xf, 0xd7, 0x0, 0x50, 0xfb, 0x6b, 0x20, 0x4a, 0x4f, 0xbb, 0x81, 0xb, 0x4f, 0x85, 0x87, 0xa4, 0x9f, 0xa3, 0xb8, 0xe0, 0x58, 0x5c, 0x9, 0x71, 0x64, 0x11, 0xa4, 0x26, 0x38, 0xee, 0x6, 0x2f, 0x39, 0x15, 0xbe, 0x7f, 0x23, 0x4e, 0x36, 0x42, 0x62, 0x9f, 0x80, 0x4a, 0x3b, 0xe0, 0xd0, 0xc, 0x48, 0x89, 0xbe, 0xe7, 0xf7, 0x5c, 0x7e, 0x3c, 0x7e, 0x95, 0x1f, 0x37, 0x56, 0x7e, 0xbc, 0x1e, 0xc0, 0x7e, 0x1f, 0x64, 0xa6, 0x97, 0xa4, 0xd, 0xa1, 0xeb, 0xb, 0xa5, 0x24, 0xd8, 0x23, 0x2b, 0x4, 0x56, 0xf4, 0x3f, 0xb0, 0x6e, 0x9d, 0x19, 0xfd, 0xab, 0xe7, 0x0, 0xd0, 0x5a, 0xe, 0x38, 0xa, 0x0, 0xd8, 0x80, 0x70, 0x51, 0xcd, 0x1, 0xa6, 0x34, 0x7c, 0xcb, 0x1, 0x94, 0xd2, 0x6a, 0x0, 0xb6, 0x9, 0x10, 0x19, 0xcf, 0x42, 0xfc, 0x51, 0x8, 0x17, 0xf5, 0x1c, 0x40, 0x72, 0x7, 0x50, 0xfa, 0x6a, 0x0, 0xbb, 0x59, 0x20, 0xad, 0x1, 0xfb, 0x97, 0x64, 0x6, 0xd8, 0x76, 0x80, 0xc4, 0xb1, 0xe3, 0x0, 0x4a, 0xad, 0x2e, 0x7a, 0xf7, 0x2f, 0x13, 0xd, 0xf0, 0xc6, 0xf, 0x9c, 0x7, 0x10, 0x85, 0xef, 0x59, 0x80, 0xf8, 0xb5, 0xef, 0xaa, 0x1, 0xc9, 0xff, 0xb5, 0x55, 0x3b, 0xad, 0x38, 0x0, 0x1e, 0xa2, 0xcd, 0xb2, 0xe8, 0x0, 0x16, 0x1d, 0x40, 0xef, 0xb, 0xe5, 0xe1, 0x15, 0x8d, 0x3, 0x80, 0x49, 0x3d, 0x0, 0x9b, 0x4, 0x54, 0x24, 0x60, 0xbe, 0x12, 0x92, 0x8f, 0x5b, 0xb5, 0x93, 0x53, 0x3d, 0x4d, 0x53, 0x7c, 0xe7, 0x1c, 0x80, 0x8e, 0x14, 0x61, 0xe, 0xf0, 0xbd, 0x21, 0xac, 0x2, 0x6e, 0xfa, 0x69, 0x26, 0xd5, 0xcd, 0x37, 0x72, 0xbf, 0x67, 0x42, 0x4f, 0x1, 0x50, 0x10, 0xc0, 0x31, 0xf, 0xcc, 0x8b, 0xa, 0x90, 0x72, 0xc0, 0xe4, 0xf5, 0x9b, 0x67, 0x47, 0x21, 0x39, 0x20, 0x2, 0x80, 0xe, 0x20, 0x78, 0xfe, 0xb5, 0x3, 0x68, 0x98, 0x66, 0x76, 0x67, 0xe, 0x43, 0x10, 0x4a, 0x61, 0x66, 0x88, 0xf6, 0x3f, 0xe9, 0x0, 0xb5, 0xbf, 0x13, 0x80, 0x79, 0x40, 0xba, 0xfb, 0x21, 0x64, 0xf9, 0x3d, 0xbf, 0xd3, 0x1a, 0x98, 0x9c, 0x3, 0xa6, 0x79, 0xe9, 0x80, 0xe9, 0xa1, 0x9, 0x2, 0x7, 0x98, 0xb9, 0xcf, 0x1c, 0x60, 0x75, 0x72, 0x7e, 0xe7, 0x8c, 0xb1, 0x3b, 0x65, 0xc0, 0x9, 0x94, 0xfb, 0x17, 0x3c, 0xe4, 0x8, 0xef, 0x95, 0x93, 0x7d, 0xc0, 0x13, 0x41, 0xa4, 0xb3, 0x2, 0x11, 0x2b, 0xe, 0xb0, 0xaa, 0xe1, 0xa, 0xc0, 0xe, 0x98, 0x1e, 0xf, 0xec, 0x0, 0x62, 0x57, 0x1, 0x43, 0xe, 0x30, 0x9a, 0x3c, 0x0, 0x4b, 0x0, 0x85, 0x1d, 0xd3, 0x41, 0xc3, 0x4a, 0xb1, 0x46, 0x78, 0x29, 0x80, 0x92, 0xff, 0x41, 0x15, 0x30, 0xea, 0x1d, 0x85, 0xdf, 0x4b, 0x7, 0x4c, 0x8f, 0x9f, 0x59, 0xe, 0xb0, 0xfa, 0x6f, 0x20, 0xf, 0x20, 0x7, 0x68, 0x2, 0x2c, 0xe5, 0xfb, 0xf0, 0xd5, 0x8c, 0x44, 0x2, 0xf6, 0xe9, 0x75, 0xbf, 0x21, 0x42, 0x8a, 0xb3, 0x9f, 0xf5, 0x1, 0xd3, 0xe4, 0xb3, 0x40, 0xee, 0x0, 0xab, 0x3f, 0xab, 0x2, 0x56, 0x39, 0xce, 0x1, 0x8, 0x80, 0x4e, 0x4, 0xd1, 0x4, 0x71, 0x9, 0x58, 0x2, 0xcc, 0x3, 0xb9, 0xbf, 0x12, 0xc0, 0x7a, 0x80, 0xca, 0x97, 0xd4, 0x43, 0x7, 0x4c, 0x11, 0x0, 0x72, 0x0, 0x5c, 0xff, 0xc9, 0x1, 0x77, 0x7e, 0x67, 0x2e, 0x46, 0x3, 0xc1, 0xe8, 0xe4, 0xfa, 0x22, 0xd5, 0x57, 0x4d, 0x7a, 0x64, 0x96, 0x0, 0xd, 0xa5, 0x81, 0xb6, 0x0, 0x0, 0x3b, 0xa0, 0xd0, 0x7, 0x3c, 0x26, 0xaf, 0x1f, 0x57, 0x81, 0xac, 0xa, 0x5a, 0x7, 0x18, 0xad, 0x56, 0xfb, 0x68, 0x82, 0x79, 0x4, 0xcc, 0xe6, 0x3d, 0x66, 0x8b, 0x83, 0x1e, 0x61, 0xdc, 0x3c, 0x3b, 0x4f, 0x34, 0xe7, 0x80, 0x65, 0x27, 0xa8, 0x1, 0xe8, 0x48, 0x2b, 0x20, 0xf5, 0x1, 0x8e, 0x1, 0x74, 0x80, 0x59, 0x0, 0x76, 0xea, 0x1d, 0x80, 0xd1, 0x65, 0x43, 0x53, 0x16, 0xdd, 0x9a, 0xb0, 0x34, 0xdc, 0x65, 0xc6, 0x9a, 0x1, 0x0, 0xbb, 0xbf, 0x69, 0xb1, 0x17, 0xb0, 0x16, 0x80, 0xfa, 0x1f, 0x68, 0x17, 0x8, 0x3c, 0x30, 0xd8, 0x7a, 0xc7, 0x92, 0x1, 0x3c, 0x2, 0x7, 0x40, 0xbf, 0xb4, 0x4, 0x18, 0x8d, 0x70, 0x5e, 0x5a, 0x5, 0xe, 0x39, 0x60, 0x5a, 0x3a, 0xc0, 0xae, 0x81, 0x9f, 0x3f, 0x53, 0x1f, 0x18, 0xfb, 0x80, 0x54, 0x9, 0x80, 0x3, 0xa8, 0x13, 0xcf, 0x12, 0x81, 0x71, 0x34, 0xab, 0x3d, 0x2e, 0x9, 0x3, 0x4, 0xac, 0x8f, 0xc6, 0x1c, 0xe0, 0xeb, 0x20, 0xde, 0xd, 0x4a, 0xd4, 0x7, 0x7, 0x7, 0x10, 0x86, 0x38, 0x4, 0x7, 0x70, 0xb7, 0x0, 0x10, 0x0, 0xad, 0xf8, 0x4e, 0xbd, 0x66, 0x38, 0x68, 0x9, 0x34, 0xe5, 0x0, 0x33, 0xff, 0xd3, 0xe2, 0x3c, 0x20, 0xde, 0x5a, 0x87, 0xe7, 0x1, 0xd1, 0xf9, 0xb8, 0x13, 0xd4, 0x0, 0x46, 0xa8, 0xdd, 0x9, 0xbd, 0x67, 0x40, 0xe2, 0xb3, 0x26, 0xd0, 0x94, 0x3, 0x4c, 0xd, 0x4, 0x75, 0x30, 0x9e, 0x0, 0xf8, 0x1d, 0x21, 0x70, 0x40, 0x9a, 0xf7, 0xac, 0x13, 0xe4, 0x58, 0xab, 0x77, 0x0, 0x2d, 0x5c, 0xb5, 0x20, 0x58, 0x63, 0xe, 0xc0, 0xfb, 0xe2, 0x78, 0x2, 0xe2, 0xb6, 0xc2, 0xf0, 0x44, 0x88, 0x94, 0x7a, 0x1, 0x3, 0xc0, 0x7b, 0x1c, 0xaf, 0x2, 0xb3, 0xec, 0x59, 0x4a, 0x8a, 0x60, 0x98, 0xb1, 0xb6, 0x1c, 0x80, 0x2b, 0x21, 0x3c, 0x5, 0x4a, 0xf3, 0x9f, 0x9f, 0x9, 0x2, 0x7, 0x38, 0x3, 0xd8, 0x7b, 0x0, 0x5c, 0x29, 0xb8, 0xda, 0x69, 0x92, 0xef, 0x6, 0xcd, 0xb0, 0xbd, 0xd0, 0x58, 0x1f, 0xb0, 0x7a, 0x26, 0x28, 0xe5, 0x81, 0x53, 0x61, 0xdb, 0x3, 0x5a, 0xf9, 0xf6, 0x36, 0x80, 0x82, 0x33, 0xed, 0xe7, 0x3e, 0xd, 0x7a, 0x2, 0xad, 0xf5, 0x1, 0xcf, 0x9f, 0xa, 0x3, 0x7, 0xb0, 0x38, 0xff, 0x42, 0x84, 0x69, 0x66, 0xc, 0xf8, 0x3e, 0x1b, 0x6c, 0x5, 0xc0, 0x50, 0x2b, 0x8c, 0x54, 0x2f, 0x51, 0x9, 0x3f, 0xcb, 0xd1, 0x2, 0x2c, 0xb8, 0xc3, 0x1c, 0xe6, 0xbb, 0xc1, 0x77, 0x6, 0x0, 0x2c, 0x90, 0x72, 0x9f, 0xc2, 0x83, 0x76, 0xec, 0xed, 0x0, 0x30, 0xe4, 0xf2, 0x25, 0x80, 0x34, 0xe8, 0x86, 0xde, 0xd, 0xc0, 0x98, 0xe5, 0x39, 0xb8, 0x2b, 0xca, 0x33, 0x64, 0x33, 0x55, 0xa0, 0x32, 0x0, 0x58, 0xe9, 0xd8, 0x98, 0x95, 0xff, 0x54, 0x23, 0xfd, 0xf5, 0xb7, 0x3, 0xb0, 0xe8, 0x75, 0xd2, 0xc5, 0xd2, 0xd0, 0x5b, 0x2, 0x18, 0x59, 0xa9, 0x1f, 0x6, 0x6e, 0x80, 0x97, 0x4f, 0x0, 0xe0, 0xcd, 0x2, 0x28, 0x45, 0x71, 0x37, 0x70, 0x12, 0x40, 0x88, 0xff, 0x29, 0x0, 0xf4, 0xdb, 0x27, 0xed, 0x3, 0x18, 0x8b, 0x0, 0x4e, 0x6c, 0x86, 0x78, 0xe3, 0x0, 0xf8, 0xc5, 0x0, 0x82, 0x72, 0xfd, 0x93, 0x4e, 0x13, 0x78, 0xb5, 0x3, 0xb8, 0xa5, 0xc3, 0x2b, 0x0, 0x88, 0xff, 0xe3, 0x39, 0x4, 0xf5, 0x1, 0xf0, 0x92, 0xf, 0xaa, 0x3a, 0x0, 0xe9, 0x77, 0x4, 0xda, 0x4f, 0x82, 0xa5, 0x43, 0x21, 0x71, 0x2, 0x0, 0x5a, 0x71, 0x4d, 0x0, 0x60, 0x47, 0x1, 0xf0, 0x78, 0xf1, 0x9b, 0x0, 0x78, 0x1, 0x0, 0x6f, 0x0, 0x0, 0xdb, 0xd4, 0xcf, 0x16, 0xb5, 0x90, 0x89, 0x13, 0x0, 0x8c, 0x6a, 0x2f, 0xfe, 0xa4, 0x5, 0x2e, 0x1, 0xc0, 0x97, 0xfa, 0x17, 0xab, 0x80, 0x9d, 0x1, 0x30, 0x36, 0x8, 0x40, 0xac, 0x59, 0x20, 0xde, 0x2c, 0x45, 0x45, 0x81, 0x9d, 0x5, 0x50, 0x2b, 0x9, 0xd4, 0x4, 0x80, 0xef, 0x8b, 0x15, 0x20, 0x60, 0x4b, 0xbc, 0x1b, 0x0, 0x51, 0x74, 0x3a, 0x68, 0x86, 0x19, 0xba, 0x41, 0xf4, 0xa6, 0x0, 0xd6, 0xda, 0xfe, 0x11, 0xdc, 0x34, 0x8f, 0xfa, 0xbf, 0xd, 0x40, 0xb4, 0x9, 0x40, 0x6c, 0x38, 0x20, 0x35, 0x83, 0x57, 0x0, 0x10, 0x8d, 0x0, 0x10, 0xeb, 0x49, 0x60, 0x59, 0x12, 0xc4, 0xb7, 0x1, 0x98, 0xc0, 0xfa, 0x45, 0xb, 0x8d, 0x90, 0x88, 0x4, 0x16, 0x47, 0x61, 0xb8, 0x0, 0x24, 0xfd, 0x27, 0x0, 0x8, 0xb8, 0xcb, 0x6e, 0x8, 0x80, 0xc8, 0x6e, 0xb, 0xa6, 0xc3, 0xb0, 0x45, 0x5, 0xa8, 0x0, 0x20, 0x39, 0x40, 0x34, 0x3, 0x40, 0x30, 0x78, 0x27, 0x14, 0x78, 0x80, 0x15, 0xf4, 0x9f, 0x1, 0x20, 0x2a, 0xe9, 0xaf, 0xe, 0x40, 0x88, 0xc2, 0xbd, 0x60, 0x96, 0x5d, 0x13, 0x15, 0x0, 0x88, 0x3a, 0xfa, 0xeb, 0x1, 0xb8, 0x3a, 0xa, 0x9f, 0xbd, 0x82, 0xfc, 0xff, 0x36, 0x80, 0xb6, 0x4e, 0x85, 0x3b, 0x80, 0xe, 0xa0, 0x3, 0xe8, 0x0, 0x3a, 0x80, 0xe, 0xa0, 0x3, 0xe8, 0x0, 0x3a, 0x80, 0xe, 0xa0, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xd1, 0xa3, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xd1, 0xa3, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xd1, 0xa3, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xd1, 0xa3, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xf1, 0x66, 0xf1, 0xf, 0x98, 0xa1, 0x5a, 0x77, 0x3b, 0x19, 0x70, 0xd6, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "terrain.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x51, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xee, 0xa1, 0x5c, 0x1, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x0, 0xb3, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0xd5, 0xc9, 0x11, 0xc3, 0x20, 0xc, 0x0, 0x40, 0x8a, 0x40, 0xea, 0xbf, 0xd4, 0x80, 0x20, 0x49, 0x7, 0xe2, 0xb3, 0x1b, 0xcf, 0x38, 0x3f, 0x5d, 0x58, 0x8c, 0x71, 0x64, 0xc6, 0xfe, 0xad, 0x67, 0x3c, 0xb1, 0x82, 0xe7, 0x3c, 0x39, 0xe4, 0x93, 0xf8, 0x27, 0xf8, 0x7d, 0xbd, 0xa8, 0x7f, 0x7, 0x8f, 0xb9, 0x67, 0x30, 0x1f, 0x4c, 0x61, 0x15, 0x7e, 0xe7, 0x1f, 0xf7, 0x28, 0xf4, 0xc6, 0x3f, 0x81, 0x2b, 0x89, 0x3d, 0x82, 0x1a, 0x44, 0x6b, 0xfd, 0xb7, 0xfa, 0x4a, 0xa1, 0x32, 0x59, 0x7f, 0x66, 0x6b, 0xfd, 0x19, 0x37, 0x89, 0x59, 0x4f, 0x54, 0x22, 0x9d, 0x7, 0x60, 0xde, 0xb3, 0x77, 0xf, 0x40, 0xce, 0xd5, 0x82, 0xec, 0x6d, 0xc1, 0x9, 0x5d, 0x33, 0x38, 0x83, 0x68, 0x3e, 0x4, 0xab, 0x5, 0xd5, 0x87, 0xdd, 0x8d, 0xfa, 0x12, 0xa2, 0x7b, 0xd, 0xfc, 0x3f, 0xc3, 0xa8, 0x14, 0xba, 0x17, 0x41, 0xdc, 0x4d, 0xf0, 0xbd, 0xf, 0xfa, 0x57, 0x61, 0x7c, 0xd7, 0x40, 0x2d, 0x82, 0x47, 0x97, 0x51, 0x9c, 0x65, 0xfc, 0xea, 0x3a, 0xfc, 0xd, 0x20, 0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0xaf, 0xf, 0xcd, 0x67, 0x16, 0xb4, 0xc6, 0x5a, 0xdb, 0x3, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "trees.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x7e, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x51, 0x0, 0x0, 0xe0, 0x5a, 0x40, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xc2, 0x5d, 0x52, 0x38, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x1, 0x58, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0x96, 0x31, 0x96, 0xc4, 0x30, 0x8, 0x43, 0xdd, 0xc6, 0x66, 0x9d, 0x14, 0x79, 0x3c, 0xdd, 0xff, 0xa4, 0xb, 0xd8, 0x9e, 0x9d, 0xbd, 0x0, 0x34, 0xfa, 0x53, 0x24, 0x9d, 0x64, 0x21, 0x93, 0x69, 0x8d, 0x10, 0xf2, 0x41, 0x55, 0x6b, 0xe5, 0x5f, 0x55, 0xc0, 0x5c, 0x5c, 0xf6, 0x2b, 0xd0, 0x37, 0x79, 0x85, 0xbb, 0xf0, 0x97, 0xec, 0x2c, 0xd4, 0x85, 0x2f, 0xdd, 0x8f, 0x2b, 0xdd, 0x80, 0xbe, 0xaf, 0x86, 0xec, 0xbb, 0x3c, 0x78, 0x12, 0xe9, 0xe9, 0xeb, 0xd1, 0x77, 0x80, 0xe4, 0x1, 0x5c, 0x71, 0xe8, 0x88, 0x1e, 0xa1, 0x7f, 0xa7, 0xf7, 0x2f, 0xea, 0xbf, 0x8c, 0xc0, 0xf4, 0xef, 0xec, 0xa, 0x9a, 0xee, 0xbb, 0xc3, 0xf, 0xfd, 0xec, 0x11, 0x7c, 0xe9, 0x1f, 0x17, 0xa9, 0x6, 0x80, 0xd5, 0x0, 0x7f, 0xae, 0x3d, 0x90, 0xda, 0x82, 0x50, 0xfc, 0x5c, 0x80, 0x5d, 0xc5, 0x54, 0x3, 0x4b, 0xf3, 0x6f, 0xe, 0x11, 0xc5, 0x9d, 0x1a, 0xc0, 0xda, 0x3d, 0x71, 0x7c, 0x9f, 0x86, 0xd7, 0xf0, 0x4e, 0xe, 0x0, 0xff, 0x8f, 0x9f, 0x3b, 0x81, 0x68, 0xfd, 0xf9, 0x18, 0x61, 0x5f, 0xc4, 0xc4, 0x4b, 0x10, 0x8a, 0xe7, 0x63, 0x1c, 0xc7, 0xcf, 0x5e, 0x44, 0xcd, 0x16, 0xcf, 0x5e, 0x2, 0x35, 0xfa, 0x9f, 0x10, 0xd6, 0x16, 0x2e, 0xd0, 0xdf, 0xe, 0xce, 0xf8, 0xb, 0xf4, 0x8d, 0x31, 0xaa, 0xe2, 0x3f, 0xf4, 0x61, 0x26, 0xee, 0x3a, 0x7d, 0x77, 0xd0, 0x87, 0x8c, 0xa7, 0xf0, 0x7f, 0x71, 0xef, 0xcf, 0x23, 0x5a, 0x6a, 0x60, 0x8, 0xa, 0xf5, 0xdb, 0xf3, 0x74, 0x33, 0xa0, 0xa5, 0x25, 0x10, 0x94, 0x8e, 0x60, 0xd8, 0x8, 0x50, 0x5c, 0xc2, 0x62, 0x3, 0x43, 0x26, 0xb4, 0xce, 0x83, 0x1b, 0xb0, 0x16, 0xa2, 0xca, 0x41, 0x1f, 0xb6, 0x6, 0x64, 0x7d, 0x10, 0x4a, 0xf4, 0x45, 0xbc, 0x85, 0x13, 0x55, 0x16, 0xc6, 0x9c, 0x5d, 0xa6, 0x6c, 0x7, 0x15, 0x1, 0x0, 0x5d, 0x7e, 0x96, 0x3, 0xcd, 0x37, 0xa0, 0x98, 0xf0, 0x4d, 0x8c, 0x39, 0xc5, 0x5e, 0xb, 0x12, 0xb0, 0xd8, 0xc7, 0x84, 0x77, 0x70, 0xce, 0xa2, 0x16, 0x7a, 0x1, 0x2b, 0x6f, 0x81, 0xd, 0xa1, 0x29, 0x50, 0x67, 0xe0, 0x6b, 0x16, 0x8d, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x10, 0x42, 0x8, 0x21, 0x84, 0x64, 0xf2, 0xb, 0x44, 0x57, 0x2c, 0xc3, 0x96, 0x39, 0x8c, 0xa5, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
		{name: "wolf.png", b: []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x80, 0x8, 0x3, 0x0, 0x0, 0x0, 0xf4, 0xe0, 0x91, 0xf9, 0x0, 0x0, 0x0, 0x78, 0x50, 0x4c, 0x54, 0x45, 0x0, 0x0, 0x0, 0x52, 0x1e, 0x2e, 0x78, 0x18, 0x26, 0xa5, 0x2d, 0x27, 0xc9, 0x6d, 0x45, 0xcc, 0xa5, 0x62, 0xcf, 0xc2, 0x81, 0xe3, 0xe7, 0xab, 0xf7, 0xe1, 0xbc, 0xd9, 0xb6, 0x91, 0xc7, 0x97, 0x7d, 0x8d, 0x5b, 0x57, 0x57, 0x32, 0x22, 0x28, 0x22, 0x1f, 0x41, 0x33, 0x25, 0x63, 0x49, 0x2c, 0x83, 0x6b, 0x3f, 0xa2, 0x9d, 0x6b, 0x97, 0xb5, 0x8a, 0x65, 0x84, 0x5c, 0x3a, 0x59, 0x41, 0x20, 0x37, 0x2d, 0x20, 0x2c, 0x11, 0x3a, 0x3b, 0x3d, 0x5f, 0x63, 0x67, 0x86, 0x95, 0x98, 0xac, 0xc2, 0xc3, 0xd4, 0xed, 0xed, 0xac, 0xdc, 0xf1, 0x7c, 0xa8, 0xd5, 0x4c, 0x6e, 0xad, 0x3e, 0x3d, 0x8a, 0x37, 0x2a, 0x5e, 0x33, 0x14, 0x4d, 0x58, 0x26, 0x4f, 0x72, 0x40, 0x4a, 0x85, 0x5b, 0x69, 0xb1, 0x8b, 0x9a, 0xdf, 0xb9, 0xca, 0xea, 0xdc, 0xe6, 0xaf, 0x44, 0x2c, 0x26, 0x0, 0x0, 0x0, 0x1, 0x74, 0x52, 0x4e, 0x53, 0x0, 0x40, 0xe6, 0xd8, 0x66, 0x0, 0x0, 0x0, 0xca, 0x49, 0x44, 0x41, 0x54, 0x78, 0xda, 0xed, 0xd6, 0xb1, 0xe, 0xc2, 0x30, 0xc, 0x84, 0x61, 0x24, 0x66, 0x66, 0x9f, 0xc5, 0xfb, 0xbf, 0x27, 0x4e, 0xd2, 0xb4, 0x88, 0x1d, 0x1f, 0x12, 0xff, 0x37, 0x75, 0xf3, 0xd5, 0x71, 0xac, 0xdc, 0x6e, 0x0, 0x0, 0xe0, 0x27, 0xc9, 0x1d, 0x20, 0xc3, 0xf7, 0xef, 0x99, 0x8a, 0x90, 0x5c, 0x9, 0xaa, 0xb4, 0x9e, 0x8f, 0x7b, 0x5, 0x28, 0x96, 0x0, 0x23, 0x41, 0xd5, 0xd7, 0x8a, 0x10, 0xa6, 0x1e, 0x68, 0x27, 0xe8, 0x3f, 0x89, 0x58, 0x72, 0x4c, 0xc2, 0x8c, 0xe0, 0xb9, 0x3, 0x55, 0x7e, 0xb6, 0x60, 0x44, 0xf0, 0x5c, 0xcd, 0x8c, 0x38, 0xcf, 0xe1, 0xc, 0xa0, 0xd6, 0x93, 0x48, 0x9d, 0x4d, 0x70, 0xb4, 0x40, 0x63, 0xc, 0xa4, 0x9c, 0x11, 0x76, 0x13, 0x7a, 0x3b, 0x70, 0x4d, 0xc2, 0x5e, 0x9, 0x9d, 0x47, 0x70, 0xfc, 0xf3, 0x15, 0x61, 0x2e, 0xc7, 0xce, 0x9, 0xdc, 0x5f, 0xeb, 0x14, 0xd6, 0x2c, 0x64, 0xe3, 0x0, 0xbc, 0x5d, 0xc8, 0x35, 0x87, 0xea, 0xac, 0xff, 0xd9, 0x8d, 0xb1, 0xf, 0xab, 0xbe, 0x29, 0xc0, 0xde, 0x8d, 0x32, 0xbe, 0x10, 0x8e, 0xed, 0xec, 0x7b, 0x21, 0xec, 0xc, 0xde, 0x27, 0x92, 0xbb, 0x7e, 0x25, 0xf8, 0xef, 0xfa, 0xe9, 0xfe, 0x7f, 0xd9, 0x27, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xbe, 0xe8, 0x5, 0xe6, 0xb, 0x15, 0xaf, 0x79, 0xe2, 0x6f, 0xc1, 0x0, 0x0, 0x0, 0x0, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}},
	} {
		decoded, err := png.Decode(bytes.NewReader(file.b))
		if err != nil {
			panic(fmt.Sprintf("png.Decode %s: %v", file.name, err))
		}
		Images[file.name] = decoded
	}
}
