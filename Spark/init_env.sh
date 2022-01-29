# /bin/bash
if ! type java &> /dev/null ; then
    apt-get install -y default-jre default-jdk
fi
if [ ! -d "/usr/local/spark" ]; then
	wget https://d3kbcqa49mib13.cloudfront.net/spark-2.2.0-bin-hadoop2.7.tgz && \
	tar -xvf spark-2.2.0-bin-hadoop2.7.tgz && \
	mv spark-2.2.0-bin-hadoop2.7 /usr/local/spark
fi
echo "環境初始化完畢"