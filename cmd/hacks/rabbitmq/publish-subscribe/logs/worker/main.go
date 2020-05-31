package main

import (
  "github.com/streadway/amqp"
  "log"
)

func failOnError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  // Make a connection
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err)
  defer conn.Close()

  // connect to a channel
  ch, err := conn.Channel()
  failOnError(err)

  err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
  failOnError(err)

  q, err := ch.QueueDeclare("", false, false, false, false, nil)
  failOnError(err)

  err = ch.QueueBind(q.Name, "", "logs", false, nil) // subscribe
  failOnError(err)

  msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
  failOnError(err)

  forever := make(chan bool)

  go func() {
    for d := range msgs {
      log.Printf("[x]: %s", d.Body)
    }
  }()
  log.Printf("Waiting for messages. To exist press CTRL+C")
  <-forever
}


