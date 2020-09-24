package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const (
	Path = "docs/paper-benchmarks/cleaned"
)

type Benchmark struct {
	name string
	logs []Log
}

type Log struct {
	validator int
	shard     int
	data      []Entry
	mean      float64
}

type Entry struct {
	time int64
	tx   int
	cpu  float64
}

func main() {
	benchmarks := parse()

	/*
		- For each benchmark, calculate the mean of every validators.
	*/

	for _, benchmark := range benchmarks {
		// For each log file, calculate the mean where cpu > 80%
		for i, log := range benchmark.logs {
			numEntries := 0
			sum := 0
			for _, entry := range log.data {
				if entry.cpu > 0.8 {
					numEntries++
					sum += entry.tx
				}
			}
			benchmark.logs[i].mean = float64(sum) / float64(numEntries)
		}

		// For each validator, add the means of its shards.
		means := make(map[int]float64)
		for _, log := range benchmark.logs {
			means[log.validator] += log.mean
		}

		var numValidators float64 = 0
		var sum float64 = 0
		for validator := range means {
			numValidators++
			sum += means[validator]
		}

		fmt.Printf("%s: %v\n", benchmark.name, sum / numValidators)
	}

}

func parse() []Benchmark {
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		log.Fatal(err)
	}

	var benchmarks []Benchmark

	for _, dir := range files {
		name := dir.Name()
		if !strings.HasPrefix(name, "final") {
			continue
		}

		var logs []Log
		bench := Benchmark{name, logs}

		logFiles, err := ioutil.ReadDir(path.Join(Path, name))
		if err != nil {
			log.Fatal(err)
		}

		for _, logFile := range logFiles {
			lName := logFile.Name()
			ident := strings.Split(lName, ".")[0]
			idents := strings.Split(ident, "-")
			validator, err := strconv.Atoi(idents[0])
			if err != nil {
				log.Fatal(err)
			}
			shard, err := strconv.Atoi(idents[1])
			if err != nil {
				log.Fatal(err)
			}

			var entries []Entry
			logItem := Log{
				validator: validator,
				shard:     shard,
				data:      entries,
			}

			content, err := ioutil.ReadFile(path.Join(Path, name, lName))
			if err != nil {
				log.Fatal(err)
			}
			contentStr := strings.TrimSpace(string(content))
			lines := strings.Split(contentStr, "\n")

			for _, line := range lines {
				match, _ := regexp.MatchString("\\d+,\\d+,\\d+\\.\\d*", line)

				if !match {
					continue
				}

				spl := strings.Split(line, ",")
				time, err := strconv.ParseInt(spl[0], 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				tx, err := strconv.Atoi(spl[1])
				if err != nil {
					log.Fatal(err)
				}
				cpu, err := strconv.ParseFloat(spl[2], 32)
				if err != nil {
					log.Fatal(err)
				}
				entry := Entry{time, tx, cpu}
				logItem.data = append(logItem.data, entry)
			}

			bench.logs = append(bench.logs, logItem)
		}
		benchmarks = append(benchmarks, bench)
	}

	return benchmarks
}
