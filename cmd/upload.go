package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func upload(managerAddr string, funcName string, sourceZip string, env string) ([]byte, error) {
	body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)

    formFile, err := writer.CreateFormFile("sourceZip", "code.zip")
    if err != nil {
        return nil, err
    }
	file, err := os.Open(sourceZip)
	if err != nil {
		return nil, err
	}
    _, err = io.Copy(formFile, file)
    if err != nil {
        return nil, err
    }

    _ = writer.WriteField("funcName", funcName)
	_ = writer.WriteField("env", env)

    err = writer.Close()
    if err != nil {
        return nil, err
	}
	
	url := "http://" + managerAddr + "/createfunction"
    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return nil, err
	}
	
    //req.Header.Set("Content-Type","multipart/form-data")
    req.Header.Add("Content-Type", writer.FormDataContentType())
	
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