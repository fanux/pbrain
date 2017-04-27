#交付时使用的基础镜像
FROM 10.1.86.51/devops/golang:1.7-alpine 
COPY pbrain $GOPATH/bin 
CMD pbrain 
