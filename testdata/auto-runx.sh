#!/bin/bash


_debug(){
	#func
	local pfx=$( printf "%-10s" $(caller 0 | cut -f2 -d" " ) 2>>/dev/null)
	#logfile
	#pid
	local pid=$(echo $(printf "%07d" "$$"))
	echo "[`date '+%a %b %d %H:%M:%S'`,$( echo $(date '+%N')|awk '{printf("%03d", $0/1000000);}' )][$pid][$pfx] INFO : $*"
}

_initALL(){

	MYSQL_PWD=""
	MYSQL_PASS=" mysql -uroot -p${MYSQL_PWD}"
	MYSQL_NPASS="mysql -uroot "
	MYSQL_BIN="${MYSQL_NPASS}"
	if [[ "${#MYSQL_PWD}" != "0" ]]; then
		#YES_ROOT_PASSWORD
		MYSQL_BIN="${MYSQL_PASS}"
	fi
}

_prepDb(){

	_debug "Start"
	for sql in "DROP DATABASE IF EXISTS  benjerry;" "DROP DATABASE IF EXISTS  benjerry_dev;" "DROP USER IF EXISTS  benjerry,benjerry_dev;" 
	{
		#dump it
		echo "$sql" | ${MYSQL_BIN}
	}
	_debug "Done"

}

_makeDb(){
	_debug "Start"
	for sql in "db_prod.sql" "dump_prod.sql" "db_dev.sql" "dump_dev.sql" 
	{
		#dump it
		_debug "Loading .... ${sql}"
		[[ ! -e "${sql:-xxxx}" ]] && {
		   _debug "Oops: NOT_FOUND: ${sql}"
		   continue	
		}
		#load it
		cat ${sql} | ${MYSQL_BIN}
	}
	_debug "Done"
}

#defaults
_initALL

_debug "Start"

#free 1st
_prepDb

#create
_makeDb

_debug "Done"
