FROM busybox:ubuntu-14.04

RUN mkdir -p /opt/aws-operator
ADD ./aws-operator /opt/aws-operator/aws-operator

RUN mkdir -p /opt/aws-operator/service
ADD service/templates/ /opt/aws-operator/service/templates

WORKDIR /opt/aws-operator

EXPOSE 8000
ENTRYPOINT ["/opt/aws-operator/aws-operator"]
