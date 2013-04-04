// --------
// fedex.go ::: fedex api
// --------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package fedex

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	// Production Account Info
	Acct = Account{
		ApiURI:      "https://ws.fedex.com:443/web-services",
		DevKey:      "BltV0ugnEcMK7KdK",
		Password:    "0yi3ftgSTxvB4TeAxfFcHJg0w",
		AcctNumber:  "119710618",
		MeterNumber: "104991078",
	}

	// Testing Account Info
	TestAcct = Account{
		ApiURI:      "https://wsbeta.fedex.com:443/web-services",
		DevKey:      "isurR8vQXGWe7NdB",
		Password:    "sLpViYx5Zem7T01fLIHRZyKAQ",
		AcctNumber:  "510087682",
		MeterNumber: "118573690",
	}

	Zoom = Contact{
		PersonName:  "",
		CompanyName: "Zoom Envlopes",
		Phone:       "0000000000",
		Address: Address{
			Street: "52 Industrial Road",
			City:   "Ephrata",
			State:  "PA",
			Zip:    "17522",
		},
	}
)

type Account struct {
	ApiURI      string
	DevKey      string
	Password    string
	AcctNumber  string
	MeterNumber string
}

type Package struct {
	Weight         float64
	Tag            string
	SequenceNumber int
	Shipment       *Shipment
}

type Shipment struct {
	Account         Account
	Shipper         Contact
	Recipient       Contact
	Packages        []*Package
	PackageCount    int
	TrackingId      string
	RequestTemplate *template.Template
	TagTemplate     *template.Template
	Timestamp       time.Time
}

type Contact struct {
	PersonName, CompanyName, Phone string
	Address                        Address
}

type Address struct {
	Street, City, State, Zip string
}

func NewShipment(shipper, recipient Contact) *Shipment {
	return &Shipment{
		Account:         Acct,
		Shipper:         shipper,
		Recipient:       recipient,
		RequestTemplate: template.Must(template.New("xml").Parse(REQUEST_SHIPMENT_XML)),
		TagTemplate:     template.Must(template.New("tag").Parse(`<img width="380" hspace="25" vspace="50" src="data:image/png;base64,{{ .tag }}">`)),
		Timestamp:       time.Now(),
	}
}

func (self *Shipment) ParsePackages(orderQty, maxCartonCnt, stockWeight int) {
	pkgs := make([]*Package, 0)
	env := (float64(stockWeight) / float64(1000))
	var i, j int
	for i = orderQty; i >= maxCartonCnt; i -= maxCartonCnt {
		pkgs = append(pkgs, &Package{(float64(maxCartonCnt) * env), "", j + 1, self})
		j++
	}
	if i != 0 {
		pkgs = append(pkgs, &Package{(float64(i) * env), "", j + 1, self})
	}
	self.PackageCount = len(pkgs)
	self.Packages = pkgs
}

func (self *Package) MakeRequest(autoParse bool) []byte {
	var xml bytes.Buffer
	self.Shipment.RequestTemplate.Execute(&xml, map[string]interface{}{"fedex": self.Shipment, "package": self})
	response, err := http.Post(self.Shipment.Account.ApiURI, "application/xml", &xml)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	xmlDat, _ := ioutil.ReadAll(response.Body)
	if autoParse {
		pngDat := ParseXmlVals(xmlDat, "Image")["Image"]
		self.Tag = self.ParseImage(pngDat)
	}
	return xmlDat
}

func (self *Package) ParseImage(pngDat string) string {
	var tag bytes.Buffer
	self.Shipment.TagTemplate.Execute(&tag, map[string]interface{}{"tag": pngDat})
	return tag.String()
}

func ParseXmlVals(xmldat []byte, s ...string) map[string]string {
	vals := make(map[string]string)
	for _, val := range s {
		dec := xml.NewDecoder(bytes.NewBufferString(string(xmldat)))
		for {
			t, err := dec.Token()
			if err != nil {
				break
			}
			switch e := t.(type) {
			case xml.StartElement:
				if e.Name.Local == val {
					b, _ := dec.Token()
					switch b.(type) {
					case xml.CharData:
						vals[val] = fmt.Sprintf("%s", b)
					}
				}
			}
		}
	}
	return vals
}

/*
func (self *Shipment) MakeShipmentRequests() bool {
	if len(self.Packages) <= 0 {
		return false
	}
	self.PackageCount = len(self.Packages)
	var xml bytes.Buffer
	self.XmlTemplate.Execute(&xml, self)
	resp, _ := http.Post(self.Account.ApiURI, "application/xml", &xml)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", body)
	return true
}
*/

var TAG *template.Template
var REQUEST_SHIPMENT *template.Template
var REQUEST_SHIPMENT_XML = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v12="http://fedex.com/ws/ship/v12">
   <soapenv:Header/>
   <soapenv:Body>

		<v12:ProcessShipmentRequest>
			
			<v12:WebAuthenticationDetail>
				<v12:UserCredential>
		      		<v12:Key>{{ .fedex.Account.DevKey }}</v12:Key>
		      		<v12:Password>{{ .fedex.Account.Password }}</v12:Password>
		    	</v12:UserCredential>
		  	</v12:WebAuthenticationDetail>
		  	
		  	<v12:ClientDetail>
		    	<v12:AccountNumber>{{ .fedex.Account.AcctNumber }}</v12:AccountNumber>
		    	<v12:MeterNumber>{{ .fedex.Account.MeterNumber }}</v12:MeterNumber>
		  	</v12:ClientDetail>
		  	
		  	<!--
		  	<v12:TransactionDetail>
		    	<v12:CustomerTransactionId>** TEST TRANSACTION**</v12:CustomerTransactionId>
		  	</v12:TransactionDetail>
		  	-->

		  	<v12:Version>
		    	<v12:ServiceId>ship</v12:ServiceId>
		    	<v12:Major>12</v12:Major>
		    	<v12:Intermediate>0</v12:Intermediate>
		    	<v12:Minor>0</v12:Minor>
		  	</v12:Version>
		  
		  	<v12:RequestedShipment>
		  		
		  		<v12:ShipTimestamp>{{ .fedex.Timestamp }}</v12:ShipTimestamp>
		    	<v12:DropoffType>REGULAR_PICKUP</v12:DropoffType>
		    	<v12:ServiceType>PRIORITY_OVERNIGHT</v12:ServiceType>
		    	<v12:PackagingType>YOUR_PACKAGING</v12:PackagingType>
		    	
		    	<v12:Shipper>
		      		<v12:Contact>
		        		<v12:PersonName>{{ .fedex.Shipper.PersonName }}</v12:PersonName>
		        		<v12:CompanyName>{{ .fedex.Shipper.CompanyName }}</v12:CompanyName>
		        		<v12:PhoneNumber>{{ .fedex.Shipper.Phone }}</v12:PhoneNumber>
		      		</v12:Contact>
		      		<v12:Address>
		        		<v12:StreetLines>{{ .fedex.Shipper.Address.Street }}</v12:StreetLines>
		        		<v12:City>{{ .fedex.Shipper.Address.City }}</v12:City>
		        		<v12:StateOrProvinceCode>{{ .fedex.Shipper.Address.State }}</v12:StateOrProvinceCode>
		        		<v12:PostalCode>{{ .fedex.Shipper.Address.Zip }}</v12:PostalCode>
		        		<v12:CountryCode>US</v12:CountryCode>
		      		</v12:Address>
		    	</v12:Shipper>
		    	
		    	<v12:Recipient>
		      		<v12:Contact>
		        		<v12:PersonName>{{ .fedex.Recipient.PersonName }}</v12:PersonName>
		        		<v12:CompanyName>{{ .fedex.Recipient.CompanyName }}</v12:CompanyName>
		        		<v12:PhoneNumber>{{ .fedex.Recipient.Phone }}</v12:PhoneNumber>
					</v12:Contact>
					<v12:Address>
						<v12:StreetLines>{{ .fedex.Recipient.Address.Street }}</v12:StreetLines>
		        		<v12:City>{{ .fedex.Recipient.Address.City }}</v12:City>
		        		<v12:StateOrProvinceCode>{{ .fedex.Recipient.Address.State }}</v12:StateOrProvinceCode>
		        		<v12:PostalCode>{{ .fedex.Recipient.Address.Zip }}</v12:PostalCode>
		        		<v12:CountryCode>US</v12:CountryCode>
		        		<!--<v12:Residential>true</v12:Residential>-->
		      		</v12:Address>
		    	</v12:Recipient>
		    
		    	<v12:ShippingChargesPayment>
		      		<v12:PaymentType>SENDER</v12:PaymentType>
		      		<v12:Payor>
		        		<v12:ResponsibleParty>
		          			<v12:AccountNumber>{{ .fedex.Account.AcctNumber }}</v12:AccountNumber>
		          			<v12:Contact/>
		        		</v12:ResponsibleParty>
		      		</v12:Payor>
		    	</v12:ShippingChargesPayment>
		    
		    	<v12:LabelSpecification>
		      		<v12:LabelFormatType>COMMON2D</v12:LabelFormatType>
		      		<v12:ImageType>PNG</v12:ImageType>
					<v12:LabelStockType>PAPER_4X6</v12:LabelStockType>
					<v12:LabelPrintingOrientation>TOP_EDGE_OF_TEXT_FIRST</v12:LabelPrintingOrientation>
		    	</v12:LabelSpecification>
		    
		    	<v12:RateRequestTypes>ACCOUNT</v12:RateRequestTypes>
		    	{{ if .fedex.TrackingId }}
			    	<v12:MasterTrackingId>
			    		<v12:TrackingIdType>FEDEX</v12:TrackingIdType>
			    		<v12:FormId></v12:FormId>
			    		<v12:TrackingNumber>{{ .fedex.TrackingId }}</v12:TrackingNumber>
			    	</v12:MasterTrackingId>
			    {{ end }}
		    	<v12:PackageCount>{{ .fedex.PackageCount }}</v12:PackageCount>
		    	
		    	<v12:RequestedPackageLineItems>
		      		<v12:SequenceNumber>{{ .package.SequenceNumber }}</v12:SequenceNumber>
		      		<v12:Weight>
		        		<v12:Units>LB</v12:Units>
		        		<v12:Value>{{ .package.Weight }}</v12:Value>
		      		</v12:Weight>
		    	</v12:RequestedPackageLineItems>

		  	</v12:RequestedShipment>

		</v12:ProcessShipmentRequest>

	</soapenv:Body>
</soapenv:Envelope>`
