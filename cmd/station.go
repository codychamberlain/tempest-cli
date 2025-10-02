/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type Station struct {
	Stations []struct {
		Capabilities []struct {
			Capability      string  `json:"capability"`
			DeviceID        int     `json:"device_id"`
			Environment     string  `json:"environment"`
			Agl             float64 `json:"agl,omitempty"`
			ShowPrecipFinal bool    `json:"show_precip_final,omitempty"`
		} `json:"capabilities"`
		CreatedEpoch int `json:"created_epoch"`
		Devices      []struct {
			DeviceID   int `json:"device_id"`
			DeviceMeta struct {
				Agl             float64 `json:"agl"`
				Environment     string  `json:"environment"`
				Name            string  `json:"name"`
				WifiNetworkName string  `json:"wifi_network_name"`
			} `json:"device_meta"`
			DeviceType       string `json:"device_type"`
			FirmwareRevision string `json:"firmware_revision"`
			HardwareRevision string `json:"hardware_revision"`
			SerialNumber     string `json:"serial_number"`
			DeviceSettings   struct {
				ShowPrecipFinal bool `json:"show_precip_final"`
			} `json:"device_settings,omitempty"`
		} `json:"devices"`
		IsLocalMode       bool    `json:"is_local_mode"`
		LastModifiedEpoch int     `json:"last_modified_epoch"`
		Latitude          float64 `json:"latitude"`
		LocationID        int     `json:"location_id"`
		Longitude         float64 `json:"longitude"`
		Name              string  `json:"name"`
		PublicName        string  `json:"public_name"`
		State             int     `json:"state"`
		StationID         int     `json:"station_id"`
		StationItems      []struct {
			DeviceID       int    `json:"device_id"`
			Item           string `json:"item"`
			LocationID     int    `json:"location_id"`
			LocationItemID int    `json:"location_item_id"`
			Sort           int    `json:"sort"`
			StationID      int    `json:"station_id"`
			StationItemID  int    `json:"station_item_id"`
		} `json:"station_items"`
		StationMeta struct {
			Elevation   float64 `json:"elevation"`
			ShareWithWf bool    `json:"share_with_wf"`
			ShareWithWu bool    `json:"share_with_wu"`
		} `json:"station_meta"`
		Timezone              string `json:"timezone"`
		TimezoneOffsetMinutes int    `json:"timezone_offset_minutes"`
	} `json:"stations"`
	Status struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"status"`
}

// stationCmd represents the station command
var stationCmd = &cobra.Command{
	Use:   "station",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := getAPIToken()
		params := url.Values{}
		params.Add("token", apiToken)
		sid := cmd.Flag("station").Value.String()
		if sid == "" {
			baseURL := "https://swd.weatherflow.com/swd/rest/stations"
			params.Add("station_id", sid)
			req, err := http.NewRequest("GET", baseURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			req.URL.RawQuery = params.Encode()
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			fmt.Println("No station ID provided, listing all stations:")
			var stations Station
			err = json.Unmarshal(body, &stations)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}
			for _, s := range stations.Stations {
				fmt.Println("*---------------------------------------------------------------*")
				fmt.Printf("| Station data for ID %v\t\t\t\t\t|\n", s.StationID)
				fmt.Println("|---------------------------------------------------------------|")
				fmt.Printf("| Name\t\t\t| %s\t\t\t|\n", s.Name)
				fmt.Printf("| Latitude\t\t| %f\t\t\t\t|\n", s.Latitude)
				fmt.Printf("| Longitude\t\t| %f\t\t\t\t|\n", s.Longitude)
				fmt.Printf("| Elevation\t\t| %.2f\t\t\t\t|\n", s.StationMeta.Elevation)
				fmt.Printf("| Timezone\t\t| %s\t\t\t|\n", s.Timezone)
				fmt.Printf("| Public Name\t\t| %s\t\t|\n", s.PublicName)
				fmt.Println("|---------------------------------------------------------------|")
				fmt.Printf("| %v devices connected to station\t\t\t\t|\n", len(s.Devices))
				for _, device := range s.Devices {
					fmt.Printf("|  -> Device Name\t| %s\t\t\t\t|\n", device.DeviceMeta.Name)
					fmt.Printf("|  --- Device Type\t| %s\t\t\t\t\t|\n", device.DeviceType)
					fmt.Printf("|  --- Firmware Rev.\t| %s\t\t\t\t\t|\n", device.FirmwareRevision)
					fmt.Printf("|  --- Hardware Rev.\t| %s\t\t\t\t\t|\n", device.HardwareRevision)
					fmt.Printf("|  --- Serial Number\t| %s\t\t\t\t|\n", device.SerialNumber)
				}
				fmt.Println("*---------------------------------------------------------------*")
			}
		} else {
			baseURL := "https://swd.weatherflow.com/swd/rest/stations" + "/" + sid
			params.Add("station_id", sid)
			req, err := http.NewRequest("GET", baseURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			req.URL.RawQuery = params.Encode()
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			if cmd.Flag("output").Value.String() == "JSON" {
				fmt.Println(string(body))
			} else {
				var s Station
				err = json.Unmarshal(body, &s)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					return
				}

				fmt.Println("*---------------------------------------------------------------*")
				fmt.Printf("| Station data for ID %s\t\t\t\t\t|\n", sid)
				fmt.Println("|---------------------------------------------------------------|")
				fmt.Printf("| Name\t\t\t| %s\t\t\t|\n", s.Stations[0].Name)
				fmt.Printf("| Latitude\t\t| %f\t\t\t\t|\n", s.Stations[0].Latitude)
				fmt.Printf("| Longitude\t\t| %f\t\t\t\t|\n", s.Stations[0].Longitude)
				fmt.Printf("| Elevation\t\t| %.2f\t\t\t\t|\n", s.Stations[0].StationMeta.Elevation)
				fmt.Printf("| Timezone\t\t| %s\t\t\t|\n", s.Stations[0].Timezone)
				fmt.Printf("| Public Name\t\t| %s\t\t|\n", s.Stations[0].PublicName)
				fmt.Println("|---------------------------------------------------------------|")
				fmt.Printf("| %v devices connected to station\t\t\t\t|\n", len(s.Stations[0].Devices))
				for _, device := range s.Stations[0].Devices {
					fmt.Printf("|  -> Device Name\t| %s\t\t\t\t|\n", device.DeviceMeta.Name)
					fmt.Printf("|  --- Device Type\t| %s\t\t\t\t\t|\n", device.DeviceType)
					fmt.Printf("|  --- Firmware Rev.\t| %s\t\t\t\t\t|\n", device.FirmwareRevision)
					fmt.Printf("|  --- Hardware Rev.\t| %s\t\t\t\t\t|\n", device.HardwareRevision)
					fmt.Printf("|  --- Serial Number\t| %s\t\t\t\t|\n", device.SerialNumber)
				}
				fmt.Println("*---------------------------------------------------------------*")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(stationCmd)
}
