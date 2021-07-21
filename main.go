package main

import (
  "encoding/json"
  "fmt"
  "flag"
  "os"
  "time"
  "crypto/tls"  
  "errors"
  "bufio"
  "sync"
  "net/smtp"
  "io/ioutil"
  "github.com/fatih/color"
)

type Config struct {
  Email       string `json:"sender_email"` 
  Password    string `json:"sender_password"`
  Host        string `json:"smtp_host"`
  Port        string `json:"smtp_port"`
  Receivers   []string `json:"receivers"`
  AlertBefore int `json:"alert_before"`
}

var (
  domain = flag.String("d", "", "Single domain to be checked, in case you don't have a file containing domain names")
  filename = flag.String("f", "", "Route to a file with domains to check")
  left = flag.Bool("l", false, "Is going to return number of days left until the certificate expires instead of the expiry date")
  email = flag.String("e", "", "Route to a JSON file with email configuration")
)

func sslCheck(url string) (time.Time, error) {
  conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", url), nil)
  if err != nil {
    return time.Time{}, errors.New("Couldn't establish a connection with the site. It's possible that you misspelled it's name or it has no certificate")
  }
  err = conn.VerifyHostname(url)
  if err != nil {
    return time.Time{}, errors.New("Couldn't find ssl certificate on the site.")
  }
  expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
  return expiry, nil
}

func main() {
  var wg sync.WaitGroup
  flag.Parse()
  var emailJson Config
  if *domain != "" && *filename != "" {
    fmt.Println("You can't choose both -d and -f flag at the same time.")
    os.Exit(2)
  } else if *domain != "" {
    func(domain string, left bool) {
      ex, err := sslCheck(domain) 
      if err != nil {
        fmt.Printf("%s: %s\n", color.RedString(domain), err)
        os.Exit(2)
      } else {
        l := int(ex.Sub(time.Now()).Hours() / 24)
        if *email != "" {
          emailJson = parseJsonFile(*email)
          if l <= emailJson.AlertBefore {
            sendMail(emailJson.Email, emailJson.Password, emailJson.Host, emailJson.Port, domain, emailJson.Receivers, l)
          }
        }
        if left {
          fmt.Printf("%s: %s | %d days left\n", color.GreenString(domain), color.YellowString(ex.String()), l)
        } else {
          fmt.Printf("%s: %s\n", color.GreenString(domain), color.YellowString(ex.String()))
        }
      }
    }(*domain, *left)
  } else if *domain == "" && *filename == "" {
    fmt.Println("You have to provide a single domain (-d flag) or a file that contain domains (-f flag)")
  } else if *filename != "" {
    fmt.Println("Reading from file:", color.MagentaString(*filename)) 
    readFile, err := os.Open(*filename)
    if err != nil {
      fmt.Println("No file named:", color.RedString(*filename))
    }
    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)
    var lines []string
    for fileScanner.Scan() {
      lines = append(lines, fileScanner.Text())
    }
    readFile.Close()
    wg.Add(len(lines))
    for _, line := range lines {
      go func(domain string, left bool) {
        defer wg.Done() 
        ex, err := sslCheck(domain) 
        if err != nil {
          fmt.Printf("%s: %s\n", color.RedString(domain), err)
        } else {
          l := int(ex.Sub(time.Now()).Hours() / 24)
          if *email != "" {
            emailJson = parseJsonFile(*email)
            if l <= emailJson.AlertBefore {
              sendMail(emailJson.Email, emailJson.Password, emailJson.Host, emailJson.Port, domain, emailJson.Receivers, l)
            }
          }
          if left {
            fmt.Printf("%s: %s | %d days left\n", color.GreenString(domain), color.YellowString(ex.String()), l)
          } else {
            fmt.Printf("%s: %s\n", color.GreenString(domain), color.YellowString(ex.String()))
          }
        }
      }(line, *left)
    }
    wg.Wait()
  }
}

func sendMail(sender, password, host, port, domain string, receiver []string, left int) string {
  auth := smtp.PlainAuth("", sender, password, host)
  msg := fmt.Sprintf("Certificate on %s is going to expire in %d days", domain, left)
  body := []byte(msg)

  err :=  smtp.SendMail(host + ":" + port, auth, sender, receiver, body)
  if err != nil {
    return "Couldn't send the email"
  }
  return "Email was sent successfully"
}

func parseJsonFile(file string) Config {
  readFile, err := os.Open(file)
  defer readFile.Close()
  if err != nil {
    fmt.Println("Couldn't find the email configuration file.")
    os.Exit(2)
  }
  byteValue, _ := ioutil.ReadAll(readFile)
  var config Config
  
  if err := json.Unmarshal(byteValue, &config); err != nil {
    fmt.Println("Failed to parse the email configuration file.")
    os.Exit(2)
  }
  return config
}
