package main

import (
	"bufio"
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/topologies/topology_util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Configuration
const (
	Interactive = false
	Verbose     = false
	Topology    = "merkleAWS28_4"
	StartupTime = 3 * time.Minute
	TestTime    = 12 * time.Minute

	/* Validators */
	NumShards               = 4
	Pooling                 = true
	PoolSize                = 64
	ValidatorLaunchTemplate = "lt-0e530c80ccc3334ab"

	/* Agents */
	AgentsLaunchTemplate = "lt-0664087648feb48a8" // us-west-1
	NumAgentsInstances   = 3 * 28 * NumShards
	NumWorkers           = 2048
)

type Validator struct {
	id    int
	shard int
	ip    string
}

var instanceIDs []string

func main() {
	topology_util.GenerateAll()
	topologiesFile := fmt.Sprintf("topologies/%s.json", Topology)
	common.InitAddresses(topologiesFile)

	numValidators := common.GetNumServers()

	log.Printf("Launching %d validators with %d shards (total: %d)\n", numValidators, NumShards, numValidators*NumShards)

	var validators []Validator

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

			id, ip := parse(out)
			instanceIDs = append(instanceIDs, id)
			injectIP(ip)
			validators = append(validators, Validator{
				id:    i,
				shard: j,
				ip:    ip,
			})
		}
	}

	pushChanges()

	if Interactive {
		log.Printf("→ wait until validators have started up.\n")
		fmt.Print("Press enter when done.")
		bufio.NewScanner(os.Stdin).Scan()
	} else {
		log.Println("waiting for startup...")
		time.Sleep(StartupTime)
	}

	/* Launch agents */
	log.Println("launching agents...")
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
	out := execute(args)
	ids := parseIDs(out)
	instanceIDs = append(instanceIDs, ids...)
	os.Remove(topologiesFile)

	if Interactive {
		log.Printf("→ wait until the test is finished.\n")
		log.Print("Press enter when done.")
		bufio.NewScanner(os.Stdin).Scan()
	} else {
		log.Println("testing...")
		time.Sleep(TestTime)
	}

	/* Retrieve logs from validators */
	log.Println("retrieving logs...")
	outputFolder := fmt.Sprintf("logs/%s-%d", Topology, time.Now().Unix())
	os.MkdirAll(outputFolder, os.ModePerm)

	for _, validator := range validators {
		cmd := exec.Command(
			"ssh",
			"-i",
			"~/.ssh/makesxi-us-west-1",
			"-oStrictHostKeyChecking=no",
			fmt.Sprintf("ec2-user@%s", validator.ip),
			"docker logs acfts_server_1",
		)
		dir, _ := os.Getwd()
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			log.Fatalf("could not get logs.. %s\n", err)
		}
		filename := fmt.Sprintf("%s/%d-%d.log", outputFolder, validator.id, validator.shard)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		f.Write(out)
		f.Close()
	}

	/* Terminate instances */
	log.Println("terminating instances...")
	args = []string{
		"ec2",
		"terminate-instances",
		"--instance-ids",
	}
	args = append(args, instanceIDs...)
	execute(args)

	log.Println("done")
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

func parse(output string) (string, string) {
	line := strings.Split(output, "\n")[1]
	id := strings.Split(line, "\t")[7]

	for i := 0; i < 10; i++ {
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
		if out != "" {
			fmt.Println(out)
			return strings.TrimSpace(id), out
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	panicIDs(instanceIDs)
	panic("after 10 tries, didn't manage to get ip. :(")
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
		panicIDs(instanceIDs)
		fmt.Println(out)
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}

func info(output string) {
	if Verbose {
		log.Println(output)
	}
}

var injectIPMutex sync.Mutex

func injectIP(ip string) {
	path := fmt.Sprintf("./topologies/%s.json", Topology)
	args := []string{
		"-i",
		fmt.Sprintf("0,/localhost/{s/localhost/%s/}", ip),
		path,
	}
	injectIPMutex.Lock()
	cmd := exec.Command("gsed", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	injectIPMutex.Unlock()
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

func parseIDs(out string) []string {
	lines := strings.Split(out, "\n")
	var ids []string
	for _, line := range lines {
		e := strings.Split(line, "\t")
		if e[0] == "INSTANCES" {
			ids = append(ids, strings.TrimSpace(e[7]))
		}
	}
	return ids
}

func panicIDs(ids []string) {
	fmt.Print("aws ec2 terminate-instances --instance-ids ")
	for _, id := range ids {
		fmt.Printf(" %s", id)
	}
	fmt.Println()
}
