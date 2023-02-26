#!/bin/bash

echo 'mysql-server mysql-server/root_password password your_password' | sudo debconf-set-selections
echo 'mysql-server mysql-server/root_password_again password your_password' | sudo debconf-set-selections
sudo apt-get -y install mysql-server
sleep 4

mysql -u root -pyour_password -e "create database new_db_kbv;"
mysql -u root -pyour_password -e "use new_db_kbv;"

export DBUSER=root
export DBPASS=your_password

# chmod +x gomysql

# ./gomysql
