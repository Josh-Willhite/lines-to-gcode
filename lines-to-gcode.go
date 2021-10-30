package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func LinesInFile(fileName string) []string {
	f, _ := os.Open(fileName)
	scanner := bufio.NewScanner(f)
	result := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

type Line struct {
	Start []float64
	End   []float64
}

func (l *Line) Full() bool {
	if len(l.Start) == 3 && len(l.End) == 3 {
		return true
	}
	return false
}

func NewLine() Line {
	return Line{[]float64{}, []float64{}}
}

type Lines struct {
	lines []Line
}

func (l *Lines) Validate() {
	// make sure we have a closed shape
}

func (l *Lines) OriginLineIndex() int {
	origin := []float64{0, 0, 0}
	closest := 0
	for idx, line := range l.lines[1:] {
		if Distance(line.Start, origin) < Distance(l.lines[closest].Start, origin) {
			closest = idx
		}
	}
	return closest
}

func Distance(a, b []float64) float64 {
	return math.Hypot(a[0]-b[0], a[1]-b[1])
}

func NextLine(end []float64, l *Line) bool {
	matchedStart, matchedEnd := true, true
	if Distance(end, l.Start) > 0.01625 {
		matchedStart = false
	}
	if Distance(end, l.End) > 0.01625 {
		matchedEnd = false
	}

	if matchedStart {
		return true
	}
	if matchedEnd {
		l.Start, l.End = l.End, l.Start
		return true
	}
	return false
}

func (l *Lines) Swap(x, y int) {
	l.lines[x], l.lines[y] = l.lines[y], l.lines[x]
}

func (l *Lines) OrderLines() {
	startLineIdx := l.OriginLineIndex()
	l.Swap(0, startLineIdx)
	for i := 1; i < len(l.lines); i++ {
		for j := i; j < len(l.lines); j++ {
			if NextLine(l.lines[i-1].End, &l.lines[j]) {
				l.Swap(i, j)
			}
		}
	}

}

func ParseLines(fileLines []string) Lines {
	lines := Lines{lines: []Line{}}
	var line Line
	for idx, l := range fileLines {
		if strings.Contains(l, "LINE") {
			line = NewLine()
		}
		switch {
		case !line.Full() && strings.Contains(l, " 10"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.Start = append(line.Start, num)
		case !line.Full() && strings.Contains(l, " 20"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.Start = append(line.Start, num)
		case !line.Full() && strings.Contains(l, " 30"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.Start = append(line.Start, num)
		case !line.Full() && strings.Contains(l, " 11"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.End = append(line.End, num)
		case !line.Full() && strings.Contains(l, " 21"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.End = append(line.End, num)
		case !line.Full() && strings.Contains(l, " 31"):
			num, err := strconv.ParseFloat(fileLines[idx+1], 64)
			if err != nil {
				log.Fatalf("Failed to parse with: %s", err)
			}
			line.End = append(line.End, num)
			lines.lines = append(lines.lines, line)
		}
	}
	return lines
}

func GenerateGCode(lines Lines, layerDepth float64, layerCount int, feedRate float64) []string {
	gcode := []string{}
	gcode = append(gcode, "G90")
	for i := 1; i <= layerCount; i++ {
		gcode = append(gcode, fmt.Sprintf("G1 Z-%.3f F%.3f", layerDepth*float64(i), feedRate/2))
		gcode = append(gcode, fmt.Sprintf("G1 X%.3f Y%.3f F%.3f", lines.lines[0].Start[0], lines.lines[0].Start[1], feedRate))
		for _, line := range lines.lines {
			gcode = append(gcode, fmt.Sprintf("G1 X%.3f Y%.3f F%.3f", line.End[0], line.End[1], feedRate))
		}
	}

	// Raise the tool above the surface.
	gcode = append(gcode, fmt.Sprintf("G1 Z2.000 F%.3f", feedRate))
	gcode = append(gcode, "M2")
	return gcode
}

func main() {
	feedRate := float64(1)
	layerDepth := float64(0.4375) / float64(3)
	layerCount := 3
	fileLines := LinesInFile("/home/joshw/Documents/plant-stand-a.dxf")
	lines := ParseLines(fileLines)
	lines.OrderLines()
	gcode := GenerateGCode(lines, layerDepth, layerCount, feedRate)
	for _, l := range gcode {
		fmt.Println(l)
	}
}
