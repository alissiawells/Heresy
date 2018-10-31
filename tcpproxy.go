package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
)

const SO_ORIGINAL_DST = 80

type Proxy struct {
	from    string
	fromTCP *net.TCPAddr
	done    chan struct{}
}

func NewProxy(from string) *Proxy {
	return &Proxy{
		from: from,
		done: make(chan struct{}),
	}
}

func (p *Proxy) Start() error {
	var err error
	p.fromTCP, err = net.ResolveTCPAddr("tcp", p.from)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", p.fromTCP)
	if err != nil {
		return err
	}
	go p.run(*listener)
	return nil
}

func (p *Proxy) Stop() {
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func getOriginalDst(clientConn *net.TCPConn) (ipv4 string, port uint16, newTCPConn *net.TCPConn, err error) {
	if clientConn == nil {
		err = errors.New("ERR: clientConn is nil")
		return
	}

	remoteAddr := clientConn.RemoteAddr()
	if remoteAddr == nil {
		err = errors.New("ERR: clientConn.fd is nil")
		return
	}

	newTCPConn = nil

	/*
    net.TCPConn.File() would cause the receiver's socket to be placed in a blocking mode.
    In order to place the new TCPConn in non-blocking mode, we need to take the file returned by .File(), 
    get the original destination by getsockopt(), create a new *net.TCPConn by net.TCPConn.FileConn(). 
    */
	clientConnFile, err := clientConn.File()
	if err != nil {
		return
	} else {
		clientConn.Close()
	}

	// Get original destination
	addr, err := syscall.GetsockoptIPv6Mreq(int(clientConnFile.Fd()), syscall.IPPROTO_IP, SO_ORIGINAL_DST)

	if err != nil {
		return
	}
	newConn, err := net.FileConn(clientConnFile)
	if err != nil {
		return
	}
	if _, ok := newConn.(*net.TCPConn); ok {
		newTCPConn = newConn.(*net.TCPConn)
		clientConnFile.Close()
	} else {
		errmsg := fmt.Sprintf("ERR: newConn is not a *net.TCPConn, instead it is: %T (%v)", newConn, newConn)
		err = errors.New(errmsg)
		return
	}

	ipv4 = itod(uint(addr.Multiaddr[4])) + "." +
		itod(uint(addr.Multiaddr[5])) + "." +
		itod(uint(addr.Multiaddr[6])) + "." +
		itod(uint(addr.Multiaddr[7]))
	port = uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])

	return
}

func (p *Proxy) run(listener net.TCPListener) {
	for {
		select {
		case <-p.done:
			return
		default:
			connection, err := listener.AcceptTCP()
			la := connection.LocalAddr()
			if la == nil {
				panic("Connection lost!")
			}
			fmt.Printf("Connectoin from %s\n", la.String())

			if err == nil {
				go p.handle(*connection)
			} else {
			}
		}
	}
}

func (p *Proxy) handle(connection net.TCPConn) {

	defer connection.Close()

	var clientConn *net.TCPConn
	ipv4, port, clientConn, err := getOriginalDst(&connection)
	if err != nil {
		panic(err)
	}
	connection = *clientConn

	dest := ipv4 + ":" + fmt.Sprintf("%d", port)
	addr, err := net.ResolveTCPAddr("tcp", dest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connectoin to %s\n", dest)
	remote, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}
	defer remote.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(*remote, connection, wg)
	go p.copy(connection, *remote, wg)
	wg.Wait()

}

func (p *Proxy) copy(from, to net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(&to, &from); err != nil {
			p.Stop()
			return
		}
	}
}

func itod(i uint) string {
	if i == 0 {
		return "0"
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; i > 0; i /= 10 {
		bp--
		b[bp] = byte(i%10) + '0'
	}

	return string(b[bp:])
}

func main() {
	NewProxy("localhost:1111").Start()
	fmt.Println("Server started.")
	select {}
}
