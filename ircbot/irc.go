package main

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "net"
  "strings"
)

type Irc struct {
  host string
  port int
  conn net.Conn
  input *bufio.Scanner
}


func (irc *Irc) Connect(nick, name string) {
  addr := fmt.Sprintf("%s:%d", irc.host, irc.port)
  conn, err := net.Dial("tcp", addr)
  irc.conn = conn
  if err != nil {
    log.Fatal(err)
  }
  println("Connection established")
  irc.Write(fmt.Sprintf("NICK %s", nick))
  irc.Write(fmt.Sprintf("USER %s host server :%s", nick, name))

  irc.input = bufio.NewScanner(conn)
}


func (i *Irc) Write(msg string) {
  fmt.Printf("< %s\n", msg)
  msg = fmt.Sprintf("%s\r\n", msg)
  total := len(msg)
  for done := 0; done < total; {
    n, err := i.conn.Write([]byte(msg))
    if err != nil {
      log.Fatal(err)
    }
    done += n
  }
}


func (irc *Irc) ProcessLine(line string) {
  parts := strings.SplitN(line, " ", 3)
  prefix := parts[0]
  command := parts[1]
  params := ""
  if len(parts) <= 2 {
    params = parts[2]
  }

  if strings.ToUpper(command) == "PRIVMSG" {
    p := strings.SplitN(prefix, "!", 2)
    nick := p[0]
    p = strings.SplitN(parts[2], " ", 2)
    dest := p[0]
    pm := dest[0] != []byte("#")[0]
    if pm {
      dest = nick
    }
    msg := strings.Trim(p[1][1:], " ")

    p = strings.Split(msg, " ")
    if strings.ToLower(p[0]) == "!bot" {
      irc.Write(fmt.Sprintf("PRIVMSG %s :I am a bot written in Go", dest))
    } else if strings.ToLower(p[0]) == "!join" && len(p) == 2 {
      irc.Write(fmt.Sprintf("JOIN :%s", p[1]))
    } else if strings.ToLower(p[0]) == "!part" && !pm {
      irc.Write(fmt.Sprintf("PART :%s", dest))
    }
  }

  if prefix == command && prefix == params {
    println("Nope")
  }
}


func main() {
  var host = flag.String("host", "localhost", "Host to connect to")
  var port = flag.Int("port", 6666, "Port to use")
  flag.Parse()
  irc := Irc{host: *host, port: *port}
  irc.Connect("testbott", "TestBot")

  for irc.input.Scan() {
    line := irc.input.Text()
    line = strings.TrimLeft(line, ":")
    fmt.Printf("> %s\n", line)

    if line[0:4] == "PING" {
      server := strings.TrimLeft(line[4:], " :")
      irc.Write(fmt.Sprintf("PONG %s", server))
      continue
    }

    irc.ProcessLine(line)
  }
}

// vim:ts=2:sw=2:expandtab
