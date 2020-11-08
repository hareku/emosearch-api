go 1.15

module github.com/hareku/emosearch-api

replace github.com/hareku/emosearch-api => ./

require (
	cloud.google.com/go/firestore v1.3.0 // indirect
	firebase.google.com/go v3.13.0+incompatible
	github.com/aquasecurity/lmdrouter v0.3.0
	github.com/aws/aws-lambda-go v1.15.0
	github.com/aws/aws-sdk-go v1.35.23
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.6.0
	github.com/google/uuid v1.1.2
	github.com/guregu/dynamo v1.10.0
	github.com/urfave/cli v1.22.1 // indirect
	google.golang.org/api v0.34.0
)
