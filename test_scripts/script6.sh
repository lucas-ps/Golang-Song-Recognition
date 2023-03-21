#!/bin/sh
ID="a1"
AUDIO=`base64 -i "$ID".wav`
URL=localhost:3002/cooltown
echo "{ \"Audio\":\"$AUDIO\" }" > input
curl -v -X POST -d @input $URL
