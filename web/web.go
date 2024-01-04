package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func HTTPGetBody(url string) (string, error) {

	// Send an HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Convert the response body to a string
	return string(body), nil

}

func HTTPGetCode(url string) (int, error) {

	response, err := http.Head(url)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	status := response.Status
	statusCode := response.StatusCode

	fmt.Printf("URL: %s\n", url)
	fmt.Printf("HTTP Status: %s\n", status)
	fmt.Printf("Status Code: %d\n", statusCode)
	return statusCode, nil
}

func HTTPostJson(url string, jsond string) (string, error) {

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsond)))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	return string(body), nil
}
