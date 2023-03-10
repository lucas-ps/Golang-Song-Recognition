#!/bin/sh
ID="Everybody+(Backstreet's+Back)+(Radio+Edit)"
ESCAPED=`perl -e "use URI::Escape; print uri_escape(\"$ID\")"`
RESOURCE=localhost:3000/tracks/$ESCAPED
curl -v -X DELETE $RESOURCE
