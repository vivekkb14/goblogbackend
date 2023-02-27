# goblogbackend
Backend for blog operations

To build this repo, Please do `make`, it will generate executable binary.

start.sh is containing the instructions to install mysql and start blogging server. 
If you're not using docker-compose for starting server, Please uncomment [last two lines] of `start.sh`

To start using docker-compose, execute `start.sh` and run `docker-compose up`

To start using docker , execute `start.sh` and run command `make dockerise` and 
`docker run --network="host" --env DBUSER=root --env DBPASS=your_password goblog:latest`

## Test cases for goblogbackend

API integration test have been writtern using ginkgo (https://onsi.github.io/ginkgo/)
install ginkgo and run `ginkgo command` in httpserver
