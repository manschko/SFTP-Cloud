package sftp

import (
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Connection represents an SFTP connection with credentials
type Connection struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewConnection creates a new SFTP connection configuration
func NewConnection(host string, port int, username, password string) *Connection {
	return &Connection{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// Connect establishes a connection to the SFTP server
func (c *Connection) Connect() (*sftp.Client, error) {
	// Configure SSH client
	sshConfig := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not recommended for production
		Timeout:         15 * time.Second,
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}

	return client, nil
}

// TestConnection attempts to connect to verify credentials
func (c *Connection) TestConnection() error {
	client, err := c.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}
