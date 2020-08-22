FROM openjdk:8-jre

MAINTAINER impiggy (ideshenghe@gmail.com)

WORKDIR /c-deployer
COPY bin bin
COPY conf conf
COPY lib lib
RUN mkdir logs

RUN rm -f /c-deployer/bin/canal.pid

RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN rm -Rf /var/lib/apt/lists/*
RUN apt update && apt install -y file
RUN ls -a /c-deployer/bin

EXPOSE 11110 11111 11112 9100

RUN chmod +x /c-deployer/bin/startup.sh
RUN chmod +x /c-deployer/bin/cmd.sh
ENTRYPOINT ["sh", "/c-deployer/bin/cmd.sh"]
