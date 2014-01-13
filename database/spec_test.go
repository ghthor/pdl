package database

import (
	"flag"
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

var cfg config.Config

func init() {
	var err error

	buildDeps := flag.Bool("rebuild-deps", true, "run Make and build the test's dependencies")
	flag.Parse()

	if *buildDeps {
		//Build Test Dependencies
		cmd := exec.Command("make", "test-deps")

		cout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		cerr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			if _, e := io.Copy(os.Stdout, cout); e != nil {
				log.Fatal(e)
			}
		}()
		go func() {
			if _, e := io.Copy(os.Stderr, cerr); e != nil {
				log.Fatal(e)
			}
		}()

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}

	}
	// Read the Database Config file
	cfg, err = config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeAddFileExecutor)
	r.AddSpec(DescribeInstallAppExecutors)

	gospec.MainGoTest(r, t)
}
