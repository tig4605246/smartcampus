#!/bin/sh


case "$1" in

    startmeter)

        echo -n "Starting smartcampus meter agent: \n"
        cd /home/ntust/smartcampus/

        ./smartcampus -meter -macfile -gwserial=${2}  &
        #./smartcampus -test -diskpath=/dev/sda1 -macfile -gwserial=${2}


    ;;

    startchiller)

        echo -n "Starting smartcampus chiller agent: \n"
        cd /home/bmw/smartcampus/

        ./smartcampus -chiller -pastmac=${3} -gwid=${2} &


    ;;
        

    stop)

        echo -n "Stoping smartcampus agent: \n"
        smartcampus_PID=`cat /tmp/smartcampus_PID`
        kill -9 ${smartcampus_PID}
                

    ;;

esac

return 0
