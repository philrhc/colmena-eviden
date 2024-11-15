######################################
FROM python:3.12-slim AS compile-image

RUN apt-get update
RUN apt-get install -y --no-install-recommends build-essential gcc

COPY . .

RUN pip install --user prometheus_client 
RUN pip install --user eclipse-zenoh 


######################################
FROM python:3.12-slim AS build-image

COPY --from=compile-image /root/.local /root/.local

EXPOSE 8999

ENV PATH=/root/.local/bin:$PATH

CMD ["python", "./main.py"]

