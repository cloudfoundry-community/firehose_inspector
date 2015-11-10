package screen_outputs

import (
  "fmt"

  "github.com/cloudfoundry/sonde-go/events"
  "github.com/cloudfoundry-community/firehose_inspector/pages"
  "github.com/cloudfoundry-community/firehose_inspector/utility"
  "github.com/nsf/termbox-go"
)

type NullDisplay struct {
  page *pages.Page
}

func (r *NullDisplay) Init() {

}

func (r *NullDisplay) Setup(page *pages.Page) {
  r.page = page

  line := 5
  for y := 0; y < 4; y++ {
		for x := 0; x < 8; x++ {
			for z := 0; z < 8; z++ {
				c1 := termbox.Attribute(256 - y*64 - x*8 - z)
				c2 := termbox.Attribute(1 + y*64 + z*8 + x)
				c3 := termbox.Attribute(256 - y*64 - z*8 - x)
				c4 := termbox.Attribute(1 + y*64 + x*4 + z*4)
        utility.WriteString(fmt.Sprintf("%d,%d,%d", y, x, z), 0, line, c1, r.page.Background)
        utility.WriteString(fmt.Sprintf("%d,%d,%d", y, z, x), 10, line, c2, r.page.Background)
        utility.WriteString(fmt.Sprintf("%d,%d,%d", y, z, x), 20, line, c3, r.page.Background)
        utility.WriteString(fmt.Sprintf("%d,%d,%d", y, x, z), 30, line, c4, r.page.Background)
        line ++
			}
		}
	}

  termbox.Flush()
}

func (r *NullDisplay) Update(env *events.Envelope) {
}

func (r *NullDisplay) KeyEvent(key termbox.Key) {

}
