package screen_outputs

import (
  "fmt"
  "sort"
  "strings"
  "github.com/cloudfoundry/sonde-go/events"
  "github.com/cloudfoundry-community/firehose_inspector/pages"
  "github.com/cloudfoundry-community/firehose_inspector/utility"
  "github.com/nsf/termbox-go"
)


type LogLine struct {
  Line string
  Source string
  Indent bool
}

type LogsDisplay struct {
  MarginPos int
  LogOrigins map[string]bool

  lineBuffer []LogLine
  selectedOrigin int
  originKeys []string
  page *pages.Page
  width int
  height int
}

func (r *LogsDisplay) padRight(s string, padStr string, overallLen int) string{
	var padCountInt int
	padCountInt = 1 + ((overallLen-len(padStr))/len(padStr))
	var retStr =  s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func (r *LogsDisplay) Init() {
  r.lineBuffer = []LogLine{}
  r.selectedOrigin = 0
  r.LogOrigins = make(map[string]bool)
}

func (r *LogsDisplay) Setup(page *pages.Page) {
  r.page = page
  r.width, r.height = termbox.Size()

  utility.WriteString ("Filters: ", 2, 3, page.Foreground, page.Background)

  termbox.SetCell(r.MarginPos, 2, rune(pages.LinesVerticalJoinDown), page.Foreground, page.Background)
  termbox.SetCell(r.MarginPos, r.height-2, rune(pages.LinesVerticalJoinUp), page.Foreground, page.Background)

  for i := 3; i < r.height - 2; i++ {
    termbox.SetCell(r.MarginPos, i, rune(pages.LinesVLineStd), page.Foreground, page.Background)
  }

  r.DrawOriginList()
  r.DrawLogBuffer()
  termbox.Flush()
}

func (r *LogsDisplay) KeyEvent(key termbox.Key) {

  if key == termbox.KeyArrowUp { if r.selectedOrigin > 0 { r.selectedOrigin -- } }

  if key == termbox.KeyArrowDown {
    if r.selectedOrigin < (len(r.originKeys) - 1) { r.selectedOrigin ++ }
  }

  if key == termbox.KeySpace {
    originKey := r.originKeys[r.selectedOrigin]
    r.LogOrigins[originKey] = !r.LogOrigins[originKey]
  }

  r.DrawOriginList()
  termbox.Flush()
}

func (r *LogsDisplay) DrawOriginList() {

  utility.WipeArea(0, 5, r.MarginPos, r.height - 2, r.page.Foreground, r.page.Background)
  index := 0

  for _, origin := range r.originKeys {

    var bg termbox.Attribute
    fg := r.page.Foreground

    if r.LogOrigins[origin] {
      bg = termbox.ColorGreen
    } else {
      bg = termbox.ColorRed
    }

    if index == r.selectedOrigin {
      origin = fmt.Sprintf("> %s", origin)
      fg = termbox.ColorWhite | termbox.AttrBold
    } else {
      origin = fmt.Sprintf("  %s", origin)
    }

    origin = r.padRight(origin, " ", r.MarginPos)

    for i := 0; i < len(origin); i++ {
      termbox.SetCell(i, index + 5, rune(origin[i]), fg, bg)
    }
    index ++
  }
}

func (r *LogsDisplay) GetLinesFromLogBytes(logBytes []byte) []string {
  logMsg := string(logBytes)
  return strings.Split(logMsg, "\n")
}

func (r *LogsDisplay) DrawLogBuffer() {
  utility.WipeArea(r.MarginPos + 1, 3, r.width, r.height - 2, r.page.Foreground, r.page.Background)

  bufferStart := 3
  bufferEnd := r.height - 2

  if len(r.lineBuffer) < (bufferEnd - bufferStart) {
    bufferStart = bufferEnd - len(r.lineBuffer)
  }

  offset := len(r.lineBuffer) - (bufferEnd - bufferStart)

  for i := bufferEnd; i > bufferStart; i-- {
    line := r.lineBuffer[(i - bufferStart) + offset - 1]

    cursor := r.MarginPos + 1

    fmt.Printf("%d - %v : %s\n", len(line.Line), line.Indent, line.Source)

    if line.Source != "" {
      source := fmt.Sprintf("[%s] ", line.Source)
      utility.WriteString(source, cursor, i - 1, utility.ColorBrightGreen, r.page.Background)
      cursor += len(source)
    }

    if line.Indent {
      utility.WriteString(">> ", cursor, i - 1, utility.ColorBrightGreen, r.page.Background)
      cursor += 3
    }

    utility.WriteString(line.Line, cursor, i - 1, r.page.Foreground, r.page.Background)
  }
}

func (r *LogsDisplay) Update(env *events.Envelope) {

  origin := *env.Origin
  eventType := *env.EventType

  if eventType != events.Envelope_LogMessage {
    return
  }

  origin = fmt.Sprintf("%s/%s/%s", origin, *env.LogMessage.SourceInstance, *env.LogMessage.SourceType)
  _, ok := r.LogOrigins[origin]

  if !ok {
    r.LogOrigins[origin] = true

    r.originKeys = []string{}

    for k := range r.LogOrigins {
      r.originKeys = append(r.originKeys, k)
    }

    sort.Strings(r.originKeys)
    r.DrawOriginList()
  }

  // skip it if disabled
  if !r.LogOrigins[origin] { return }

  // add line to buffer
  lines := r.GetLinesFromLogBytes(env.LogMessage.Message)
  displayLineWidth := r.width - r.MarginPos

  croppedLines := []string{}

  // split to screen_width
  for _, line := range lines {
    fmt.Printf("%s\n--------------\n", line)
    for _, croppedLine := range utility.SplitString(line, displayLineWidth - 10) {
      croppedLines = append(croppedLines, croppedLine)
    }
  }

  // format the logs
  for i, line := range croppedLines {
    if len(line) == 0 { continue }

    var logLine LogLine
    if i == 0 {
      logLine = LogLine {
        Line: line,
        Source: fmt.Sprintf("%s/%s", *env.LogMessage.SourceType, *env.LogMessage.SourceInstance),
        Indent: false,
      }
    } else {
      logLine = LogLine {
        Line: line,
        Indent: true,
      }
    }

    r.lineBuffer = append(r.lineBuffer, logLine)
  }

  r.DrawLogBuffer()

  termbox.Flush()
}
