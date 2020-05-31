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

  ch, err := conn.Channel()
  failOnError(err)
  defer ch.Close()

  q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
  failOnError(err)

  body := "Hello World!"
  err = ch.Publish("", q.Name, false, false,
    amqp.Publishing{
      ContentType: "text/plain",
      Body: []byte(body),
    })
  failOnError(err)

}
