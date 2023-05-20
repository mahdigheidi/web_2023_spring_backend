# web_2023_spring_backend
The following repository is for the backend assignment of the web development course I did in SUT in 2023 Spring, lectured by Dr. JafariNezhad.

### To generate the self-signed certificate
We use the command below:
`sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/nginx-selfsigned.key -out certs/nginx-selfsigned.crt`

This command generates the key & cert needed by nginx to add the https directory.

### To generate the proto go files cd into the service folder & use the commands below
`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
`protoc -I=. --go_out=. ./services.proto`