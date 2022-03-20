package main

import (
	"flag"
	"fmt"
	"go-image-water-case/pkg"
	"log"
	"strconv"
	"sync"
	"time"
)

var num int
var concurrent int
var sweep int

func main() {
	flag.IntVar(&num, "n", 50, "图片数量")
	flag.IntVar(&concurrent, "c", 5, "并发数量")
	flag.IntVar(&sweep, "s", 1, "清理生成出来的水印")
	flag.Parse()

	defer func() {
		pkg.Sweep(sweep)
	}()

	fmt.Printf("水印图：%d 张\n", num)

	syncGen()
	goroutineGen()
	goroutineCtrlGen()
}

// 同步
func syncGen() {
	startT := time.Now()

	// 打开水印图
	water, err := pkg.OpenPngImage("./test_water.png")
	if err != nil {
		log.Fatal(err)
	}

	// 这个循环是遍历原图
	imgPath := "./test.png"
	for i := 1; i <= num; i++ {
		id := i + 1000 // 模拟 id
		if err = pkg.Generate(imgPath, id, water); err != nil {
			log.Println("水印图生成失败，id=" + strconv.Itoa(i))
		}
	}

	// 计算耗时
	tc := time.Since(startT)

	fmt.Printf("同步执行时间 = %v\n", tc)
}

// 无限制并发
func goroutineGen() {
	startT := time.Now()

	water, err := pkg.OpenPngImage("./test_water.png")
	if err != nil {
		log.Fatal(err)
	}

	imgPath := "./test.png"

	wg := sync.WaitGroup{}

	// 这个循环是遍历目录下的原图
	for i := 1; i <= num; i++ {
		wg.Add(1)
		id := i + 1000 // 模拟 id

		go func() {
			defer wg.Done()
			if err = pkg.Generate(imgPath, id, water); err != nil {
				log.Println("水印图生成失败，id=" + strconv.Itoa(i))
			}
		}()
	}

	wg.Wait()
	// 计算耗时
	tc := time.Since(startT)
	fmt.Printf("无并发数量控制时间 = %v\n", tc)
}

// 控制并发数量
func goroutineCtrlGen() {
	startT := time.Now()
	wg := sync.WaitGroup{}
	// 指定缓冲数量
	ch := make(chan struct{}, concurrent)

	water, err := pkg.OpenPngImage("./test_water.png")
	if err != nil {
		log.Fatal(err)
	}

	// 模拟数据量
	// 这个循环是遍历目录下的原图
	for i := 1; i <= num; i++ {
		// 模拟读取的量
		imgPath := "./test.png"
		wg.Add(1)

		id := i + 1000 // 模拟 id
		pkg.Go(func() {
			defer wg.Done()
			ch <- struct{}{}

			if err = pkg.Generate(imgPath, id, water); err != nil {
				log.Println("水印图生成失败，id=" + strconv.Itoa(id))
			}
			<-ch
		})
	}

	wg.Wait()
	// 计算耗时
	tc := time.Since(startT)
	fmt.Printf("同时并发 %d 个时间 = %v\n", concurrent, tc)
}
