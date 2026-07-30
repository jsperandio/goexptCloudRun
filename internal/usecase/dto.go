package usecase

type GetWeatherByZipcodeInput struct {
	Zipcode string `validate:"required,len=8,number"`
}

type GetWeatherByZipcodeOutput struct {
	TempC float64
	TempF float64
	TempK float64
}
