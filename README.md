Here lies the piece of tsunginator that will: Collect metadata and store it using Consul k/v, injest a config.json file for use in configuring the test xml to the available nodes, and other things really really cool.


Scratch notes on writing a consul service check:

./empd -names
epmd: up and running on port 4369 with data:

#check if erlang dist communications are running
#check if in use (any beams) *issue warning*
    retrieve test time
    set timer
    if > ttl then issue warning 1
    if < ttl then issue fail 2
    else issue passing 0
#is the master node key value in /etc/hosts

#register hostname in key/value store?
#status updates, running, 
