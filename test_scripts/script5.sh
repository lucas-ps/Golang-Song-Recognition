#!/bin/sh
ID="b"
AUDIO=`base64 -i "$ID".wav`
RESOURCE=localhost:3001/search
echo "{ \"Audio\":\"$AUDIO\" }" > input
curl -v -X POST -d @input $RESOURCE
