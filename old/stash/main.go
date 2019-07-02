package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

var (
	basedir string
	Usage   func()
	details bool
)

func init() {
	flag.StringVar(&basedir, "dir", ".", "Use path as a Store, cwd by default")
	flag.BoolVar(&details, "details", false, "Show details when displaying info")

	Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// Create a writer to pass into our command handler
	wr := os.Stdout

	// Decide what to do based on cmd line flags and args
	if len(flag.Args()) < 1 {
		Usage()
	}

	args := flag.Args()
	if len(args) < 1 {
		displayStoreList(wr, []string{"."})
	}

	switch args[0] {
	case "ls":
		displayStoreList(wr, args[1:])
	default:
		log.Errorf("unsupported cmd %s", args[0])
	}
}

func displayStoreList(w io.Writer, dirs []string) {
	var (
		st  *store.Store
		err error
	)

	fmt.Fprintln(w, "Store Meta: ")
	tabs := ""
	for _, dir := range dirs {
		if st, err = store.UseStore(dir); err != nil {
			fmt.Fprintf(w, "%sFailed to open store %v", tabs, err)
			return
		}
		fmt.Fprintln(w, tabs+st.String())
		if details {
			tabs := tabs + "\t"
			idx := st.Index()
			for i, _ := range idx {
				fmt.Fprintf(w, "%s%s\n", tabs, i)
			}
		}
	}
}
