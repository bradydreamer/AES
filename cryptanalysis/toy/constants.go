package toy

import (
	"github.com/OpenWhiteBox/primitives/encoding"
	"github.com/OpenWhiteBox/primitives/gfmatrix"
	"github.com/OpenWhiteBox/primitives/matrix"
	"github.com/OpenWhiteBox/primitives/number"
)

var (
	// subBytes is the linear part of of AES's SubBytes transformation.
	subBytes = encoding.NewByteLinear(matrix.Matrix{
		matrix.Row{0xF1}, // 0b11110001
		matrix.Row{0xE3}, // 0b11100011
		matrix.Row{0xC7}, // 0b11000111
		matrix.Row{0x8F}, // 0b10001111
		matrix.Row{0x1F}, // 0b00011111
		matrix.Row{0x3E}, // 0b00111110
		matrix.Row{0x7C}, // 0b01111100
		matrix.Row{0xF8}, // 0b11111000
	})

	// blocks maps an 8-by-8 matrix to its MixColumns coefficient, c. The matrix is $[c] \circ subBytes$ (from above).
	// This is used for compressing affine layers to matrices over F_{2^8}.
	blocks = map[string]number.ByteFieldElem{
		"0000000000000000": 0,
		"f1e3c78f1f3e7cf8": 1,
		"f809e33f771f3e7c": 2,
		"09ea24b068214284": 3,
	}

	// smallRound is a compressed round matrix, without ShiftRows (so essentially just MixColumns four times). This is
	// used in the search algorithm for unpermuting affine layers.
	smallRound = gfmatrix.Matrix{
		gfmatrix.Row{2, 3, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{1, 2, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{1, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{3, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 2, 3, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 1, 2, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 1, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 3, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 2, 3, 1, 1, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 1, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 2, 3, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 3, 1, 1, 2, 0, 0, 0, 0},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 3, 1, 1},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 1},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 2, 3},
		gfmatrix.Row{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 1, 1, 2},
	}

	// round is the linear part of an affine layer of AES.
	round = encoding.NewBlockLinear(matrix.Matrix{
		matrix.Row{0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1},
		matrix.Row{0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3},
		matrix.Row{0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7},
		matrix.Row{0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f},
		matrix.Row{0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f},
		matrix.Row{0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e},
		matrix.Row{0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c},
		matrix.Row{0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8},
		matrix.Row{0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1},
		matrix.Row{0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3},
		matrix.Row{0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7},
		matrix.Row{0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f},
		matrix.Row{0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f},
		matrix.Row{0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e},
		matrix.Row{0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c},
		matrix.Row{0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8},
		matrix.Row{0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09},
		matrix.Row{0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea},
		matrix.Row{0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24},
		matrix.Row{0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0},
		matrix.Row{0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68},
		matrix.Row{0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21},
		matrix.Row{0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42},
		matrix.Row{0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84},
		matrix.Row{0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8},
		matrix.Row{0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09},
		matrix.Row{0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3},
		matrix.Row{0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f},
		matrix.Row{0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77},
		matrix.Row{0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f},
		matrix.Row{0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e},
		matrix.Row{0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c},
		matrix.Row{0x00, 0x00, 0x00, 0xf1, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xe3, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xc7, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x8f, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x1f, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x3e, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x7c, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xf8, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xf1, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xe3, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xc7, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x8f, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x1f, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x3e, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x7c, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xf8, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x09, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xea, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x24, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xb0, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x68, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x21, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x42, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x84, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xf8, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x09, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0xe3, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x3f, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x77, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x1f, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x3e, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00},
		matrix.Row{0x00, 0x00, 0x00, 0x7c, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00},
		matrix.Row{0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00},
		matrix.Row{0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00},
		matrix.Row{0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0xf8, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x09, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0xe3, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x3f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x77, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x1f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x3e, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x7c, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xf1, 0xf1, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0xe3, 0xe3, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x00, 0xc7, 0xc7, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x8f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x1f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x3e, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x7c, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0xf8, 0xf8, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x09, 0xf1, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0xea, 0xe3, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x24, 0xc7, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x8f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x00, 0x68, 0x1f, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x21, 0x3e, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x42, 0x7c, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x84, 0xf8, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x09, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x00, 0x00, 0x00, 0x00, 0x09, 0xea, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xc7, 0x00, 0x00, 0x00, 0x00, 0xe3, 0x24, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x8f, 0x00, 0x00, 0x00, 0x00, 0x3f, 0xb0, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x77, 0x68, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x21, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x42, 0x00, 0x00, 0x00},
		matrix.Row{0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x7c, 0x84, 0x00, 0x00, 0x00},
	})
)
