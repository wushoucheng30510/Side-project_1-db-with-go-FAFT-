#!/usr/bin/env bash
# Run the tests and put the result in DB
# EX: 

user=$(pwd | cut -d '/' -f3)

for i in {1..$2}
	do
	for test in $1
	do
		rm ~/chromiumos/src/scripts/run/tmp/log.txt
		rm ~/chromiumos/src/scripts/run/tmp/result.json
		echo -e "\n\n                       ========> Run test: \033[5m\033[31m$test\033[0m <========\n"
		echo 

		year=`date +%Y`
		month=`date +%m`
		day=`date +%d`
		hour=`date +%H`
		minute=`date +%M`
		second=`date +%S`
		empty=" "
		time=$year-$month-$day$empty$hour:$minute:$second
		fileHeader=$year$day$hour$minute$second
		echo $time
		case=$test
		path="/home/$user/chromiumos/src/scripts/run/tmp"
		
		tast run -var=servo=localhost:9998 10.240.102.203 firmware.$test | tee -a $path/log.txt 
		
		echo -e "\n               =============Finish============\n\n"
		
		echo "Start parse the data in shell"
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

		if [ $PorF == "Pass" ] || [ $PorF == "Fail" ]
		then 

			dutVersion=$(cat $path/log.txt | grep "Primary DUT version" | cut -d ' ' -f5)
			logPath=$(cat $path/log.txt | grep Writing | cut -d ' ' -f5)
			model=$(cat "/home/$user/chromiumos/chroot$logPath/dut-info.txt" | grep model | cut -d '"' -f2)
			board=$(echo "$dutVersion" | cut -d '/' -f1 | cut -d '-' -f1)
			version=$(echo "$dutVersion" | cut -d '/' -f2)

			echo "{
			\"time\":\"$time\",
			\"tester\":\"$user\",
			\"name\":\"$test\",
			\"board\":\"$board\",
			\"model\":\"$model\",
			\"version\":\"$version\",
			\"logPath\":\"$fileHeader-$user.txt\",
			\"result\":\"$PorF\"

			}" > $path/result.json
		

			echo -e "\033[31mStart to process data on DB\033[0m\n"
			cd ~/chromiumos/src/scripts/run/backend/database
			go run main.go
			scp $path/log.txt ubuntu@10.240.102.16:/home/ubuntu/backend/logDB/$fileHeader-$user.txt
		fi
	done
done


