package remote

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	config *SSHConfig
	client *ssh.Client
}

func NewSSHClient(config *SSHConfig) *SSHClient {
	return &SSHClient{
		config: config,
	}
}

func (c *SSHClient) Connect() error {
	sshConfig := &ssh.ClientConfig{
		User: c.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For simplicity - should use proper verification in production
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	c.client = client
	return nil
}

func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *SSHClient) DownloadFile(remotePath, localPath string) error {
	if c.client == nil {
		return fmt.Errorf("not connected to server")
	}

	// Create SFTP session
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Create local directory if it doesn't exist
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Use cat command to read remote file content
	cmd := fmt.Sprintf("cat %s", remotePath)
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Copy content from remote to local
	_, err = io.Copy(localFile, stdout)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func (c *SSHClient) ListLogFiles(logDir string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("ls -la %s/*.log 2>/dev/null || echo 'No log files found'", logDir)
	output, err := session.Output(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return []string{string(output)}, nil
}

func (c *SSHClient) CheckConnection() error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Simple command to test connection
	_, err = session.Output("echo 'connection test'")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

func TestConnection(config *SSHConfig) error {
	client := NewSSHClient(config)
	
	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()

	return client.CheckConnection()
}