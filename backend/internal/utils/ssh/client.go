package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// ========================================
// Types
// ========================================

type ConnectParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type CommandParams struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Client struct {
	Host     string
	Port     int
	Username string
	Password string

	Client  *ssh.Client
	Session *ssh.Session
}

// ========================================
// Errors
// ========================================
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %s: %s - %v", e.Code, e.Message, e.Cause)
}

const (
	ErrConnectionFailed = "SSH_CONNECTION_FAILED"
	ErrSessionFailed    = "SSH_SESSION_FAILED"
	ErrCommandFailed    = "SSH_COMMAND_FAILED"
)

// ========================================
// Instance
// ========================================

func NewClient(params ConnectParams) *Client {
	return &Client{
		Host:     params.Host,
		Port:     params.Port,
		Username: params.Username,
		Password: params.Password,
	}
}

// ========================================
// Methods
// ========================================

// Connect establishes an SSH connection to the specified host using the provided credentials.
func (c *Client) Connect() error {
	var err error

	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.Client, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		return &Error{
			Code:    ErrConnectionFailed,
			Message: "Failed to connect to SSH server",
			Cause:   err,
		}
	}

	return nil
}

// Execute executes a command on the remote SSH server and returns the output or error.
func (c *Client) Execute(params CommandParams) (*CommandResponse, error) {
    var err error

    if c.Client == nil {
        return nil, &Error{
            Code:    ErrConnectionFailed,
            Message: "SSH client is not connected",
        }
    }

    session, err := c.Client.NewSession()
    if err != nil {
        return nil, &Error{
            Code:    ErrSessionFailed,
            Message: "Failed to create SSH session",
            Cause:   err,
        }
    }
    defer session.Close()

    output, err := session.CombinedOutput(params.Command)
    if err != nil {
        return &CommandResponse{
            Error: string(output),
        }, &Error{
            Code:    ErrCommandFailed,
            Message: "Failed to execute command on SSH server",
            Cause:   err,
        }
    }

    return &CommandResponse{
        Output: string(output),
    }, nil
}

func (c *Client) Close() error {
    if c.Client != nil {
        return c.Client.Close()
    }
    return nil
}
