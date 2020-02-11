package main

import (
	"fmt"
	"os"
)

// printUsage of dtn-tool and exit with an error code afterwards.
func printUsage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s create|show|serve-dir:\n\n", os.Args[0])

	_, _ = fmt.Fprintf(os.Stderr, "%s create sender receiver -|filename -|bundle-name\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "  Creates a new Bundle, addressed from sender to receiver with the stdin (-)\n")
	_, _ = fmt.Fprintf(os.Stderr, "  or the given file (filename) as payload. This Bundle will be written to the\n")
	_, _ = fmt.Fprintf(os.Stderr, "  stdout (-) or saved as bundle-name.\n\n")

	_, _ = fmt.Fprintf(os.Stderr, "%s show filename\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "  Prints a human-readable version of the given Bundle.\n\n")

	_, _ = fmt.Fprintf(os.Stderr, "%s serve-dir websocket endpoint-id directory\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "  %s registeres itself as an agent on the given websocket and writes\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "  incoming Bundles in the directory. If the user dropps a new Bundle in the\n")
	_, _ = fmt.Fprintf(os.Stderr, "  directory, it will be sent to the server.\n\n")

	os.Exit(1)
}

// printFatal of an error with a short context description and exits afterwards.
func printFatal(err error, msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s errored: %s\n  %v\n", os.Args[0], msg, err)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	switch os.Args[1] {
	case "create":
		createBundle(os.Args[2:])

	case "show":

	case "serve-dir":

	default:
		printUsage()
	}
	//
	//api := "ws://localhost:8080/ws"
	//eid := "dtn://foo/bar"
	//
	//wac, err := agent.NewWebSocketAgentConnector(api, eid)
	//if err != nil {
	//	log.WithError(err).Fatal("Creating WebSocket agent errored")
	//}
	//
	//b := createBundle(eid, "dtn://uff/", []byte("hello world"))
	//if err := wac.WriteBundle(b); err != nil {
	//	log.WithError(err).Fatal("Sending Bundle errored")
	//}
	//
	//wac.Close()
}
