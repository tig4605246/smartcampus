#!/bin/sh


case "$1" in

  start)

    echo -n "Starting smartcampus agent: \n"
          cd /home/ntust/

          ./smartcampus -meter -gwserial=03 &


 	;;

  stop)

    echo -n "Stoping smartcampus agent: \n"
          smartcampus_PID=`cat /tmp/smartcampus_PID`
          kill -9 ${smartcampus_PID}
         

 	;;

esac

return 0
