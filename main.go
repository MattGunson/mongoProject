package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type statePopulation struct {
	State      string `json:"state"`
	Population int    `json:"population"`
}

type stateVaccine struct {
	State        string `json:"state"`
	Date         string `json:"time"`
	Vaccinations int    `json:"vaccinated"`
}

type stateCases struct {
	Date   string `json:"time"`
	State  string `json:"state"`
	Cases  int    `json:"cases"`
	Deaths int    `json:"deaths"`
}

// Run this with `go run .` as long as
func main() {
	// SUMLEV,REGION,DIVISION,STATE,NAME,POPESTIMATE2019,POPEST18PLUS2019,PCNT_POPEST18PLUS
	popMapString, _ := readCsvFile("data/state_pops.csv", "STATE", addStatePopRecord)
	populationArr := make([]statePopulation, 0, len(popMapString))
	for k, v := range popMapString {
		population, _ := strconv.Atoi(v[5])
		populationArr = append(populationArr, statePopulation{k, population})
	}
	printPopulationWithJson(populationArr)

	// Entity,Code,Day,people_vaccinated
	stateVaccineMap, _ := readCsvFile("data/us-covid-19-total-people-vaccinated.csv", "Entity", addStateVaccineRecord)
	vaccineArr := make([]stateVaccine, 0, len(stateVaccineMap))
	for k, v := range stateVaccineMap {
		vaccinations, _ := strconv.Atoi(v[3])
		vaccineArr = append(vaccineArr, stateVaccine{k, v[2], vaccinations})
	}
	printVaccineWithJson(vaccineArr)

	// date,county,state,fips,cases,deaths
	casesRecords, _ := readCsvFile("data/us-counties.csv", "county", addCaseRecord)
	caseMap := make(map[string]stateCases)
	for _, rec := range casesRecords {
		cases, _ := strconv.Atoi(rec[4])
		deaths, _ := strconv.Atoi(rec[5])
		if v, found := caseMap[rec[2]]; !found {
			caseMap[rec[2]] = stateCases{rec[0], rec[2], cases, deaths}
		} else {
			caseStruct := v
			caseStruct.Cases += cases
			caseStruct.Deaths += deaths
			caseMap[rec[2]] = caseStruct
		}
	}
	caseArr := make([]stateCases, 0, len(caseMap))
	for _, v := range caseMap {
		caseArr = append(caseArr, v)
	}
	printCasesWithJson(caseArr)
}

func readCsvFile(filePath string, key string, addRecord func(int, []string, map[string][]string)) (map[string][]string, []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	records := make(map[string][]string)
	csvReader := csv.NewReader(f)
	keys, err := csvReader.Read()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	var keyIndex int
	for i, v := range keys {
		if v == key {
			keyIndex = i
		}
	}
	loop := true
	for loop {
		vals, err := csvReader.Read()
		if vals == nil {
			loop = false
		} else if err != nil {
			log.Fatal("Unale to parse file as CSV for "+filePath, err)
		} else {
			addRecord(keyIndex, vals, records)
		}
	}
	return records, keys
}

func addStatePopRecord(mapKey int, vals []string, records map[string][]string) {
	state := vals[4]
	if _, found := records[state]; !found {
		records[state] = vals
	} else {
		log.Println("duplicate record for " + state)
		log.Println(vals)
	}
}

func addStateVaccineRecord(mapKey int, vals []string, records map[string][]string) {
	state := vals[0]
	if v, found := records[state]; !found {
		records[state] = vals
	} else if v[2] < vals[2] {
		records[state] = vals
	}
}

func addCaseRecord(mapKey int, vals []string, records map[string][]string) {
	county := vals[mapKey]
	if v, found := records[county]; !found {
		records[county] = vals
	} else if v[0] < vals[0] {
		records[county] = vals
	}
}

func printVaccineWithJson(arr []stateVaccine) {
	dat, err := json.MarshalIndent(arr, "", "\t")
	if err != nil {
		return
	}
	err = os.WriteFile("json/vaccineArrQuotes.json", dat, 0644)
	if err != nil {
		return
	}
}

func printPopulationWithJson(arr []statePopulation) {
	dat, err := json.MarshalIndent(arr, "", "\t")
	if err != nil {
		return
	}
	err = os.WriteFile("json/populationArrQuotes.json", dat, 0644)
	if err != nil {
		return
	}
}

func printCasesWithJson(arr []stateCases) {
	dat, err := json.MarshalIndent(arr, "", "\t")
	if err != nil {
		return
	}
	err = os.WriteFile("json/caseArrQuotes.json", dat, 0644)
	if err != nil {
		return
	}
}
