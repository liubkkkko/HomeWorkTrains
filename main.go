package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

const (
	priceCriteria         = "price"
	arrivalTimeCriteria   = "arrival-time"
	departureTimeCriteria = "departure-time"
	lenResoult            = 3
)

var (
	unsupportedCriteria      = errors.New("unsupported criteria")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badArrivalStationInput   = errors.New("bad arrival station input")
	badDepartureStationInput = errors.New("bad departure station input")
)

func main() {
	var departureStation, arrivalStation, criteria string

	fmt.Print("Enter departure station : ")
	fmt.Scanf("%s", &departureStation)
	fmt.Print("Enter arrival station : ")
	fmt.Scanf("%s", &arrivalStation)
	fmt.Print("Enter criteria : ")
	fmt.Scanf("%s", &criteria) 

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
		return
	} 

	for _, res := range result {
		fmt.Printf("% +v\n", res)
	} 

}

//Пошук поїздів за заданими критеріями
func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	jsonFile, err := os.Open("data.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonByte, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var trains Trains

	err = json.Unmarshal(jsonByte, &trains)
	if err != nil {
		return nil, err
	}

	departureStationValid, arrivalStationValid, criteriaValid, err := parameterValidity(departureStation, arrivalStation, criteria)
	if err != nil {
		return nil, err
	}

	selectionTrains, err := trains.selection(departureStationValid, arrivalStationValid)
	if err != nil {
		return nil, err
	}
	selectionTrains.sorting(criteriaValid)

	if len(selectionTrains) > lenResoult {
		selectionTrains = selectionTrains[:lenResoult]
	}
	return selectionTrains, nil
}

// Відбір потрібних нам поїздів
func (t Trains) selection(departureStation int, arrivalStation int) (Trains, error) {
	var selectionTrains Trains

	for _, value := range t {
		if value.DepartureStationID == departureStation {
			if value.ArrivalStationID == arrivalStation {
				selectionTrains = append(selectionTrains, value)
			}
		}
	}

	if len(selectionTrains) == 0 {
		return nil, nil
	}

	return selectionTrains, nil
}

//Функція сортування за введеною критерією
func (t Trains) sorting(criteria string) {

	if criteria == priceCriteria {
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].Price < t[j].Price
		})
	}
	if criteria == arrivalTimeCriteria {
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		})
	}
	if criteria == departureTimeCriteria {
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		})
	}

}

// Кастомний анмаршал для формату time.Time
func (t *Train) UnmarshalJSON(data []byte) error {
	type TrainClone Train
	layout := "15:04:05"

	str := &struct {
		ArrivalTime   string
		DepartureTime string
		*TrainClone
	}{
		TrainClone: (*TrainClone)(t),
	}

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	t.ArrivalTime, err = time.Parse(layout, str.ArrivalTime)
	if err != nil {
		return err
	}

	t.DepartureTime, err = time.Parse(layout, str.DepartureTime)
	if err != nil {
		return err
	}

	return nil
}

// Валідація вхідних даних
func parameterValidity(departureStation, arrivalStation, criteria string) (departureStationValid int, arrivalStationValid int, criteriaValid string, err error) {
	const firstNaturalNumber = 1
	if departureStation == "" {
		return departureStationValid, arrivalStationValid, criteriaValid, emptyDepartureStation
	}
	if arrivalStation == "" {
		return departureStationValid, arrivalStationValid, criteriaValid, emptyArrivalStation
	}
	intDepartureStation, err := strconv.Atoi(departureStation)
	if intDepartureStation < firstNaturalNumber || err != nil {
		return departureStationValid, arrivalStationValid, criteriaValid, badDepartureStationInput
	}
	intArrivalStation, err := strconv.Atoi(arrivalStation)
	if intArrivalStation < firstNaturalNumber || err != nil {
		return departureStationValid, arrivalStationValid, criteriaValid, badArrivalStationInput
	}
	if criteria != priceCriteria && criteria != arrivalTimeCriteria && criteria != departureTimeCriteria {
		return departureStationValid, arrivalStationValid, criteriaValid, unsupportedCriteria
	}
	departureStationValid, _ = strconv.Atoi(departureStation)
	arrivalStationValid, _ = strconv.Atoi(arrivalStation)
	criteriaValid = criteria
	return departureStationValid, arrivalStationValid, criteriaValid, nil
}
