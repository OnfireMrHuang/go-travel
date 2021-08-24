package func_options

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Serve struct {
	Addr string
	Port int
	Acl *tls.Config
	Timeout time.Duration
	Protocol string
	MaxConns int
}

type Option func(*Serve)

func Protocol(p string) Option {
	return func(serve *Serve) {
		serve.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(serve *Serve) {
		serve.Timeout = timeout
	}
}

func MaxConns(maxConns int) Option {
	return func(serve *Serve) {
		serve.MaxConns = maxConns
	}
}

func TLS(tls *tls.Config) Option {
	return func(serve *Serve) {
		serve.Acl = tls
	}
}

func NewServe(addr string,port int,options ...func(*Serve)) (*Serve,error) {
	srv := Serve{
		Addr: addr,
		Port: port,
		Protocol: "tcp",
		Timeout: 30 * time.Second,
		MaxConns: 1000,
		Acl: nil,
	}
	for _,option := range options {
		option(&srv)
	}
	return &srv,nil
}

func Demo()  {
	s1,_ := NewServe("localhost",1024)
	fmt.Println(s1)
	s2,_ := NewServe("localhost",1024,Protocol("udp"))
	fmt.Println(s2)
	s3,_ := NewServe("localhost",8080,Timeout(300 * time.Second),MaxConns(1000))
	fmt.Println(s3)
}


