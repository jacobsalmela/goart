/*
MIT License

Copyright (c) 2021 Jacob Salmela <me@jacobsalmela.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package cmd

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jacobsalmela/goart/sketch"
	"github.com/spf13/cobra"
)

var (
	sourceImgName                                                                                     string
	outputImgName                                                                                     string
	strokeRatio, initialAlpha, strokeReduction, alphaIncrease, strokeInversionThreshold, strokeJitter float64
	destWidth, destHeight, minEdge, maxEdge, totalCycleCount                                          int
	keepSourceDimensions                                                                              bool
)

// generateCmd generates a new image from a source image
var generateCmd = &cobra.Command{
	Use:   "generate FILE",
	Short: "Generate a art from a source image",
	Long: `Generates art using a source image as the starting point.
	
	The generated image can be further manipulated with flags to this
	command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// The source image is the first argument
		sourceImgName = args[0]
		// Get just the filename
		sourceBasename := filepath.Base(sourceImgName)
		// Trim the extension and append the new extension
		sourceBasename = strings.TrimSuffix(sourceBasename, filepath.Ext(sourceBasename)) + "-goart.png"
		// Save to the directry of the source image by default
		sourceDir := filepath.Dir(sourceImgName)
		// Save to the same director as the source image
		outputImgName = filepath.Join(sourceDir, sourceBasename)
		fmt.Println("Generating...", outputImgName)
		// Generate the art
		err := generateArt(sourceImgName, outputImgName)
		if err != nil {
			log.Panicln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.DisableAutoGenTag = true

	generateCmd.Flags().Float64VarP(&strokeRatio, "stroke-ratio", "r", 0.75, "size of the initial stroke compared to that of the final result")
	generateCmd.Flags().BoolVarP(&keepSourceDimensions, "keep-source-dimensions", "K", true, "generate a new image with the same dimensions as the source image")
	generateCmd.Flags().IntVarP(&destWidth, "width", "W", 2000, "width of the generated image")
	generateCmd.Flags().IntVarP(&destHeight, "height", "H", 2000, "height of the generated image")
	generateCmd.Flags().Float64VarP(&initialAlpha, "initial-alpha", "a", 0.1, "beginning stroke transparency")
	generateCmd.Flags().Float64VarP(&strokeReduction, "stroke-reduction", "R", 0.002, "the initial stroke size gets minimized by this amount on each iteration")
	generateCmd.Flags().Float64VarP(&alphaIncrease, "alpha-increase", "A", 0.06, "the step of transparency increase at each iteration")
	generateCmd.Flags().Float64VarP(&strokeInversionThreshold, "stroke-inversion-threshold", "t", 0.05, "the minimum stroke size")
	generateCmd.Flags().Float64VarP(&strokeJitter, "stroke-jitter", "j", 0.1, "deviation of the colored stroke from its projected position in the original image")
	generateCmd.Flags().IntVarP(&minEdge, "min-edge", "m", 3, "minimum stroke will be a n-edge polygon (3 is a triangle)")
	generateCmd.Flags().IntVarP(&maxEdge, "max-edge", "M", 4, "maximum stroke will be a n-edge polygon (4 is a square)")
	generateCmd.Flags().IntVarP(&totalCycleCount, "total-cycles", "T", 5000, "Copy any discovered k8s squashfs images from SRC to DEST")
}

// generateArt generates a new image from a source image
func generateArt(sourcePath string, destPath string) error {

	// load the source image
	img, err := loadImage(sourcePath)

	if err != nil {
		log.Panicln(err)
	}

	// if the user wants to keep the source dimensions,
	if keepSourceDimensions {
		// set them appropriately
		destWidth = img.Bounds().Dx()
		destHeight = img.Bounds().Dy()
	}

	// StrokeRatio determines the size of the initial stroke compared
	// to that of the final result. Our idea is to start big and reduce
	// the stroke size upon each iteration. This is why an initial
	// StrokeRatio of 3/4, or 0.75 should be a reasonable default value
	//
	// DestWidth and DestHeight are the dimensions of our end sketch.
	// Regardless of how large or small our source image is, we can
	// produce a sketch of arbitrary dimensions, as long as our
	// computers have enough capacity to render it.
	//
	// InitialAlpha controls the stroke transparency at the start.
	// To achieve the desired effect, stroke transparency is low at
	// the beginning and increases at every step. Alpha is a floating
	// point number that varies between 0 and 255. We can set it to a
	// pretty low value: 0.1
	//
	// StrokeReduction is directly related to StrokeRatio. The reduction
	// is the step, with which the initial stroke size gets minimized
	// upon each loop iteration
	//
	// AlphaIncrease is the step of transparency increase at each
	// iteration. This one may differ from setup to setup, so feel
	// free to play with it. In this case, I have set it to 0.06.
	//
	// StrokeInversionThreshold determines the minimum stroke size,
	// beyond which we start adding borders around the stroke, for
	// more contrast. To deepen the contrast even more. borders will
	// be the negative of the brightness of the stroke’s color. A
	// bright stroke will get a dark border, and vice versa. To make
	// and effect of depth-of-field in my sketches, I keep this value
	// low - 0.05, so that the algorithm starts adding borders only
	// when the stroke is at 5% of its original size, but not before.
	//
	// StrokeJitter is the deviation of our colored stroke from its
	// projected position in the original image. If you recall the
	// algorithm, a color gets picked at certain coordinates (X, Y)
	// in the original image an a polygon with that color gets “drawn”
	// onto the destination. The jitter helps to make the final image
	// more unpredictable. By adding an element of randomness, the
	// jitter would shake the hand of our algorithm, so that the
	// stroke ends near where it should have been in the original
	// image, but not exactly. I have computed this value to be a
	// certain portion of the output image’s width:
	// int(0.1 * float64(destWidth))
	//
	// MinEdgeCount and MaxEdgeCount are the minimum and maximum number
	// of edges of each stroke, respectively. Each stroke will be a
	// n-edge polygon. The defaults here I leave totally to the
	// imagination of the reader. If you prefer simpler shapes, you can
	// play with drawing just triangles and squares (3, and 4,
	// respectively), or go for more complex shapes.
	//
	s := sketch.NewSketch(img, sketch.UserParams{
		StrokeRatio:              strokeRatio,
		DestWidth:                destWidth,
		DestHeight:               destHeight,
		InitialAlpha:             initialAlpha,
		StrokeReduction:          strokeReduction,
		AlphaIncrease:            alphaIncrease,
		StrokeInversionThreshold: strokeInversionThreshold,
		StrokeJitter:             int(strokeJitter * float64(destWidth)),
		MinEdgeCount:             minEdge,
		MaxEdgeCount:             maxEdge,
	})

	rand.Seed(time.Now().Unix())

	// for each iteration defined,
	for i := 0; i < totalCycleCount; i++ {
		// update the sketch
		s.Update()
	}

	// save the sketch once it is generated
	err = saveOutput(s.Output(), destPath)

	if err != nil {
		return err
	}

	return nil
}

// loadImage loads an image from the given path
func loadImage(src string) (image.Image, error) {
	file, _ := os.Open(sourceImgName)
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

// saveOutput saves the output image to the specified path
func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level then save to a file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
