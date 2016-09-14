# may be need -v on production
docker run --name shipyard-db \
    -e POSTGRES_USER=shipyard \
    -e POSTGRES_DB=shipyard \
    -e POSTGRES_PASSWORD=111111 \
    -p 5432:5432 -d postgres
