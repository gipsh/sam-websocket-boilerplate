# sam-websocket-boilerplate

Example websockets application using lambdas and api gateway on AWS


![alt text](https://github.com/gipsh/sam-websocket-boilerplate/blob/master/aws-ws-bp.png)


## How it works

The app defines 3 lambda functions for handling the main events of a websocket service: 

- Connections: when a new client connects the `connectHandler` is called and the connection is added
on the dynamodb table 

- Disconnections: when a client disconnects the `disconnectHandler` removes the entry from dynamodb

- Message received: when a message arrives to the service the `defaultHandler` is called and first uses
the api-gateway inerface to answert to the client and then retrieves the connections from dynamodb and updates
the array with new messages. 


## Deploy

Just use `serverless depoy` 



