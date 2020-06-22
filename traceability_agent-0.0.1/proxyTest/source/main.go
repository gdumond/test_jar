package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jcollins-axway/socks-proxy-check/proxy"
	"github.com/jcollins-axway/socks-proxy-check/socks"
)

var proxyAddress, hostURL, username, password string

const (
	NoAuth               = 0
	GSSAPIAuth           = 1
	UsernamePasswordAuth = 2
)

func main() {
	flag.StringVar(&proxyAddress, "proxy", "", "The socks 5 proxy address host[:port], port defaults to 1080")
	flag.StringVar(&hostURL, "url", "", "The URL to hit through the proxy")
	flag.StringVar(&username, "username", "username", "The username to use for auth via username/password, default is username")
	flag.StringVar(&password, "password", "password", "The password to use for auth via username/password, default is password")

	flag.Parse()

	// Check that variables were supplied
	if proxyAddress == "" {
		fmt.Fprintln(os.Stderr, "proxy is a required argument")
	}
	if hostURL == "" {
		fmt.Fprintln(os.Stderr, "url is a required argument")
	}

	// Append default port if one was not given
	if proxyParts := strings.Split(proxyAddress, ":"); len(proxyParts) == 1 {
		proxyAddress += ":1080"
	}

	// create a socks5 dialer, no auth
	fmt.Println("Attempting connection through Proxy without authentication")
	authType := NoAuth
	dialer := createSocksDialer(authType)
	fmt.Println("Connection to proxy made, connecting to url now")
	testConnection(authType, dialer)

	// create a socks5 dialer, username/password auth
	fmt.Println("Attempting connection through Proxy with username/password")
	authType = UsernamePasswordAuth
	dialer = createSocksDialer(authType)
	fmt.Println("Connection to proxy made, connecting to url now")
	testConnection(authType, dialer)

	// create a socks5 dialer, gssapi auth
	fmt.Println("Attempting connection through Proxy with gss-api")
	authType = GSSAPIAuth
	dialer = createSocksDialer(authType)
	fmt.Println("Connection to proxy made, connecting to url now")
	testConnection(authType, dialer)
}

// createSocksDialer - creates a socks dialer and returns the connectio
// authType
//		0 - no auth
//		2 - username/password
func createSocksDialer(authType int) proxy.Dialer {
	var dialer proxy.Dialer
	var err error

	// Connect with the auth type selected
	switch authType {
	case NoAuth:
		dialer, err = proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	case UsernamePasswordAuth:
		dialer, err = proxy.SOCKS5("tcp", proxyAddress, &proxy.Auth{Type: socks.AuthMethodUsernamePassword, User: username, Password: password}, proxy.Direct)
	case GSSAPIAuth:
		dialer, err = proxy.SOCKS5("tcp", proxyAddress, &proxy.Auth{Type: socks.AuthMethodGSSAPI}, proxy.Direct)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	return dialer
}

func testConnection(authType int, dialer proxy.Dialer) {
	// Attempt connection via proxy, then close
	conn, err := dialer.Dial("tcp", hostURL)
	if err != nil {
		switch authType {
		case NoAuth:
			fmt.Fprintln(os.Stderr, "connection without authentication failed", err)
		case UsernamePasswordAuth:
			fmt.Fprintln(os.Stderr, "connection with username/password authentication failed", err)
		case GSSAPIAuth:
			fmt.Fprintln(os.Stderr, "connection with gss-api authentication failed", err)
		}
	}

	if conn != nil {
		switch authType {
		case NoAuth:
			fmt.Println("connection without authentication was successful")
		case UsernamePasswordAuth:
			fmt.Println("connection with username/password authentication was successful")
		case GSSAPIAuth:
			fmt.Println("connection with gss-api authentication was successful")
		}
		conn.Close()
	}
}
