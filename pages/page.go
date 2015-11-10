package pages

import (
  "fmt"
  // "log"
  "github.com/nsf/termbox-go"
  "github.com/cloudfoundry/sonde-go/events"
)

type Page struct {
  Title       string
  HelpText    string
  Outputs     []Output
  Foreground  termbox.Attribute
  Background  termbox.Attribute
}

type Output interface {
  Init()
  Setup(*Page)
  Update(*events.Envelope)
  KeyEvent(termbox.Key)
}

const LinesHLineStd = '─'
const LinesVLineStd = '│'
const LinesHLineBold = '━'
const LinesVerticalJoinUp = '┴'
const LinesVerticalJoinDown = '┬'

func (p *Page) Draw(pages []Page) {
  w, h := termbox.Size()

  termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
  // draw divider

  for lx := 0; lx < w; lx++ {
    termbox.SetCell(lx, h - 2, LinesHLineStd, p.Foreground, p.Background)
    termbox.SetCell(lx, 2, LinesHLineStd, p.Foreground, p.Background)
  }

  for i := 0; i < len(p.Title); i++ {
    termbox.SetCell(i + 2, 1, rune(p.Title[i]), p.Foreground, p.Background)
  }

  // tabWidth := w / len(pages)
  start := 0

  for _, page := range pages {

    title := fmt.Sprintf("   %s   ", page.Title)

    var fg termbox.Attribute
    var bg termbox.Attribute

    if page.Title == p.Title {
      fg = termbox.ColorWhite
      bg = termbox.ColorRed
    } else {
      fg = p.Foreground
      bg = p.Background
    }

    for i := 0; i < len(title); i++ {
      termbox.SetCell(start + i, h - 1, rune(title[i]), fg, bg)
    }

    start = start + len(title) + 2

  }

  for _, output := range p.Outputs {
    output.Setup(p)
  }

  termbox.Flush()
}
