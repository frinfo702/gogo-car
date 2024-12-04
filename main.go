package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func Routers() *mux.Router {
	var port = "8080"
	r := mux.NewRouter() // router instance

	flag.Parse()
	log.SetFlags(0)

	// register handlers
	r.HandleFunc("/", adminHandler)
	r.HandleFunc("/observe", observeSitutationHandler)
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/home", home)

	fmt.Println("starting server on localhost: 8080")

	http.ListenAndServe(":"+port, r)
	log.Fatal(http.ListenAndServe(*addr, nil))
	return r
}

func adminHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "admin endpoint is running well")
}

func observeSitutationHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this endppoint will display collected logs")
}
func main() {
	// port := 8000

	// welcome messeage to check if the server is running well
	fmt.Println("Go WebSockets is running")
	// call setupRoutes function
	Routers()
	// make router

}

// server
var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default option

func echo(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	defer conn.Close()

	// use conn to send and receive message
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message: ", err)
			break
		}
		log.Printf("received: %s\n", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write: ", err)
			break
		}

	}
}

func home(w http.ResponseWriter, req *http.Request) {
	homeTemplate.Execute(w, "ws://"+req.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))

// client
// 1秒ごとにメッセージを送信する

