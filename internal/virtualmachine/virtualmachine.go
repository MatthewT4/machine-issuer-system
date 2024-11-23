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
	filePath   = "/Users/mixalight/.ssh/id_ed25519" //change on another vm
	port       = "22"
	passPhrase = "cloud" //remove if another vm
)

type Metrics struct {
	uptime int64
	cpu    float64
	ram    float64
	mem    int64 //mb
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	publicKeyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	publicKey, err := ssh.ParsePrivateKeyWithPassphrase(publicKeyBytes, []byte(passPhrase))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return ssh.PublicKeys(publicKey), nil
}

func createConnection(ip string) (*ssh.Session, error) {
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
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	return session, nil
}

func parseAnswer(ans string) Metrics {
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
		uptime: time.Now().Unix() - uptime.Unix(),
		cpu:    cpu,
		ram:    ram,
		mem:    mem,
	}
	return metrics
}

func requestAndProcessMetrics(session *ssh.Session) (Metrics, error) {
	commands := []string{
		"uptime -s", //time
		"top -bn1 | grep \"Cpu(s)\" | sed \"s/.*, *\\([0-9.]*\\)%* id.*/\\1/\" | awk '{print 100 - $1}'", //CPU
		"free | awk 'NR==2{printf \"%.2f\", $3*100/$2 }'; echo",                                          //RAM
		"df -m / | awk 'NR==2{print $4}'",                                                                //MEM
	}
	script := ""
	for _, cmd := range commands {
		script += cmd + "\n"
	}
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err := session.Run(script)
	if err != nil {
		panic(fmt.Errorf("failed to execute script: %w", err))
	}
	ans := parseAnswer(stdout.String())
	return ans, nil
}

func GetMetrics(ip string) {
	session, err := createConnection(ip)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	metrics, err := requestAndProcessMetrics(session)
	if err != nil {
		panic(err)
	}
	fmt.Println(metrics.uptime, metrics.cpu, metrics.ram, metrics.mem)
}
