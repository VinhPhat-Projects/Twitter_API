package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	api "twitterapi/twitterapi"

	"github.com/ChimeraCoder/anaconda"
	"github.com/rivo/uniseg"
)

const ERR_CORD = 1

func ReadTextFile(path string) (string, error) {
	bf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("text file is invalid. \n%v", err)
	}

	text := string(bf)
	gr := uniseg.GraphemeClusterCount(text)
	// fmt.Println("text length gr:", gr, text)
	if 240 < gr {
		return "", fmt.Errorf("text limit over. 240 < %v", gr)
	}

	return text, nil
}

func ReadImageFiles(path string, api *anaconda.TwitterApi) (url.Values, error) {
	if path == "" {
		return nil, nil
	}

	image, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Image file read error. \n%v", err)
	}

	base64str := base64.StdEncoding.EncodeToString(image)
	media, err := api.UploadMedia(base64str)
	if err != nil {
		return nil, fmt.Errorf("Image file Upload error. \n%v", err)
	}

	v := url.Values{}
	v.Add("media_ids", media.MediaIDString)

	return v, nil
}

func main() {
	text := flag.String("t", "", "Plane text")
	filePath := flag.String("f", "", "Text File Path")
	imagePath := flag.String("i", "", "Image File Path")
	env := flag.String("c", "", `Define Twitter API Keys File Path
Format:
{
	"API_KEYS"           : "your api key",
	"API_SECRET"         : "your api secret",
	"ACCESS_TOKEN"       : "your access token",
	"ACCESS_TOKEN_SECRET": "your token secret",
	"BEARER_TOKEN"       : "your bearer token"
}`)
	flag.Parse()

	tweet := ""
	if *filePath != "" {
		var e error
		tweet, e = ReadTextFile(*filePath)
		if e != nil {
			fmt.Println(e)
			os.Exit(ERR_CORD)
		}
	}

	if *text != "" {
		tweet = *text
	}

	if *env == "" || tweet == "" {
		fmt.Println("args invalid error. Please set args.")
		fmt.Println(`For example: 
	go run Tweet.go -c API_Key.txt -t Tweet-Description
	go run Tweet.go -c API_Key.txt -f tweet.txt

	with Image:
	go run Tweet.go -c API_Key.txt -t Tweet-Description -i img.png
	go run Tweet.go -c API_Key.txt -f tweet.txt -i img.jpg
	`)
		os.Exit(ERR_CORD)
	}

	api, err := api.NewAPI(*env)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERR_CORD)
	}

	twitter := api.GetTwitterAPI()
	media, err := ReadImageFiles(*imagePath, twitter)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERR_CORD)
	}

	tweeted, err := twitter.PostTweet(tweet, media)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERR_CORD)
	}

	fmt.Printf("tweeted: %v\n", tweeted.Text)
}
