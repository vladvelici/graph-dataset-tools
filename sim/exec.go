package sim

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gonum/matrix/mat64"
)

var Path = "mlscript/"
var ScriptName = "train.sh"

// Run the script at Path with arguments: inputPath, outputPath, mu, k.
func EigenRaw(inputPath, outputPath string, mu float64, k int) error {
	attr := &os.ProcAttr{
		Dir:   Path,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	proc, err := os.StartProcess(ScriptName, []string{ScriptName, inputPath, fmt.Sprint(mu), strconv.Itoa(k), outputPath}, attr)
	if err != nil {
		return err
	}

	state, err := proc.Wait()
	if err != nil {
		return err
	}

	if !state.Success() {
		return fmt.Errorf("Matlab (eigen) process finished with errors (non-zero code).")
	}

	return nil
}

// Run the script, generate a temporary file for output, parse it, delete it.
func Eigen(inputPath string, mu float64, k int) (*mat64.Dense, *mat64.Dense, error) {
	file, err := ioutil.TempFile("", "eigen_output")
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		os.Remove(file.Name())
	}()
	outputPath := file.Name()
	err = file.Close()
	if err != nil {
		return nil, nil, err
	}

	err = EigenRaw(inputPath, outputPath, mu, k)
	if err != nil {
		return nil, nil, err
	}

	return ParseEigenOutput(outputPath)
}

// Parse the output file of the eigen algorithm.
// The format is:
// Qr Qc
// q1 q2 q3 ...
// Zr Zc
// z1 z2 z3 ...
//
// Where:
// Qr - rows of Q;
// Qc - columns of Q;
// q1, q1, qn - elements of Q, ordered as:
//
// Matrix:
// 1 2
// 3 4
// 5 6
//
// Is written as:
// 1 3 5 2 4 5
func ParseEigenOutput(path string) (*mat64.Dense, *mat64.Dense, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	var qr, qc, zr, zc int
	_, err = fmt.Fscanf(f, "%d %d ", &qr, &qc)
	if err != nil {
		return nil, nil, err
	}

	qdata := make([]float64, qr*qc)
	for i := 0; i < len(qdata); i++ {
		var tmp float64
		_, err = fmt.Fscanf(f, "%f ", &tmp)
		fmt.Print(tmp, "\t")
		qdata[i] = tmp
		if err != nil {
			return nil, nil, err
		}
	}

	q := mat64.NewDense(qr, qc, qdata)

	_, err = fmt.Fscanf(f, "%d %d ", &zr, &zc)
	if err != nil {
		return nil, nil, err
	}

	zdata := make([]float64, zr*zc)
	for i := 0; i < len(zdata); i++ {
		_, err = fmt.Fscanf(f, "%f ", &zdata[i])
		if err != nil {
			return nil, nil, err
		}
	}

	z := mat64.NewDense(zr, zc, zdata)
	return q, z, nil
}
