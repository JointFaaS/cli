package cmd

import (
	"net/http"
	"io"
	"io/ioutil"
)

func invoke(managerAddr string, funcName string, args io.Reader) ([]byte, error) {
	url := "http://" + managerAddr + "/invoke?funcName=" + funcName
    req, err := http.NewRequest("POST", url, args)
    if err != nil {
        return nil, err
	}
	
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return content, nil
}