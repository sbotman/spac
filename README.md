# Simple Port Application Checker
Checking if ports are in use by other applications

### About
The usage is simple: offloading the port checking of your application to another process,
which checks your application internally by querying the availability of the port.  
  
This way your application log will not be polluted by health-check connections that are  
being triggered externally and you do not have to open your application to the world,  
because you can run spac on any port you like.

### Overview



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


