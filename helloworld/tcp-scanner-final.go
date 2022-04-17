package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

//Worker creator using input and ouptut given channels so no need for waitinggroup
func worker(ports, results chan int) {
	//the work is specified here
	//the work is: just take the port value passed in the channel and try to connect to it
	//since there are only 100 workers to process 1024 jobs, each worker
	//will have near 10 values in its channel; to process all of them we use "range"
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p

	}
}

func main() {

	//Workload is received here
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the ports to scan:")
	portsText, _ := reader.ReadString('\n')
	// convert CRLF to LF
	portsText = strings.Trim(portsText, "\n\r")
	//portsText := "80,443,8080,21-25"
	//portsText := "1-100"

	portsIndicators := strings.Split(portsText, ",")
	var ports2check []int
	for _, portsIndicator := range portsIndicators {
		if strings.Contains(portsIndicator, "-") {
			ports2check = addRangeToPorts2Check(ports2check, portsIndicator)
		} else {
			currentPort, _ := strconv.Atoi(portsIndicator)
			ports2check = append(ports2check, currentPort)
		}
	}

	//This is like a foreach
	sort.Ints(ports2check)
	fmt.Println("Ports to check:")
	for _, port := range ports2check {
		fmt.Printf("%d ", port)
	}
	fmt.Println()

	//a buffered channel is created here
	ports := make(chan int, 300)
	results := make(chan int) //not buffered, so a blocking is made untile sender receives
	var openports []int       //array containg the list of open ports

	//here we initialize workers
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	//this is the work that will be executed concurrently in goroutines through the ports channel
	//execution will flow automatically to the next for where it will be blocked
	go func() {
		for _, port := range ports2check {
			//the job parameter is sent to the channel: port, stating tyhe port number to check
			//this value will land in the buffer of the channel of a given worker in the pool
			//that worker will be processing all the buffer through the "range" keyword
			ports <- port
			//some worker in the pool will receive the data in its channel and process it
			//then the result will be transferred to the results channel; since it is not buffered
			//it will make the execution blocked for each result received
		}
	}()

	//This will compensate the async execution because n results are expected
	for i := 0; i < len(ports2check); i++ {
		//this will block the execution until a result is received
		port := <-results
		//only nonzero ports will be added to the array of open ports
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	//Open ports are sorted
	sort.Ints(openports)
	//This is like a foreach
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

}

func addRangeToPorts2Check(ports2check []int, portsIndicator string) []int {
	limits := strings.Split(portsIndicator, "-")
	lowerLimit, _ := strconv.Atoi(limits[0])
	upperLimit, _ := strconv.Atoi(limits[1])
	for i := lowerLimit; i <= upperLimit; i++ {
		ports2check = append(ports2check, i)
	}
	return ports2check
}
