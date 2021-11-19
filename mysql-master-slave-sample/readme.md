## Introduction
This repo will help you to run a mysql master-slave in local via docker.

## Quick Start
1. use `docker-compose up -d` to start the service
1. use `docker network inspect mysql-master-slave-sample_myweb` to see the master's IP address.
	* the pattern for network is `foldername` + `_` + `network name`, you can see the network name in `docker-compose.yml`
	* you can also use `docker network ls` to check the network name.
1. enter the master container `docker exec -it <master-container> /bin/bash` to check the master node.
	1. can use `docker ps` to see the containerID.
	1. login the mysql instance and check master status.
	```shell
	mysql> show variables like '%server_id%';
	+----------------+-------+
	| Variable_name  | Value |
	+----------------+-------+
	| server_id      | 1     |
	| server_id_bits | 32    |
	+----------------+-------+
	2 rows in set (0.00 sec)

	mysql> show master status;
	+------------------+----------+--------------+------------------+-------------------+
	| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
	+------------------+----------+--------------+------------------+-------------------+
	| mysql-bin.000009 |      759 |              | mysql            |                   |
	+------------------+----------+--------------+------------------+-------------------+
	1 row in set (0.01 sec)
	```
1. enter the slave container to config the slave. `docker exec -it <slave-container> /bin/bash`
	```shell
	CHANGE MASTER TO
    MASTER_HOST='172.29.0.2', // the IP address in step <2>, 
    MASTER_USER='root',
    MASTER_PASSWORD='root',
    MASTER_LOG_FILE='mysql-bin.000009', // the file in <show master status>
    MASTER_LOG_POS=759; // the position in <show master status>
	```
	here's an examle
	```shell
	mysql> CHANGE MASTER TO
    ->     MASTER_HOST='172.29.0.2',
    ->     MASTER_USER='root',
    ->     MASTER_PASSWORD='root',
    ->     MASTER_LOG_FILE='mysql-bin.000009',
    ->     MASTER_LOG_POS=759;
	Query OK, 0 rows affected, 2 warnings (0.06 sec)
	```
1. use `slave start` to start the slave.
	```shell
	mysql> start slave;
	Query OK, 0 rows affected (0.01 sec)
	```
1. use `show slave status` to check the slave status, if you see the `Slave_IO_Running` and `Slave_SQL_Running` is `Yes`, that means config success.
	```shell
	mysql> show slave status\G;
	*************************** 1. row ***************************
				Slave_IO_State: Waiting for master to send event
					Master_Host: 172.29.0.2
					Master_User: root
					Master_Port: 3306
					Connect_Retry: 60
				Master_Log_File: mysql-bin.000009
			Read_Master_Log_Pos: 1043
				Relay_Log_File: f29ade3cf371-relay-bin.000002
					Relay_Log_Pos: 604
			Relay_Master_Log_File: mysql-bin.000009
				Slave_IO_Running: Yes
				Slave_SQL_Running: Yes
				Replicate_Do_DB:
			Replicate_Ignore_DB: mysql
			Replicate_Do_Table:
		Replicate_Ignore_Table:
		Replicate_Wild_Do_Table:
	Replicate_Wild_Ignore_Table:
					Last_Errno: 0
					Last_Error:
					Skip_Counter: 0
			Exec_Master_Log_Pos: 1043
				Relay_Log_Space: 818
				Until_Condition: None
				Until_Log_File:
					Until_Log_Pos: 0
			Master_SSL_Allowed: No
			Master_SSL_CA_File:
			Master_SSL_CA_Path:
				Master_SSL_Cert:
				Master_SSL_Cipher:
				Master_SSL_Key:
			Seconds_Behind_Master: 0
	Master_SSL_Verify_Server_Cert: No
					Last_IO_Errno: 0
					Last_IO_Error:
				Last_SQL_Errno: 0
				Last_SQL_Error:
	Replicate_Ignore_Server_Ids:
				Master_Server_Id: 1
					Master_UUID: 82e97859-4775-11ec-986b-0242ac130002
				Master_Info_File: /var/lib/mysql/master.info
						SQL_Delay: 0
			SQL_Remaining_Delay: NULL
		Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
			Master_Retry_Count: 86400
					Master_Bind:
		Last_IO_Error_Timestamp:
		Last_SQL_Error_Timestamp:
				Master_SSL_Crl:
			Master_SSL_Crlpath:
			Retrieved_Gtid_Set:
				Executed_Gtid_Set:
					Auto_Position: 0
			Replicate_Rewrite_DB:
					Channel_Name:
			Master_TLS_Version:
	1 row in set (0.00 sec)
	```

Finally, you can create the DB and table in master, and check the status.
