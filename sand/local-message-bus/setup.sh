#!/bin/sh

# restart rabbitmq
systemctl restart rabbitmq-server

# enable management plugin
# create a dashboard to monitor + configure the broker
# accessible under http://127.0.0.1:15672
# don't forget to add "-p 5672:5672 -p 15672:15672" when creating container
rabbitmq-plugins enable rabbitmq_management

# create an admin user
# outside localhost (which would be the container), "guest" "guest" has no access
# hence new user with name + password "jeff"
# https://www.youtube.com/watch?v=AfIOBLr1NDU
rabbitmqctl add_user "jeff" "jeff"
rabbitmqctl set_permissions -p "/" "jeff" ".*" ".*" ".*"
rabbitmqctl set_user_tags "jeff" "administrator"

# leave the container running
tail -f /dev/null