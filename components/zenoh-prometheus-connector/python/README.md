# ZENOH-SUSCRIBER / PROMETHEUS-EXPORTER

URLs: 
- https://github.com/eclipse-zenoh/zenoh-python


## REQUIREMENTS AND DEPLOYMENT

```
pip install prometheus_client

pip install eclipse-zenoh 
```


Build the Docker image:

```bash
docker build -t zenoh-prometheus-connector .

docker tag 76d68172d12e zenoh-prometheus-conn:<version>
```

Run the container:

```bash
docker run -d -ti -p 8999:8999 zenoh-prometheus-conn:latest
```

----------------------------------------------

## ZENOH 

1. Put a key/value into Zenoh:

```bash
curl -X PUT -H "content-type:application/json" -d "122" http://192.168.137.47:8000/tests/example/test2
```

2. Get values

```bash
curl http://192.168.137.47:8000/tests/**
```

--------------------------------------------------------------

# TESTS: PUT METRICS

```bash
curl -X PUT -H "content-type:application/json" -d "2" http://192.168.137.47:8000/tests/planta01/habitacion01

curl -X PUT -H "content-type:application/json" -d "3" http://192.168.137.47:8000/tests/planta01/habitacion02

curl -X PUT -H "content-type:application/json" -d "1" http://192.168.137.47:8000/tests/planta02/habitacion01

curl -X PUT -H "content-type:application/json" -d "2" http://192.168.137.47:8000/tests/planta02/habitacion02

curl -X PUT -H "content-type:application/json" -d "3" http://192.168.137.47:8000/tests/planta02/habitacion03
```


--------------------------------------------------------------

# TESTS: READ METRICS IN PROMETHEUS


```
colmena_total_people{metric_name="tests"}

colmena_total_people{metric_name="tests", label1="planta01"}

colmena_total_people{metric_name="tests", label1="planta02"}

sum by (metric_name, label1) (colmena_total_people{metric_name="tests", label1="planta01"})
```