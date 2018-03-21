package Utils

import (
	"net/http"
	"bytes"
	"io"
	"log"
	"io/ioutil"
	"github.com/Daniele122898/Project-Aegis-Website/src/config"
)

func AdminPostRequest(url string, data []byte) ([]byte, error){
	req, err := http.NewRequest("POST",url , bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Get().AdminToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		log.Println("Failed to get response from post")
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func AdminPostRequestLong(url string, data []byte, expectedCode int) ([]byte,bool, error){
	req, err := http.NewRequest("POST",url , bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Get().AdminToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		log.Println("Failed to get response from post")
		return nil, false,err
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedCode{
		return nil, false, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, true, nil
}

func AdminGetRequestLong(url string, expectedCode int) ([]byte, bool ,error){
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.Get().AdminToken)
	if err != nil{
		return nil, false, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil{
		return nil, false, err
	}
	//close the body at the end
	defer resp.Body.Close()

	if resp.StatusCode != expectedCode{
		return nil, false, nil
	}

	buf := bytes.NewBuffer(nil)

	_, err = io.Copy(buf, resp.Body)
	if err != nil{
		return nil, false, err
	}

	return buf.Bytes(), true, nil
}


func AdminGetRequest(url string) ([]byte, error){
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.Get().AdminToken)
	if err != nil{
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	//close the body at the end
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)

	_, err = io.Copy(buf, resp.Body)
	if err != nil{
		return nil, err
	}

	return buf.Bytes(), nil
}
