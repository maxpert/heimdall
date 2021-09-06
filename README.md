# Heimdall

Listen to Postgres changes and publish them over 
various channels (Kafka, HTTP endpoint). This repo is 
WIP early PoC.


## TODO

 - Complete HTTP hook implementation including connection reuse/pooling
 - Introduce a scripting engine
 - Work on backpressure system to make sure slow upstream won't cause memory issues to server
 - Introduce Kafka support


# Why?

I am starting off with Postgres (but will expand), and I want to 
introduce a fully scriptable engine rather than a limited versions of 
transformers (like in [debezium](https://debezium.io/)), and with popular 
sinks it can be used in combination with [benthos](https://www.benthos.dev/) 
to compose a more complex systems.
