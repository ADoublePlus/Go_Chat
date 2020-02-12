package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/fatih/color"
)

func main()  {
	// Both the host and port can be set from the command line
	// Default values are localhost:8080
	host := flag.String("host", "localhost", "hostname/ip")
	port := flag.String("port", ":8080", "chat server port")
	flag.Parse()

	/* Connect to server through TCP */
	conn, err := net.Dial("tcp", *host + *port)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	/* Display a welcome message */
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	welcome := color.New(color.Bold, color.FgCyan).PrintlnFunc()
	welcome("===============================================================")
	welcome("Connected to the server...")
	welcome("After you have sent your first message you'll be assigned a client ID.")
	welcome("Press", red("enter"), "to send a message.")
	welcome("Type ", red(":exit"), "to quit.")
	welcome("===============================================================")

	for {
		/* Read a message typed in the console */
		newMsg, err := stdinReader.ReadString('\n')

		// Each client implementation can have a custom form of exiting
		// Even if the server requires a specific set of commands,
		// the client can translate their commands to match that of the servers
		exitCmd, err := regexp.MatchString(":exit", newMsg)

		if exitCmd {
			log.Printf("Connection closed by the user.")
			os.Exit(0)
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		// Send the typed message
		fmt.Fprintf(conn, white(newMsg))

		/* Output messages sent by clients */
		go func() {
			for {
				reply, err := connReader.ReadString('\n')

				if err != nil {
					log.Printf("Connection closed unexpectedly by the server.")
					os.Exit(1)
				}

				fmt.Printf(green(reply))
			}
		}()
	}
}