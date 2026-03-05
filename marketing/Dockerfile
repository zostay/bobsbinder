FROM hugomods/hugo:latest

WORKDIR /src

COPY . .

EXPOSE 1313

CMD ["hugo", "server", "--bind", "0.0.0.0", "--baseURL", "http://localhost/", "--appendPort=false", "--disableFastRender"]
