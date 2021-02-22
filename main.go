package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	_ "unsafe" // for go:linkname
)

// Implemented in the syscall package.
//go:linkname fcntl syscall.fcntl
func fcntl(fd int, cmd int, arg int) (int, error)

func main() {
	// launchd_shim <target> <args>
	args := os.Args[1:]
	target := os.Args[1]

	verbose := os.Getenv("LAUNCHD_SHIM_VERBOSE") == "1"

	// LISTEN_PID should be set to the current PID.
	err := os.Setenv("LISTEN_PID", strconv.Itoa(os.Getpid()))
	if err != nil {
		log.Fatal(err)
	}

	// LISTEN_FDNAMES should be provided to specify the socket names.
	listenFdnames := os.Getenv("LISTEN_FDNAMES")
	i := 0
	for _, name := range strings.Split(listenFdnames, ":") {
		fds, err := activateSocket(name)
		if err != nil {
			log.Fatal(err)
		}

		if len(fds) != 1 {
			log.Fatalf("require exactly one socket for %s", name)
		}

		// fds should start at 3 and be incrementing.
		if fds[0] != i+3 {
			log.Fatalf("fd for %s must be %d", name, i+3)
		}

		// Clear the close on exec flag so that the fd persists to the target.
		fcntl(fds[0], syscall.F_SETFD, 0)
		i++
	}

	// Set LISTEN_FDS for the target.
	err = os.Setenv("LISTEN_FDS", strconv.Itoa(i))
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		log.Printf("LISTEN_PID=%s\n", os.Getenv("LISTEN_PID"))
		log.Printf("LISTEN_FDNAMES=%s\n", os.Getenv("LISTEN_FDNAMES"))
		log.Printf("LISTEN_FDS=%s\n", os.Getenv("LISTEN_FDS"))
	}

	// Inherit modified environment.
	env := os.Environ()
	err = syscall.Exec(target, args, env)
	if err != nil {
		log.Fatal(err)
	}
}
