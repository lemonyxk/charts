/**
* @program: charts
*
* @description:
*
* @author: lemo
*
* @create: 2022-06-02 05:41
**/

package main

import (
	"fmt"

	"github.com/lemonyxk/charts"
	"github.com/lemonyxk/charts/example/data"
	"github.com/olekukonko/ts"
)

func main() {

	var size, err = ts.GetSize()
	if err != nil {
		panic(err)
	}
	//
	// println(strings.Repeat("︿", size.Col()/text.RuneCount("︿")))

	var t, p = data.GetData()

	_, _ = t, p
	// res[0] = -958

	// res[0] = 60
	// res[26] = 90
	//
	// graph := asciigraph.Plot(
	// 	res,
	// 	asciigraph.Width(size.Col()-8),
	// 	asciigraph.Height(size.Row()-3),
	// 	// asciigraph.Caption(),
	// )
	//
	// fmt.Print(graph)
	//
	// fmt.Println()

	var t1 = []string{"hello", "world", "lemo", "hah", "xixix"}
	var t2 = []float64{1, 2, 3, 4, 5}
	var t3 = []int{1, 2}
	var t4 = []int8{1}
	var t5 = []string{"9:30", "13:00", "15:00"}

	_, _, _, _, _ = t1, t2, t3, t4, t

	var l = charts.New(t5, p, 240)

	l.RenderXBorder = func(isEmpty bool, x string) string {
		return "━"
	}

	l.RenderSymbol = func(lastValue float64, isLastEmpty bool, value float64, isEmpty bool, symbol string) string {
		return "┃"
	}

	l.RenderEmpty = func(lastValue float64, isLastEmpty bool, value float64, isEmpty bool, empty string) string {
		return " "
	}

	l.SetYPrecision(2)

	l.SetSize(size.Col(), size.Row())

	var r = l.Render()

	_ = r

	fmt.Println(r)

}
