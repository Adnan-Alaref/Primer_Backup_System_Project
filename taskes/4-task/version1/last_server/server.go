/*
.
.
SERVER CODE
.
.
*/
package main
 
import (
    "fmt"
    "log"
    "net"
)
 
func main() {
    fmt.Println("server listening on 3000")
    listener, err := net.Listen("tcp", "localhost:3000")
    if err != nil {
        log.Fatalln(err)
    }
    defer listener.Close()
    
    // listening for incoming connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println("new connection")
        
        // listen to connections in another gorutine
        go listenConnection(conn)
    }
}
 
// listening for messages from connection
func listenConnection(conn net.Conn) {
    for {
        buffer := make([]byte, 1400)
        dataSize, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("connection closed")
            return
        }
        
        // the actual message
        data := buffer[:dataSize]
        fmt.Println("received message: ", string(data))
        
        // echoing the message back out
        _, err = conn.Write(data)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println("Message sent: ", string(data))
    }
}


