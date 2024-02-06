package im2ascii

import (
	"fmt"
	"image"
	"math"
	"testing"
)

func TestLoadImage(t *testing.T) {
	imagePath := "examples/pikachu.png"
	_, err := LoadImage(imagePath)
	if err != nil {
		t.Errorf("Failed to load image: %v", err)
	}
}

func TestLoadImageInvalidPath(t *testing.T) {
	imagePath := "examples/invalid.png"
	_, err := LoadImage(imagePath)
	if err == nil {
		t.Errorf("Failed to handle invalid image path")
	}
}

func TestResizeImage(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	newWidth := 50
	newHeight := 50
	resizedImage := resizeImage(src, newWidth, newHeight)
	if resizedImage.Bounds().Dx() != newWidth || resizedImage.Bounds().Dy() != newHeight {
		t.Errorf("Failed to resize image to the expected dimensions")
	}
}

func TestResizeImageAspectRatio(t *testing.T) {
	aspectRatios := []struct {
		width  int
		height int
	}{
		{100, 50},
		{50, 100},
	}

	tolerance := 0.001 // Define a tolerance value for comparing aspect ratios

	for _, aspectRatio := range aspectRatios {
		t.Run(fmt.Sprintf("Width: %d, Height: %d", aspectRatio.width, aspectRatio.height), func(t *testing.T) {
			src := image.NewRGBA(image.Rect(0, 0, aspectRatio.width, aspectRatio.height))
			newWidth := 50
			newHeight := 50
			resizedImage := resizeImage(src, newWidth, newHeight)
			aspectRatioExpected := float64(aspectRatio.width) / float64(aspectRatio.height)
			aspectRatioActual := float64(resizedImage.Bounds().Dx()) / float64(resizedImage.Bounds().Dy())
			if math.Abs(aspectRatioActual-aspectRatioExpected) > tolerance {
				t.Errorf("Failed to resize image to the expected aspect ratio")
			}
		})
	}
}

func TestPrintImageASCII(t *testing.T) {
	imagePath := "examples/pikachu.png"
	width := 80
	height := 80
	img, err := LoadImage(imagePath)
	if err != nil {
		t.Errorf("Failed to load image: %v", err)
	}

	if asciiImage, err := CreateASCIIImage(img, width, height); err != nil {
		t.Errorf("Failed to print image: %v", err)
	} else {
		for _, line := range asciiImage {
			fmt.Println(line)
		}
	}

}
