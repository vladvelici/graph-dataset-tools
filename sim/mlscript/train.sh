#!/bin/bash

## Set matlab path to this variable
MATLAB_PATH=/Applications/MATLAB_R2014b.app/bin/matlab

MATLAB_OPTIONS="-nodisplay -nojvm -r"

MATLAB_COMMAND="train('"$1"', '"$2"', '"$3"', '"$4"'); exit;"

echo $MATLAB_COMMAND
$MATLAB_PATH $MATLAB_OPTIONS "$MATLAB_COMMAND"

