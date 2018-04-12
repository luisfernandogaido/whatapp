package fs

import (
	"io/ioutil"
	"os"
	"strings"
	"archive/zip"
	"path/filepath"
	"io"
	"bufio"
)

func Zips(folder string) ([]string, error) {
	fi, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	zips := make([]string, 0)
	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		zips = append(zips, f.Name())
	}
	return zips, nil
}

func Unzip(file string) error {
	folder := file[:strings.LastIndex(file, ".")]
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, 0644)
		if err != nil {
			return err
		}
	}
	r, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		f, err := os.OpenFile(filepath.Join(folder, f.Name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
		rc.Close()
		f.Close()
	}
	return nil
}

func Descompacta(dataFolder string) error {
	zips, err := Zips(dataFolder)
	if err != nil {
		return err
	}
	for _, z := range zips {
		fullPath := filepath.Join(dataFolder, z)
		if err := Unzip(fullPath); err != nil {
			return err
		}
		if err = os.Remove(fullPath); err != nil {
			return err
		}
	}
	return nil
}

func MesclaTxt(dataFolder string) error {
	fis, err := ioutil.ReadDir(dataFolder)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		fis2, err := ioutil.ReadDir(filepath.Join(dataFolder, fi.Name()))
		if err != nil {
			return err
		}
		for _, fi2 := range fis2 {
			if !strings.HasSuffix(fi2.Name(), ".txt") || !strings.HasPrefix(fi2.Name(), "Conversa do WhatsApp") {
				continue
			}
			nomeWriter := filepath.Join(
				dataFolder,
				fi.Name(),
				strings.Replace(fi2.Name(), "Conversa do WhatsApp com ", "", 1),
			)
			f, err := os.OpenFile(nomeWriter, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			nomeReader := filepath.Join(
				dataFolder,
				fi.Name(),
				fi2.Name(),
			)
			f2, err := os.Open(nomeReader)
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(f2)
			for scanner.Scan() {
				linha := scanner.Text()
				_, err = f.WriteString(linha+ "\n")
				if err != nil {
					return err
				}
			}
			f.Close()
			f2.Close()
			err = os.Remove(nomeReader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
