package sim

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

// In result, unlike in Matlab, the node IDs start from 0.
// In matlab and the CSV files provided, they start from 1.
type Result struct {
	q *mat64.Dense
	z *mat64.Dense
}

// FromQZ creates a result object form a q and a z matrix.
// Obtain those matrices from Eigen().
func FromQZ(q, z *mat64.Dense) *Result {
	return &Result{
		q: q,
		z: z,
	}
}

// Len returns the number of nodes in the graph behind this Result.
func (r *Result) Len() int {
	rows, _ := r.z.Dims()
	return rows
}

func printdim(a mat64.Matrix, name string) {
	rows, cols := a.Dims()
	fmt.Printf("%s (%d, %d)\n", name, rows, cols)
}

// Multiplies a * m * b'.
func multipl(a *mat64.Vector, m mat64.Matrix, b *mat64.Vector) float64 {
	printdim(a, "a")
	printdim(m, "m")
	printdim(b, "b")
	p1 := new(mat64.Dense)
	p1.MulTrans(a, true, m, false)
	res := new(mat64.Dense)
	printdim(p1, "p1")
	res.Mul(p1, b)
	return res.At(0, 0)
}

// Sim returns the similarity metric between two nodes.
func (r *Result) DistanceSim(from, to int) float64 {
	fromRow := r.z.RowView(from)
	toRow := r.z.RowView(to)

	// matlab code to port:
	// norma = z(a,:)*q*z(a,:)';
	// normb = z(b,:)*q*z(b,:)';
	// similarity = norma + normb - 2 * (z(a,:)*q*z(b,:)');

	similarity := multipl(fromRow, r.q, fromRow) + multipl(toRow, r.q, toRow) + 2*multipl(fromRow, r.q, toRow)
	return similarity
}

// Could have more similarity measures here...
