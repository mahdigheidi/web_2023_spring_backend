# web_2023_spring_backend
The following repository is for the backend assignment of the web development course I did in SUT in 2023 Spring, lectured by Dr. JafariNezhad.

### To generate the self-signed certificate
We use the command below:
`sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/nginx-selfsigned.key -out certs/nginx-selfsigned.crt`

This command generate the key & cert needed by nginx to add the https directory.