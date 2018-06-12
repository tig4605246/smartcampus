#!/bin/sh


case "$1" in

  start)

    echo -n "Starting smartcampus agent: \n"
          cd /home/ntust/

          ./smartcampus -me:7080/gateway/v1/rawdata -aemurl=http://beta2-api.dforcepro.com:7080/gateway/v1/rawdata &


 	;;

  stop)

    echo -n "Stoping smartcampus agent: \n"
          smartcampus_PID=`cat /tmp/smartcampus_PID`
          kill -9 ${smartcampus_PID}
         

 	;;

esac

return 0
