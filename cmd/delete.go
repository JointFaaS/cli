package cmd

import (
	"net/http"
	"io"
	"io/ioutil"
)

func deleteFunc(managerAddr string, funcName string) ([]byte, error) {
	url := "http://" + managerAddr + "/delete?funcName=" + funcName
    req, err := http.NewRequest("POST", url, nil)
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