package main

import (
  "crypto/tls"
	"fmt"
  "os"
  // "log"

	"github.com/cloudfoundry/cli/plugin"
  "github.com/cloudfoundry/cli/cf/configuration/config_helpers"
  "github.com/cloudfoundry/cli/cf/configuration/core_config"
  "github.com/cloudfoundry/noaa"
  "github.com/cloudfoundry/sonde-go/events"

  so "github.com/cloudfoundry-community/firehose_inspector/screen_outputs"
  p "github.com/cloudfoundry-community/firehose_inspector/pages"
  "github.com/nsf/termbox-go"

)

const firehoseSubscriptionId = "firehose-a"

type FirehoseInspector struct {
  currentPage p.Page
  pages []p.Page
}

type ConsoleDebugPrinter struct{}

func (fi *FirehoseInspector) ChangePage(index int) {
  fi.currentPage = fi.pages[index]
  fi.currentPage.Draw(fi.pages)
}

func (fi *FirehoseInspector) Run(cliConnection plugin.CliConnection, args []string) {

	if args[0] == "firehose-inspector" {

    confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)

    dopplerEndpoint := confRepo.DopplerEndpoint();
    accessToken := confRepo.AccessToken();

    connection := noaa.NewConsumer(dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)

    msgChan := make(chan *events.Envelope)

    fi.pages = []p.Page {
      p.Page {
        Title: "Logs",
        Outputs: []p.Output{
          &so.LogsDisplay {
            MarginPos: 32,
          },
        },
        Foreground: termbox.ColorWhite,
        Background: termbox.ColorDefault,
      },
      p.Page {
        Title: "Page 1",
        Outputs: []p.Output{
          &so.NullDisplay{},
        },
        Foreground: termbox.ColorWhite,
        Background: termbox.ColorDefault,
      },
      p.Page {
        Title: "Page 2",
        Outputs: []p.Output{
          &so.NullDisplay{},
        },
        Foreground: termbox.ColorWhite,
        Background: termbox.ColorDefault,
      },
    }

    // call init on all Outputs
    for _, page := range fi.pages {
      for _, output := range page.Outputs {
        output.Init()
      }
    }

    fi.currentPage = fi.pages[0]

    err := termbox.Init()
  	if err != nil {
  		panic(err)
  	}
  	defer termbox.Close()
    termbox.SetOutputMode(termbox.Output256)
    fi.currentPage.Draw(fi.pages)

    go func() {

      defer close(msgChan)
      errorChan := make(chan error)
      go connection.Firehose(firehoseSubscriptionId, accessToken, msgChan, errorChan)

      for err := range errorChan {
        fmt.Fprintf(os.Stderr, "%v\n", err.Error())
      }

    }()

    go func() {
      for msg := range msgChan {
        for _, output := range fi.currentPage.Outputs {
          output.Update(msg)
        }
      }
    }()

    loop:

  	for {
  		switch ev := termbox.PollEvent(); ev.Type {
  		case termbox.EventKey:
        for _, output := range fi.currentPage.Outputs {
          output.KeyEvent(ev.Key)
        }
  			switch ev.Key {
  			case termbox.KeyEsc:
  				break loop
  			case termbox.KeyF1:
          fi.ChangePage(0)
        case termbox.KeyF2:
          fi.ChangePage(1)
        case termbox.KeyF3:
          fi.ChangePage(2)

  			// 	switch_output_mode(1)
  			// 	draw_all()
  			// case termbox.KeyArrowDown, termbox.KeyArrowLeft:
  			// 	switch_output_mode(-1)
  			// 	draw_all()
  			}
  		case termbox.EventResize:
        fi.currentPage.Draw(fi.pages)
        // draw_all()
  		}
  	}
	}
}

func (c *FirehoseInspector) GetMetadata() plugin.PluginMetadata {

	return plugin.PluginMetadata{
		Name: "Firehose Inspector",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "firehose-inspector",
				HelpText: "Firehose Inspector's help text",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "firehose-inspector\n   cf firehose-inspector",
				},
			},
		},
	}
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {

  logOut, _ := os.OpenFile("/tmp/firehose.stdout.log", os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
  logErr, _ := os.OpenFile("/tmp/firehose.stderr.log", os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
	os.Stdout = logOut;
	os.Stderr = logErr;

  // if err != nil {
  //     t.Fatalf("error opening file: %v", err)
  // }
  defer logOut.Close()
  defer logErr.Close()

	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "github.com/cloudfoundry/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
  fmt.Print("Starting plugin\n")

	plugin.Start(new(FirehoseInspector))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}
