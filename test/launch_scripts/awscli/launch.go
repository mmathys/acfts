package main

import (
	"bufio"
	"fmt"
	"github.com/mmathys/acfts/common"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Configuration
const (
	Verbose  = true
	Topology = "merkleAWS"

	/* Validators */
	NumShards               = 1
	Pooling                 = true
	PoolSize                = 2
	ValidatorLaunchTemplate = "lt-0e530c80ccc3334ab" // us-west-1
	//ValidatorLaunchTemplate = "lt-046c9ee55757a7a33" // eu-central-1

	/* Agents */
	AgentsLaunchTemplate = "lt-0664087648feb48a8" // us-west-1
	NumAgentsInstances   = 2
	NumWorkers           = 1
)

func main() {
	common.InitAddresses(fmt.Sprintf("topologies/%s.json", Topology))

	numValidators := common.GetNumServers()

	fmt.Printf("Launching %d validators which %d shards (total: %d)\n", numValidators, NumShards, numValidators*NumShards)

	/* Launch validators */
	for i := 0; i < numValidators; i++ {
		for j := 0; j < NumShards; j++ {
			info(fmt.Sprintf("Launching validator %d, shard %d...", i, j))
			filename := fmt.Sprintf("validator-%s-%d-%d", Topology, i, j)
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			f.WriteString(validatorConfig(i, j))

			args := []string{
				"ec2",
				"run-instances",
				"--launch-template",
				fmt.Sprintf("LaunchTemplateId=%s", ValidatorLaunchTemplate),
				"--user-data",
				fmt.Sprintf("file://%s", filename),
			}
			info(fmt.Sprintf("%v", args))
			out := execute(args)

			fmt.Printf("%d,%d: %s\n", i, j, getIp(out))
		}
	}

	fmt.Printf("→ Update %s.json validator addresses accordingly and push, then wait until validators have started up.\n", Topology)
	fmt.Print("Press enter when done.")
	bufio.NewScanner(os.Stdin).Scan()

	/* Launch agents */
	info("Launching agents...\n")
	filename := fmt.Sprintf("agents-%s", Topology)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	f.WriteString(agentsConfig())

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
	fmt.Println(out)
	fmt.Println("done :)")
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
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if re.MatchString(output) {
		return re.FindString(output)
	} else {
		panic("could not find any ip addresses in output")
	}
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
