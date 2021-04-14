package main

import (
	"context"
	"flag"
	serial2 "github.com/MaksimSvinin/com-port-monitor/serial"
	"github.com/MaksimSvinin/com-port-monitor/tui"
	"log"
)

var port = flag.String("p", "", "port")
var baudRate = flag.Uint("b", 9600, "BaudRate")
var dataBits = flag.Uint("d", 8, "DataBits")
var stopBits = flag.Uint("s", 1, "StopBits")

func main() {
	flag.Parse()
	if *port == "" {
		log.Fatal("port flag -p not found")
	}

	ctx := context.Background()
	inputCh := make(chan string)
	writeCh := make(chan string, 10)

	sidebarData := tui.SidebarData{
		PortName: *port,
		BaudRate: *baudRate,
		DataBits: *dataBits,
		StopBits: *stopBits,
	}

	serialConn, err := serial2.NewSerial(*port, *baudRate, *dataBits, *stopBits)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer serialConn.Close()

	ui := tui.NewUi(sidebarData, inputCh, writeCh)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			m := <-inputCh
			err = serialConn.Write(m)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			buf, err := serialConn.Read()
			if err != nil {
				log.Fatal(err)
			}
			if len(buf) != 0 {
				writeCh <- string(buf)
			}
		}
	}()

	ui.Run(ctx)
}
