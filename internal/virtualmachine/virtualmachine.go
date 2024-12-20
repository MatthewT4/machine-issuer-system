package virtualmachine

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	port       = "22"
	passPhrase = "cloud" //remove if another vm
	loginPass  = "5625bf05e16df92105b9ada132add04a"
)

type Metrics struct {
	Uptime int64
	CPU    float64
	RAM    float64
	MEM    int64 //mb
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	publicKeyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	publicKey, err := ssh.ParsePrivateKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return ssh.PublicKeys(publicKey), nil
}

func CreateConnection(ip string, filePath string) (*ssh.Session, error) {
	publicKey, err := publicKeyFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh public key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: "mixalight",
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password(loginPass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		conn, err = ssh.Dial("tcp", ip+":"+port, config)
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	return session, nil
}

func parseMetricsAnswer(ans string) Metrics {
	ans = strings.TrimSpace(ans)
	lines := strings.Split(ans, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	uptime, _ := time.Parse("2006-01-02 15:04:05", result[0])
	cpu, _ := strconv.ParseFloat(result[1], 2)
	ram, _ := strconv.ParseFloat(result[2], 2)
	mem, _ := strconv.ParseInt(result[3], 10, 32)
	metrics := Metrics{
		Uptime: time.Now().Unix() - uptime.Unix(),
		CPU:    cpu,
		RAM:    ram,
		MEM:    mem,
	}
	return metrics
}

func GetMetrics(session *ssh.Session, commands []string) (Metrics, error) {
	script := ""
	for _, cmd := range commands {
		script += cmd + "\n"
	}
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err := session.Run(script)
	if err != nil {
		return Metrics{}, fmt.Errorf("failed to execute script: %w", err)
	}
	ans := parseMetricsAnswer(stdout.String())
	return ans, nil
}

func ExecuteCommandsOnVirtualMachine(session *ssh.Session, commands []string) error {
	script := ""
	for _, cmd := range commands {
		script += cmd + "\n"
	}
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err := session.Run(script)
	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}
	return nil
}
