#!/usr/bin/env bash
# Run the tests and put the result in DB
# EX: 

user=$(pwd | cut -d '/' -f3)

for test in $1
do
        echo -e "\n\n========Start to test \033[5m\033[31m$test\033[0m========\n"
        echo 
        time=`date +%Y%m%d%H%M`
        year=`date +%Y`
        month=`date +%m`
        day=`date +%d`
        date=$year-$month-$day
        case=$test
        path="/home/$user/chromiumos/database"
        # hide IP
        tast run -var=servo=localhost:9998 xx.xxx.xxx.$2 firmware.$test | tee -a $path/log.txt

        echo -e "\n=============Finish============\n\n"
        log=$(cat $path/log.txt)

        P=$(grep -c "\[ PASS \]" $path/log.txt)
        F=$(grep -c "\[ FAIL \]" $path/log.txt)

        if [ $F -gt 0 ]
        then
                cat $path/log.txt | grep "\[ FAIL \]" > $path/failReason.txt
                PorF="Fail"
        fi

        if [ $P -gt 0 ]
        then
                PorF="Pass"
        fi

        dutVersion=$(cat $path/log.txt | grep "Primary DUT version" | cut -d ' ' -f5)
        logPath=$(cat $path/log.txt | grep Writing | cut -d ' ' -f5)
        model=$(cat "/home/$user/chromiumos/chroot$logPath/dut-info.txt" | grep model | cut -d '"' -f2)
        board=$(echo "$dutVersion" | cut -d '/' -f1 | cut -d '-' -f1)
        version=$(echo "$dutVersion" | cut -d '/' -f2)

        echo "{
        \"time\":\"$date\",
        \"tester\":\"$user\",
        \"name\":\"$test\",
        \"board\":\"$board\",
        \"model\":\"$model\",
        \"version\":\"$version\",
        \"logPath\":\"$logPath\",
        \"result\":\"$PorF\"

        }" > $path/result.json

        echo -e "\033[31mStart to process data on DB\033[0m\n"
        cd ~/chromiumos/database
        go run main.go

done