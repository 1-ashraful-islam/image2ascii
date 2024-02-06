package im2ascii

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

func CreateASCIIImage(img image.Image, width, height int) ([]string, error) {

	// Resize the image to half the height as we're going to use half-blocks to represent two vertical pixels.
	resizedImg := resizeImage(img, width, height) // Adjust width and height as needed

	asciiImage := make([]string, 0)

	for y := resizedImg.Bounds().Min.Y; y < resizedImg.Bounds().Max.Y; y += 2 {
		var line string
		for x := resizedImg.Bounds().Min.X; x < resizedImg.Bounds().Max.X; x++ {
			topPixelColor := resizedImg.At(x, y)
			bottomPixelColor := color.RGBA{255, 255, 255, 255} // Default to white if bottom pixel exceeds image bounds
			if y+1 < resizedImg.Bounds().Max.Y {
				bottomPixelColor, _ = resizedImg.At(x, y+1).(color.RGBA)
			}

			// Extract RGBA components for top and bottom pixels
			r1, g1, b1, _ := topPixelColor.RGBA()
			r2, g2, b2, _ := bottomPixelColor.RGBA()

			// Print the half-block with top color as foreground and bottom color as background
			line += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m", r1>>8, g1>>8, b1>>8, r2>>8, g2>>8, b2>>8)
		}
		asciiImage = append(asciiImage, line)
	}
	return asciiImage, nil
}

func LoadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func resizeImage(src image.Image, newWidth, newHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Calculate the aspect ratio of the source image
	srcAspectRatio := float64(srcWidth) / float64(srcHeight)

	// Calculate the target width and height while maintaining the aspect ratio
	var targetWidth, targetHeight int
	if float64(newWidth)/float64(newHeight) > srcAspectRatio {
		targetWidth = int(float64(newHeight) * srcAspectRatio)
		targetHeight = newHeight
	} else {
		targetWidth = newWidth
		targetHeight = int(float64(newWidth) / srcAspectRatio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(dst, dst.Rect, src, srcBounds, draw.Over, nil)
	return dst
}
