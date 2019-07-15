# README

## About

Little helper script to generate load on the test ptrac deployment.
Originally developed to help test the cloudwatch monitoring stats

## Usage

`$ go run main.go -w 10 -min 60 -max 120`

## IMO file generation

```
ssh $PTRACIP "cd /opt/purpletrac; python3 manage.py shell -c 'from dashboard.models.ship import Ship; x = {s.imo_number for s in Ship.objects.all()}; print(x)'" >imos
```
