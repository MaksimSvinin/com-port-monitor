package serial

import (
	"errors"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type Serial struct {
	options serial.OpenOptions
	conn    io.ReadWriteCloser
}

func NewSerial(port string, baudRate, dataBits, stopBits uint) (*Serial, error) {
	options := serial.OpenOptions{
		PortName:               port,
		BaudRate:               baudRate,
		DataBits:               dataBits,
		StopBits:               stopBits,
		MinimumReadSize:        0,
		InterCharacterTimeout:  100,
		ParityMode:             serial.PARITY_NONE,
		Rs485Enable:            false,
		Rs485RtsHighDuringSend: false,
		Rs485RtsHighAfterSend:  false,
	}

	f, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	return &Serial{
		options: options,
		conn:    f,
	}, nil
}

func (s *Serial) Write(data string) error {
	_, err := s.conn.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func (s *Serial) Read() ([]byte, error) {
	buf := make([]byte, 100)
	n, err := s.conn.Read(buf)
	if err != nil && err != io.EOF {
		return nil, errors.New(fmt.Sprintf("Error reading from serial port: %s", err))
	}
	buf = buf[:n]
	return buf, nil
}

func (s *Serial) Close() {
	s.conn.Close()
}
