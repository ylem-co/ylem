package postgresql

import (
	"database/sql/driver"
	"fmt"
	"net"
	"net/url"
	"ylem_integrations/config"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

type Driver struct {
	Client *ssh.Client
}

type SshConnectionsPool struct {
	sync.RWMutex
	Drivers map[string]*Driver
}

func (s *SshConnectionsPool) Register(key string, driver *Driver) {
	s.Lock()
	s.Drivers[key] = driver
	s.Unlock()
}

func (s *SshConnectionsPool) Unregister(key string) error {
	s.Lock()
	if Driver, ok := s.Drivers[key]; ok {
		err := Driver.Client.Close()
		if err != nil {
			s.Unlock()

			return err
		}

		delete(s.Drivers, key)
		s.Unlock()

		return nil
	}

	fmt.Println("No such Driver with the key " + key)
	s.Unlock()

	return nil
}

var SshPool SshConnectionsPool

func InitSshPool() {
	SshPool = SshConnectionsPool{}
	SshPool.Drivers = make(map[string]*Driver, 0)
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	parsedDsnString, err := url.Parse(dsn)

	if err != nil {
		return nil, fmt.Errorf("invalid pgsql dsn: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(config.Cfg().Ssh.PrivateKey)
	if err != nil {
		return nil, err
	}

	sshHost := parsedDsnString.Query()["sshHost"][0]
	sshPort := parsedDsnString.Query()["sshPort"][0]
	sshUser := parsedDsnString.Query()["sshUser"][0]
	sshDriverId := parsedDsnString.Query()["sshDriverId"][0]
	sslMode := parsedDsnString.Query()["sslmode"][0]

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // @todo MUST HAVE: Should be fixed one
		Timeout: time.Duration(10) * time.Second, // probably should add the channel. Probably? Because it doesn't work otherwise XD
	}

	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), sshConfig)
	if err != nil {
		return nil, err
	}

	newDriver := &Driver{Client: sshcon}

	SshPool.Register(sshDriverId, newDriver)

	// Sanitize DSN from our alien query params
	dsn = strings.Replace(dsn, parsedDsnString.RawQuery, fmt.Sprintf("sslmode=%s", sslMode), 1)

	return pq.DialOpen(newDriver, dsn)
}


// Dial make socket connection via SSH.
func (d *Driver) Dial(network, address string) (net.Conn, error) {
	return d.Client.Dial(network, address)
}

// DialTimeout make socket connection via SSH.
func (d *Driver) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return d.Client.Dial(network, address)
}
