@echo off
chcp 65001 >nul
echo hold on right installing the pie-rum sdk 😃
go mod tidy
echo just a sec running the server 🤗
go build -o app.exe
echo now running the file 🌟
app.exe
echo the pie-rum server started 🤩