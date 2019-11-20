/*
	Source code example for TDC 2019
	Author: Guilherme Ribeiro <email@guiguis.net>
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	flagOp     string // Valid operations are generate-config, update-backends
	flagCommit bool   // Whether changes must be commit or only shown
)

// CustomWriter main object
type CustomWriter struct{}

// Write prints informational log event
func (c *CustomWriter) Write(bytes []byte) (int, error) {
	var out string
	if strings.HasPrefix(string(bytes), "WARN") {
		out = color.YellowString(time.Now().Format(time.RFC3339) + " " + string(bytes))
	} else {
		out = color.GreenString(time.Now().Format(time.RFC3339) + " " + string(bytes))
	}
	return fmt.Print(out)
}

func main() {

	log.SetFlags(0)
	log.SetOutput(&CustomWriter{})

	flag.Parse()

	if flag.NFlag() == 0 {
		log.Print("INFO No flags given. Exiting.")
		os.Exit(0)
	}

	if flagCommit {
		log.Print("WARN Commit is actived.")
	}

	op := &Op{}

	switch flagOp {
	case "generate-config":
		log.Print("INFO Creating new basic configuration")
		op.GenerateConfig()
	case "update-backends":
		log.Print("INFO Updating backends accordingly to consul cluster")
		op.UpdateBackends()
	case "clear":
		log.Print("INFO Cleaning up old files...")
		op.Clear()
	case "show":
		log.Print("INFO Getting previously generate configs to show...")
		op.Show()
	default:
		log.Print("WARN Invalid operation!")
	}

}

func init() {
	flag.StringVar(&flagOp, "op", "", "Inform a valid operation: generate-config, update-backends, clear")
	flag.BoolVar(&flagCommit, "commit", true, "Whether changes must be commited or only shown")

}
