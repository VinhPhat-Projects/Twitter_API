package twitterapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ChimeraCoder/anaconda"
)

type api struct {
	API_KEYS            string
	API_SECRET          string
	ACCESS_TOKEN        string
	ACCESS_TOKEN_SECRET string
	BEARER_TOKEN        string
}

func NewAPI(path string) (*api, error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File can't read error:", path)
		return nil, err
	}

	api := new(api)
	err = json.Unmarshal(file, api)
	if err != nil {
		fmt.Println("Unmarshal error:", path)
		return nil, err
	}

	return api, nil
}

func (api *api) GetTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(api.API_KEYS)
	anaconda.SetConsumerSecret(api.API_SECRET)
	a := anaconda.NewTwitterApi(api.ACCESS_TOKEN, api.ACCESS_TOKEN_SECRET)
	return a
}
