module github.com/garfeng/tiledToRM

go 1.18

require (
	github.com/disintegration/imaging v1.6.2
	github.com/fsnotify/fsnotify v1.6.0
	github.com/lafriks/go-tiled v0.11.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/lafriks/go-tiled v0.11.0 => github.com/garfeng/go-tiled v0.11.1
