import React, { Component } from 'react';
import { DateTime } from "luxon";
import prettyMilliseconds from 'pretty-ms';

import Accordion from 'react-bootstrap/Accordion';
import Spinner from "react-bootstrap/Spinner";
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import Tooltip from '@mui/material/Tooltip';
import Circle from '@mui/icons-material/Circle';

import Datepicker from "../datepicker.component";

import TasksService, {TASK_TYPE_AGGREGATOR} from "../../services/task.service";
import StatService from "../../services/stat.service";
import {PIPELINE_TYPE_GENERIC} from "../../services/pipeline.service";

import {decodeOutput, isJson} from "./pipelineRun.component";

class PipelineLogs extends Component {
  constructor(props) {
      super(props);
      this.handleGetTasks = this.handleGetTasks.bind(this);
      this.handleGetLogs = this.handleGetLogs.bind(this);
      this.handlePrepareLogs = this.handlePrepareLogs.bind(this);
      this.getTaskNameByUuid = this.getTaskNameByUuid.bind(this);

      this.state = {
        organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
        item: this.props.item,
        dateFrom: this.props.dateFrom,
        dateTo: this.props.dateTo,
        tasks: null,
        logs: null,
        preparedLogs: null,
      };
  }

  componentDidMount = async() => {
      await this.handleGetTasks(this.props.item.uuid);
      await this.handleGetLogs(
        this.props.item.uuid,
        this.props.dateFrom,
        this.props.dateTo
      );
  };

  UNSAFE_componentWillReceiveProps(props) {
      this.setState({
          item: props.item,
          dateFrom: props.dateFrom,
          dateTo: props.dateTo,
      });

      this.handleGetLogs(
        props.item.uuid,
        props.dateFrom,
        props.dateTo
      );
  }

  promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

  handleGetTasks = async(uuid) => {
      let tasks = this.state.tasks;

      if (
          tasks === null
          || tasks.length === 0
      ) {
          tasks = TasksService.getTasksByPipeline(uuid);

          await Promise.resolve(tasks)
              .then(async(tasks) => {
                  if (tasks.data && tasks.data.items !== null) {
                      var items = tasks.data.items;
                        
                      await this.promisedSetState({tasks: items});
                  } else {
                      this.setState({tasks: []});
                  }
              });
        } 
    };

  getTaskNameByUuid = (uuid) => {
    let task = this.state.tasks.find(o => o.uuid === uuid);

    return task ? task.name : "Undefined";
  };

  handleGetLogs = async(uuid, dateFrom, dateTo) => {
    let logs = StatService.getLastRunsLog(uuid, dateFrom, dateTo);

    await Promise.resolve(logs)
        .then(async(logs) => {
            if (logs.data) {
                var items = logs.data;       
                this.setState({logs: items});
                this.handlePrepareLogs(items);
            } else {
                this.setState({logs: [], preparedLogs: []});
            }
        })
        .catch(() => {
            this.setState({logs: []});
        });
    };

  handlePrepareLogs = (rawLogs) => {
    let logs = [];

    for(var i = 0; i < rawLogs.length; i++) {
      var pipelineRun = logs.find(
        o => o.pipeline_run_uuid === rawLogs[i].pipeline_run_uuid
      );

      var output = atob(rawLogs[i].output);

      if(isJson(output)) {
        output = decodeOutput(output, false, false);
        //output = JSON.stringify(JSON.parse(output), null, 4);
      }

      if (!pipelineRun) {
        logs.push({
          "pipeline_run_uuid": rawLogs[i].pipeline_run_uuid,
          "task_count": 1,
          "succesful_task_runs": rawLogs[i].is_successful ? 1 : 0,
          "duration": rawLogs[i].duration,
          "executed_at": DateTime.fromISO(rawLogs[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', millisecond: '3-digit' }),
          "is_successful": rawLogs[i].is_successful,
          "value": rawLogs[i].metric_value,
          "task_runs": [{
            "is_successful": rawLogs[i].is_successful,
            "duration": rawLogs[i].duration,
            "task_uuid": rawLogs[i].task_uuid,
            "task_type": rawLogs[i].task_type,
            "task_name": this.getTaskNameByUuid(rawLogs[i].task_uuid),
            "executed_at": DateTime.fromISO(rawLogs[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', millisecond: '3-digit' }),
            "output": output,
          }]
        });
      } else {
        pipelineRun.task_runs.push({
          "is_successful": rawLogs[i].is_successful,
          "duration": rawLogs[i].duration,
          "task_uuid": rawLogs[i].task_uuid,
          "task_type": rawLogs[i].task_type,
          "task_name": this.getTaskNameByUuid(rawLogs[i].task_uuid),
          "executed_at": DateTime.fromISO(rawLogs[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', millisecond: '3-digit' }),
          "output": output,
        });
        pipelineRun.task_count++;
        pipelineRun.duration += rawLogs[i].duration;

        if (
          this.state.item.type !== PIPELINE_TYPE_GENERIC
          && rawLogs[i].task_type === TASK_TYPE_AGGREGATOR
        ) {
          pipelineRun.is_successful = rawLogs[i].is_successful;
          pipelineRun.value = rawLogs[i].metric_value;
        }

        if (rawLogs[i].is_successful) {
          pipelineRun.succesful_task_runs++;
        }
      }
    }

    this.setState({preparedLogs: logs});
  }

  render() {
    const { item } = this.props;

    const { dateFrom, dateTo, showDatePicker, tasks, preparedLogs } = this.state;


    return (
      <>
        <div className="text-right">
          <div className="float-right">
            <Datepicker
              dateTo={dateTo}
              dateFrom={dateFrom}
              showDatePicker={showDatePicker}
              changeDates={this.props.changeDatesHandler}
            />
          </div>
          <div className="dropdownLabel">Date range: </div>
          <div className="clearfix"></div>
        </div>
        {tasks !== null 
          && preparedLogs !== null
          ?
            preparedLogs.length === 0
            ? 
              <div className="text-center">
                No run logs are available
              </div>
            : 
              <div>
                <h1 className="px-3">Runs: {preparedLogs.length}</h1><br/><br/>
                <Accordion defaultActiveKey="pipeline_run-0">
                {
                  item.type === PIPELINE_TYPE_GENERIC
                  ?
                    preparedLogs.map((value, index) => (
                    <Accordion.Item eventKey={"pipeline_run-" + index}>
                      <Accordion.Header>
                        <Row className="w-98 logTable">
                          <Col xs={1} className="pt-2">
                            <Tooltip title={value.task_count === value.succesful_task_runs ? "Success" : "Failure"} placement="right">
                              <Circle 
                                className={value.task_count === value.succesful_task_runs ? "icon_online" : "icon_offline"}
                                alt={value.task_count === value.succesful_task_runs ? "Success" : "Failure"}
                              />
                            </Tooltip>
                          </Col>
                          <Col xs={4}>
                            <span className="statValue">{value.executed_at}</span><br/>
                            <span className="statKey">Executed at</span>
                          </Col>
                          <Col xs={2}>
                            <span className="statValue">{value.task_count}</span><br/>
                            <span className="statKey">Tasks</span>
                          </Col>
                          <Col xs={2}>
                            <span className="statValue">{value.succesful_task_runs}</span><br/>
                            <span className="statKey">Successful tasks</span>
                          </Col>
                          <Col xs={3}>
                            <span className="statValue">
                              {
                                value.duration > 0
                                  ? prettyMilliseconds(value.duration, {verbose: true})
                                  : "< 1 millisecond"
                              }
                            </span><br/>
                            <span className="statKey">Duration</span>
                          </Col>
                        </Row>
                      </Accordion.Header>
                      <Accordion.Body>
                        <Accordion>
                        {
                          value.task_runs.map((run, index) => (
                            <Accordion.Item eventKey={"task_run-" + index}>
                              <Accordion.Header>
                                <Row className="w-98 logTable">
                                  <Col xs={1} className="pt-2">
                                    <Tooltip title={run.is_successful ? "Success" : "Failure"} placement="right">
                                      <Circle 
                                        className={run.is_successful ? "icon_online" : "icon_offline"}
                                        alt={run.is_successful ? "Success" : "Failure"}
                                      />
                                    </Tooltip>
                                  </Col>
                                  <Col xs={5}>
                                    {run.task_type &&
                                      <div className={"dropdownTaskType " + run.task_type + "Node draggableNode"}>
                                        {run.task_type[0].charAt(0).toUpperCase()}
                                      </div>
                                    }
                                    <div>
                                      <span className="statValue">{run.task_name}</span><br/>
                                      <span className="statKey">UUID: {run.task_uuid}</span>
                                    </div>
                                  </Col>
                                  <Col xs={3}>
                                    <span className="statValue">{run.executed_at}</span><br/>
                                    <span className="statKey">Executed at</span>
                                  </Col>
                                  <Col xs={3}>
                                    <span className="statValue">
                                      {
                                        run.duration > 0
                                          ? prettyMilliseconds(run.duration, {verbose: true})
                                          : "< 1 millisecond"
                                      }
                                    </span><br/>
                                    <span className="statKey">Duration</span>
                                  </Col>
                                </Row>
                              </Accordion.Header>
                              <Accordion.Body>
                                {
                                  <code className="taskOutputCode">
                                      { run.output }
                                  </code>
                                }
                              </Accordion.Body>
                            </Accordion.Item>
                          ))
                        }
                        </Accordion>
                      </Accordion.Body>
                    </Accordion.Item>
                  ))
                  :
                  preparedLogs.map((value, index) => (
                    <Accordion.Item eventKey={"pipeline_run-" + index} className="noExpandAccordion">
                      <Accordion.Header>
                        <Row className="w-98 logTable">
                          <Col xs={1} className="pt-2">
                            <Tooltip title={value.is_successful ? "Success" : "Failure"} placement="right">
                              <Circle 
                                className={value.is_successful ? "icon_online" : "icon_offline"}
                                alt={value.is_successful ? "Success" : "Failure"}
                              />
                            </Tooltip>
                          </Col>
                          <Col xs={4}>
                            <span className="statValue">{value.executed_at}</span><br/>
                            <span className="statKey">Executed at</span>
                          </Col>
                          <Col xs={4}>
                            <span className="statValue">
                              {
                                value.duration > 0
                                  ? prettyMilliseconds(value.duration, {verbose: true})
                                  : "< 1 millisecond"
                              }
                            </span><br/>
                            <span className="statKey">Duration</span>
                          </Col>
                          <Col xs={3}>
                            <span className="statValue">{value.value || "-"}</span><br/>
                            <span className="statKey">Value</span>
                          </Col>
                        </Row>
                      </Accordion.Header>
                    </Accordion.Item>
                  ))}
                </Accordion>
              </div>
          : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
        }
      </>
    );
  }
};

export default PipelineLogs;
