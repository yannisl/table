package main

import (
	"fmt"
	//"strings"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"os"
)

func ExampleNew() {
	df := dataframe.New(
		series.New([]string{"b", "a"}, series.String, "COL.1"),
		series.New([]int{1, 2}, series.Int, "COL.2"),
		series.New([]float64{3.0, 4.0}, series.Float, "COL.3"),
	)
	fmt.Println(df)
}

func ExampleReadCSV() {
	counter := 0.00

	f, _ := os.Open("ps.csv")
	defer f.Close()
	df := dataframe.ReadCSV(f, dataframe.DetectTypes(true))
	fmt.Println(df)

	psnames := df
	psnames = df.Select([]int{0, 1})
	for i := 0; i < psnames.Nrow(); i++ {
		fmt.Printf("%s\t& %s\\\\\n", psnames.Elem(i, 0), psnames.Elem(i, 1))
	}

	return
	sub := df.Subset([]int{0, 1})
	sel := df.Select([]int{0, 3})
	sel = df.Select([]int{0, 3, 4})
	fmt.Println(df)

	fmt.Println(df.Names())

	fmt.Println(sub)

	fmt.Println(sel)

	for i := 1; i < df.Nrow(); i++ {
		fmt.Println(df.Elem(i, 1), df.Elem(i, 3))
		counter = counter + df.Elem(i, 3).Float() + df.Elem(i, 4).Float()
		fmt.Println(counter)
	}

	df = dataframe.LoadRecords(
		[][]string{
			[]string{"A", "B", "C", "D"},
			[]string{"a", "4", "5.1", "true"},
			[]string{"b", "4", "6.0", "true"},
			[]string{"c", "3", "6.0", "false"},
			[]string{"a", "2", "7.1", "false"},
		},
	)
	sorted := df.Arrange(
		dataframe.Sort("A"),
		dataframe.Sort("D"))
	fmt.Println(sorted)

}

func main() { //ExampleNew()

	ExampleReadCSV()
}
