import React, { Component } from 'react';
import { DateTime } from "luxon";
import { Link } from 'react-router-dom';
import prettyMilliseconds from 'pretty-ms';

import Accordion from 'react-bootstrap/Accordion';
import Spinner from "react-bootstrap/Spinner";
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Dropdown from "react-bootstrap/Dropdown";

import Tooltip from '@mui/material/Tooltip';
import Circle from '@mui/icons-material/Circle';

import Datepicker from "../datepicker.component";

import TasksService, {TASKS} from "../../services/task.service";
import StatService from "../../services/stat.service";
import PipelineService, {PIPELINE_TYPE_GENERIC} from "../../services/pipeline.service";

import {decodeOutput} from "./pipelineRun.component";

class PipelineSlowTaskRuns extends Component {
  constructor(props) {
      super(props);
      this.handleGetTask = this.handleGetTask.bind(this);
      this.handleGetPipeline = this.handleGetPipeline.bind(this);
      this.handleGetLogs = this.handleGetLogs.bind(this);
      this.handlePrepareLogs = this.handlePrepareLogs.bind(this);

      this.state = {
        organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
        item: this.props.item,
        dateFrom: this.props.dateFrom,
        dateTo: this.props.dateTo,
        threshold: this.props.threshold,
        taskType: this.props.taskType,
        tasks: [],
        pipelines: [],
        logs: null,
        preparedLogs: null,
      };
  }

  componentDidMount = async() => {
      await this.handleGetLogs(
        this.props.dateFrom,
        this.props.dateTo,
        this.props.threshold,
        this.props.taskType
      );
  };

  promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

  UNSAFE_componentWillReceiveProps(props) {
      this.setState({
          dateFrom: props.dateFrom,
          dateTo: props.dateTo,
          threshold: props.threshold,
          taskType: props.taskType,
      });

      this.handleGetLogs(
        props.dateFrom,
        props.dateTo,
        props.threshold,
        props.taskType
      );
  }

  handleGetTask = async(uuid, pipelineUuid) => {
      let tasks = this.state.tasks;

      let task = TasksService.getTask(uuid, pipelineUuid);

      let result = await Promise.resolve(task)
          .then(async(task) => {
              if (task.data && task.data !== null) {
                  var item = task.data;

                  tasks.push(item)
                        
                  await this.promisedSetState({tasks});

                  return item;
              } else {
                  return null;
              }
          })
          .catch((error) => {
              return null;
          });

      return result;
  };

  handleGetPipeline = async(uuid) => {
      let pipelines = this.state.pipelines;

      let pipeline = pipelines.find(o => o.uuid === uuid);

      if (!pipeline) {
        pipeline = PipelineService.getPipeline(uuid);

        let result = await Promise.resolve(pipeline)
            .then(async(pipeline) => {
                if (pipeline.data && pipeline.data !== null) {
                    var item = pipeline.data;

                    pipelines.push(item)
                          
                    await this.promisedSetState({pipelines});

                    return item;
                } else {
                  return null;
                }
            })
            .catch((error) => {
                return null;
            });

        return result;
      } else {
        return pipeline;
      }
  };

  handleGetLogs = async(dateFrom, dateTo, threshold, taskType) => {
    let logs = StatService.getSlowTaskRuns(
      this.state.organization.uuid,
      dateFrom, 
      dateTo, 
      threshold, 
      taskType
    );

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

  handlePrepareLogs = async(rawLogs) => {
    let logs = [];

    for(var i = 0; i < rawLogs.length; i++) {
      var task = logs.find(
        o => o.task_uuid === rawLogs[i].task_uuid
      );

      if (!task) {
        task = {
          "slow_runs_count": 1,
          "succesful_runs": rawLogs[i].is_successful ? 1 : 0,
          "cumulative_duration": rawLogs[i].duration,
          "task_uuid": rawLogs[i].task_uuid,
          "pipeline_uuid": rawLogs[i].pipeline_uuid,
          "task_type": rawLogs[i].task_type,
          "task": false,
          "pipeline": false,
          "task_runs": [{
            "pipeline_run_uuid": rawLogs[i].pipeline_run_uuid,
            "is_successful": rawLogs[i].is_successful,
            "duration": rawLogs[i].duration,
            "executed_at": DateTime.fromISO(rawLogs[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }),
            "output": atob(rawLogs[i].output),
          }]
        };

        logs.push(task);
      } else {
        task.task_runs.push({
          "pipeline_run_uuid": rawLogs[i].pipeline_run_uuid,
          "is_successful": rawLogs[i].is_successful,
          "duration": rawLogs[i].duration,
          "executed_at": DateTime.fromISO(rawLogs[i].executed_at, { zone: 'UTC' }).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }),
          "output": atob(rawLogs[i].output),
        });
        task.slow_runs_count++;
        task.cumulative_duration += rawLogs[i].duration;

        if (rawLogs[i].is_successful) {
          task.succesful_runs++;
        }
      }
    }

    logs.sort((a, b) => b.cumulative_duration - a.cumulative_duration);

    await this.promisedSetState({preparedLogs: logs});

    for(i = 0; i < logs.length; i++) {
      let t = await this.handleGetTask(logs[i].task_uuid, logs[i].pipeline_uuid);
      let w = await this.handleGetPipeline(logs[i].pipeline_uuid);

      logs[i].task = t;
      logs[i].pipeline = w;

      this.setState({preparedLogs: logs});
    }
  }

  render() {
    const { logs, dateFrom, dateTo, threshold, taskType, showDatePicker, preparedLogs } = this.state;

    return (
      <>
        <Row>
          <Col xs={4}>
            <div className="dropdownLabel">Type:</div>
            <Dropdown size="lg" className="mb-4 float-left">
                <Dropdown.Toggle 
                    variant="light" 
                    id="dropdown-basic"
                    className="tasksDropdown"
                >
                    <div className={"dropdownTaskType " + taskType + "Node draggableNode"}>
                        {taskType[0].charAt(0).toUpperCase()}
                    </div>
                    <div className="dropdownTaskTitle float-left">{taskType[0].charAt(0).toUpperCase() + taskType.slice(1)}</div>
                </Dropdown.Toggle>

                <Dropdown.Menu className="tasks">
                    {TASKS.map((value, index) => (
                        <Dropdown.Item
                            value={value}
                            key={"task-" + index}
                            active={this.state.taskType !== null && value === this.state.taskType}
                            onClick={(e) => this.props.changeTaskTypeHandler(value)}
                        >
                            <div className={"dropdownTaskType " + value + "Node draggableNode"}>
                                {value.charAt(0).toUpperCase()}
                            </div>
                            <div className="dropdownTaskTitle">{value[0].charAt(0).toUpperCase() + value.slice(1)}</div>
                        </Dropdown.Item>
                    ))}
                </Dropdown.Menu>
            </Dropdown>
          </Col>
          <Col xs={4}>
            <div className="dropdownLabel">Threshold (milliseconds):</div>
            <input
                className="form-control form-control-lg"
                type="text"
                placeholder="Threshold"
                autoComplete="threshold"
                name="threshold"
                value={threshold}
                onChange={this.props.changeThresholdHandler}
            />
          </Col>
          <Col xs={4}>
            <div>
              <div>
                <div className="dropdownLabel">Date range: </div>
                <Datepicker
                  dateTo={dateTo}
                  dateFrom={dateFrom}
                  showDatePicker={showDatePicker}
                  changeDates={this.props.changeDatesHandler}
                />
              </div>
              <div className="clearfix"></div>
            </div>
          </Col>
        </Row>
        {logs !== null 
          && preparedLogs !== null
          ?
            preparedLogs.length === 0
            ? 
              <div className="text-center">
                No slow tasks with the selected conditions
              </div>
            : 
              <div>
                <h1 className="px-3">Slow tasks: {preparedLogs.length}</h1><br/><br/>
                <Accordion>
                {
                  preparedLogs.map((value, index) => (
                    <Accordion.Item eventKey={"pipeline_run-" + index}>
                      <Accordion.Header>
                        <Row className="w-98 logTable">
                          <Col xs={5}>
                            {value.task_type &&
                              <div className={"dropdownTaskType " + value.task_type + "Node draggableNode"}>
                                {value.task_type[0].charAt(0).toUpperCase()}
                              </div>
                            }
                            <div>
                              {
                                value.task !== false
                                ?
                                  (value.task !== null
                                  ?
                                    <div>
                                      <span className="statValue">{value.task.name}</span><br/>
                                      <span className="statKey">UUID: {value.task_uuid}</span>
                                    </div>
                                  :
                                    <div>
                                      <span className="statValue">Undefined</span><br/>
                                      <span className="statKey">UUID: {value.task_uuid}</span>
                                    </div>
                                  )
                                : <div className="text-center"><Spinner animation="border" className="spinner-primary"/></div>
                              }
                            </div>
                          </Col>
                          <Col xs={3}>
                            <span className="statValue">
                              {
                                value.pipeline !== false
                                ?
                                (
                                  value.pipeline !== null
                                  ?
                                    <div>
                                      <Link to={
                                        '/' + (value.pipeline.type === PIPELINE_TYPE_GENERIC ? "Pipelines" : "metrics")
                                        + '/folder/'
                                        + (value.pipeline.folder_uuid || 'root')
                                        + '/details/' + value.pipeline.uuid
                                      }>
                                        {value.pipeline.name}
                                      </Link><br/>
                                      <span className="statKey">Pipeline</span>
                                    </div>
                                  :
                                    <div>
                                      <span className="statValue">Undefined</span><br/>
                                      <span className="statKey">Pipeline</span>
                                    </div>
                                )
                                : <div className="text-center"><Spinner animation="border" className="spinner-primary"/></div>
                              }
                            </span>
                          </Col>
                          <Col xs={2}>
                            <span className="statValue">{value.slow_runs_count + '/' + value.succesful_runs}</span><br/>
                            <span className="statKey">Runs/Successful</span>
                          </Col>
                          <Col xs={2}>
                            <span className="statValue">
                              {
                                value.cumulative_duration > 0
                                  ? prettyMilliseconds(value.cumulative_duration, {verbose: true})
                                  : "< 1 millisecond"
                              }
                            </span><br/>
                            <span className="statKey">Cumulative duration</span>
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
                                  <Col xs={7}>
                                    <span className="statValue">{run.executed_at}</span><br/>
                                    <span className="statKey">Executed at</span>
                                  </Col>
                                  <Col xs={4}>
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
                                    {decodeOutput(run.output, false, false)}
                                  </code>
                                }
                              </Accordion.Body>
                            </Accordion.Item>
                          ))
                        }
                        </Accordion>
                      </Accordion.Body>
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

export default PipelineSlowTaskRuns;
