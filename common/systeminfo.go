package common

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadAvg() []string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		log.Fatal(err)
	}
	loadavg := strings.Split(string(data), " ")
	return loadavg
}

func MemoryInfo() map[string]int {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	memInfo := make(map[string]int)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) == 2 {
			num, _ := strconv.Atoi(strings.TrimSpace(strings.Split(fields[1], "kB")[0]))
			memInfo[fields[0]] = num / 1024
		}
	}

	return memInfo
}

func UpdateSystemInfo(meminfo1 *map[string]int, loadavg1 *[]string) {
	for {
		meminfo := MemoryInfo()
		loadavg := LoadAvg()
		*meminfo1 = meminfo
		*loadavg1 = loadavg
		time.Sleep(time.Second * 1)
	}
}
