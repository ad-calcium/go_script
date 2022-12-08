package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func readFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	i := make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil
			}
		}
		i = append(i, strings.Trim(string(line), ""))
	}
	return i
}

func ssh_comm(ips string, pass string, comm string) {
	sshHost := ips
	sshUser := "root"
	sshPassword := pass
	sshType := "password" //password 或者 key
	sshPort := 22

	config := &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {

		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	}

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Println("[-] "+ips+" 创建ssh client 失败", err)
		return
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
  
	combo, err := session.CombinedOutput(comm)
	if err != nil {
		return
	}

	if len(combo) >= 10 {
		log.Println("[+] "+ips+" 命令执行成功，命令执行结果:\n", string(combo))
	} else {
		log.Println("[+] "+ips+" 命令执行成功，命令执行结果:", string(combo))
	}

}

func main() {
	var target string
	var comm string
	var pass string
	flag.StringVar(&pass, "p", "", "ssh密码")
	flag.StringVar(&comm, "c", "whoami", "执行的命令")
	flag.StringVar(&target, "f", "", "文件名")
	flag.Parse()
	if target != "" {
		ip := readFile(target)
		for _, ips := range ip {
			ssh_comm(ips, pass, comm)
			//fmt.Println(v)
		}
	} else {
		flag.PrintDefaults()
	}
}
