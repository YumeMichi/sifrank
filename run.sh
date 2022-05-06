#!/bin/bash
BIN_SERVICES=(
	"sifrank"
)

SERVICES=(
	${BIN_SERVICES[@]}
)

# 初始化目录
ROOT_PATH=`pwd`
BIN_PATH="$ROOT_PATH/bin"
TEMP_PATH="$ROOT_PATH/temp"

# create temp director
mkdir -p $TEMP_PATH

# 重启
function jar_restart()
{
	jar_stop
	jar_start
}

# 结束进程树
function treekill()
{
	local parent=$1
	# 获取所有孩子结点
	childs=(`ps -ef | awk -v parent=$parent 'BEGIN{ORS=" ";} $3==parent{print $2;}'`)
	# 结束自己
	sudo kill -9 $parent
	# 结束子进程
	if [ ${#childs[@]} -ne 0 ]
	then
		for child in ${childs[*]}
		do
			# 结束子进程
			treekill $child
		done
	fi
}

# 结束原来的进程
function jar_stop()
{
	# 枚举进程的中的服务
	for svr in ${SERVICES[@]}; do
		# 读取文件中存储的pid
		path_1="${TEMP_PATH}/${svr}"
		# 判断文件是否存在
		if [ ! -f "${path_1}.pid" ]; then
			continue
		fi
		# 根据pid结束进程
		read pid < "${path_1}.pid"
		path_2=`readlink /proc/${pid}/exe`
#		lock=`flock -x -n ${path_1}.pid echo 1`
#		# 比较存储的pid是否与路径相符
#		if [ "1" = $lock ]; then
#			continue
#		fi
		{
			flock -xn 6
			if [ "$?" -eq "1" ]; then
				echo "Stoping ${svr} (${pid})"
				treekill $pid
			fi
		} 6<>"${path_1}.pid"
		# 删除pid文件
		rm -f "${path_1}.pid"
	done
	# 完成
	# echo "jar stop success"
}

# 后台启动进程
function jar_start()
{
	chmod +x "${BIN_PATH}/binsvr.sh"
	for svr in ${BIN_SERVICES[@]}; do
		echo "Starting ${svr}"
		nohup "${BIN_PATH}/binsvr.sh" $svr ${ROOT_PATH} > /dev/null 2>&1 &
	done
}

if [ $# -lt 1 ]; then
	echo "use restart|start|stop"
	exit 0
fi

case $1 in
"start")
	jar_restart;;
"stop")
	jar_stop;;
"restart")
	jar_restart;;
*)
echo "use restart|start|stop"
exit 0
;;
esac
