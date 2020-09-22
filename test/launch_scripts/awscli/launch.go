package main

import (
	"bufio"
	"fmt"
	"github.com/mmathys/acfts/common"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Configuration
const (
	Verbose     = false
	Topology    = "merkleAWS4"
	Interactive = false
	StartupTime = 3 * time.Minute
	TestTime    = 6 * time.Minute

	/* Validators */
	NumShards               = 1
	Pooling                 = true
	PoolSize                = 64
	ValidatorLaunchTemplate = "lt-0e530c80ccc3334ab" // us-west-1
	//ValidatorLaunchTemplate = "lt-046c9ee55757a7a33" // eu-central-1

	/* Agents */
	AgentsLaunchTemplate = "lt-0664087648feb48a8" // us-west-1
	NumAgentsInstances   = 16
	NumWorkers           = 4096
)

func main() {
	topologiesFile := fmt.Sprintf("topologies/%s.json", Topology)
	common.InitAddresses(topologiesFile)

	numValidators := common.GetNumServers()

	fmt.Printf("Launching %d validators which %d shards (total: %d)\n", numValidators, NumShards, numValidators*NumShards)

	var validatorIPs []string

	/* Launch validators */
	for i := 0; i < numValidators; i++ {
		for j := 0; j < NumShards; j++ {
			info(fmt.Sprintf("Launching validator %d, shard %d...", i, j))
			filename := fmt.Sprintf("aws-validator-%s-%d-%d", Topology, i, j)
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			f.WriteString(validatorConfig(i, j))
			f.Close()

			args := []string{
				"ec2",
				"run-instances",
				"--launch-template",
				fmt.Sprintf("LaunchTemplateId=%s", ValidatorLaunchTemplate),
				"--user-data",
				fmt.Sprintf("file://%s", filename),
			}
			out := execute(args)

			ip := getIp(out)
			injectIP(ip)
			validatorIPs = append(validatorIPs, ip)
		}
	}

	pushChanges()

	if Interactive {
		fmt.Printf("→ wait until validators have started up.\n")
		fmt.Print("Press enter when done.")
		bufio.NewScanner(os.Stdin).Scan()
	} else {
		fmt.Println("waiting for startup...")
		time.Sleep(StartupTime)
	}

	/* Launch agents */
	fmt.Println("Launching agents...")
	filename := fmt.Sprintf("aws-agents-%s", Topology)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	f.WriteString(agentsConfig())
	f.Close()

	args := []string{
		"ec2",
		"run-instances",
		"--launch-template",
		fmt.Sprintf("LaunchTemplateId=%s", AgentsLaunchTemplate),
		"--count",
		strconv.Itoa(NumAgentsInstances),
		"--user-data",
		fmt.Sprintf("file://%s", filename),
	}
	execute(args)
	os.Remove(topologiesFile)

	if Interactive {
		fmt.Printf("→ wait until the test is finished.\n")
		fmt.Print("Press enter when done.")
		bufio.NewScanner(os.Stdin).Scan()
	} else {
		fmt.Println("testing...")
		time.Sleep(TestTime)
	}

	/* Retrieve logs from validators */
	fmt.Println("retrieving logs...")
	os.Mkdir("logs", os.ModePerm)
	for _, ip := range validatorIPs {
		cmd := exec.Command(
			"ssh",
			"-i",
			"~/.ssh/makesxi-us-west-1",
			"-oStrictHostKeyChecking=no",
			fmt.Sprintf("ec2-user@%s", ip),
			"docker logs acfts_server_1",
		)
		dir, _ := os.Getwd()
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("could not get logs.. %s\n", err)
		}
		filename := fmt.Sprintf("logs/%s.log", ip)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		f.Write(out)
		f.Close()
	}

	fmt.Println("done")
}

func validatorConfig(validator int, shard int) string {
	base, err := ioutil.ReadFile("test/launch_scripts/awscli/base.sh")
	if err != nil {
		panic(err)
	}
	baseStr := string(base)
	baseStr += "# machine specific config\n"
	addr := common.GetServers()[validator]
	baseStr += export("ADDRESS", fmt.Sprintf("%x", addr))
	baseStr += export("TOPOLOGY", Topology)
	if Pooling {
		baseStr += export("POOLING", "true")
	} else {
		baseStr += export("POOLING", "false")
	}
	baseStr += export("POOL_SIZE", strconv.Itoa(PoolSize))
	baseStr += export("INSTANCE", strconv.Itoa(shard))
	baseStr += "/usr/local/bin/docker-compose up --detach --build server\n"
	return baseStr
}

func agentsConfig() string {
	base, err := ioutil.ReadFile("test/launch_scripts/awscli/base.sh")
	if err != nil {
		panic(err)
	}
	baseStr := string(base)
	baseStr += "# machine specific config\n"
	baseStr += export("TOPOLOGY", Topology)
	baseStr += export("NUM_WORKERS", strconv.Itoa(NumWorkers))
	baseStr += "/usr/local/bin/docker-compose up --detach --build agents\n"
	return baseStr
}

func export(key string, value string) string {
	return fmt.Sprintf("export %s=%s\n", key, value)
}

func getIp(output string) string {
	line := strings.Split(output, "\n")[1]
	id := strings.Split(line, "\t")[7]

	time.Sleep(time.Second)
	args := []string{
		"ec2",
		"describe-instances",
		"--query",
		"Reservations[*].Instances[*].PublicIpAddress",
		"--instance-ids",
		id,
	}
	out := execute(args)
	out = strings.TrimSpace(out)
	fmt.Println(out)
	return out
}

func execute(args []string) string {
	info(fmt.Sprintf("args: %s", args))
	cmd := exec.Command("aws", args...)
	// cmd config
	dir, _ := os.Getwd()
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, "AWS_DEFAULT_OUTPUT=text")
	cmd.Env = append(cmd.Env, "AWS_PAGER=")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}

func info(output string) {
	if Verbose {
		fmt.Println(output)
	}
}

func injectIP(ip string) {
	path := fmt.Sprintf("./topologies/%s.json", Topology)
	args := []string{
		"-i",
		fmt.Sprintf("0,/localhost/{s/localhost/%s/}", ip),
		path,
	}
	cmd := exec.Command("gsed", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func pushChanges() {
	cmd := exec.Command("git", "add", "topologies/*")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("git add failed with %s\n", err)
	}

	cmd = exec.Command("git", "commit", "-m", "updated topology")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("git commit failed with %s\n", err)
	}

	cmd = exec.Command("git", "push")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("git push failed with %s\n", err)
	}
}
