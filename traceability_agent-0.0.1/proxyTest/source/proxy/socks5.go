// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"context"
	"io"
	"net"

	"github.com/jcollins-axway/socks-proxy-check/socks"
)

// SOCKS5 returns a Dialer that makes SOCKSv5 connections to the given
// address with an optional username and password.
// See RFC 1928 and RFC 1929.
func SOCKS5(network, address string, auth *Auth, forward Dialer) (Dialer, error) {
	d := socks.NewDialer(network, address)
	if forward != nil {
		if f, ok := forward.(ContextDialer); ok {
			d.ProxyDial = func(ctx context.Context, network string, address string) (net.Conn, error) {
				return f.DialContext(ctx, network, address)
			}
		} else {
			d.ProxyDial = func(ctx context.Context, network string, address string) (net.Conn, error) {
				return dialContext(ctx, forward, network, address)
			}
		}
	}
	if auth != nil {
		if auth.Type == socks.AuthMethodUsernamePassword {
			up := socks.UsernamePassword{
				Username: auth.User,
				Password: auth.Password,
			}
			d.AuthMethods = []socks.AuthMethod{
				socks.AuthMethodUsernamePassword,
			}
			d.Authenticate = up.Authenticate
		}
		if auth.Type == socks.AuthMethodGSSAPI {
			d.AuthMethods = []socks.AuthMethod{
				socks.AuthMethodGSSAPI,
			}
			d.Authenticate = func(ctx context.Context, rw io.ReadWriter, auth socks.AuthMethod) error {
				b := []byte{0x01}
				// TODO(mikio): handle IO deadlines and cancelation if
				// necessary
				if _, err := rw.Write(b); err != nil {
					return err
				}
				// if _, err := io.ReadFull(rw, b); err != nil {
				// 	return err
				// }
				return nil
			}
		}
	}

	return d, nil
}
