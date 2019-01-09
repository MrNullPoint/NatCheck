package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const STUN_SERVER = "118.89.173.65"
const STUN_PORT = 3478
const READ_TIME_OUT = 10

var NatType = map[string]string{
	"UDP_BLOCKED":            "Firewall blocks UDP",
	"PUBLIC_IP":              "Your IP is public on the Internet",
	"SYMMETRIC_UDP_FIREWALL": "Firewall that allows UDP out, and responses have to come back to the source of the request",
	"FULL_CONE":              "Full Cone Nat",
	"SYMMETRIC":              "Symmetric NAT",
	"PORT_RESTRICT":          "Port Restrict Cone NAT",
	"ADDR_RESTRICT":          "(Address) Restrict Cone NAT",
}

func connect(ip string, port int) *net.UDPConn {
	client := net.UDPAddr{}
	client.IP = net.ParseIP("0.0.0.0")

	server := net.UDPAddr{}

	if ip == "" {
		server.IP = net.ParseIP(STUN_SERVER)
	} else {
		server.IP = net.ParseIP(ip)
	}

	if port == 0 {
		server.Port = STUN_PORT
	} else {
		server.Port = port
	}

	conn, err := net.DialUDP("udp", &client, &server)
	check(err)

	return conn
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func read(conn *net.UDPConn) []byte {
	err := conn.SetReadDeadline(time.Now().Add(time.Duration(READ_TIME_OUT * time.Second)))
	check(err)

	data := make([]byte, 4096)
	n, _ := conn.Read(data)

	return data[:n]
}

func write(conn *net.UDPConn, b []byte) error {
	if _, err := conn.Write(b); err != nil {
		fmt.Println("Send Msg Failed: " + err.Error())
		return errors.New("Send Data Error")
	}
	return nil
}

func TestI(conn *net.UDPConn) *Message {
	fmt.Println("Running TestI: Check UDP Response")

	msg := NewBindRequest()
	err := write(conn, msg.ToBytes())
	check(err)

	var resp Message
	resp.FromBytes(read(conn))

	fmt.Println("- Local Address:", conn.LocalAddr().String())
	fmt.Println("- Mapped Address:", resp.GetMappedAddress().String())
	fmt.Println("- Changed Address:", resp.GetChangedAddress().String())

	return &resp
}

func TestII(conn *net.UDPConn) *Message {
	fmt.Println("Running TestII: Check Changing IP & Port")

	msg := NewChangeRequest(true, true)
	err := write(conn, msg.ToBytes())
	check(err)

	var resp Message
	resp.FromBytes(read(conn))

	fmt.Println("- Receive Changing Response")

	return &resp
}

func TestIII(conn *net.UDPConn) *Message {
	fmt.Println("Running TestIII: Check Connect Changed Address & Change Port")

	msg := NewChangeRequest(false, true)
	err := write(conn, msg.ToBytes())
	check(err)

	var resp Message
	resp.FromBytes(read(conn))

	fmt.Println("- Receive Changing Response")

	return &resp
}

func RunCheck(ip string) string {
	conn1 := connect(ip, 0)

	resp := TestI(conn1)
	caddr := resp.GetChangedAddress()
	m1addr := resp.GetMappedAddress()

	if resp.Len() == 0 {
		return NatType["UDP_BLOCKED"]
	}

	if conn1.LocalAddr().String() == resp.GetMappedAddress().String() {
		resp = TestII(conn1)
		if resp.Len() == 0 {
			return NatType["SYMMETRIC_UDP_FIREWALL"]
		}
		return NatType["PUBLIC_IP"]
	}

	resp = TestII(conn1)
	if resp.Len() > 0 {
		return NatType["FULL_CONE"]
	}

	conn2 := connect((net.IP)(caddr.Ip).String(), int(BytesToUint16(caddr.Port)))

	resp = TestI(conn2)
	if resp.Len() == 0 {
		return NatType["UDP_BLOCKED"]
	}

	m2addr := resp.GetMappedAddress()

	if m1addr.String() != m2addr.String() {
		return NatType["SYMMETRIC"]
	}

	resp = TestIII(conn1)

	if resp.Len() == 0 {
		return NatType["PORT_RESTRICT"]
	} else {
		return NatType["ADDR_RESTRICT"]
	}
}
