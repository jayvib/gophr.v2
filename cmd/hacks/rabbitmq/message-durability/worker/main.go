package main

import (
  "bytes"
  "github.com/streadway/amqp"
  "log"
  "time"
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

  q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
  failOnError(err)

  msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
  failOnError(err)

  forever := make(chan bool)

  go func() {
    for d := range msgs {
      log.Printf("Received a message: %s", d.Body)
      dot_count := bytes.Count(d.Body, []byte("."))
      t := time.Duration(dot_count)
      time.Sleep(t * time.Second)
      log.Printf("Done")
      err := d.Ack(true)
      if err != nil {
        log.Printf("Error: %v\n", err)
      }
    }
  }()
  log.Printf("Waiting for messages. To exist press CTRL+C")
  <-forever
}


