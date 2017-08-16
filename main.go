
package main

import (
    "net/smtp"
    "io/ioutil"
    "log"
    "sync"
    "time"
)

var (
    relay = "localhost:2525"
    sessions = 100
    mail_per_session = 1000
    total = 100000
    sent = 0
    username = "detect"
    password = "123123"
)

var mt sync.Mutex

func stat() {
    for {
		time.Sleep(time.Second)
		log.Println("sent", sent, "emails")
	}
}

func worker(b []byte, jobs <-chan int, rslt chan<- int) {
    c, err := smtp.Dial(relay)

    if err != nil {
        log.Fatal("dial", err)
    }

    var auth = smtp.CRAMMD5Auth(username, password)

    if err := c.Auth(auth); err != nil {
        log.Fatal("auth: ", err)
    }

    done := 0
    for j := range jobs {

        // Set the sender and recipient first
        if err := c.Mail("someone@example.org"); err != nil {
            log.Fatal("mail from: ", err)
        }

        if err := c.Rcpt("someone@example.org"); err != nil {
            log.Fatal("rcpt: ", err)
        }

        // Send the email body.
        wc, err := c.Data()
        if err != nil {
            log.Fatal("data: ", err)
        }

        _, err = wc.Write(b)
        if err != nil {
            log.Fatal("data write: ", err)
        }

        err = wc.Close()
        if err != nil {
            log.Fatal("close: ", err)
        }

        mt.Lock()
        sent++
        mt.Unlock()

        rslt <- j
        done++

        if done >= mail_per_session {

            err = c.Quit()
            if err != nil {
                log.Fatal("quit: ", err)
            }

            c, err = smtp.Dial(relay)
            if err != nil {
                log.Fatal("dial: ", err)
            }

            if err := c.Auth(auth); err != nil {
                log.Fatal("auth: ", err)
            }

            done = 0
        }

    }

    // Send the QUIT command and close the connection.
    err = c.Quit()
    if err != nil {
        log.Fatal("quit: ", err)
    }

}

func main() {
    jobs := make(chan int, total)
    rslt := make(chan int, total)

    b, err := ioutil.ReadFile("email.txt")
    if err != nil {
        log.Fatal(err)
    }

    go stat()

    for w := 1; w <= sessions; w++ {
        go worker(b, jobs, rslt)
    }

    for j := 1; j <= total; j++ {
        jobs <- j
    }

    close(jobs)

    log.Println("wait...")

    for a := 1; a <= total; a++ {
        <- rslt
    }

    log.Println("sent", sent, "emails")

}
