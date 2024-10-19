<p align="center">
  <img width="748" title="Ylem. The open-source data streaming platform" alt="Ylem. The open-source data streaming platform" src="https://github.com/user-attachments/assets/385e8209-a4fe-4846-8499-e156ce2d2403">
</p>

<div align="center">

  ![Static Badge](https://img.shields.io/badge/checks-passing-green)
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

# Installation

## 1. Install and configure dependencies

To run Ylem, the following dependencies must be first installed and configured.

### Docker 4

If you don't yet have Docker 4 installed, [install](https://www.docker.com/products/docker-desktop/) it from their official website for your OS.

### Apache Kafka

Ylem uses Apache Kafka to process pipelines and tasks. [Install](https://kafka.apache.org/) it from their official website or skip this step if you already have it installed.

We also recommend that you install the Apache Kafka GUI software to manage its topics and subscriptions.

## 2. Create Kafka topics

To be able to work correctly, Ylem requires the following Apache Kafka topics to be created:

* task_runs
* task_runs_load_balanced
* task_run_results
* query_task_run_results
* notification_task_run_results

## 3. Configure environment variables in .env files

Some particular integrations might require extra steps and using `.env` files. 

More information is in [our documentation](https://docs.ylem.co/open-source-edition/configuring-integrations-with-.env-variables).

## 4. Install from pre-build containers

The best way to install Ylem is to clone the repository https://github.com/ylem-co/ylem-installer and run `docker compose up` or `docker compose up -d` from it. It will install Ylem from the latest version of pre-build containers stored on Docker Hub.

## Or build and install from the source

If you want to compile Ylem from the source, run `docker compose up` or `docker compose up -d` from this repository. It will compile the code and run all the necessary containers.

### To rebuild a particular container

If you want to rebuild a particular container from source locally, run:

``` bash
docker compose build --no-cache %%CONTAINER_NAME%%
```

E.g.

``` bash
docker compose build --no-cache ylem_users
```

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

# Key contributors

* [olschaefer](https://github.com/olschaefer)
* [schneekatze](https://github.com/schneekatze)
* [lunoshot](https://github.com/lunoshot)
* [Ardem](https://github.com/Ardem)
