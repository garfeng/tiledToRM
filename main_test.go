package main

import "testing"

func TestRenderOneMap(t *testing.T) {
	*dstRoot = "./tmp"
	*separateGroup = true
	renderOneMap("C:/Users/jerri/Documents/RMMZ/XiuXian/maps/7.tmx")
}
