# go-channels-monitor
go-channels-monitor is a utility and service for keeping track of channel properties and reporting on them. 

##Install
```bash
$ go get github.com/longshotsyndicate/go-channels-monitor
```

##Usage
```go

  foo := make(chan bool, 100)
  
  // add your channel to be monitored, optionally specifying a string instance id for the channel.
  //in this example we are registering one of potentially many log-writer channels, each with an id like "instance-10"
  //instances need to be unique for the channel name but not globally unique. 
  monitor.AddNamed("log-writer", "instance-"+strconv.Itoa(nLogWriters++), channel)
  
  properties := monitor.Get("log-writer")
  
  log.Printf("log-writer len: %d, cap: %d, instance: %s", properties.Len, properties.Cap, properties.Instance)
  
```
## Service
The service type allows you to query for channel properties.
```go
  
  //async error reporting
  errc := make(chan error)
  go func() {
    for _, err := range errc {
      log.Printf("Error: %v", err)
    }
  }()
  
  //create and start the service that responds with channel properties.
  service.New("my-service", ":9999", errc).Start()
  
```

###Query Response
Making a GET request on `your-ip:9999/channels` will result in the following reponse:
```json
{"service":"my-service",
  "channels":{
    "foo":{"length":0,"capacity":100, "instance": "10"}
  }
}
```








