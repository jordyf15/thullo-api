package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func ResizeImage(reader NamedFileReader, width int) (*os.File, error) {
	fileExt := GetFileExtension(reader.Name())
	tmpFile, err := ioutil.TempFile(os.TempDir(), "dp-*."+fileExt)
	if err != nil {
		return nil, err
	}

	reader.Seek(0, io.SeekStart)
	if _, err = io.Copy(tmpFile, reader); err != nil {
		return nil, err
	}

	image, err := imaging.Open(tmpFile.Name(), imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	dst := imaging.Fill(image, width, width, imaging.Center, imaging.Lanczos)
	return tmpFile, imaging.Save(dst, tmpFile.Name())
}

func GetImageFromURL(url string) (*os.File, error) {
	if len(url) == 0 {
		return nil, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "img-*.jpg")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}

func GenerateAvatar(initials string, avatarSize int, fontSize float64) (*os.File, error) {
	rect := image.Rect(0, 0, avatarSize, avatarSize)
	img := image.NewRGBA(image.Rect(0, 0, avatarSize, avatarSize))

	draw.Draw(img, rect, &image.Uniform{color.RGBA{R: 0x4A, G: 0xA2, B: 0x48}},
		image.Point{}, draw.Src)

	ctx := freetype.NewContext()
	fg := image.NewUniform(color.White)

	ctx.SetDst(img)
	ctx.SetSrc(fg)

	fontPath := filepath.Join("fonts", "Roboto-Medium.ttf")
	fontFile, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, err
	}

	ctx.SetFont(font)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())

	textWidth, textHeight := getTextSize(initials, font, fontSize)
	pt := freetype.Pt((avatarSize-textWidth)/2, (avatarSize+textHeight)/2)
	if _, err := ctx.DrawString(initials, pt); err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.jpg")
	if err != nil {
		return nil, err
	}

	if err := jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}

	return tmpFile, err
}

func getTextSize(text string, font *truetype.Font, fontSize float64) (int, int) {
	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})

	width := 0
	height := 0

	for _, x := range text {
		bounds, aWidth, ok := face.GlyphBounds(x)
		if !ok {
			return width, height
		}

		width += int(float64(aWidth) / 64)
		runeHeight := int(math.Abs(float64(bounds.Min.Y)) / 64)
		if height < runeHeight {
			height = runeHeight
		}
	}

	return width, height
}
