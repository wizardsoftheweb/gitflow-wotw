package main

import (
	"fmt"
	"time"
)

// func shell(args ...string) (string, error) {
// 	process := exec.Command(args[0], args[1:]...)
// 	combined, err := process.CombinedOutput()
// 	if nil != err {
// 		logrus.Fatal(err)
// 	}
// 	fmt.Println(string(combined))
// 	return string(combined), err
// }

// func pollConfig() []string {
// 	result, _ := shell("git", "config", "--global", "--list")
// 	return strings.Split(result, "\n")
// }

// func ConfigGrabber() {

// }

// func explodeLine(line string) []string {
// 	exploded := strings.Split(line, "=")
// 	value := strings.TrimSpace(strings.Join(exploded[1:len(exploded)-1], ""))
// 	exploded = strings.Split(strings.Join(exploded[:len(exploded)-1], ""), ".")
// 	section := strings.TrimSpace(exploded[0])
// 	key := strings.TrimSpace(exploded[len(exploded)-1])
// 	subsection := strings.TrimSpace(strings.Join(exploded[1:len(exploded)-1], ""))
// 	return []string{section, subsection, key, value}
// }

// func ConfigParser() {

// }

// func main() {
// 	config := pollConfig()
// 	explodeLine(config[0])

// 	// jobs := [][]string{
// 	// 	[]string{"git", "config", "--get", "gitflow.prefix.feature"},
// 	// 	[]string{"git", "config", "--get", "gitflow.prefix.hotfix"},
// 	// 	[]string{"git", "config", "--get", "gitflow.prefix.release"},
// 	// 	[]string{"git", "config", "--get", "gitflow.prefix.support"},
// 	// }
// 	// parallelStart := time.Now()
// 	// parallel := make(chan string)
// 	// go func() {
// 	// 	for _, job := range jobs {
// 	// 		ConfigInParallel(parallel, job)
// 	// 	}
// 	// }()
// 	// serialStart := time.Now()
// 	// serial := make(chan string)
// 	// go func() {
// 	// 	ConfigInSerial(serial, jobs)
// 	// }()
// 	// for {
// 	// 	select {
// 	// 	case output := <-parallel:
// 	// 		fmt.Println("Parallel", fmt.Sprintf("time since start: %s", time.Since(parallelStart)))
// 	// 		fmt.Println("Parallel", output)
// 	// 	case output := <-serial:
// 	// 		fmt.Println("Serial", fmt.Sprintf("time since start: %s", time.Since(serialStart)))
// 	// 		fmt.Println("Serial", output)
// 	// 	}
// 	// }
// }

func ConfigInParallel(parallel chan string, job []string) {
	fmt.Println(fmt.Sprintf("%s: %s", "parallel", time.Now()))

	result, _ := shell(job...)
	parallel <- result
}

func ConfigInSerial(serial chan string, jobs [][]string) {
	fmt.Println(fmt.Sprintf("%s: %s", "serial", time.Now()))
	for _, job := range jobs {
		result, _ := shell(job...)
		serial <- result
	}
}
