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

// InsertFile creates a new file in Drive from the given file and details
func InsertFile(d *drive.Service, title string, description string,
	parentId string, mimeType string, filename string) (*drive.File, error) {
	m, err := os.Open(filename)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	f := &drive.File{Title: title, Description: description, MimeType: mimeType}
	if parentId != "" {
		p := &drive.ParentReference{Id: parentId}
		f.Parents = []*drive.ParentReference{p}
	}
	r, err := d.Files.Insert(f).Media(m).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return r, nil
}

// PrintFile fetches and displays the given file.
func PrintFile(d *drive.Service, fileId string) error {
	f, err := d.Files.Get(fileId).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	fmt.Printf("Title: %v", f.Title)
	fmt.Printf("Description: %v", f.Description)
	fmt.Printf("MIME type: %v", f.MimeType)
	return nil
}

// DownloadFile downloads the content of a given file object
func DownloadFile(d *drive.Service, t http.RoundTripper, f *drive.File) (string, error) {
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

// DeleteFile deletes a file, skipping the trash
func DeleteFile(d *drive.Service, fileId string) error {
	err := d.Files.Delete(fileId).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	return nil
}

// IsFileInFolder tests whether the given file is in the given folder
func IsFileInFolder(d *drive.Service, folderId string,
	fileId string) (bool, error) {
	_, err := d.Children.Get(folderId, fileId).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

// InsertFileIntoFolder adds a given file to a given folder
func InsertFileIntoFolder(d *drive.Service, folderId string, fileId string) error {
	c := &drive.ChildReference{Id: fileId}
	_, err := d.Children.Insert(folderId, c).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	return nil
}

// RemoveFileFromFolder removes the given file from the given folder
func RemoveFileFromFolder(d *drive.Service, folderId string, fileId string) error {
	err := d.Children.Delete(folderId, fileId).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	return nil
}
