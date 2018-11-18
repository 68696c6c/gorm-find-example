FROM golang:1.9

RUN echo 'alias ll="ls -lahG"' >> ~/.bashrc

RUN go get -u github.com/golang/dep/cmd/dep
