package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMatrixReturnCorrectMatrix(t *testing.T) {
	row := 8
	col := 10
	m := NewMatrix(row, col)
	expectedValue := make([][]float64, row)
	for i := 0; i < row; i++ {
		expectedValue[i] = make([]float64, col)
	}
	assert.EqualValues(t, m, Matrix(expectedValue))
}

func Test_sameSizeReturnTrueIfMatrixHasSameSize(t *testing.T) {
	a := NewMatrix(2, 3)
	b := NewMatrix(2, 3)
	c := NewMatrix(3, 3)
	assert.True(t, sameSize(a, b))
	assert.False(t, sameSize(a, c))
}

func Test_deepCopyCreateADeepCopyOfMatrix(t *testing.T) {
	a := NewMatrix(3, 5)
	aClone := deepCopy(a)
	assert.Exactly(t, a, aClone)
}

func Test_AddReturnSumOfMatrices(t *testing.T) {
	a := Matrix{
		[]float64{1, 2},
		[]float64{3, 4},
	}
	b := Matrix{
		[]float64{1, 2},
		[]float64{3, 4},
	}
	c := Matrix{
		[]float64{2, 4},
		[]float64{6, 8},
	}
	d := Matrix{
		[]float64{2, 4},
	}
	sum, err := Add(a, b)
	assert.NoError(t, err)
	assert.EqualValues(t, c, sum)

	_, err = Add(a, d)
	assert.Error(t, err)
}
