#!/bin/sh

ROOT_PATH="$2"
BIN_PATH="$ROOT_PATH/bin"
TEMP_PATH="$ROOT_PATH/temp"
EXE="$BIN_PATH/$1"

# 记录下pid
echo $$ > "${TEMP_PATH}/$1.pid"

# 判断文件是否被锁定,再执行动作
{
	flock -xn 6
	if [ "$?" -eq "1" ]; then
		exit -1
	fi
	# 进入文件循环
	while (true)
	do
		chmod +x $EXE
		nohup sudo $EXE > /dev/null 2>&1
		sleep 1s
	done
} 6<>"${TEMP_PATH}/$1.pid"
