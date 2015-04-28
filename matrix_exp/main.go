package main

import (
	"fmt"
	"os"

	"github.com/vladvelici/graph-dataset-tools/sim"
)

func main() {
	outputPath := "../../output.txt"
	inputPath := "/Users/vlad/Projects/uni/project/datasets/fb/egonets/test/a_g_fb.csv"
	mu := 0.5
	k := 20

	sim.Path = "./sim/" + sim.Path
	err := sim.EigenRaw(inputPath, outputPath, mu, k)
	if err != nil {
		fmt.Println(err)
		return
	}
	q, z, err := sim.ParseEigenOutput("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, cols := q.Dims()

	o, _ := os.Create("q.txt")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Fprintf(o, "%.4f\t", q.At(i, j))
		}
		fmt.Fprintln(o)
	}
	o.Close()

	//		fmt.Println(z)

	r := sim.FromQZ(q, z)

	fmt.Println("Q(1,2) golang:", q.At(1, 2))
	fmt.Println("Q(2,3) golang:", q.At(2, 3))
	fmt.Println("Q(3,4) golang:", q.At(3, 4))

	dst := r.DistanceSim(2, 3)
	fmt.Println("sim 2,3 golang: ", dst)

}
