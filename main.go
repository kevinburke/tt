package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"
)

func getMocha() string {
	wd, err := os.Getwd()
	if err != nil {
		return "mocha"
	}
	path := filepath.Join(wd, "node_modules", ".bin", "mocha")
	s, err := os.Stat(path)
	if err != nil {
		return "mocha"
	}
	if s.Mode()&0100 > 0 {
		return path
	}
	return "mocha"
}

var bail = flag.Bool("bail", false, "bail on tests")

func main() {
	flag.Parse()
	mocha := getMocha()
	var cmd *exec.Cmd
	args := []string{"--slow", "2"}
	if *bail {
		args = append(args, "--bail")
	}
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		args = append(args, "--colors")
	}
	if flag.NArg() > 0 {
		args = append(args, "")
		copy(args[1:], args[:])
		args[0] = flag.Arg(0)
		cmd = exec.Command(mocha, args...)
	} else {
		cmd = exec.Command(mocha, args...)
	}
	cmd.Env = append(cmd.Env, "TZ=UTC")
	cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, "NODE_ENV=test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
