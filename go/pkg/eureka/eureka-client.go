package eureka

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/procyon-projects/chrono"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// ? ==================== Structs ==================== ?

type AppRegistrationBody struct {
	Instance InstanceDetails `json:"instance"`
}

// * =========== *

type InstanceDetails struct {
	InstanceId       string         `json:"instanceId"`
	HostName         string         `json:"hostName"`
	App              string         `json:"app"`
	VipAddress       string         `json:"vipAddress"`
	SecureVipAddress string         `json:"secureVipAddress"`
	IpAddr           string         `json:"ipAddr"`
	Status           string         `json:"status"`
	Port             Port           `json:"port"`
	SecurePort       Port           `json:"securePort"`
	HealthCheckUrl   string         `json:"healthCheckUrl"`
	StatusPageUrl    string         `json:"statusPageUrl"`
	HomePageUrl      string         `json:"homePageUrl"`
	DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
}

// * =========== *

type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

// * =========== *

type Port struct {
	Port    int    `json:"$"`
	Enabled string `json:"@enabled"`
}

// ? ==================== Functions ==================== ?

// ScheduleHeartbeat sends a heartbeat every 25 seconds to keep track of the instance on the Eureka server
func ScheduleHeartbeat(appName string, appId string) chrono.ScheduledTask {

	taskScheduler := chrono.NewDefaultTaskScheduler()

	task, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		sendHeartbeat(appName, appId)
	}, 25*time.Second)

	if err != nil {
		log.Fatalln(err)
	}

	return task

}

// * =========== *

// RegisterApp register the instance on the Eureka server
func RegisterApp(appName string, appId string, port int) {

	log.Println("registering app on Eureka server")

	body := buildBody(appName, appId, port, "STARTING")

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)

	if err != nil {
		log.Fatalln(err)
	}

	server := viper.GetString("eureka.client.service-url.defaultZone")
	if server == "" {
		server = "http://localhost:8761/eureka"
	}

	resp, err := http.Post(server+"/apps/"+appName, "application/json", &buf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	responseBody, parseErr := io.ReadAll(resp.Body)

	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	if string(responseBody) != "" {
		log.Println(string(responseBody))
	}

}

// * =========== *

// UpdateAppStatus updates the status of the instance on the Eureka server
func UpdateAppStatus(appName string, appId string, port int, status string) {

	log.Println("updating app status")

	body := buildBody(appName, appId, port, status)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)

	if err != nil {
		log.Fatalln(err)
	}

	server := viper.GetString("eureka.client.service-url.defaultZone")
	if server == "" {
		server = "http://localhost:8761/eureka"
	}

	req, err := http.Post(server+"/apps/"+appName, "application/json", &buf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(req.Body)

	responseBody, parseErr := io.ReadAll(req.Body)

	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	if string(responseBody) != "" {
		log.Println(string(responseBody))
	}

}

// * =========== *

// DeleteApp deletes the Eureka server instance
func DeleteApp(appName string, appId string) {

	log.Println("deleting app from Eureka server")

	server := viper.GetString("eureka.client.service-url.defaultZone")
	if server == "" {
		server = "http://localhost:8761/eureka"
	}

	req, err := http.NewRequest(http.MethodDelete, server+"/apps/"+appName+"/"+appId, nil)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	responseBody, parseErr := io.ReadAll(resp.Body)

	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	if string(responseBody) != "" {
		log.Println(string(responseBody))
	}

}

// * =========== *

// sendHeartbeat sends a heartbeat to keep track of the instance on the Eureka server
func sendHeartbeat(appName string, appId string) {

	server := viper.GetString("eureka.client.service-url.defaultZone")
	if server == "" {
		server = "http://localhost:8761/eureka"
	}

	req, err := http.NewRequest("PUT", server+"/apps/"+appName+"/"+appId, nil)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	responseBody, parseErr := io.ReadAll(resp.Body)

	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	if string(responseBody) != "" {
		log.Println(string(responseBody))
	}

}

// * =========== *

// buildBody constructs the body of the request to register the instance on the Eureka server
func buildBody(appName string, appId string, port int, status string) *AppRegistrationBody {

	hostname := viper.GetString("server.hostname")
	httpport := port

	if hostname == "" {

		hostname = "localhost"
		httpport = 8080

	}

	homePageUrl := "http://" + hostname + ":" + strconv.Itoa(httpport)
	statusPageUrl := "http://" + hostname + ":" + strconv.Itoa(httpport) + "/swagger/index.html"
	healthCheckUrl := "http://" + hostname + ":" + strconv.Itoa(httpport) + "/swagger/index.html"

	dataCenterInfo := DataCenterInfo{Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo", Name: "MyOwn"}

	basePort := Port{Port: httpport, Enabled: "true"}
	securePort := Port{Port: httpport, Enabled: "false"}

	instance := InstanceDetails{InstanceId: appId, HostName: hostname, App: appName, VipAddress: appName, SecureVipAddress: appName, IpAddr: hostname, Status: status, Port: basePort, SecurePort: securePort, HealthCheckUrl: healthCheckUrl, StatusPageUrl: statusPageUrl, HomePageUrl: homePageUrl, DataCenterInfo: dataCenterInfo}

	requestBody := &AppRegistrationBody{Instance: instance}

	return requestBody

}

// * =========== *

// Init initializes Eureka communication
func Init(appName string, appId string, port int) chrono.ScheduledTask {

	log.Println("initializing Eureka")

	// Initializes the Eureka server

	RegisterApp(appName, appId, port)
	UpdateAppStatus(appName, appId, port, "UP")
	return ScheduleHeartbeat(appName, appId)

}

// * =========== *

// Stop stops Eureka communication
func Stop(appName string, appId string, port int, task chrono.ScheduledTask) {

	log.Println("stopping Eureka server")

	// Stops communication with the Eureka server.
	task.Cancel()

	UpdateAppStatus(appName, appId, port, "DOWN")
	time.Sleep(5 * time.Second)
	DeleteApp(appName, appId)

}

// * =========== *

// StartClient Eureka's client starts
func StartClient(appName string, appId string, port int) {

	log.Println("starting Eureka client")

	// Initialize the Eureka client

	task := Init(appName, appId, port)

	// Create channel for receiving interrupt signals

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	// Wait for interruption signal
	go func() {
		select {
		case sgn := <-ch:
			_ = sgn
			Stop(viper.GetString("application.name"), appId, port, task) // Stop Eureka service
			os.Exit(1)
		}
	}()

}
