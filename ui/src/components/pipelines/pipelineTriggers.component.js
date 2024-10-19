import React from 'react';

import ContentCopy from '@mui/icons-material/ContentCopy';
import Tooltip from '@mui/material/Tooltip';

import Card from 'react-bootstrap/Card';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from "react-bootstrap/Button";

import { PIPELINE_PAGE_TRIGGERS } from "../../services/pipeline.service";

function copyLink(copiedLink){
  navigator.clipboard.writeText(copiedLink);
}

export const PipelineTriggers = ({ item, openPipelineForm, openScheduleForm }) => {
  return (
    <div className="pb-4">
      <Card className="onboardingCard withHeader noBorder mb-5">
        <Card.Header>
          Schedule
        </Card.Header>
        <Card.Body className="p-4">
          {
            item.schedule !== ""
            ?
            <div className="mb-3">
              Current pipeline schedule:
              <div className="code">
                {item.schedule}
              </div>
            </div>
            :
            <div className="mb-3">
              This pipeline is not scheduled yet.
            </div>
          }

          <Button
            variant="primary"
            className="mx-0"
            onClick={() => openScheduleForm(item, PIPELINE_PAGE_TRIGGERS)}
          >
            Edit pipeline schedule
          </Button>                           
        </Card.Body>
      </Card>

      <Card className="onboardingCard withHeader noBorder mb-5">
        <Card.Header>
          Apache Kafka / RabbitMQ / Google Cloud PubSub / AWS S3
        </Card.Header>
        <Card.Body className="p-4">
          <h4>Apache Kafka</h4>
          This pipeline can listen to updates in <strong>Apache Kafka</strong> topics and run automatically in real time. To do so, add an <a href="https://docs.ylem.co/workflows/tasks-ip#external-trigger" target="_blank" rel="noreferrer">External trigger</a> task to the beginning of the pipeline and install our open-source <a href="https://github.com/ylem-co/ylem-kafka-trigger" target="_blank" rel="noreferrer">kafka-trigger</a> library. More information is in our <a href="https://docs.ylem.co/datamin-api/api-endpoints" target="_blank" rel="noreferrer">documentation</a>.
          <br/><br/>
          <h4>RabbitMQ</h4>
          The same way data can be transferred to Ylem from RabbitMQ. We even provide you an open-source docker container with <a href="https://github.com/ylem-co/ylem-rabbitmq-consumer" target="_blank" rel="noreferrer">RabbitMQ consumer</a> written on Golang here.
          <br/><br/>
          <h4>Google Cloud PubSub</h4>
          The integration with Google Cloud PubSub is even easier. To trigger this pipeline just call from there the API Endpoint you see below.
          <br/><br/>
          <h4>AWS S3</h4>
          Coming soon.
          <br/><br/>
          <Button
            variant="primary"
            className="mx-0"
            onClick={() => openPipelineForm(item, PIPELINE_PAGE_TRIGGERS)}
          >
            Edit pipeline
          </Button>                            
        </Card.Body>
      </Card>

      <Card className="onboardingCard withHeader noBorder mb-5">
        <Card.Header>
          API
        </Card.Header>
        <Card.Body className="p-4">
          This pipeline can be triggered through Ylem API using the following endpoint:

          <div className="code">
            <Row>
              <Col xs={10}>
                {
                  "POST https://api.datamin.io/v1/pipelines/" + item.uuid + "/runs/" 
                }
              </Col>
              <Col xs={2} className="text-right">
                <Tooltip title="Click to copy link to clipboard" placement="left">
                  <ContentCopy
                    onClick={() => copyLink("https://api.datamin.io/v1/pipelines/" + item.uuid + "/runs/")}
                  />
                </Tooltip>
              </Col>
            </Row>
          </div>

          More information about Ylem API is in our <a href="https://docs.ylem.co/datamin-api/api-endpoints" target="_blank" rel="noreferrer">documentation</a>.                      
        </Card.Body>
      </Card>

      <Card className="onboardingCard withHeader noBorder mb-5">
        <Card.Header>
          Manual on-demand trigger
        </Card.Header>
        <Card.Body className="p-4">
          Every pipeline can be run manually for on-demand and debugging purposes from the pipeline editing layout.<br/><br/>
          <Button
            variant="primary"
            className="mx-0"
            onClick={() => openPipelineForm(item, PIPELINE_PAGE_TRIGGERS)}
          >
            Edit pipeline
          </Button>                            
        </Card.Body>
      </Card>
    </div>
  );
};
