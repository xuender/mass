package app

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/xuender/kit/los"
	"github.com/xuender/kit/types"
)

func Csv(titles []string, data []map[string]any) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	los.Must0(writer.Write(titles))

	for _, row := range data {
		val := make([]string, len(titles))

		for num, key := range titles {
			col := row[key]
			if col != nil {
				val[num] = fmt.Sprintf("%v", row[key])
			}
		}

		los.Must0(writer.Write(val))
	}
}

func Grid(titles []string, data []map[string]any) {
	titles = append([]string{"#"}, titles...)

	vals := [][]string{}

	for idx, row := range data {
		val := make([]string, len(titles))

		for num, key := range titles {
			if num == 0 {
				val[0] = types.Itoa(idx + 1)
			} else {
				val[num] = fmt.Sprintf("%v", row[key])
			}
		}

		vals = append(vals, val)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(titles)
	// footer := make([]string, len(titles))
	// footer[0] = fmt.Sprintf("%d", len(data))
	// table.SetFooter(footer)
	table.AppendBulk(vals)
	table.SetBorder(false)
	table.Render()
}
