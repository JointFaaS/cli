package cmd

import (
	"archive/zip"
	"io"
	"io/ioutil"
    "path"
	"os"
)

func compressDir(dir string, dst io.Writer) (error) {
	w := zip.NewWriter(dst)
	defer w.Close()
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