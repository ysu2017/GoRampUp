package matrix

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"strconv"
)

type Comparison int

const (
	Unknown Comparison = iota
	Less
	Equals
	Greater
	LessOrEquals
	GreaterOrEquals
)

type Matrix [][]float64

func NewMatrix(row int, col int) Matrix {
	matrix := make([][]float64, row)
	for i := range matrix {
		matrix[i] = make([]float64, col)
	}
	return Matrix(matrix)
}

func sameSize(a Matrix, b Matrix) bool {
	return a.Col() == b.Col() && a.Row() == b.Row()
}

func deepCopy(m Matrix) Matrix {
	row := m.Row()
	col := m.Col()
	matrix := make([][]float64, row)
	for i := range matrix {
		matrix[i] = make([]float64, col)
		for j := range matrix[i] {
			matrix[i][j] = m[i][j]
		}
	}
	return Matrix(matrix)
}

func NewMatrixFromGreyImage(grey *image.Gray) Matrix {
	bounds := grey.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	matrix := make([][]float64, dx)
	for i := range matrix {
		matrix[i] = make([]float64, dy)
		for j := range matrix[i] {
			matrix[i][j] = float64(grey.GrayAt(i, j).Y)
		}
	}

	return Matrix(matrix)
}

func Ones(row int, col int) Matrix {
	matrix := make([][]float64, row)
	for i := range matrix {
		matrix[i] = make([]float64, col)
		for j := range matrix[i] {
			matrix[i][j] = 1
		}
	}
	return Matrix(matrix)
}

func Add(a Matrix, matrices ...Matrix) (Matrix, error) {
	result := deepCopy(a)
	if len(matrices) == 0 {
		return result, nil
	}

	col := a.Col()
	row := a.Row()

	for _, m := range matrices {
		if !sameSize(m, a) {
			return nil, fmt.Errorf("added matrices does not have same dimention")
		}

		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				result[i][j] = result[i][j] + m[i][j]
			}
		}
	}

	return result, nil
}

func ElementWiseMultiplication(a Matrix, matrices ...Matrix) (Matrix, error) {
	result := deepCopy(a)
	if len(matrices) == 0 {
		return result, nil
	}

	col := a.Col()
	row := a.Row()

	for _, m := range matrices {
		if !sameSize(m, a) {
			return nil, fmt.Errorf("element wise multiplicated matrices does not have same dimention")
		}

		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				result[i][j] = result[i][j] * m[i][j]
			}
		}
	}

	return result, nil
}

func (m Matrix) Row() int {
	return len(m)
}

func (m Matrix) Col() int {
	if len(m) == 0 {
		return 0
	}

	return len(m[0])
}

func (m Matrix) String() string {
	row := m.Row()
	col := m.Col()
	result := ""
	for i := 0; i < row; i++ {
		row := "["
		for j := 0; j < col; j++ {
			row += strconv.FormatFloat(m[i][j], 'f', -1, 64) + " "
		}
		result += row + "]\n"
	}
	return result
}

func (m Matrix) Ctranspose() Matrix {
	newCol := m.Row()
	newRow := m.Col()
	newMatrix := NewMatrix(newRow, newCol)
	for i := 0; i < newRow; i++ {
		for j := 0; j < newCol; j++ {
			newMatrix[i][j] = m[j][i]
		}
	}
	return newMatrix
}

func (m Matrix) Scale(s float64) Matrix {
	col := m.Col()
	row := m.Row()
	newMatrix := NewMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			newMatrix[i][j] = m[i][j] * s
		}
	}
	return newMatrix
}

func (m Matrix) Sum() float64 {
	var result float64
	for i := 0; i < m.Row(); i++ {
		for j := 0; j < m.Col(); j++ {
			result += m[i][j]
		}
	}
	return result
}

func (m Matrix) Power(p float64) Matrix {
	row := m.Row()
	col := m.Col()
	newMatrix := NewMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			newMatrix[i][j] = math.Pow(m[i][j], p)
		}
	}
	return newMatrix
}

func (m Matrix) Compare(compared Matrix, symbol Comparison) (Matrix, error) {
	if !sameSize(m, compared) {
		return nil, fmt.Errorf("compared matrix have a different dimension")
	}

	row := m.Row()
	col := m.Col()
	newMatrix := NewMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			switch symbol {
			case Less:
				if m[i][j] < compared[i][j] {
					newMatrix[i][j] = m[i][j]
				} else {
					newMatrix[i][j] = 0
				}
			case Equals:
				if m[i][j] == compared[i][j] {
					newMatrix[i][j] = m[i][j]
				} else {
					newMatrix[i][j] = 0
				}
			case Greater:
				if m[i][j] > compared[i][j] {
					newMatrix[i][j] = m[i][j]
				} else {
					newMatrix[i][j] = 0
				}
			case LessOrEquals:
				if m[i][j] <= compared[i][j] {
					newMatrix[i][j] = m[i][j]
				} else {
					newMatrix[i][j] = 0
				}
			case GreaterOrEquals:
				if m[i][j] >= compared[i][j] {
					newMatrix[i][j] = m[i][j]
				} else {
					newMatrix[i][j] = 0
				}
			default:
				return nil, fmt.Errorf("Invalid comparison symbol")
			}
		}
	}
	return newMatrix, nil
}

func (m Matrix) subMatrix(x int, y int, halfSize int) (Matrix, error) {
	col := m.Col()
	row := m.Row()
	if x-halfSize < 0 || x+halfSize >= row || y-halfSize < 0 || y+halfSize >= col {
		return nil, fmt.Errorf("Error creating sub matrix")
	}
	matrix := NewMatrix(2*halfSize+1, 2*halfSize+1)
	for i := x - halfSize; i <= x+halfSize; i++ {
		for j := y - halfSize; j <= y+halfSize; j++ {
			matrix[i-x+halfSize][j-y+halfSize] = m[i][j]
		}
	}
	return matrix, nil
}

func (m Matrix) Filter(filter Matrix) (Matrix, error) {
	col := m.Col()
	row := m.Row()
	filterCol := filter.Col()
	filterRow := filter.Row()
	if col < filterCol || row < filterRow {
		return nil, fmt.Errorf("filter row %d and col %d is larger than matrix row %d and col %d",
			filterRow, filterCol, row, col)
	}

	if filterCol != filterRow {
		return nil, fmt.Errorf("filter row must be equal to its col, but got row %d, col %d", filterRow, filterCol)
	}

	if filterCol%2 != 1 {
		return nil, fmt.Errorf("filter row and col must be a odd number, currently got %d", filterCol)
	}

	halfSize := filterCol / 2
	newMatrix := deepCopy(m)

	channel := make(chan int, runtime.NumCPU())
	errChannel := make(chan error)
	for i := halfSize; i < row-halfSize; i++ {
		for j := halfSize; j < col-halfSize; j++ {
			channel <- 1
			go func(i, j int) {
				var err error
				defer func() {
					if err != nil {
						errChannel <- err
					}
				}()
				defer func() { <-channel }()
				computedSubMatrix, err := m.subMatrix(i, j, halfSize)
				if err != nil {
					return
				}
				matrixProduct, err := ElementWiseMultiplication(filter, computedSubMatrix)
				if err != nil {
					return
				}
				newMatrix[i][j] = matrixProduct.Sum()
			}(i, j)

			select {
			case err := <-errChannel:
				return nil, err
			default:
				// do nothing
			}
		}
	}
	return newMatrix, nil
}
