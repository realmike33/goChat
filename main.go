//Fun fact: Go is built using Go

//All Go files require a package, the name doesn't matter for most cases, except for main packages
//All main packages require a main function
package main

//Import is like a mixture of bower components/npm modules and the require statement
//You would also require your Go files in this import statement if this file is dependent on them
//Remember: They are based on your GOPATH/GOROOT, in my case, my dev environment is nested under my GOPATH
//However, Go will check your GOPATH files for these files before checking the GOROOT.
//In reality, you could download all of your dependencies to your GOROOT and import from there
//Haven't played around too much to see if the pros/cons to dev environments like that
import (
	//command line flag parsing. Look down below for more information
	"flag"
	//fmt is a format package. Lots of different methods you can use to format
	//More info below
	"fmt"
	//custom web socket built for Golang, this is the actual location for the file, just like a require state but
	//absolute path from GOPATH/GOROOT which you define either during Go's installation or in your terminal
	//you obtain these files form running 'go get' in your terminal, they will be fetched and downloaded to your pkg
	//file in your GOPATH/GOROOT
	//You'll need to read their documentation on how to use their library of course
	"github.com/gorilla/websocket"
	//A logger package. Don't think of this like Console.log. It's not quite like that. There are methods like console.log
	//More info below
	"log"
	//Let's Golang perform like a server, you would also use this client side if you were in need to make a http/https request
	//More info below
	"net/http"
)

//Note: Go doesn't require var statements, however, if you use the var statement, it actually only works are a pointer and not the actual value
//This line is basically creating a variable called connections which points to a map, or hash table
//which is setting the websocket.Conn as the key and boolean as the value
var connections map[*websocket.Conn]bool

//Note: Go allows MULTIPLE return statements. Which is really cool and a bit frustrating if you aren't too careful
//Note: Go allows you to declare what type, from int to string to float32, the function returns.
//Hoisting is a different in go. Read: http://stratus3d.com/blog/2015/03/07/variable-hoisting-in-golang/

//This is a function declaration, looks familiar to function nameOfFunc(){} in JavaScript
//This function takes a collection of bytes: NO MORE TYPE CHECKS! YAY!
func sendAll(msg []byte) {
	//setting a for loop. For loops looks almost the sane as they do in JavaScript generally they look like
	//for i := 0; i < len(Array); i++, the function len is the same as Array.length in JavaScript
	//the range keyword is just a short handed way of a for loop. Like forEach but with no call back
	//We are also creating a conn value here. Are you can see, I am not using the var statement, this is because
	//conn is not being point but as the value
	for conn := range connections {
		//So a lot is going on on this single line. I am setting the value err to the return value of the conn.WriteMessage function
		//If a function has a check for errors, it will either return the error or nil
		//So this line is running the conn.WriteMessage method and then checking to see if it returned an error
		//Once again, this can be written a different way, but this is just a short handed way.
		//I think of this like Golang's ternary
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			//if there was an error, delete the connections hashtable and the current websocket connection and return
			delete(connections, conn)
			return
		}
	}
}

//another function declaration, however, this will be our callback function to our route
//I am setting w to the response and r to the request, all done over http
func wsHandler(w http.ResponseWriter, r *http.Request) {
	//See, told you! This function return two values, the actual connection and an err/or nil if no error
	//Most function return the actual value as the first return and the err as the second value, but this isn't always the case
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	//Another if statement, checking to see if connection to the websocket happened
	//when you see a _, that means the returning value isn't important, since you have to put something there, the go to symbol is the underscore
	if _, ok := err.(websocket.HandshakeError); ok {
		//if there was an error, send a 400 with a message
		http.Error(w, "Not a websocket handshake", 400)
	} else if err != nil {
		//This log statement will print the output. The standard logger, no formatting
		//If err, log error and return
		log.Println(err)
		return
	}
	//The defer statement, think of the final thing to happen before the function actually stops running
	//I think of this as like my personal garbage collection.
	//No matter what the outcome in this function is, from an error to a websocket connection, this will close that connection
	defer conn.Close()
	//This is the map we created up top, we're setting the value to the current connection and the value to true since it's an actual connection
	//As mentioned above, no need for type checking, I know my map connections will consist of connections as the key and a boolean as the value and will never be anything other than that
	connections[conn] = true
	//Infinite for loop. Will keep running unless told otherwise
	for {
		//This function returned three values, don't forget to read the docs for these things! Trust me.
		_, msg, err := conn.ReadMessage()
		//If there is an error
		if err != nil {
			//delete the connection
			delete(connections, conn)
			return
		}
		//else, log the message, which will show up in your terminal
		//string(msg), just turns the msg value into a string before logging it
		log.Println(string(msg))
		//sendAll is a websocket function but basically posts the message to the client
		sendAll(msg)
	}
}

func main() {
	//setting port as a flag integer, flag.Int(Name of flag, number, what is flag for)
	//First argument is a string, second argument is a number, third is also a string. All required.
	//setting port to the value, more info below
	port := flag.Int("port", 8000, "port to serve on")
	//setting dir as a flag string, flag.String(Name of flag, some string, what is flag for)
	//All arguments are strings. The argument values can be an empty strings but all arguments are required.
	//setting dir to the value, more info below
	dir := flag.String("directory", "web/", "directory of client files")
	//parses all flags into actual defined flags.
	//Parse allows you to use flags
	flag.Parse()

	//creating a connections hashtable like the one we created above
	//Note: The * is a key symbol in Go. It's a pointer.
	//You can read more here: http://golang.org/ref/spec#Pointer_types
	connections = make(map[*websocket.Conn]bool)

	//Setting fs as the returning value
	//which as basically the location to the static files of our application
	fs := http.Dir(*dir)
	//Go's fileserver, seems to work like node's FS.
	handler := http.FileServer(fs)
	//serves static assets
	http.Handle("/", handler)
	//links users to websockets via get request
	http.HandleFunc("/ws", wsHandler)

	//This log in your terminal. %d logs the second argument, *port(Which is our current port and not anything different).
	log.Printf("Running on port %d\n", *port)
	//formats this string to include our port. Just a different was of putting "localhost:8000"
	//NOTE: Although log and fmt looks to be the same thing, they are two separate pkgs and shouldn't be confused.
	addr := fmt.Sprintf("127.0.0.1:%d", *port)

	//This is where the server actually starts
	err := http.ListenAndServe(addr, nil)

	//Will print if there was an error when starting the server
	fmt.Println(err.Error())
}
