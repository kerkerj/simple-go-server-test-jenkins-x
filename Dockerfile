FROM scratch
EXPOSE 8080
ENTRYPOINT ["/simple-go-server-test"]
COPY ./bin/ /