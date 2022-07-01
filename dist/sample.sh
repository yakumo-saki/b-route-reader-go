#!/bin/bash

echo "JSONファイル=$1"

WATT=`jq -r '."E7"' $1`
AMP=`jq -r '."E8"' $1`
DELTA=`jq -r '."E0"' $1`
DELTA_R=`jq -r '."E3"' $1`

if [ "$WATT" != "null" ]; then
    echo "消費電力は$WATT W"
    $SENDER 
fi
if [ "$AMP" != "null" ]; then
    echo "消費電流は$AMP A"
fi
if [ "$DELTA" != "null" ]; then
    echo "積算電力（正方向）は$DELTA kWh"
fi
if [ "$DELTA" != "null" ]; then
    echo "積算電力（逆方向）は$DELTA_R kWh"
fi
