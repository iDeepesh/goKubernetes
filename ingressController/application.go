package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

var rClient redis.Client

func sayHello(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	name := strings.TrimPrefix(url, "/")

	var count int64
	v, err := rClient.Get(name).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("Key does not exist", name)
			count = 0
		} else {
			message := fmt.Sprintln("Can't connect to redis\n", err)
			w.Write([]byte(message))
			return
		}
	} else {
		count, _ = strconv.ParseInt(v, 10, 64)
	}

	count++

	rClient.Set(name, count, 0)
	message := fmt.Sprintf("Hello %s. I am Latest. This is the visit number %d for you.\n", name, count)
	w.Write([]byte(message))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Just fine"))
}

func main() {
	rClient = *(redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Addr: "192.168.99.100:7001",
	}))
	pong, err := rClient.Ping().Result()
	fmt.Println(pong, err)

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":7080", nil); err != nil {
		panic(err)
	}
}
