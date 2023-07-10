package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png" // 引入 PNG 解码器
	"math"
	"os"
)

func dctII2(x [][]float64) [][]float64 {
	N := len(x)
	M := len(x[0])
	X := make([][]float64, N)

	for u := 0; u < N; u++ {
		X[u] = make([]float64, M)
		for v := 0; v < M; v++ {
			var sum float64
			for i := 0; i < N; i++ {
				for j := 0; j < M; j++ {
					sum += x[i][j] * math.Cos((2*float64(i)+1)*float64(u)*math.Pi/2/float64(N)) * math.Cos((2*float64(j)+1)*float64(v)*math.Pi/2/float64(M))
				}
			}
			X[u][v] = sum
		}
	}
	return X
}

// dctII implements the discrete cosine transform (DCT-II).
func dctII(x [][]float64) [][]float64 {
	N := len(x)
	X := make([][]float64, N)

	for i := 0; i < N; i++ {
		X[i] = make([]float64, N)
		for u := 0; u < N; u++ {
			for v := 0; v < N; v++ {
				for i := 0; i < N; i++ {
					for j := 0; j < N; j++ {
						X[u][v] += x[i][j] * math.Cos((2*float64(i)+1)*float64(u)*math.Pi/16) * math.Cos((2*float64(j)+1)*float64(v)*math.Pi/16)
					}
				}
				X[u][v] *= 0.25 * (func() float64 {
					if u == 0 {
						return 1.0 / math.Sqrt(2)
					}
					return 1.0
				}()) * (func() float64 {
					if v == 0 {
						return 1.0 / math.Sqrt(2)
					}
					return 1.0
				}())
			}
		}
	}
	return X
}

func normalizeAndSave(dct [][]float64, filename string) {
	minVal, maxVal := dct[0][0], dct[0][0]
	for _, row := range dct {
		for _, val := range row {
			minVal = math.Min(minVal, val)
			maxVal = math.Max(maxVal, val)
		}
	}

	width, height := len(dct[0]), len(dct)
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			normalized := (dct[y][x] - minVal) / (maxVal - minVal) // Normalize to [0, 1]
			pixelVal := uint8(normalized * 255)                    // Map to [0, 255]
			img.Set(x, y, color.Gray{pixelVal})
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func normalize(dct [][]float64) [][]uint8 {
	minVal, maxVal := dct[0][0], dct[0][0]
	for _, row := range dct {
		for _, val := range row {
			minVal = math.Min(minVal, val)
			maxVal = math.Max(maxVal, val)
		}
	}

	width, height := len(dct[0]), len(dct)
	img := make([][]uint8, height)
	for y := 0; y < height; y++ {
		img[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			normalized := (dct[y][x] - minVal) / (maxVal - minVal) // Normalize to [0, 1]
			img[y][x] = uint8(normalized * 255)                    // Map to [0, 255]
		}
	}
	return img
}

func main() {
	// 读取图像
	imgFile, err := os.Open("/Users/firshme/Desktop/work/go-png2ascii/ascii/img.png")
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	// 对图像的每个 8x8 块进行 DCT
	bounds := img.Bounds()
	dctImage := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 8 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 8 {
			// 提取 8x8 块
			block := make([][]float64, 8)
			for i := range block {
				block[i] = make([]float64, 8)
			}
			if len(block) != 8 || len(block[0]) != 8 {
				fmt.Println("Block size must be 8x8")
				continue
			}
			for j := 0; j < 8; j++ {
				for i := 0; i < 8; i++ {
					r, _, _, _ := img.At(x+i, y+j).RGBA()
					// 像素值需要归一化到 [0, 1]
					block[j][i] = float64(r) / 65535.0
				}
			}
			// 对块进行 DCT
			dct := dctII2(block)
			normalized := normalize(dct)
			// 将 DCT 结果复制到新图像
			for j := 0; j < 8; j++ {
				for i := 0; i < 8; i++ {
					dctImage.SetGray(x+i, y+j, color.Gray{normalized[j][i]})
				}
			}
			fmt.Printf("DCT of block at (%d, %d):\n", x, y)
		}
	}
	// 保存 DCT 结果图像
	outFile, err := os.Create("dct.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, dctImage)
	if err != nil {
		panic(err)
	}
}
