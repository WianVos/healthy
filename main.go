package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

var port string
var proc string

func init() {
	processCmd.Flags().StringVarP(&port, "port", "p", "18080", "port to run the heatlh check on")

	healthy.AddCommand(processCmd)
}

func main() {

	healthy.Execute()

}

func checkProcess(w http.ResponseWriter, r *http.Request, p string) {

	pid := strconv.Itoa(os.Getpid())

	cmd := "ps -ef |grep -v " + pid + "| grep -i " + p + "|grep -v grep"
	log.Println(cmd)
	//_, err := exec.Command("sh", "-c", cmd).CombinedOutput()

	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		// log.Println(err)
		// log.Println(out)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Process: "+proc+" unavailable")
		os.Exit(1)
	}
	log.Println(err)
	log.Println(string(out))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Healthy")
}

func monitorProcess(cmd *cobra.Command, args []string) {
	for _, a := range args {
		fmt.Println(a)
		findProcessByName(a)
	}

	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/health", checkProcess, "blah").Methods("GET")
	http.ListenAndServe(":"+port, router)
}

func findProcessByName(n string) []process.Process {

	procs := make([]process.Process, [...])

	cmd := "ps -ef |grep -v " + ownPid() + "| grep -i " + n + "|grep -v grep"

	out, _ := exec.Command("sh", "-c", cmd).Output()
	s := string(out)

	for _, l := range strings.Split(s, "\n") {
		spid := strings.Split(l, " ")
		//fmt.Println(pid[3])
		ipid, _ := strconv.Atoi(spid[2])
		procs = append(procs, process.NewProcess(int32(ipid)))
	}

} 

func ownPid() string {
	return strconv.Itoa(os.Getpid())
}
