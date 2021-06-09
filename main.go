package main

import (
  "fmt"
  "flag"
  "os"
  "time"
  "crypto/tls"  
  "errors"
  "bufio"
  "sync"
)

var (
  email = flag.String("e", "", "Email address for sending expiry date of provided domains")
  domain = flag.String("d", "", "Single domain to be checked, in case you don't have a file containing domain names")
  filename = flag.String("f", "", "Route to a file with domains to check")
  left = flag.Bool("l", false, "Is going to return number of days left until the certificate expires instead of the expiry date")
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

  if *domain != "" && *filename != "" {
    fmt.Println("You can't choose both -d and -f flag at the same time.")
    os.Exit(2)
  } else if *domain != "" {
    func(domain string, left bool) {
      ex, err := sslCheck(domain) 
      if err != nil {
        fmt.Printf("%s: %s\n", domain, err)
        os.Exit(2)
      } else {
        if left {
          left := ex.Sub(time.Now()) 
          fmt.Printf("%s: %s | %d days left\n", domain, ex, int(left.Hours()) / 24)
        } else {
          fmt.Printf("%s: %s\n", domain, ex)
        }
      }
    }(*domain, *left)
  } else if *domain == "" && *filename == "" {
    fmt.Println("You have to provide a single domain (-d flag) or a file that contain domains (-f flag)")
  } else if *filename != "" {
    fmt.Println("Reading from file:", *filename) 
    readFile, err := os.Open(*filename)
    if err != nil {
      fmt.Println("Wrong filename was provided.")
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
          fmt.Printf("%s: %s\n", domain, err)
        } else {
          if left {
            left := ex.Sub(time.Now()) 
            fmt.Printf("%s: %s | %d days left\n", domain, ex, int(left.Hours()) / 24)
          } else {
            fmt.Printf("%s: %s\n", domain, ex)
          }
        }
      }(line, *left)
    }
    wg.Wait()
  }
}
