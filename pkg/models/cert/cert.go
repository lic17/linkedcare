/*

 Copyright 2019 The Linkedcare Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var files []string

type CertTime struct {
	Name string
	Time string
}

func init() {
	dir := "/etc/kubernetes/pki"
	files, _ = walkDir(dir, "crt")

}
func walkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return files, err
}

func getCertInfo(file string) (string, error) {
	expiretime := ""
	certCerFile, err := os.Open(file)
	if err != nil {
		return expiretime, err
	}

	chunks, err := ioutil.ReadAll(certCerFile)
	certCerFile.Close()
	if len(chunks) <= 10 {
		return "-1", errors.New("the file is nil")
	}
	block, _ := pem.Decode(chunks)
	if block == nil {
		return expiretime, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return expiretime, err
	}

	now := time.Now().Unix()

	noTime := cert.NotAfter.Unix() - now
	if noTime <= 0 {
		return "0", nil
	}

	noTime /= 24 * 3600

	expiretime = strconv.FormatInt(noTime, 10)
	return expiretime, nil

}

func GetCertTime() ([]*CertTime, error) {

	var certsTime []*CertTime
	var fileErr error
	fileErr = nil
	for _, f := range files {
		tim, err := getCertInfo(f)
		if err != nil {
			fileErr = errors.New(f + err.Error())
			break
		}
		certsTime = append(certsTime, &CertTime{
			Name: f,
			Time: tim,
		})
	}
	return certsTime, fileErr
}
