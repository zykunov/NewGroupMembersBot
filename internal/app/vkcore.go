package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zykunov/courseGoFirst/vkApiBot/models"
)

type Response struct {
	Response Data `json:"response"`
}
type Data struct {
	Count    int            `json:"count"`
	Items    []GroupMembers `json:"items"`
	NextFrom string         `json:"next_from"`
}

type GroupMembers struct {
	Id        int    `json:"id"`
	Sex       int    `json:"sex"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Access    bool   `json:"can_access_closed"`
	IsClosed  bool   `json:"is_closed"`
}

func getVkToken() (token string) {
	vkToken := os.Getenv("vk_token")
	if vkToken == "" {
		log.Fatalf("set token for vk!")
	}
	return vkToken
}

func GetGroupMembers(groupId string, userId int) []models.User {

	token := getVkToken()

	var members Response
	var result []models.User
	var offset int
	var i int = 1

	for {

		url := fmt.Sprintf("https://api.vk.com/method/groups.getMembers?group_id=%s&offset=%d&fields=sex&access_token=%s&v=5.131", groupId, offset, token)
		res, err := http.Get(url)

		if err != nil {
			log.Printf("error with vk api")
		}
		defer res.Body.Close()

		byteValue, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(byteValue, &members)

		for _, m := range members.Response.Items {
			result = append(result, models.User{VkgroupId: groupId, VkId: m.Id, UserID: userId, Sex: m.Sex, IsNew: true})
		}

		offset = offset + 1000

		time.Sleep(210 * time.Millisecond) // не больше 5 запросов в секунду
		log.Println("number of iteration", i)
		i++

		if members.Response.NextFrom == "" {
			break
		}

	}

	return result
}
