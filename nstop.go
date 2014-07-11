package main
import (
    "fmt"
    "time"
    "os/exec"
    "encoding/json"
    "github.com/ric03uec/nstop/cli"
)

func main() {

    myMap := make([]string, 3)
    myMap = append(myMap, "h")
    myMap = append(myMap, "e")
    myMap = append(myMap, "l")
    fmt.Println(myMap)

    c := make([]string, len(myMap))
    copy(c, myMap)
    fmt.Println("E:", c)

    date := exec.Command("date")
    dateOut, _ := date.Output()
    fmt.Println(string(dateOut))

    bleh, _ := json.Marshal(true)
    fmt.Println(string(bleh))

    a := 5
    p := &a
    fmt.Println("pointer: ", &a)
    fmt.Println("pointer: ", *p)

    cn := container{}
    cn.id = "asdfsfdf"

    fmt.Println(cn.getId())
    cp := &cn
    fmt.Println(cp.getId())

    message := make (chan string, 2)
    go func(){
      message <- "hello from other end "
      message <- "bye bye from other end "
    }()

    msg := <-message
    msg1 := <-message
    fmt.Println(msg)
    fmt.Println(msg1)

    fmt.Printf("Hello worldddd")

    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ {
      go worker(w, jobs, results)
    }
    for j := 1; j <= 9; j++{
      jobs <- j
    }
    close(jobs)

    for a := 1; a <=9; a++ {
      <-results
    }

    cli.PrintSth()
}

func worker (id int, jobs <-chan int, results chan<- int){
  for j := range jobs {
    fmt.Println("Worker:", id, " processing job: ", j)
    time.Sleep(time.Second)
    results <- j*2
  }
}

type container struct {
  id string
}

func (c container) getId() string {
  return c.id
}

func (c *container) getIdPtr() string {
  return c.id
}
