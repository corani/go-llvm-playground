#!/bin/bash
echo "building"
go build -o compile .

echo "example 1"
./compile -i <( echo "123-321" ) -o calc && ./calc

echo "example 2"
./compile -i example.ex -o calc && ./calc

rm calc compile
