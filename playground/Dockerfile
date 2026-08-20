FROM debian:buster-slim

WORKDIR /corewar

RUN apt-get update && apt-get install -y bsdmainutils

COPY . . 

ENTRYPOINT ["/bin/bash"]
