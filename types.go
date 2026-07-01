package main

type Config struct {
	LLM struct {
		Endpoint	string 	`json:"endpoint"`
		Token			string 	`json:"token"`
	} `json:"llm"`
	Anytype struct {
		Endpoint	string	`json:"endpoint"`
		Token			string	`json:"token"`
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
