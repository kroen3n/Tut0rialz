package main

//a very simple
//tcp listener

import (
    "fmt"
    "net"
    "os"
)


func main() {

    // create server
    service := ":80"
    listener, err := net.Listen("tcp", service)

    if err != nil {
                fmt.Fprintln(os.Stdout, err)
                os.Exit(2)
    }
    

    fmt.Println("Listening...")

    for {
    
        // Listening for incoming connection.

        conn, err := listener.Accept()

        if err != nil {
	    fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
            os.Exit(1)
        }
	
        // Handle connections 
	      // Documentation: https://golang.org/pkg/net/
	
        go handleConnection(conn)
    }
}

// function to handle incoming requests 
func handleConnection(handleconn net.Conn) {

  defer handleconn.Close()

  handleconn.Write([]byte("'twas a success!Bye!\n"))

  handleconn.Close()

}
