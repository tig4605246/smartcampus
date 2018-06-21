# Gateway Daemon for Aemdra and CPM70 Meters (DAM)

A smart campus gw daemon. This program also contains the test function for airbox.

## Build Version

go version go1.9.2 linux/amd64

## Build for Beaglebone

````
env GOOS=linux GOARCH=arm go build
````

## Usage of Flags:

````
-aemurl (string)
    post url for aem-dra (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
-chillerurl (string)
    post url for chiller (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
-cpmurl (string)
    post url for cpm (default "https://beta2-api.dforcepro.com/gateway/v1/rawdata")
-cpupath (string)
    cpu path (default "/proc/stat")
-diskpath (string)
    disk path (default "/dev/mmcblk0p1")
-gwserial (string)
    Declare GW serial number (default "03")
-help
    Show information
-macfile
    Use macFile to set up mac serial numbers 
-meter (bool)
    run in smart meter mode
-chiller (bool)
    run in chiller mode
-airbox
    post to airbox for testing
-version
    Check the version
````

## Example Usage

````
# Run post test
$ ./smartcampus -airbox

# Run in smart meter mode
# Set gw serial to 05
# Set disk file path to /dev/sda1
$ ./smartcampus -meter -gwserial=05 -diskpath=/dev/sda1

# Run in chiller mode
# Set GW ID to chiller_02
# Set post mac to aa:bb:03:01:01:02
# Set disk file path to /dev/sda1
$ ./smartcampus -chiller -diskpath=/dev/sda1 -postmac=aa:bb:03:01:01:02 gwid=chiller_02

# Run smart meter inside wood house (Available after 1.2)
$ ./smartcampus -meter -woodhouse
````

## Known Issues

*  CPU check only process the first core
*  Disk check always return 0