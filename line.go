/**
* @program: charts
*
* @description:
*
* @author: lemo
*
* @create: 2022-06-02 05:56
**/

package charts

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/text"
	"github.com/olekukonko/ts"
)

func New[T comparable](x []T, y []float64) *Line[T] {
	return &Line[T]{X: x, Y: y}
}

type Line[T comparable] struct {
	X []T
	Y []float64

	width  int
	height int

	matrix     [][]int
	xOffset    int
	yOffset    int
	yMaxCount  int
	yMin       float64
	yMax       float64
	size       ts.Size
	xMaxCount  int
	yPrecision int
}

func (l *Line[T]) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *Line[T]) SetYPrecision(precision int) {
	l.yPrecision = precision
}

func (l *Line[T]) Render() string {
	l.init()

	if l.width == 0 || l.height == 0 {
		return ""
	}

	// var xMin, xMax = l.xMinAndMax()
	var yMin, yMax = l.yMinAndMax()

	var yRange = yMax - yMin
	if yRange == 0.0 {
		yRange = 1.0
	}

	l.yMax = yMax
	l.yMin = yMin

	var lY = len(l.Y)

	var xScale = float64(l.width) / float64(lY)

	var mMap = make(map[int]bool)
	var yScale = float64(l.height) / yRange

	for i := 0; i < lY; i++ {
		var x = int(float64(i) * xScale)
		if xScale >= 1 {
			x = i
		}
		var y = int((l.Y[i] - yMin) * yScale)

		if x > l.width-1 {
			x = l.width - 1
		}
		if y > l.height-1 {
			y = l.height - 1
		}

		if mMap[x] {
			continue
		}

		l.matrix[x][y] = y
		mMap[x] = true
	}

	if xScale >= 1 {
		return l.outPut()
	}

	var next = 0
	for i := 0; i < l.width; i++ {
		var n = getNextY(l.matrix[i])
		if n == math.MinInt {
			l.matrix[i][next] = next
		} else {
			next = n
		}
	}

	return l.outPut()
}

// output
func (l *Line[T]) outPut() string {
	var buf bytes.Buffer
	for i := l.height - 1; i >= 0; i-- {
		var count = 0
		for j := 0; j < l.width+l.xOffset-l.yMaxCount; j++ {

			if count == l.yMaxCount {
				j--
				count++
				if i == 0 {
					buf.WriteString("┃")
					continue
				}

				if i%2 == 1 {
					buf.WriteString("┫")
				} else {
					buf.WriteString("┃")
				}
				continue
			}

			if count >= 0 && count < l.yMaxCount {
				if i%2 == 0 {
					var x = l.yMax - l.yMin
					if x == 0 {
						x = l.yMin
					}
					var v = l.yMin + (x)/float64(l.height)*float64(i)
					var s = l.parseFloatToString(v)
					var c = text.RuneCount(s)
					if count >= c {
						buf.WriteString(" ")
					} else {
						buf.WriteString(s[count : count+1])
					}
					j--
					count++
					continue
				} else {
					j--
					count++
					buf.WriteString(" ")
					continue
				}

			}

			if j >= l.width {
				buf.WriteString(" ")
				continue
			}

			if l.matrix[j][i] != math.MinInt {
				buf.WriteString("┃")
			} else {
				var n = getNextY(l.matrix[j])
				if i < n {
					buf.WriteString("┃")
					continue
				}
				buf.WriteString(" ")
			}
		}
	}

	for i := 0; i < l.size.Col(); i++ {
		if i >= 0 && i < l.yMaxCount {
			buf.WriteString(" ")
			continue
		}

		if i >= l.yMaxCount && i < l.width+l.yMaxCount {

			if i == l.yMaxCount {
				buf.WriteString("┗")
			} else {

				var n = math.MinInt
				if i < l.width+l.yMaxCount {
					n = getNextY(l.matrix[i-l.yMaxCount-1])
				}

				if n == math.MinInt {
					if i == len(l.Y)+l.yMaxCount {
						buf.WriteString("┻")
					} else {
						buf.WriteString("━")
					}
				} else {
					buf.WriteString("┻")
				}
			}
			continue
		}

		if i == l.width+l.yMaxCount && l.xOffset != -1 {
			var n = math.MinInt
			n = getNextY(l.matrix[i-l.yMaxCount-1])

			if n == math.MinInt {
				buf.WriteString("━")
			} else {
				buf.WriteString("┻")
			}
			continue
		} else {
			buf.WriteString(" ")
		}

	}

	var count = 0

	var spaceWidth = (l.width - l.xMaxCount*len(l.X)) / (len(l.X) - 1)

	for i := 0; i < l.size.Col(); i++ {
		if i >= 0 && i < l.yMaxCount {
			buf.WriteString(" ")
			continue
		}

		if i >= l.yMaxCount && i < l.width+l.yMaxCount {

			var index = i - l.yMaxCount

			if count > len(l.X)-1 {
				continue
			}

			var s = fmt.Sprintf("%v", l.X[count])
			var r = text.RuneCount(s)
			if r < l.xMaxCount {
				if count == 0 {
					s = s + strings.Repeat(" ", l.xMaxCount-r) + strings.Repeat(" ", spaceWidth)
				} else if count == len(l.X)-1 {
					s = strings.Repeat(" ", l.xMaxCount-r) + strings.Repeat(" ", l.width-index-l.xMaxCount) + s
				} else {
					s = strings.Repeat(" ", l.xMaxCount-r) + s + strings.Repeat(" ", spaceWidth)
				}
			}

			buf.WriteString(s)

			i += l.xMaxCount + spaceWidth - 1

			count++

			continue
		}

		if i == l.width+l.yMaxCount && l.xOffset != -1 {
			buf.WriteString(" ")
			continue
		} else {
			buf.WriteString(" ")
		}

	}

	return buf.String()

}

func getNextY(y []int) int {
	var v = math.MinInt
	for i := 0; i < len(y); i++ {
		if y[i] != math.MinInt {
			v = y[i]
			break
		}
	}
	return v
}

func (l *Line[T]) init() {
	var size, err = ts.GetSize()
	if err != nil {
		panic(err)
	}

	l.size = size

	if l.width == 0 || l.height == 0 {
		l.width = size.Col()
		l.height = size.Row()
	}

	l.yMaxCount = l.getMaxFloatCount(l.Y) + 1

	l.xMaxCount = getMaxRuneCount(l.X) + 1

	l.height = l.height - 1 - 1
	l.width = l.width - 1 - l.yMaxCount

	if l.height < 4 {
		l.height = 4
	}

	if l.width < 1+l.yMaxCount {
		l.width = 1 + l.yMaxCount
	}

	l.xOffset = size.Col() - l.width - 1
	l.yOffset = size.Row() - l.height - 1

	l.matrix = make([][]int, l.width)
	for i := 0; i < l.width; i++ {
		l.matrix[i] = make([]int, l.height)
	}

	for i := 0; i < len(l.matrix); i++ {
		for j := 0; j < len(l.matrix[i]); j++ {
			l.matrix[i][j] = math.MinInt
		}
	}

	var scale = float64(len(l.X)) / (float64(l.width) / float64(l.xMaxCount))
	if scale < 1 {
		scale = 1
	} else {
		scale = scale * 2 * float64(size.Col()) / float64(l.width)
	}

	var xLen = int(float64(len(l.X)) / (scale))

	var b []T
	for i := 0; i < xLen; i++ {
		b = append(b, l.X[int((scale)*float64(i))])
	}

	if len(l.X) > 1 && len(b) > 1 {
		if b[len(b)-1] != l.X[len(l.X)-1] {
			b = append(b, l.X[len(l.X)-1])
		}
	}

	if len(b) == 0 && len(l.X) > 0 {
		b = make([]T, 1)
		b[0] = l.X[0]
	}

	l.X = b
}

func getMaxRuneCount[T any](res []T) int {
	var max = 0
	for _, v := range res {
		var s = fmt.Sprintf("%v", v)
		var c = text.RuneCount(s)
		if c > max {
			max = c
		}
	}
	return max
}

func (l *Line[T]) parseFloatToString(f float64) string {
	var yPrecision = l.yPrecision
	if l.yPrecision == 0 {
		yPrecision = 1
	}
	var format = "%." + strconv.Itoa(yPrecision) + "f"
	var s = fmt.Sprintf(format, f)
	return s
}

func (l *Line[T]) getMaxFloatCount(res []float64) int {
	var max = 0
	var yPrecision = l.yPrecision
	if l.yPrecision == 0 {
		yPrecision = 1
	}
	var format = "%." + strconv.Itoa(yPrecision) + "f"
	for _, v := range res {
		var s = fmt.Sprintf(format, v)
		var c = text.RuneCount(s)
		if c > max {
			max = c
		}
	}

	if max < 4 {
		max = 4
	}

	return max
}

// get y min and max
func (l *Line[T]) yMinAndMax() (float64, float64) {
	if len(l.Y) == 0 {
		return 0, 0
	}
	var min = l.Y[0]
	var max = l.Y[0]
	for _, v := range l.Y {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}
