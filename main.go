package main

import (
	"log"
	"net"
	"time"
)

func main() {
	ip := getIP().String()
	log.Printf("external ip: %v", ip)

	go probe(ip)

	log.Printf("listen on %v:9000", "0.0.0.0")

	listener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			time.Sleep(10 * time.Second)
			conn.Close()
		}(conn)
	}
}

func probe(ip string) {
	log.Printf("probing %v", ip)

	ipAddr, err := net.ResolveIPAddr("ip4:tcp", ip)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenIP("ip4:tcp", ipAddr)
	if err != nil {
		panic(err)
	}

	pending := make(map[string]uint32)

	var buffer [4096]byte
	for {
		n, addr, err := conn.ReadFrom(buffer[:])
		if err != nil {
			log.Printf("conn.ReadFrom() error: %v", err)
			continue
		}

		pkt := &tcpPacket{}
		data, err := pkt.decode(buffer[:n])
		if err != nil {
			log.Printf("tcp packet parse error: %v", err)
			continue
		}

		if pkt.DestPort != 9000 {
			continue
		}

		log.Printf("tcp packet: %+v, flag: %v, data: %v, addr: %v", pkt, pkt.FlagString(), data, addr)

		if pkt.Flags&SYN != 0 {
			pending[addr.String()] = pkt.Seq + 1
			continue
		}
		if pkt.Flags&RST != 0 {
			panic("RST received")
		}
		if pkt.Flags&ACK != 0 {
			if seq, ok := pending[addr.String()]; ok {
				log.Println("connection established")
				delete(pending, addr.String())

				badPkt := &tcpPacket{
					SrcPort:    pkt.DestPort,
					DestPort:   pkt.SrcPort,
					Ack:        seq,
					Seq:        pkt.Ack - 100000,      // Bad: seq out-of-window
					Flags:      (5 << 12) | PSH | ACK, // Offset and Flags  oooo000F FFFFFFFF (o:offset, F:flags)
					WindowSize: pkt.WindowSize,
				}

				data := []byte("boom!!!")
				remoteIP := net.ParseIP(addr.String())
				localIP := net.ParseIP(conn.LocalAddr().String())
				_, err := conn.WriteTo(badPkt.encode(localIP, remoteIP, data[:]), addr)
				if err != nil {
					log.Printf("conn.WriteTo() error: %v", err)
				}
			}
		}
	}
}

func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
