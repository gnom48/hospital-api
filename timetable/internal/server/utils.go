package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

func IsUserInRole(roles []models.Role, roleId string) bool {
	for _, role := range roles {
		if role.Id == roleId {
			return true
		}
	}
	return false
}

func IsValidTime(from, to time.Time) bool {
	if from.After(to) {
		return false
	}

	if to.Sub(from) > 12*time.Hour {
		return false
	}

	if from.Second() != 0 || from.Minute()%30 != 0 || to.Second() != 0 || to.Minute()%30 != 0 {
		return false
	}

	return true
}

func CheckIfDoctorExists(tokenString string, doctorId string) (bool, error) {
	doctorNotFoundError := fmt.Errorf("Doctor not found")
	if tokenString == "" {
		return false, fmt.Errorf("Token is empty")
	}
	if doctorId == "" {
		return false, fmt.Errorf("Doctor Id is empty")
	}

	url := "http://account-service:8081/api/Doctors/" + doctorId
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("Request error: %w", err)
	}

	req.Header.Set("Authorization", tokenString)

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("Request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, doctorNotFoundError
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return false, doctorNotFoundError
	}
	if _, exists := responseMap["server_error"]; exists {
		return false, fmt.Errorf(responseMap["server_error"].(string))
	}

	return true, nil
}

func CheckIfHospitalExists(tokenString string, hospitalId string) (bool, error) {
	hospitalNotFoundError := fmt.Errorf("Hospital not found")
	if tokenString == "" {
		return false, fmt.Errorf("Token is empty")
	}
	if hospitalId == "" {
		return false, fmt.Errorf("Doctor Id is empty")
	}

	url := "http://hospital-service:8082/api/Hospitals/" + hospitalId
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("Request error: %w", err)
	}

	req.Header.Set("Authorization", tokenString)

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("Request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, hospitalNotFoundError
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return false, hospitalNotFoundError
	}
	if _, exists := responseMap["server_error"]; exists {
		return false, fmt.Errorf(responseMap["server_error"].(string))
	}

	return true, nil
}
