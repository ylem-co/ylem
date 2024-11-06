<p align="center">
  <img width="748" title="Ylem. The open-source data streaming platform" alt="Ylem. The open-source data streaming platform" src="https://github.com/user-attachments/assets/5f3bcbf3-92db-4bff-bdb5-68b88ddefcf4">
</p>

<div align="center">

  ![GitHub branch check runs](https://img.shields.io/github/check-runs/ylem-co/ylem/main?color=green)
  ![Static Badge](https://img.shields.io/badge/Go-1.23-black)
  ![Static Badge](https://img.shields.io/badge/React-18.3.1-black)
  <a href="https://github.com/ylem-co/ylem?tab=Apache-2.0-1-ov-file">![Static Badge](https://img.shields.io/badge/license-Apache%202.0-black)</a>
  <a href="https://github.com/ylem-co/ylem/tags">![Static Badge](https://img.shields.io/badge/tag-v0.0.1_pre_release-black)</a>
  <a href="https://ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/website-ylem.co-black)</a>
  <a href="https://docs.ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/documentation-docs.ylem.co-black)</a>
  <a href="https://join.slack.com/t/ylem-co/shared_invite/zt-2nawzl6h0-qqJ0j7Vx_AEHfnB45xJg2Q" target="_blank">![Static Badge](https://img.shields.io/badge/community-join%20Slack-black)</a>
</div>

# Ylem
The open-source data streaming platform is a one-stop-shop solution for orchestrating data streams on top of Apache Kafka, Amazon SQS, Google Pub/Sub, RabbitMQ, various APIs, and data storages.

<img width="1158" alt="Screenshot 2024-10-18 at 13 20 37" src="https://github.com/user-attachments/assets/fee384d3-bc10-4681-a3d7-a8a9f5dc8983">

|    |  |  |
| ------------- | ---------- | ----------- |
| <img  alt="User dashboard" src="https://github.com/user-attachments/assets/3740b8bd-f127-4ccb-8496-3b68888c5846"> | <img alt="Pipeline running" src="https://github.com/user-attachments/assets/00520fa1-3912-4c9a-850e-546fdd541a06"> | <img alt="Pipeline log" src="https://github.com/user-attachments/assets/628e1be5-5a45-4935-a9fa-6769edada491">  |

# Installation

## Install Docker 4

If you don't yet have Docker 4 installed, [install](https://www.docker.com/products/docker-desktop/) it from their official website for your OS.

## Install Ylem

### Option 1. Install from pre-build containers

The best way to install Ylem is to clone the repository https://github.com/ylem-co/ylem-installer and follow the installation instructions from it. It will install Ylem from the latest version of pre-build containers stored on Docker Hub.

Ylem will be available at http://localhost:7331/

### Option 2. Build and install from the source

If you want to compile Ylem from the source, run `docker compose up` or `docker compose up -d` from this repository. It will compile the code and run all the necessary containers.

Ylem is available at http://127.0.0.1:7330/

:warning: Please pay attention. Compiling from the source might take some time and will keep the resources of your machine busy.

#### To rebuild a particular container

If you want to rebuild a particular container from a source locally, run the following:

``` bash
docker compose build --no-cache %%CONTAINER_NAME%%
```

E.g.

``` bash
docker compose build --no-cache ylem_users
```

# Using your own Apache Kafka cluster

Ylem uses Apache Kafka to exchange messages for processing pipelines and tasks. By default Ylem already comes with the pre-configured Apache Kafka container, however, you might already have an Apache Kafka cluster in your infrastructure and might want to reuse it.

In this case, you need to take the steps described [here in our documentation](https://docs.ylem.co/open-source-edition/usage-of-apache-kafka).

# Configuring environment variables in .env files

Some particular integrations might require extra steps and using `.env` files. Configure them if you need to.

The list of such integrations and more information about them is in [our documentation](https://docs.ylem.co/open-source-edition/configuring-integrations-with-.env-variables).

# Folder structure in this repository

Ylem is a set of microservices. Each microservice is represented by one or more containers in the same network and communicates with each other via the API.

``` bash
|-- api                  # api microservice
|-- backend
|--|-- integrations      # integrations with external APIs, databases and other software
|--|-- pipelines         # pipelines, tasks, connectors
|--|-- statistics        # statistics of pipeline and task runs
|--|-- users             # users and organizations
|-- database             # a container for storing databases for all the microservices
|-- processor
|--|-- python_processor  # processor of the Python code written in pipelines 
|--|-- taskrunner        # task runner and load balancer
|-- server               # Nginx container in front of all the microservice APIs allowing to avoid CORS issues on the UI side
|-- ui                   # user interface
```

Each microservice has its own README file containing more information about its usage and functionality.

# Documentation

The user and developer documentation of Ylem is available at https://docs.ylem.co/.

The [open-source section](https://docs.ylem.co/open-source-edition) contains information about the [task-processing architecture](https://docs.ylem.co/open-source-edition/task-processing-architecture) and [configuration of integrations](https://docs.ylem.co/open-source-edition/configuring-integrations-with-.env-variables) using .env files and parameters.

# Explore our additional integration packages

| Integration   | Repository | Description |
| ------------- | ---------- | ----------- |
| <img width="100px" alt="Apache Kafka" title="Apache Kafka" src="https://docs.ylem.co/~gitbook/image?url=https%3A%2F%2F3180830455-files.gitbook.io%2F%7E%2Ffiles%2Fv0%2Fb%2Fgitbook-x-prod.appspot.com%2Fo%2Fspaces%252FD0FT8l3QzMrw546vOdHU%252Fuploads%252FFGLlkFHxjsJMN4bfHJFR%252Fkf.png%3Falt%3Dmedia%26token%3D47a37fc7-13df-410e-9760-85525a424fde&width=245&dpr=2&quality=100&sign=b464460b&sv=1">  | https://github.com/ylem-co/ylem-kafka-trigger | Containerized Apache Kafka listener to stream data to Ylem |
| <img width="100px" alt="RabbitMQ" title="RabbitMQ" src="https://docs.ylem.co/~gitbook/image?url=https%3A%2F%2F3180830455-files.gitbook.io%2F%7E%2Ffiles%2Fv0%2Fb%2Fgitbook-x-prod.appspot.com%2Fo%2Fspaces%252FD0FT8l3QzMrw546vOdHU%252Fuploads%252FOG2ufh5uxAi54ht1kIBm%252Frabbitmq_logo_icon_170812.png%3Falt%3Dmedia%26token%3D01fb235c-59a3-4a3d-83e3-f6a256f99dd6&width=245&dpr=2&quality=100&sign=d1eda75a&sv=1">      | https://github.com/ylem-co/ylem-rabbitmq-consumer | Containerized RabbitMQ consumer to stream data to Ylem |
| <img width="100px" alt="AWS S3" title="AWS S3" src="https://docs.ylem.co/~gitbook/image?url=https%3A%2F%2F3180830455-files.gitbook.io%2F%7E%2Ffiles%2Fv0%2Fb%2Fgitbook-x-prod.appspot.com%2Fo%2Fspaces%252FD0FT8l3QzMrw546vOdHU%252Fuploads%252F9ltBsXtDuKbyl1YlaX1M%252Fs3.png%3Falt%3Dmedia%26token%3Daec101c7-f899-4c6b-ad91-7b230b694a63&width=245&dpr=2&quality=100&sign=b3b158ab&sv=1"> | https://github.com/ylem-co/s3-lambda-trigger  | AWS Lambda function to stream data from AWS S3 to Ylem |
| <img width="100px" alt="Tableau" title="Tableau" src="https://docs.ylem.co/~gitbook/image?url=https%3A%2F%2F3180830455-files.gitbook.io%2F%7E%2Ffiles%2Fv0%2Fb%2Fgitbook-x-prod.appspot.com%2Fo%2Fspaces%252FD0FT8l3QzMrw546vOdHU%252Fuploads%252Fr5M0FiIfvC7on34tLYNN%252FTB.png%3Falt%3Dmedia%26token%3D3d29f3ad-2df5-454f-ae2c-c1cc4a1f191c&width=245&dpr=2&quality=100&sign=8ee0d519&sv=1">       | https://github.com/ylem-co/tableau-http-wrapper | Containerized HTTP wrapper to stream data from Ylem to Tableau |

# Key contributors

* [olschaefer](https://github.com/olschaefer)
* [schneekatze](https://github.com/schneekatze)
* [lunoshot](https://github.com/lunoshot)
* [Ardem](https://github.com/Ardem)
