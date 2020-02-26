package cmd

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
    "mime/multipart"
    "path"
	"net/http"
	"os"
)

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func compressDir(dir string, w *zip.Writer) (error) {
	fileInfos, _ := ioutil.ReadDir(dir)
	for _, fi := range fileInfos {
		f, err := os.Open(path.Join(dir, fi.Name()))
		if err != nil {
			return err
		}
		err = compress(f, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

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

    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }
    if fileInfo.IsDir() {
        zWriter := zip.NewWriter(formFile)
        compressDir(sourceZip, zWriter)
        zWriter.Close()
    } else {
        _, err = io.Copy(formFile, file)
        if err != nil {
            return nil, err
        }
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