package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var userkey string
var AircraftType string
var Terminal string

type AircraftItems struct {
	XMLName  xml.Name   `xml:"AircraftItems"`
	Aircraft []Aircraft `xml:"Aircraft"`
}

type Aircraft struct {
	XMLName      xml.Name `xml:"Aircraft"`
	SerialNumber string   `xml:"SerialNumber"`
	Location     string   `xml:"Location"`
	NeedsRepair  string   `xml:"NeedsRepair"`
}

type IcaoJobsFrom struct {
	XMLName    xml.Name     `xml:"IcaoJobsFrom"`
	Assignment []Assignment `xml:"Assignment"`
}

type Assignment struct {
	XMLName        xml.Name `xml:"Assignment"`
	Location       string   `xml:"Location"`
	ToIcao         string   `xml:"ToIcao"`
	Type           string   `xml:"Type"`
	AircraftId     string   `xml:"AircraftId"`
	Commodity      string   `xml:"Commodity"`
	Pay            float64  `xml:"Pay"`
	Expires        string   `xml:"Expires"`
	ExpireDateTime string   `xml:"ExpireDateTime"`
}

type Job struct {
	Location string
	ToIcao   string
	Expires  string
	Pay      int
}

type GPS_Coordinates struct {
	lat1 float64
	lon1 float64
}

type Airport struct {
	Ident        string `json:"ident"`
	Name         string `json:"name"`
	Iso_country  string `json:"iso_country"`
	Elevation_ft string `json:"elevation_ft"`
	Iata_code    string `json:"iata_code"`
	Iso_region   string `json:"iso_region"`
	Type         string `json:"type"`
	Coordinates  string `json:"coordinates"`
}

type Config struct {
	Userkey      string `json:"userkey"`
	Aircrafttype string `json:"aircrafttype"`
	Terminal     string `json:"terminal"`
}

var data []byte

func main() {
	readConfig()
	data = readAirportData()
	printTopJobs()
}

func readConfig() {
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Print(err)
	}
	var Config Config
	json.Unmarshal(config, &Config)
	userkey = Config.Userkey
	Terminal = Config.Terminal
	AircraftType = strings.Replace(Config.Aircrafttype, " ", "%20", -1)
}

func convertGPSstring(Coordinates string) GPS_Coordinates {
	var GPS_Coordinates GPS_Coordinates
	coords := strings.Split(Coordinates, ",")
	if lat1, err := strconv.ParseFloat(strings.Trim(coords[0], " "), 64); err == nil {
		GPS_Coordinates.lat1 = lat1
	}
	if lon1, err := strconv.ParseFloat(strings.Trim(coords[1], " "), 64); err == nil {
		GPS_Coordinates.lon1 = lon1
	}
	return GPS_Coordinates
}

func convertKMtoNM(KM int) int {
	return int(float64(KM) * 0.5399565)
}
func degreesToRadians(degress float64) float64 {
	return degress * math.Pi / 180
}

func distanceInKmBetweenEarthCoordinates(lat1 float64, lon1 float64, lat2 float64, lon2 float64) int {
	var earthRadiusKm = 6371.0
	var dLat = degreesToRadians(lat2 - lat1)
	var dLon = degreesToRadians(lon2 - lon1)
	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return int(earthRadiusKm * c)
}

func readAirportData() []byte {
	data, err := ioutil.ReadFile("./airport-codes_json.json")
	if err != nil {
		fmt.Print(err)
	}
	return data
}

func getAirportData(ICAO string) Airport {
	var AirportInfo Airport
	var obj []Airport
	json.Unmarshal(data, &obj)
	for i := 0; i < len(obj); i++ {
		if ICAO == obj[i].Ident {
			AirportInfo = obj[i]
		}
	}
	return AirportInfo
}

func getAircrafts() AircraftItems {

	Aircraftresponse, err := http.Get("https://server.fseconomy.net/data?userkey=" + userkey + "&format=xml&query=aircraft&search=makemodel&makemodel=" + AircraftType)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(Aircraftresponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	var AircraftItem AircraftItems
	xml.Unmarshal(responseData, &AircraftItem)
	if Terminal == "bash" {
		fmt.Println("Found " + "\033[31m" + strconv.Itoa(len(AircraftItem.Aircraft)) + "\033[0m Aircrafts of type " + strings.Replace(AircraftType, "%20", "", -1))
	} else {
		fmt.Println("Found " + strconv.Itoa(len(AircraftItem.Aircraft)) + " Aircrafts of type " + strings.Replace(AircraftType, "%20", "", -1))
	}
	return AircraftItem
}

func getAssignment(AircraftItem AircraftItems) IcaoJobsFrom {
	var Locationstr string
	var IcaoJobsFromList IcaoJobsFrom
	for i := 0; i < len(AircraftItem.Aircraft); i++ {
		// concat all airports
		if AircraftItem.Aircraft[i].Location != "In Flight" {
			if i == 0 {
				Locationstr = AircraftItem.Aircraft[i].Location
			} else {
				Locationstr = Locationstr + "-" + AircraftItem.Aircraft[i].Location
			}
		}
	}
	// get possible assignments on all airports, that have our aircraft type
	Assignmentresponse, err := http.Get("https://server.fseconomy.net/data?userkey=" + userkey + "&format=xml&query=icao&search=jobsfrom&icaos=" + Locationstr)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(Assignmentresponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(responseData, &IcaoJobsFromList)
	return IcaoJobsFromList
}

func printTopJobs() {
	var Aircrafts AircraftItems = getAircrafts()
	var IcaoJobsFromList IcaoJobsFrom = getAssignment(Aircrafts)
	var myList []Job
	// var Joblists [] Joblist
	for i := 0; i < len(IcaoJobsFromList.Assignment); i++ {
		if IcaoJobsFromList.Assignment[i].AircraftId != "0" {
			for j := 0; j < len(Aircrafts.Aircraft); j++ {
				// Only add Planes, that we searched
				if Aircrafts.Aircraft[j].SerialNumber == IcaoJobsFromList.Assignment[i].AircraftId {
					myList = append(myList, Job{
						Location: IcaoJobsFromList.Assignment[i].Location,
						Pay:      int(IcaoJobsFromList.Assignment[i].Pay),
						ToIcao:   IcaoJobsFromList.Assignment[i].ToIcao,
						Expires:  IcaoJobsFromList.Assignment[i].Expires,
					})
				}
			}
		}
	}
	sort.SliceStable(myList, func(i, j int) bool {
		return myList[i].Pay > myList[j].Pay
	})
	for i := 0; i < len(myList) && i < 10; i++ {
		var Airport Airport = getAirportData(myList[i].Location)
		var TravelDistance int = calculateDistanceNM(myList[i].Location, myList[i].ToIcao)
		var AirportInfoStr string = Airport.Ident + ", " + Airport.Name + ", " + Airport.Iso_country + ", " + Airport.Type
		if Terminal == "bash" {
			fmt.Println("\033[32m" + myList[i].Location + " > " + myList[i].ToIcao + "\033[0m" + ", " + "\033[33m" + strconv.Itoa(TravelDistance) + " NM" + "\033[0m" + ", Expires in: " + myList[i].Expires + " , " + "\033[35m" + strconv.Itoa(myList[i].Pay) + "$" + "\033[0m" + " : " + AirportInfoStr)
		} else {
			fmt.Println(myList[i].Location + " > " + myList[i].ToIcao + ", " + strconv.Itoa(TravelDistance) + " NM , Expires in: " + myList[i].Expires + " , " + strconv.Itoa(myList[i].Pay) + "$" + " : " + AirportInfoStr)
		}

	}
}

func calculateDistanceNM(FromICAO string, ToICAO string) int {
	var FROM GPS_Coordinates
	var TO GPS_Coordinates
	FROM = convertGPSstring(getAirportData(FromICAO).Coordinates)
	TO = convertGPSstring(getAirportData(ToICAO).Coordinates)
	var distance int = distanceInKmBetweenEarthCoordinates(FROM.lat1, FROM.lon1, TO.lat1, TO.lon1)
	return convertKMtoNM(distance)
}