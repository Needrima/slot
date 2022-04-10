package models

type CovidInfo struct {
	Country      string `json:"country"`
	TotalCases   int    `json:"cases"`
	TodaysCases  int    `json:"todayCases"`
	TotalDeaths  int    `json:"deaths"`
	TodaysDeaths int    `json:"todayDeaths"`
	Recovered    int    `json:"recovered"`
	Active       int    `json:"active"`
	Critical     int    `json:"critical"`
	CPM          int    `json:"casesPerOneMillion"`  //Cases per million
	DPM          int    `json:"deathsPerOneMillion"` //Deaths per million
	TotalTests   int    `json:"totalTests"`
	TPM          int    `json:"testsPerOneMillion"` //Tests per million
}

type WeatherInfo struct {
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	MinTemp     float64 `json:"temp_min"`
	MaxTemp     float64 `json:"temp_max"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
	SeaLevel    float64 `json:"sea_level"`
	GroundLevel float64 `json:"grnd_level"`
}

type Weather struct {
	City        string `json:"name"`
	WeatherInfo `json:"main"`
}
