package gotik

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/jjcinaz/gotik/proto"
)

// Client is a RouterOS API client.
type Client struct {
	Queue int

	rwc                  io.ReadWriteCloser
	serverName           string // dns name or IP address
	isTLS                bool
	useInsecureCleartext bool
	r                    proto.Reader
	w                    proto.Writer
	closing              bool
	async                bool
	nextTag              int64
	tags                 map[string]sentenceProcessor
	mu                   sync.Mutex
}

// NewClient returns a new Client over rwc. Login must be called.
func NewClient(rwc io.ReadWriteCloser) (*Client, error) {
	return &Client{
		rwc: rwc,
		r:   proto.NewReader(rwc),
		w:   proto.NewWriter(rwc),
	}, nil
}

// AllowInsecureCleartext -- With versions 6.43 or newer, RouterOS takes the user name and password (in cleartext)
// on the initial login command.  If the connection is not TLS, then this library will
// not send the cleartext password and will fall back to the old MD5 challenge
// method (only supported on versions earlier than 6.45.1).  If the target is newer than
// 6.45.1, MD5 challenge will not work.  In such cases, the connection must be TLS
// or you must explicitly enable sending cleartext passwords over non-TLS by calling
// this function with a value of true.
func (c *Client) AllowInsecureCleartext(value bool) {
	c.useInsecureCleartext = value
}

func fqRouterIP(ip string, useTLS bool) string {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if useTLS {
			return ip + ":8729"
		}
		return ip + ":8728"
	}
	return ip
}

// Dial connects and logs in to a RouterOS device.
func Dial(address, username, password string) (*Client, error) {
	var (
		conn net.Conn
		c    *Client
		err  error
	)
	conn, err = net.Dial("tcp", fqRouterIP(address, false))
	if err != nil {
		return nil, err
	}
	c, err = newClientAndLogin(conn, username, password, false)
	if err == nil {
		c.serverName = address
	}
	return c, err
}

// DialTimeout connects to and logs in to a RouterOS device.
func DialTimeout(address, username, password string, timeout time.Duration) (*Client, error) {
	var (
		conn net.Conn
		c    *Client
		err  error
	)
	conn, err = net.DialTimeout("tcp", fqRouterIP(address, false), timeout)
	if err != nil {
		return nil, err
	}
	c, err = newClientAndLogin(conn, username, password, false)
	if err == nil {
		c.serverName = address
	}
	return c, err
}

// DialTLS connects to and logs in to a RouterOS device using TLS.
func DialTLS(address, username, password string, tlsConfig *tls.Config) (*Client, error) {
	var (
		conn net.Conn
		c    *Client
		err  error
	)
	conn, err = tls.Dial("tcp", fqRouterIP(address, true), tlsConfig)
	if err != nil {
		return nil, err
	}
	c, err = newClientAndLogin(conn, username, password, true)
	if err == nil {
		c.serverName = address
	}
	return c, err
}

// DialTLSTimeout connects to and logs in to a RouterOS device using TLS with an optional timeout
func DialTLSTimeout(address, username, password string, tlsConfig *tls.Config, timeout time.Duration) (*Client, error) {
	var (
		conn net.Conn
		c    *Client
		err  error
	)
	d := new(net.Dialer)
	d.Timeout = timeout
	conn, err = tls.DialWithDialer(d, "tcp", fqRouterIP(address, true), tlsConfig)
	if err != nil {
		return nil, err
	}
	c, err = newClientAndLogin(conn, username, password, true)
	if err == nil {
		c.serverName = address
	}
	return c, err
}

func newClientAndLogin(rwc io.ReadWriteCloser, username, password string, isTLS bool) (*Client, error) {
	c, err := NewClient(rwc)
	if err != nil {
		_ = rwc.Close()
		return nil, err
	}
	c.isTLS = isTLS
	c.useInsecureCleartext = true
	err = c.Login(username, password)
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

// CurrentAddress returned the DNS name or IP address of the router to which this client is working against
func (c *Client) CurrentAddress() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.serverName
}

// Close closes the connection to the RouterOS device.
func (c *Client) Close() {
	c.mu.Lock()
	if c.closing {
		c.mu.Unlock()
		return
	}
	c.closing = true
	c.mu.Unlock()
	_ = c.rwc.Close()
}

// Login runs the /login command. Dial and DialTLS call this automatically.
func (c *Client) Login(username, password string) error {
	var (
		r   *Reply
		err error
		b   []byte
	)
	if c.isTLS || c.useInsecureCleartext {
		// RouterOS post v6.43 has new authentication method and wants the login/pass (in cleartext) on
		// the first command.  We only do this if it is a TLS connection.
		r, err = c.Run("/login", "=name="+username, "=password="+password)
	} else {
		r, err = c.Run("/login")
	}
	if err != nil {
		return err
	}
	ret, ok := r.Done.Map["ret"]
	if !ok {
		// if we didn't get a =ret= in the response, then we assume it's a new login method (post 6.45.1)
		return nil
	}
	b, err = hex.DecodeString(ret)
	if err != nil {
		return fmt.Errorf("RouterOS: /login: invalid ret (challenge) hex string received: %s", err)
	}

	r, err = c.Run("/login", "=name="+username, "=response="+c.challengeResponse(b, password))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) challengeResponse(cha []byte, password string) string {
	h := md5.New()
	h.Write([]byte{0})
	_, _ = io.WriteString(h, password)
	h.Write(cha)
	return fmt.Sprintf("00%x", h.Sum(nil))
}
