package systatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SystatusOptions struct {
	Prefix      string
	ExposeEnv   bool
	Healthcheck func(w http.ResponseWriter, r *http.Request)
}

type HealthResponse struct{}
type UptimeResponse struct {
	Systime string `json:"systime"`
	Uptime  string `json:"uptime"`
}
type CPURepsponse struct{}
type MemResponse struct{}
type EnvResponse struct {
	Env map[string]string `json:"env"`
}

func Enable(opts SystatusOptions) {
	if opts.Healthcheck == nil {
		http.HandleFunc(fmt.Sprintf("%s/health", opts.Prefix), handleHealth)
	} else {
		http.HandleFunc(fmt.Sprintf("%s/health", opts.Prefix), opts.Healthcheck)
	}

	http.HandleFunc(fmt.Sprintf("%s/uptime", opts.Prefix), handleUptime)
	http.HandleFunc(fmt.Sprintf("%s/cpu", opts.Prefix), handleCPU)
	http.HandleFunc(fmt.Sprintf("%s/mem", opts.Prefix), handleMem)
	http.HandleFunc(fmt.Sprintf("%s/disk", opts.Prefix), handleDisk)

	if opts.ExposeEnv {
		http.HandleFunc(fmt.Sprintf("%s/env", opts.Prefix), handleEnv)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {

}

func handleUptime(w http.ResponseWriter, r *http.Request) {

	res := UptimeResponse{}

	if runtime.GOOS == "windows" {
		// TODO Implement windows uptime
	} else {
		cmdoutput, err := exec.Command("/bin/uptime").Output()
		if err != nil {
			http.Error(w, "Could not exec uptime command on this machine", http.StatusInternalServerError)
			return
		}
		split := strings.Split(string(cmdoutput), " ")

		res.Systime = split[1]
		// Remove comma e.g 3:05,
		res.Uptime = strings.Split(split[4], ",")[0]

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
func handleCPU(w http.ResponseWriter, r *http.Request) {

}
func handleMem(w http.ResponseWriter, r *http.Request) {

}
func handleDisk(w http.ResponseWriter, r *http.Request) {

}
func handleEnv(w http.ResponseWriter, r *http.Request) {

	res := EnvResponse{}
	env := os.Environ()

	res.Env = make(map[string]string, len(env))

	for _, val := range env {
		split := strings.Split(val, "=")
		res.Env[split[0]] = split[1]
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(200)

}
