# applog endpoint app

## development

for cli testing of the WebSocket stream, try [wscat](http://einaros.github.io/ws/), which can be done easily on a Stackato VM as:

```
stackon-fetch http://stackato:suchDogeW0w@docker-internal.stackato.com master activestate/wscat
export TOKEN='….'  # set token here
export APPGUID='…' # set GUID here
docker run -i -t activestate/wst wscat -c "ws://logs.stackato-abcd.local/tail?token=$TOKEN&appid=$APPGUID"

```
