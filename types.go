package main

import "github.com/epheo/anytype-go"

type Config struct {
	LLM struct {
		Endpoint	string 	`json:"endpoint"`
		Token			string 	`json:"token"`
	} `json:"llm"`
	Anytype struct {
		Endpoint	string	`json:"endpoint"`
		Token			string	`json:"token"`
		TagId			string	`json:"tagid"`
	}	`json:"antrype"`
}

type AnytypeObjIds struct {
	SpaceId		string
	ObjectId	string
}

type SuggestionTags struct {
	Tag			string
	Reason	string
}

type Application struct {
	config	Config
	Anytype struct {
		client anytype.Client
		objIds AnytypeObjIds
	}
}
