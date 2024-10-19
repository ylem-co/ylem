package helpers

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"ylem_integrations/config"
	"time"
)

type SshContext struct {
	Host string
	Port uint16
	User string
}

const MySqlSSHTCPNetName = "mysql+ssh+tcp"
const MySQLTimeout = 20 * time.Second

func SSHDialContextFunc(ctx context.Context, addr string) (net.Conn, error) {
	if ctx.Value("config") == nil {
		fmt.Println("No SSH context, gonna dial directly")

		nd := net.Dialer{Timeout: MySQLTimeout}

		return nd.DialContext(ctx, "tcp", addr)
	}

	cfg := ctx.Value("config").(*SshContext)

	sshHost := cfg.Host
	sshPort := cfg.Port
	sshUser := cfg.User

	signer, err := ssh.ParsePrivateKey(config.Cfg().Ssh.PrivateKey)
	if err != nil {
		log.Errorf("could not parse a private ssh key: %s", err.Error())

		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // @todo MUST HAVE: Should be fixed one
		Timeout: time.Duration(3) * time.Second, // probably should add the channel. Probably? Because it doesn't work otherwise XD
	}

	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), sshConfig)
	if err != nil {
		return nil, err
	}

	conn, err := sshcon.Dial("tcp", addr)
	if err != nil {
		return nil, errors.New("mysql: " + err.Error())
	}

	go func() {
		select { //nolint:all
		case <-ctx.Done():
			conn.Close()
			sshcon.Close()
		}
	}()

	return conn, nil
}

func CreateSSHConfigurationContext(parent context.Context, Host string, Port uint16, User string) (context.Context, context.CancelFunc) {
	return context.WithCancel(
		context.WithValue(
			parent,
			"config", //nolint:all
			&SshContext{
				Host: Host,
				Port: Port,
				User: User,
			},
		),
	)
}
