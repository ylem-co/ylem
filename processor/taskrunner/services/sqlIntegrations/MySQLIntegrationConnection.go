package sqlIntegrations

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
	"database/sql"
	"ylem_taskrunner/config"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type MySQLSQLIntegrationConnection struct {
	Host       string
	Port       uint16
	User       string
	Password   string
	Database   string
	Context    context.Context
	CancelFunc context.CancelFunc
	DB         *sql.DB
}

const MySqlSSHTCPNetName = "mysql+ssh+tcp"
const MySQLTimeout = 20 * time.Second

func init() {
	mysql.RegisterDialContext(MySqlSSHTCPNetName, sshDialContextFunc)
}

func (s *MySQLSQLIntegrationConnection) Open() error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", s.User, s.Password, MySqlSSHTCPNetName, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}

	s.Context = context.Background()
	s.CancelFunc = nil

	return nil
}

func (s *MySQLSQLIntegrationConnection) OpenSsh(Host string, Port uint16, User string) error {
	return s.OpenSshContext(context.Background(), Host, Port, User)
}

func (s *MySQLSQLIntegrationConnection) OpenSshContext(ctx context.Context, Host string, Port uint16, User string) error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", s.User, s.Password, MySqlSSHTCPNetName, s.Host, s.Port, s.Database))
	if err != nil {
		return err
	}

	s.Context, s.CancelFunc = createSSHConfigurationContext(
		ctx,
		Host,
		Port,
		User,
	)

	return nil
}

func (s *MySQLSQLIntegrationConnection) Close() error {
	err := s.DB.Close()
	if s.CancelFunc != nil {
		s.CancelFunc()
	}

	return err
}

func (s *MySQLSQLIntegrationConnection) Test() error {
	return s.DB.PingContext(s.Context)
}

func (s *MySQLSQLIntegrationConnection) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(s.Context, query)
}

func (s *MySQLSQLIntegrationConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(s.Context, query, args...)
}

func (s *MySQLSQLIntegrationConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(s.Context, query, args...)
}

type SshContext struct {
	Host string
	Port uint16
	User string
}

func sshDialContextFunc(ctx context.Context, addr string) (net.Conn, error) {
	if ctx.Value("config") == nil {
		log.Debug("No SSH context, gonna dial directly")

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
		Timeout:         time.Second * 3,
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

func createSSHConfigurationContext(parent context.Context, Host string, Port uint16, User string) (context.Context, context.CancelFunc) {
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
