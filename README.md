# goblogbackend
Backend for blog operations

To build this repo, Please do "make", it will generate executable binary.

start.sh is containing the instructions to install mysql and start blogging server. If you're not using docker-compose for starting server, Please uncomment last two lines of "start.sh"

to start using docker-compose, execute "start.sh" and run "docker-compose up"

to start using docker , execute "start.sh" and run command "make dockerise" and "docker run --network="host" --env DBUSER=root --env DBPASS=your_password goblog:latest"