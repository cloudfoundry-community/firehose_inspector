package utility

import (
  "github.com/nsf/termbox-go"
)

const (
  ColorLightGreen = termbox.Attribute(73)
  ColorBrightGreen = termbox.Attribute(41)
)


func WipeArea(x1 int, y1 int, x2 int, y2 int, fg termbox.Attribute, bg termbox.Attribute) {
  for y := y1; y < y2; y++ {
    for x := x1; x < x2; x++ {
      termbox.SetCell(x, y, ' ', fg, bg)
    }
  }
}

func WriteString(s string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
  for i := 0; i < len(s); i++ {
    termbox.SetCell(x + i, y, rune(s[i]), fg, bg)
  }
}

func SplitString(s string, width int) []string {
  ret := []string{}
  last := 0
  for i := 0; i+ width < len(s); i+=width {
    ret = append(ret, s[i:i + width])
    last = i + width
  }

  ret = append(ret, s[last:])
  return ret
}
