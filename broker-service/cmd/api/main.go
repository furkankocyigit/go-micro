package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct{
	Rabbit *amqp.Connection
}

func main() {
	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitCon.Close()
	
	app := Config{
		Rabbit: rabbitCon,
	}

	log.Printf("Starting Broker service on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	
	//start server
	err = srv.ListenAndServe()
	if err !=nil{
		log.Panic(err)
	}
}

func connect()(*amqp.Connection, error){
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for{
		c,err :=amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil{
			fmt.Println("RabbitMQ is not ready yet...")
			counts++
		}else{
			connection = c
			log.Println("Connected to RabbitMQ")
			return connection, nil
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(backOff), 2)) * time.Second
		log.Printf("Backing off for %s", backOff)
		time.Sleep(backOff)
		continue
	}
}