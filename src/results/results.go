package results

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func SampleTable() {
	data := [][]string{
		[]string{"A", "The Good", "500"},
		[]string{"B", "The Very very Bad Man", "288"},
		[]string{"C", "The Ugly", "120"},
		[]string{"D", "The Gopher", "800"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func DrawTable(data [][]string) {
	header := []string{
		"Name",
		"Environment",
		"Instance Id",
		"Private Ip",
		"Key Name",
		"Launch Time",
		"State",
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, x := range data {
		table.Append(x)
	}
	table.Render()
}
