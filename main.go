package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"image/png"
	"os"

	"log"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	mapRoot = flag.String("map", "./maps", "tiled map dir")
	dstRoot = flag.String("dst", "./img/tilesets", "dst tilesets dir")
)

func main() {
	flag.Parse()
	ScanAndGenerateAll()
	SetupMonitor()

}

func ScanAndGenerateAll() {
	w, err := os.Open(*mapRoot)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	names, err := w.Readdirnames(-1)

	if err != nil {
		panic(err)
	}

	for _, name := range names {
		AddUpdateEvent(filepath.Join(*mapRoot, name))
	}
}

func SetupMonitor() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.Add(*mapRoot)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				fmt.Println(event)
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					AddUpdateEvent(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		for {
			e := <-events
			if !TryGenerateMap(e) {
				events <- e
			}
		}
	}()

	<-make(chan struct{})
}

var (
	generateTime = map[string]time.Time{}
	generateLock = map[string]*atomic.Int32{}
	events       = make(chan *UpdateEvent, 100)
)

type UpdateEvent struct {
	Name string
	Time time.Time
}

func AddUpdateEvent(name string) {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".tmx" {
		return
	}

	events <- &UpdateEvent{
		Name: name,
		Time: time.Now(),
	}
}

func TryGenerateMap(e *UpdateEvent) bool {
	if generateLock[e.Name] == nil {
		generateLock[e.Name] = new(atomic.Int32)
		generateLock[e.Name].Store(0)
	}

	m := generateLock[e.Name]
	log.Println("尝试处理", e.Name)
	// 正在处理中
	if !m.CompareAndSwap(0, 1) {
		return false
	}
	// 结束的时候换回来
	defer m.Store(0)

	// 最后编译时间晚于用户修改时间
	if e.Time.Sub(generateTime[e.Name]) <= 0 {
		log.Println("已处理过")
		return true
	}

	generateTime[e.Name] = time.Now()

	renderOneMap(e.Name)

	return true
}

func renderOneMap(name string) {
	tiledMap, err := tiled.LoadFile(name)
	if err != nil {
		log.Println("Err when load tilemap", name, err)
		return
	}

	renderer, err := render.NewRenderer(tiledMap)
	if err != nil {
		log.Println("Err when create render", name, err)
		return
	}

	err = renderer.RenderVisibleLayersAndObjectGroups()
	if err != nil {
		log.Println("Err when render map", name, err)
		return
	}

	err = renderer.RenderVisibleGroups()

	if err != nil {
		log.Println("Err when render map", name, err)
		return
	}
	defer renderer.Clear()

	img := renderer.Result
	buffer := bytes.NewBuffer(nil)
	err = png.Encode(buffer, img)
	if err != nil {
		log.Println("Err when encode map to png", name, err)
		return
	}

	_, dstName := filepath.Split(name)
	dstFullname := filepath.Join(*dstRoot, dstName)
	dstFullname = ReplaceExtTo(dstFullname, ".png")
	os.WriteFile(dstFullname, buffer.Bytes(), 0755)

	if err != nil {
		log.Println("Err when save png", name, err)
		return
	}

}

func ReplaceExtTo(name string, newExt string) string {
	ext := filepath.Ext(name)
	idx := strings.LastIndex(name, ext)
	return name[:idx] + newExt
}
