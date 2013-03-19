#!/bin/bash

mkdir $1
cd $1
newgo modl.go "data models"
newgo view.go "view handlers"
newgo cont.go "http controller"
mkdir templates
mkdir static
mkdir static/js
mkdir static/css
mkdir static/img
mkdir uploads
echo "project started!"