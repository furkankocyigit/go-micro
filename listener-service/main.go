package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try to connect to rabbitmq

	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitCon.Close()
	
	//start listening for messages
	fmt.Println("Listening and consuming RabbitMQ messages...")

	//create consumer

	consumer, err := event.NewConsumer(rabbitCon)
	if err != nil {
		panic(err)
	}

	//watch the queue and consume evnts
	err = consumer.Listen([]string{"log.INFO","log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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