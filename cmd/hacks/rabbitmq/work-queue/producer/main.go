package main

import (
  "github.com/streadway/amqp"
  "log"
  "os"
  "strings"
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

  body := bodyFrom(os.Args)
  err = ch.Publish("", q.Name, false, false,
    amqp.Publishing{
      DeliveryMode: amqp.Persistent,
      ContentType: "text/plain",
      Body: []byte(body),
    })
  failOnError(err)

}

func bodyFrom(args []string) string {
  var s string
  if (len(args) < 2) || os.Args[1] == "" {
    s = "hello"
  } else {
    s = strings.Join(args[1:], " ")
  }
  return s
}


