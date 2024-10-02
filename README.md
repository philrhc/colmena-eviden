# COLMENA
COLMENA Project.

## Architecture Reference:

### Component Diagram:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/9e7cab42-58b0-4001-8304-97e74799201f)

### REST APIs overview per component:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/e342be37-d2af-4ad7-904f-c520c07162dc)


### Resulting Packages:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/8339995c-e02a-4b0b-bf3d-5b2cd56de5d5)

## Sequence Diagrams:

### Deployment:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/9da74b11-0515-4874-af58-d57c5a0ad0aa)

### Update Deployment:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/3421434b-297a-4b57-8b90-414dc1bca794)

### Context Awareness Management:

![image](https://github.gsissc.myatos.net/storage/user/7653/files/8bbbd125-a746-4d80-a208-61c80ab7f20c)

### SLA/SLO Management:
![image](https://github.gsissc.myatos.net/storage/user/7653/files/5acab036-ec3b-4b42-9fad-677b61d78625)


The project is composed by the following services:

- [aggregator](aggregator-svc/README.md)
    - golang microservice
    - prometheus / thanos
    - node_exporter
- [context awareness manager](context-awareness-manager-svc/README.md)
    - golang microservice
- [sla manager](sla-manager-svc/README.md)
    - golang microservice
