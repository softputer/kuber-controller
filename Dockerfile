FROM hub.qingyuanos.com/admin/haproxy
MAINTAINER  Jayson  yjge@qingyuanos.com

COPY kuber-controller  /
COPY ./provider/haproxy/haproxy_reload  /etc/haproxy/
COPY ./provider/haproxy/*.cfg  /etc/haproxy/
RUN chmod +x /etc/haproxy/haproxy_reload
CMD ["/kuber-controller"]
