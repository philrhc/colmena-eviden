#!/usr/bin/pythonCOLMENA-DEPLOYMENT-TOOL
#  COLMENA-DEPLOYMENT-TOOL
#  Copyright © 2024 EVIDEN
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
#  This work has been implemented within the context of COLMENA project.

import sys
from prometheus_client import start_http_server, Counter, Gauge
import zenoh
import time

# Create a metric to track some event
ops_processed = Counter('myapp_processed_ops_total', 'The total number of processed events')

# Create a Prometheus gauge metric
colmena_total_people = Gauge('colmena_total_people', 'people in a room / floor', ['metric', 'path', 'desc'])

# zenoh listener
def listener(sample):
    print(f"{sample.key_expr} => {sample.payload.decode('utf-8')}")
    
    k = str(sample.key_expr)
    print(k)
    
    # Split the string by the '#' character
    arr = k.split("/")
    print(arr)
    
    vmetric_name = ""
    vlabel1 = ""
    
    if len(arr) == 1:
        vmetric_name = k
    elif len(arr) >= 2:
        vmetric_name = arr[0]
        for e in arr:
            vlabel1 = vlabel1 + "/" + e
    
    # Set the value of the Prometheus metric
    colmena_total_people.labels(metric=vmetric_name, path=vlabel1, desc='people in a room / floor').set(sample.payload.decode('utf-8'))
    
    ops_processed.inc()
    
    print("Processed an operation: [name=" + vmetric_name + ", path=" + vlabel1 + "]")

# main
def main():
    ###########################################################################
    # ZENOH
    # Initialize Zenoh
    config = zenoh.Config().from_file("./zenoh_config.json5")
    session = zenoh.open(config)
    print(session.__dict__)
    
    # suscriber
    subscriber = session.declare_subscriber('tests/**', listener)
    
    ###########################################################################
    # PROMETHEUS
    # Start up the server to expose the metrics.
    start_http_server(8999)
    
    # Simulate some work
    while True:
        try:
            print("Waiting for metrics ...")
            time.sleep(120)
        except KeyboardInterrupt:
            sys.exit()

if __name__ == '__main__':
    main()
