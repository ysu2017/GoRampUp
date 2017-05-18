package harris

import (
	"image"
	"log"
	"yi/matrix"
	"yi/point"
)

const k = 0.04

func smoothFilter() matrix.Matrix {
	result := matrix.Ones(3, 3).Scale(float64(1) / float64(9))
	return result
}

func gaussianFilter() matrix.Matrix {
	arr := [][]float64{
		[]float64{1, 4, 7, 4, 1},
		[]float64{4, 16, 26, 16, 4},
		[]float64{7, 26, 41, 26, 7},
		[]float64{4, 16, 26, 16, 4},
		[]float64{1, 4, 7, 4, 1},
	}
	newMatrix := matrix.Matrix(arr)
	return newMatrix.Scale(float64(1) / newMatrix.Sum())
}

func imageGradientX() matrix.Matrix {
	arr := [][]float64{
		[]float64{-1, 0, 1},
		[]float64{-1, 0, 1},
		[]float64{-1, 0, 1},
	}
	return matrix.Matrix(arr)
}

func imageGradientY() matrix.Matrix {
	return imageGradientX().Ctranspose()
}

func isLargestNeighbor(m matrix.Matrix, xPos int, yPos int) bool {
	row := m.Row()
	col := m.Col()

	for i := xPos - 1; i <= xPos+1; i++ {
		for j := yPos - 1; j <= yPos+1; j++ {
			if i >= 0 && j >= 0 && i < row && j < col {
				if m[xPos][yPos] < m[i][j] {
					return false
				}
			}
		}
	}

	return true
}

// compute the image gradient
func imageGradient(img matrix.Matrix) (matrix.Matrix, matrix.Matrix, error) {
	imgX, err := img.Filter(imageGradientX())
	if err != nil {
		log.Println("Error Filter image gradient x")
		return nil, nil, err
	}
	imgY, err := img.Filter(imageGradientY())
	if err != nil {
		log.Println("Error Filter image gradient y")
		return nil, nil, err
	}
	return imgX, imgY, nil
}

// compute the product of derivatives and pass them with gaussian filter
func derivatives(imgX matrix.Matrix, imgY matrix.Matrix) (matrix.Matrix, matrix.Matrix, matrix.Matrix, error) {
	ixx, err := matrix.ElementWiseMultiplication(imgX, imgX)
	if err != nil {
		log.Println("Error calculate ixx")
		return nil, nil, nil, err
	}
	iyy, err := matrix.ElementWiseMultiplication(imgY, imgY)
	if err != nil {
		log.Println("Error calculate iyy")
		return nil, nil, nil, err
	}
	ixy, err := matrix.ElementWiseMultiplication(imgX, imgY)
	if err != nil {
		log.Println("Error calculate ixy")
		return nil, nil, nil, err
	}

	gfilter := gaussianFilter()
	sxx, err := ixx.Filter(gfilter)
	if err != nil {
		log.Println("Error gaussian filter ixx")
		return nil, nil, nil, err
	}
	syy, err := iyy.Filter(gfilter)
	if err != nil {
		log.Println("Error gaussian filter iyy")
		return nil, nil, nil, err
	}
	sxy, err := ixy.Filter(gfilter)
	if err != nil {
		log.Println("Error gaussian filter ixy")
		return nil, nil, nil, err
	}

	return sxx, syy, sxy, nil
}

//Compute the R response
func rResponse(sxx matrix.Matrix, syy matrix.Matrix, sxy matrix.Matrix) (matrix.Matrix, error) {
	xyProduct, err := matrix.ElementWiseMultiplication(sxx, syy)
	if err != nil {
		log.Println("Error multiply sxx an syy")
		return nil, err
	}
	xyAdd, err := matrix.Add(sxx, syy)
	if err != nil {
		log.Println("Error adding sxx an syy")
		return nil, err
	}
	r, err := matrix.Add(xyProduct, sxy.Power(2).Scale(-1), xyAdd.Power(2).Scale(-1*k))
	if err != nil {
		log.Println("Error calculate r response")
		return nil, err
	}
	return r, nil
}

//Threshold R response
func applyThreshold(r matrix.Matrix, threshold float64) (matrix.Matrix, error) {
	thresholdMatrix := matrix.Ones(r.Row(), r.Col()).Scale(threshold)
	rThresholded, err := r.Compare(thresholdMatrix, ">")
	if err != nil {
		return nil, err
	}
	return rThresholded, nil
}

//Compute Nonmax_suppression
func nonmaxSuppression(m matrix.Matrix) {
	row := m.Row()
	col := m.Col()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if m[i][j] != 0 && !isLargestNeighbor(m, i, j) {
				m[i][j] = 0
			}
		}
	}
}

func HarrisCornerDetector(grey *image.Gray, threshold float64) (point.Points, error) {
	img := matrix.NewMatrixFromGreyImage(grey)
	img, err := img.Filter(smoothFilter())
	if err != nil {
		return nil, err
	}

	imgX, imgY, err := imageGradient(img)
	if err != nil {
		return nil, err
	}

	sxx, syy, sxy, err := derivatives(imgX, imgY)
	if err != nil {
		return nil, err
	}

	r, err := rResponse(sxx, syy, sxy)
	if err != nil {
		return nil, err
	}

	rThresholded, err := applyThreshold(r, threshold)
	if err != nil {
		return nil, err
	}

	nonmaxSuppression(rThresholded)

	// find all corner points
	var corners point.Points
	row := rThresholded.Row()
	col := rThresholded.Col()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if rThresholded[i][j] != 0 {
				corners = append(corners, point.Point{X: i, Y: j})
			}
		}
	}
	return corners, err
}
