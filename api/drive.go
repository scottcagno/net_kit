// *
// * Copyright 2014, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Google Drive API Calls
// *

package api

import (
	// required package by google
	"code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	"io/ioutil"
	"net/http"
)

// mime file type
const (
	PDF = "application/pdf"
	XLS = "application/vnd.ms-excel"
	CSV = "text/csv"
)

// google drive structure
type Drive struct {
	*drive.Service
}

// initialize and return a drive instance
func InitDrive() *Drive {
	return &Drive{}
}

// upload file
func (self *Drive) Upload(fldr, file, mime string) (*drive.File, bool) {
	m, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	defer m.Close()
	f := &drive.File{
		Title:       name,
		Description: fldr + " uploaded file",
		MimeType:    mime,
	}
	if fldr != "" {
		p := &drive.ParentReference{Id: fldr}
		f.Parents = []*drive.ParentReference{p}
	}
	r, err := self.Files.Insert(f).Media(m).Do()
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return r, true
}

// download file
func DownloadFile(fldr, file, mime string) (*drive.File, bool) {
	m, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	defer m.Close()
	f := &drive.File{
		Title:       name,
		Description: fldr + " uploaded file",
		MimeType:    mime,
	}
	if fldr != "" {
		p := &drive.ParentReference{Id: fldr}
		f.Parents = []*drive.ParentReference{p}
	}
	// t parameter should use an oauth.Transport
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		fmt.Printf("An error occurred: File is not downloadable")
		return "", nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	resp, err := t.RoundTrip(req)
	// Make sure we close the Body later
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return "", err
	}
	return string(body), nil
}
