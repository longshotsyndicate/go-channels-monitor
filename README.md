# go-channels-monitor
go-channels-monitor is a utility and service for keeping track of channel properties and reporting on them. 

##Install
```bash
$ go get github.com/longshotsyndicate/go-channels-monitor
```

##Usage
```go

  foo := make(chan bool, 100)
  
  // add your channel to be monitored.
  monitor.AddNamed("foo", channel)
  
  properties := monitor.Get("foo")
  
  log.Printf("foo len: %d cap: %d", properties.Len, properties.Cap)
  
  //add this monitor to the service
  
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

##Query Response
Making a GET request on `your-ip:9999/channels` will result in the following reponse:
```json
{"service":"my-service",
  "channels":{
    "foo":{"length":0,"capacity":100}
  }
}
```








