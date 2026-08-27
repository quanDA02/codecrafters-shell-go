package builtin

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// background jobs
type jobs struct {
	id     int
	name   string
	recent int
	status string
}

var jobMap = make(map[int]*jobs)

func Jobs(doneOnly bool) {
	if len(jobMap) < 1 {
		return
	}
	key := make([]int, 0)
	for _, job := range jobMap {
		key = append(key, job.id)
	}
	sort.Ints(key)
	recents := make([]int, 0)
	for _, job := range jobMap {
		recents = append(recents, job.recent)
	}
	sort.Ints(recents)
	for _, id := range key {
		job := jobMap[id]
		mark := " "
		if job.recent == recents[0] {
			mark = "+"
		}
		if len(jobMap) > 1 && job.recent == recents[1] {
			mark = "-"
		}
		if doneOnly {
			if job.status == "Done" {
				fmt.Printf("[%d]%s  %-24s%s\n", job.id, mark, job.status, job.name)
				delete(jobMap, job.id)
			}
		} else {
			fmt.Printf("[%d]%s  %-24s%s\n", job.id, mark, job.status, job.name)
			if job.status == "Done" {
				delete(jobMap, job.id)
			}
		}
	}
}

func IsBackground(args []string) ([]string, bool) {
	isBackground := false
	if args[len(args)-1] == "&" {
		isBackground = true
		args = args[0 : len(args)-1]
	}
	return args, isBackground
}

func CreateJob(args []string, cmd *exec.Cmd, input string) {
	jobID := 1
	for {
		_, exist := jobMap[jobID]
		if exist {
			jobID++
			continue
		}
		break
	}
	job := &jobs{
		id:     jobID,
		name:   input,
		recent: 0,
		status: "Running",
	}
	jobMap[jobID] = job
	for _, job := range jobMap {
		jobMap[job.id].recent += 1
	}
	fmt.Printf("[%d] %d\n", jobID, cmd.Process.Pid)
	go func(jobID int) {
		cmd.Wait()
		jobMap[jobID].status = "Done"
		jobMap[jobID].name = strings.Join(args, " ")
	}(jobID)
}
