package main

import (
	"fmt"
	"net"
)

func mainSerialScanner() {
	for i := 1; i <= 1024; i++ {
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			//port closed or filtered
			fmt.Printf("%d X\n", i)
			continue
		}
		conn.Close()
		fmt.Printf("%d      open\n", i)
	}
}
