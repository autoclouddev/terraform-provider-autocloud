package properties

import "os"

type Properties struct {
	ApiHost         string
	AppClientID     string
	UserPoolID      string
	AppClientSecret string
}

func LoadProperties() Properties {
	// err := godotenv.Load("./autocloud_sdk/.env")

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	return Properties{
		ApiHost:         os.Getenv("SDK_API_URL"),
		AppClientID:     os.Getenv("SDK_COGNITO_APP_CLIENT_ID"),
		UserPoolID:      os.Getenv("SDK_COGNITO_USER_POOL_ID"),
		AppClientSecret: os.Getenv("SDK_COGNITO_APP_CLIENT_SECRET"),
	}
}

/*
ApiHost:         "http://localhost:8080/graphql",  //os.Getenv("SDK_API_URL"),
		AppClientID:     "23f966jp1nes26piqfo19bmv04",     //os.Getenv("SDK_COGNITO_APP_CLIENT_ID"),
		UserPoolID:      "us-east-1_yNaj3AKmz",            //os.Getenv("SDK_COGNITO_USER_POOL_ID"),
		AppClientSecret: "jhAupWgfm3xPdBmRfBWGA6caqcsNCa", //os.Getenv("SDK_COGNITO_APP_CLIENT_SECRET"),
*/
