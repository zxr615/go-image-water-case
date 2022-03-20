package pkg

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Generate 生成水印
// imgPath 源文件路径
// id 文字水印id
// waterImg 水印
func Generate(imgPath string, id int, water image.Image) error {
	// 打开原图文件
	origin, err := OpenPngImage(imgPath)
	if err != nil {
		return err
	}

	x := origin.Bounds().Dx() // 图片长(px)
	y := origin.Bounds().Dy() // 图片宽(px)

	// 水印位置：以原图长宽 - 去水印图 - 偏移像素
	waterPosition := image.Pt(
		x-water.Bounds().Dx()-100,
		y-water.Bounds().Dy()-100,
	)

	// 新建一个画布，把原图、水印、文字叠加上去
	newImage := image.NewRGBA(origin.Bounds())
	// 叠加原图
	draw.Draw(newImage, origin.Bounds(), origin, image.Point{}, draw.Over)
	// 叠加水印
	draw.Draw(newImage, water.Bounds().Add(waterPosition), water, image.Point{}, draw.Over)

	// 加载字体
	font, err := GetFont()
	if err != nil {
		return err
	}
	fc := freetype.NewContext()   // xx
	fc.SetDPI(72)                 // dpi
	fc.SetFont(font)              // 字体
	fc.SetFontSize(48)            // 字体大小(磅)
	fc.SetClip(newImage.Bounds()) // 设置剪裁矩形以进行绘制
	fc.SetDst(newImage)           // 设置目标图像(字体加在哪里)
	//fc.SetSrc(image.Black)        // 字体颜色
	// 可根据 RGBA 设置颜色
	fc.SetSrc(image.NewUniform(color.RGBA{R: 21, G: 33, B: 57, A: 255}))

	// 绘制文字eg: id:10086
	text := "id: " + strconv.Itoa(id)
	// 字体水印位置：在水印位置的基础上 偏移 Y 轴
	fontPt := freetype.Pt(waterPosition.X, waterPosition.Y+200)
	if _, err = fc.DrawString(text, fontPt); err != nil {
		return err
	}

	// 创建水印文件
	savePath := "./water/" + strconv.Itoa(id) + ".png"
	created, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer created.Close()

	// 写入
	if err = jpeg.Encode(created, newImage, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return nil
}

// OpenPngImage 获取水印
func OpenPngImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	water, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return water, nil
}

// GetFont 获取字体
func GetFont() (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile("./font.ttf")
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(fontBytes)
}

// Sweep 清理图片
func Sweep(sweep int) {
	if sweep == 0 {
		return
	}

	fs, err := ioutil.ReadDir("./water")
	if err != nil {
		log.Printf("水印图清理失败：%+v", err)
	}

	for _, v := range fs {
		if v.Name() == ".gitkeep" {
			continue
		}
		if err := os.Remove("./water/" + v.Name()); err != nil {
			log.Printf("水印图清理失败：%+v", err)
		}
	}
}

// Go 避免 go func(){} 如果方法中抛出 panic 无法被捕获到
// 或者是每在每个 go 前面都 recover() 一次，造成的代码混乱不可维护
func Go(f func()) {
	defer func() {
		if err := recover(); err != nil {
			// 记录日志
			log.Println(err)
		}
	}()

	go f()
}
