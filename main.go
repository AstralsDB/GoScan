package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zan8in/masscan"
)

func main() {
	context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(os.Args[1]),
		masscan.SetParamPorts("0-9999,u0-9999"),
		masscan.SetParamExclude("255.255.255.255"),
		masscan.SetParamWait(0),
		masscan.SetParamRate(100000000),
		masscan.WithContext(context),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		panic(err)
	}

	stdout := scanner.GetStdout()
	stderr := scanner.GetStderr()

	go func() {
		fmt.Printf("Address\t\tPort\n")

		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			// fmt.Println(srs.IP, srs.Port)
			fmt.Printf("%s\t\t%v\n", srs.IP, srs.Port)
			scannerResult = append(scannerResult, srs)
		}
	}()

	go func() {
		for stderr.Scan() {
			// fmt.Println("err: ", stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("%v Hosts Found\n", len(scannerResult))
	fmt.Printf("PID %v Exited", scanner.GetPid())
}
