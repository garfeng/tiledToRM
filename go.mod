module github.com/garfeng/tiledToRM

go 1.18

require (
	github.com/fsnotify/fsnotify v1.6.0
	github.com/lafriks/go-tiled v0.11.0
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)

replace github.com/lafriks/go-tiled v0.11.0 => github.com/garfeng/go-tiled v0.11.4
