package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

var healthy = &cobra.Command{
	Short: "healty is a quick solution if u need to provide a healthcheck inside a docker container",
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "monitor linux process",
	Long:  "monitor a linux process and serve the healthy message as long as this process is up and running", //TODO document the long output,
	Run:   monitorProcess,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays version information on this happy",
	Run:   showVersion,
}

var (
	port string
	proc string

	CommitHash string
	VersionTag string
	BuildTime  string

	hc        HealthCheck
	arguments []string
)

func init() {
	processCmd.Flags().StringVarP(&port, "port", "p", "18080", "port to run the heatlh check on")

	healthy.AddCommand(versionCmd)
	healthy.AddCommand(processCmd)
}

func main() {
	healthy.Execute()
}

//HealthCheck represents all info needed to execute a health check on the system
// this struct will get initialized upon first run and used all consecutive runs
type HealthCheck struct {
	processes []*process.Process
}

func showVersion(cmd *cobra.Command, args []string) {
	log.Printf("The version is: %s; the commit hash is: %s. Build time is: %s", VersionTag, CommitHash, parseBuildTime(BuildTime).Format(time.RFC1123))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	//now run the healthcheck
	w.Header().Set("Content-Type", "application/json")

	//if the length of the healthcheck processes is 0 then we have not initialized the health check
	if len(hc.processes) == 0 {
		fmt.Println(len(hc.processes))
		// initialize the healthcheck
		initHealthCheck(arguments)
		// if the number of hc.processes is still 0 .. we dies
		if len(hc.processes) == 0 {

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(getJSONResponse("unable to find process"))

			os.Exit(1)
			return
		}
	}
	fmt.Println(len(hc.processes))

	// if the healthcheck returns true.. report success
	if hc.runHealthCheck() {

		w.WriteHeader(http.StatusOK)
		w.Write(getJSONResponse("Healthy"))

	} else {

		w.Write(getJSONResponse("Dead in the Water"))
		w.WriteHeader(http.StatusServiceUnavailable)

		os.Exit(1)

	}
}

//monitorProcess sets up the muxer
func monitorProcess(cmd *cobra.Command, args []string) {

	// copy the args to a globally declared variable
	arguments = args

	// setup the router
	router := mux.NewRouter().StrictSlash(true)

	// add the health route
	router.HandleFunc("/health", handleHealthCheck).Methods("GET")

	// start the listener
	http.ListenAndServe(":"+port, router)

}

//initHealthCheck will initialize the health check upon first run
func initHealthCheck(gs []string) {

	// range over the arguments to healthy
	for _, g := range gs {

		// compose the grep command
		cmd := "ps -ef |grep -v " + ownPid() + "| grep -i " + g + "|grep -v grep"

		//execute the grep command and catch the output .. we are ignoring the stderr here :-)
		out, _ := exec.Command("sh", "-c", cmd).Output()

		// convert the output to a string
		s := string(out)

		// split the output by line and range over the lines
		for _, l := range strings.Split(s, "\n") {
			// if we encounter an empty line .. ignore
			if len(l) != 0 {
				// split the output line on empty space
				spid := strings.Split(l, " ")
				// get the pid
				ipid, _ := strconv.Atoi(spid[3])
				// initialize the process object
				p, err := process.NewProcess(int32(ipid))
				// error handling
				if err != nil {
					fmt.Println(err)
				}

				//add the process to the health check
				hc.processes = append(hc.processes, p)
			}
		}
	}
}

//runHealthCheck returns true if  all the processes in the healthcheck are still alive and kicking
// guess what happens if they are not
func (hc *HealthCheck) runHealthCheck() bool {

	// loop over the processes
	for _, p := range hc.processes {
		// get the status and act accordingly
		// if we get an error of the process status is Zombie then we will fail
		if s, err := p.Status(); err != nil {
			return false
		} else if s == "Z" {
			fmt.Printf("%d has status %s", p.Pid, s)
			return false
		}
	}

	// all other responses are ok
	return true

}

//returns the pid for Healthy
func ownPid() string {
	return strconv.Itoa(os.Getpid())
}

//returns a results json object
func getJSONResponse(s string) []byte {
	result := struct {
		Status string
	}{
		Status: s,
	}

	payload, _ := json.Marshal(result)

	return payload
}

func parseBuildTime(BuildTime string) time.Time {
	// See https://pauladamsmith.com/blog/2011/05/go_time.html
	// See https://golang.org/pkg/time/#pkg-constants

	layout := "2006-01-02T15:04:05-0700"
	t, err := time.Parse(layout, BuildTime)

	if err != nil {
		log.Println(err)
	}

	return t
}
