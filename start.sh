#!/bin/sh


case "$1" in

  start)

    echo -n "Starting smartcampus agent: \n"
          cd /home/bmw/smarcampus/

          ./smartcampus -chiller -diskpath=/dev/sda1 -postmac=aa:bb:03:01:01:02 gwid=chiller_02 &


 	;;

  stop)

    echo -n "Stoping smartcampus agent: \n"
          smartcampus_PID=`cat /tmp/smartcampus_PID`
          kill -9 ${smartcampus_PID}
         

 	;;

esac

return 0
