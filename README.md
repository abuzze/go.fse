# FS Economy All In Profit Finder

## Purpose

The purpose of this tool is to find the top 10 paying "Airline Pilot for Hire" jobs for
a specific aircraft.

## Install

Copy the binary, the config.sample.json and the airport-codes_json.json to the same folder.
go.fse is the MacOS version. go.fse.exe is for Windows.
Run the binary in a terminal like iTerm (MacOS) or cmd (Windows)

## Config

You need to rename config.sample.json to config.json and edit your userkey and the
desired aircraft type. You will find you userkey on the (FSEconomy Website)[https://server.fseconomy.net/] via
Home > Datafeeds in the Access Key field. It looks like this 3EF7A5E5746994GF.
The config.json must be in the some folder as the go.fse tool.

The aircraft model should be the same like in the "Search Airport" dropdown menu "Airports that have this aircraft".

Set terminal to cmd for Windows or to bash, if you like colors in your terminal.

From and to can contain country iso codes like "US", "UK" or "DE" to narrow to the search for only a particular list of country.
Both can be a list of multiple countries. If you want to search for all jobs that start from China or the US use this list

  "from":["US","CN"],

The same applies for the to list. If you don't want to restrict it, use ["all"]. Please note that the square brackets surround the list, even
if it contains only one item.

  "to":["all"],

You can limit the number of jobs shown be setting the searchlimit.

  "searchlimit":4

## Airport Data

The Airport data was downloaded from (datahub.io)[https://datahub.io/core/airport-codes] You can update the file, if you keep
the name.


## Licence

Free to use.
