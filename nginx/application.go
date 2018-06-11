package main

import (
	"fmt"
	"net/http"
	"os"
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
	host := os.Getenv("REDIS_SERVICE_SERVICE_HOST")
	port := os.Getenv("REDIS_SERVICE_SERVICE_PORT")
	addr := host + ":" + port
	fmt.Println("Redis service address retrived from the k8s generated env variable is:", addr)
	fmt.Println("Redis service address retrieved from the custom env variable is:", os.Getenv("SVC_DISC_REDIS_HOST"))
	rClient = *(redis.NewClient(&redis.Options{
		Addr: addr,
	}))
	pong, err := rClient.Ping().Result()
	fmt.Println(pong, err)

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":7080", nil); err != nil {
		panic(err)
	}
}
