FROM openjdk:8-jre

MAINTAINER impiggy (ideshenghe@gmail.com)

WORKDIR /c-admin
COPY bin bin
COPY conf conf
COPY lib lib
COPY logs logs

RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN rm -Rf /var/lib/apt/lists/*
RUN apt update && apt install -y file

RUN chmod +x /c-admin/bin/startup.sh
ENTRYPOINT ["sh", "/c-admin/bin/startup.sh"]
