package main

import "testing"

func TestRenderOneMap(t *testing.T) {
	*dstRoot = "./tmp"
	renderOneMap("C:/Users/jerri/Documents/RMMZ/XiuXian/maps/Test1.tmx")
}
