# Simple Port Application Checker
Small HTTP Proxy to remotely check if ports are in use by other applications.

### About
The usage is simple: offloading the port checking of your application to another process,  
which checks your application internally by querying the availability of the port.  
  
This way your application log will not be polluted by health-check connections that are  
being triggered externally and you do not have to open your application to the world,  
because you can run spac on any port you like.

### Example
Let's take a simple ftp service as an example.  
We configure ftp access (port 21) based on security group and ip restrictions.  
We configure spac to check the port 21 internally and open port 9000 for the world.  
We configure the route 53 health-checks to check ftp via spac on port 9000 

```
                    +-----------------+        +----------------------+
                    |   route 53      |        | port access based on |
                    | health checkers |        |    security group    |
                    +--------+--------+        +----------+-----------+
                             |                            |
   http result == 200 ok     |                            |
   http result >= 400 error  |                            |
                             v                            v
                    +-------------------------------------------------+
                    | +-----------------+           +---------------+ |
                    | |  spac on port   |           |FTP Application| |
                    | |      9000       |           |  on port 21   | |
                    | |                 |           |               | |
                    | +--------+--------+           +-------+-------+ |
                    |          |                            ^         |
                    |          |                            |         |
                    |          +----------------------------+         |
                    |          checks availability of port 21         |
                    +-------------------------------------------------+
```

### Config

In this setup, the spac configuration would look like:
```
{
  "server": "0.0.0.0",
  "port": 9000,
  "apps": {
    "ftp": ":21"
  }
}
```
The Route53 health-check machines can query the status on:  
`http://my.example.com:9000/services/ftp`
