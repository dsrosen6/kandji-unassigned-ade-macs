package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ADEDevices struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Result    `json:"results"`
}

type Result struct {
	BlueprintID                string     `json:"blueprint_id"`
	MdmDevice                  MdmDevice  `json:"mdm_device"`
	UserID                     string     `json:"user_id"`
	DepAccount                 DepAccount `json:"dep_account"`
	AssetTag                   string     `json:"asset_tag"`
	Color                      string     `json:"color"`
	Description                string     `json:"description"`
	DeviceAssignedBy           string     `json:"device_assigned_by"`
	DeviceAssignedDate         time.Time  `json:"device_assigned_date"`
	DeviceFamily               string     `json:"device_family"`
	Model                      string     `json:"model"`
	OS                         string     `json:"os"`
	ProfileAssignTime          time.Time  `json:"profile_assign_time"`
	ProfilePushTime            time.Time  `json:"profile_push_time"`
	ProfileStatus              string     `json:"profile_status"`
	SerialNumber               string     `json:"serial_number"`
	ID                         string     `json:"id"`
	LastAssignmentStatus       string     `json:"last_assignment_status"`
	FailedAssignmentAttempts   int64      `json:"failed_assignment_attempts"`
	AssignmentStatusReceivedAt time.Time  `json:"assignment_status_received_at"`
	Blueprint                  string     `json:"blueprint"`
	User                       string     `json:"user"`
}

type DepAccount struct {
	ID         string `json:"id"`
	ServerName string `json:"server_name"`
}

type MdmDevice struct {
	ID               string    `json:"id"`
	EnrolledAt       time.Time `json:"enrolled_at"`
	Name             string    `json:"name"`
	EnrollmentStatus int64     `json:"enrollment_status"`
	DeferredInstall  bool      `json:"deferred_install"`
	IsMissing        bool      `json:"is_missing"`
	IsRemoved        bool      `json:"is_removed"`
}

func main() {
	key, err := getApiKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	url, err := getApiURL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d, err := getAdeDeviceList(key, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ua, err := extractUnassignedDevices(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(ua.Results) < 1 {
		fmt.Println("No available devices found.")
		return
	}

	fmt.Printf("Found %d Unassigned Devices\n\n", len(ua.Results))
	for _, device := range ua.Results {
		fmt.Printf("%s\n%s\n\n", device.SerialNumber, device.Model)
	}
}

// creates filtered struct of Macs with no assigned user (aka, available devices)
func extractUnassignedDevices(adeResp *ADEDevices) (*ADEDevices, error) {
	var unassignedDevices ADEDevices
	for _, device := range adeResp.Results {
		if device.User == "" && strings.Contains(device.Model, "Mac") {
			unassignedDevices.Results = append(unassignedDevices.Results, device)
		}
	}

	return &unassignedDevices, nil
}

// gets unfiltered struct of ADE devices from Kandji
func getAdeDeviceList(key, url string) (*ADEDevices, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var adeResp ADEDevices
	if err := json.Unmarshal(respBody, &adeResp); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &adeResp, nil
}

// gets API key from environment for use with device list call
func getApiKey() (string, error) {
	key := os.Getenv("KANDJI_API_KEY")
	if key == "" {
		return "", fmt.Errorf("no valid Kandji API key obtained\nmake sure you have a key set as an environmental variable with the name 'KANDJI_API_KEY'")
	}

	return key, nil
}

// gets API subdomain from environment and creates full url
func getApiURL() (string, error) {
	s := os.Getenv("KANDJI_API_SUBDOMAIN")
	if s == "" {
		return "", fmt.Errorf("no valid Kandji API subdomain obtained\nmake sure you have a key set as an environmental variable with the name 'KANDJI_API_SUBDOMAIN'")
	}

	url := fmt.Sprintf("https://%s.api.kandji.io/api/v1/integrations/apple/ade/devices", s)
	return url, nil
}
