package main

type Config struct {
	LLM struct {
		Host			string 	`json:"host"`
		Port			string 	`json:"port"`
		IsHttps		bool		`json:"isHttps"`
		Endpoint	string 	`json:"endpoint"`
		Key				string 	`json:"key"`
	} `json:"llm"`
}

type AnytypeObjIds struct {
	SpaceId		string
	ObjectId	string
}
