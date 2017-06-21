package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
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
	processCmd.Flags().StringVarP(&proc, "proc", "P", "", "process to check for by name")

	healthy.AddCommand(processCmd)
}

func main() {

	healthy.Execute()

}

func checkProcess(w http.ResponseWriter, r *http.Request) {

	pid := strconv.Itoa(os.Getpid())

	cmd := "ps -ef |grep -v " + pid + "| grep -i " + proc + "|grep -v grep"
	// log.Println(cmd)
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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", checkProcess).Methods("GET")
	http.ListenAndServe(":"+port, router)
}
