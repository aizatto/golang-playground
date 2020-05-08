module github.com/aizatto/golang-playground/auth0-test

go 1.14

replace auth => ./auth

replace session => ./session

require (
	auth v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/urfave/negroni v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	session v0.0.0-00010101000000-000000000000 // indirect
)
