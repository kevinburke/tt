package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const Version = "0.4"

func getMocha(wd string) string {
	if wd == "" {
		return "mocha"
	}
	for i := 0; i < 10; i++ {
		path := filepath.Join(wd, "node_modules", ".bin", "mocha")
		s, err := os.Stat(path)
		if err != nil {
			wd = filepath.Dir(wd)
			continue
		}
		if s.Mode()&0100 > 0 {
			return path
		}
		break
	}
	return "mocha"
}

var bail = flag.Bool("bail", false, "bail on tests")
var grep = flag.String("grep", "", "Only run tests matching this string")
var timeout = flag.Uint("timeout", 2000, "Test timeout in ms")
var slow = flag.Uint("slow", 2, "Mark tests slower than this as slow")
var verbose = flag.Bool("verbose", false, "Print info about run commands")
var version = flag.Bool("version", false, "Print the current version and exit")

func main() {
	flag.Parse()
	if *version {
		fmt.Fprintf(os.Stderr, "tt version %s\n", Version)
		os.Exit(2)
	}
	// TODO walk upwards to git root
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mocha := getMocha(wd)
	var cmd *exec.Cmd
	args := []string{
		"--slow", strconv.FormatUint(uint64(*slow), 10),
		"--timeout", strconv.FormatUint(uint64(*timeout), 10),
	}
	if *bail {
		args = append(args, "--bail")
	}
	if *grep != "" {
		args = append(args, "--grep", *grep)
	}
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		args = append(args, "--colors")
	}
	if flag.NArg() > 0 {
		args = append(args, flag.Args()...)
		cmd = exec.Command(mocha, args...)
	} else {
		path := filepath.Join(wd, "test")
		_, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		files := make([]string, 0)
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			idx := strings.Index(path, "/test")
			if strings.Contains(path[idx:], "node_modules") {
				return nil
			}
			if filepath.Ext(path) == ".js" {
				files = append(files, path)
			}
			return nil
		})
		args = append(args, files...)
		cmd = exec.Command(mocha, args...)
	}
	cmd.Env = append(cmd.Env, "TZ=UTC")
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, "NODE_ENV=test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if *verbose {
		cmd.Env = append(cmd.Env, "VERBOSE=true")
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(cmd.Args, " "))
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
