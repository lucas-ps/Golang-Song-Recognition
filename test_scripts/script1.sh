#!/bin/sh
ID="a"
ESCAPED=`perl -e "use URI::Escape; print uri_escape(\"$ID\")"`
AUDIO=`test`
RESOURCE=localhost:3000/tracks/$ESCAPED
echo "{ \"Id\":\"$ID\", \"Audio\":\"$AUDIO\" }" > input
curl -v -X PUT -d @input $RESOURCE  