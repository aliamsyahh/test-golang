package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Define structure for User
type User struct {
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Email    string   `json:"email"`
	Age      int      `json:"age"`
	Phone    string   `json:"phone"`
	Cell     string   `json:"cell"`
	Picture  []string `json:"picture"`
}

type RandomUserResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City     string `json:"city"`
			State    string `json:"state"`
			Country  string `json:"country"`
			Postcode int    `json:"postcode"`
		} `json:"location"`
		Email string `json:"email"`
		Dob   struct {
			Age int `json:"age"`
		} `json:"dob"`
		Phone   string `json:"phone"`
		Cell    string `json:"cell"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
	} `json:"results"`
}

func FetchUser(c *gin.Context) {
	resultsParam := c.DefaultQuery("results", "10")
	pageParam := c.DefaultQuery("page", "1")

	results, err := strconv.Atoi(resultsParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'results' parameter"})
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
		return
	}

	url := fmt.Sprintf("https://randomuser.me/api?results=%d&page=%d", results, page)
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from API"})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var apiResponse RandomUserResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode JSON response"})
		return
	}

	var users []User
	for _, result := range apiResponse.Results {
		users = append(users, User{
			Name:     fmt.Sprintf("%s, %s %s", result.Name.Title, result.Name.First, result.Name.Last),
			Location: fmt.Sprintf("%d %s, %s, %s, %s", result.Location.Street.Number, result.Location.Street.Name, result.Location.City, result.Location.State, result.Location.Country),
			Email:    result.Email,
			Age:      result.Dob.Age,
			Phone:    result.Phone,
			Cell:     result.Cell,
			Picture:  []string{result.Picture.Large, result.Picture.Medium, result.Picture.Thumbnail},
		})
	}

	// Return the response in the desired format
	c.JSON(http.StatusOK, users)
}
