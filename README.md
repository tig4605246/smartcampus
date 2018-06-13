## Gateway Daemon for Aemdra and CPM70 Meters (GDACM)

A backup plan for smart campus gw daemon. This repo also contains the test function for airbox.

Build for Beaglebone:
````
env GOOS=linux GOARCH=arm go build
````

Usage of flags:
````
  -aemurl string
        post url for aem-dra (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
  -chiller
        run in chiller mode
  -chillerurl string
        post url for chiller (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
  -cpmurl string
        post url for cpm (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
  -cpupath string
        cpu path (default "/proc/stat")
  -diskpath string
        disk path (default "/dev/mmcblk0p1")
  -gwserial string
        Declare GW serial number (default "03")
  -help
        a bool
  -macfile
        Use macFile to set up mac serial numbers 
  -meter
        run in smart meter mode
  -test
        test all functions
````