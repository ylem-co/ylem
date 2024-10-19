import React from 'react';

import Spinner from "react-bootstrap/Spinner";

import PipelineService, { PIPELINE_TYPE_METRIC } from "../../services/pipeline.service";

const RUN_STATE_PENDING = "pending"
const RUN_STATE_EXECUTED = "executed"

const MAX_ATTEMPTS = 30;

export function decodeOutput(decodedOutput, trimmedVersion = false, usePrefix = true) {
    let output = "";
    let prefix = usePrefix === true ? 'And returned:\n' : '';

    if (isJson(decodedOutput)) {
        decodedOutput = JSON.parse(decodedOutput);
        if (Array.isArray(decodedOutput)) {
            if (decodedOutput.length > 0) {
                if (trimmedVersion === true) {
                    const slicedArray = decodedOutput.slice(0, 10);
                    output = output + prefix + JSON.stringify(slicedArray, null, 4) + '\n';
                        
                    if (decodedOutput.length > 10) {
                        let cnt = decodedOutput.length - 10;
                        output = output + 'And ' + cnt + ' other rows. See the full log\n';
                    }
                } else {
                    output = output + prefix + JSON.stringify(decodedOutput, null, 4) + '\n';
                }
            } else {
                output = output + 'And returned 0 rows\n';
            }
        } else if (decodedOutput === null || decodedOutput === 'null') {
            output = output + 'And returned: null';
        } else if (decodedOutput.hasOwnProperty('result')) {
            output = output + prefix + decodedOutput.result + '\n';
        } else {
            output = output + prefix + JSON.stringify(decodedOutput, null, 4) + '\n';
        }
    } else {
        output = output + prefix + decodedOutput + '\n';
    }

    return output;
}

export function isJson(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

export function handleErrors(errors) {
    var output = "";
    if (errors.length > 0) {
        for(var i = 0; i < errors.length; i++) {
            output = output + errors[i].message + "\n";
        }
    }

    return output;
};

class PipelineRun extends React.Component {
    constructor(props) {
        super(props);
        this.handleRunningOutput = this.handleRunningOutput.bind(this);
        this.handlePipelineRun = this.handlePipelineRun.bind(this);
        //this.handleErrors = this.handleErrors.bind(this);
        this.handlePipelineRunResults = this.handlePipelineRunResults.bind(this);
        this.getTaskNameByUuid = this.getTaskNameByUuid.bind(this);
        this.getType = this.getType.bind(this);
        //this.decodeOutput = this.decodeOutput.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            item: this.props.item,
            elements: this.props.elements,
            run: null,
            finished: null,
            output: this.getType(this.props.item) + " started...\n",
            fixedOutput: null,
            attempts: 0,
        };
    }

    componentDidMount = async() => {
        if (this.props.showLogWithoutRunning !== true) {
            this.mounted = true;
            this.handlePipelineRun(this.props.item.uuid);
        } else {
            await this.promisedSetState({
                output: this.props.oldOutput
            });
        }
    }

    componentWillUnmount() {
        this.mounted = false;
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    getType(item) {
        return item.type === PIPELINE_TYPE_METRIC ? "Metric" : "Pipeline";
    }

    getTaskNameByUuid = (uuid) => {
        var element = this.state.elements.find(o => o.id === uuid);

        if (element) {
            return element.data.name;
        }

        return "";
    };

    handlePipelineRun = async(uuid) => {
        var run = PipelineService.runPipeline(uuid);

        await Promise.resolve(run)
            .then(async(run) => {
                var output = this.state.output;
                if (run.data) {
                   await this.promisedSetState({
                        run: run.data,
                        output: output + this.getType(this.state.item) + " is running...\n"
                   });
                   this.handlePipelineRunResults(uuid)
                } else {
                    this.setState({
                        finished: true,
                        output: output + "\nSomething went wrong. Please make sure " + this.getType(this.state.item) + " is correct and try again.\n"
                    });
                    this.props.outputHandler([], true, output, false, true);
                }
            })
            .catch((error) => {
                var output = this.state.output;
                output = output + "\nSomething went wrong. Please make sure you have correct access to it, and the " + this.getType(this.state.item) + " is correct, and try again.\n\nIf you are not an administrator of the organization, please ask an administrator for a help";
                this.setState({
                    finished: true,
                    output,
                });
                this.props.outputHandler([], true, output, false, true);
            });
        ;
    };

    handlePipelineRunResults = async(uuid) => {
        var results = PipelineService.getPipelineRunResults(uuid);

        await Promise.resolve(results)
            .then(async(results) => {
                var output = this.state.output;
                if (results.data) {
                    if (this.state.fixedOutput === null) {
                        await this.promisedSetState({
                            fixedOutput: output
                        });
                    }
                    this.handleRunningOutput(results.data);
                } else {
                    this.setState({
                        finished: true,
                        output: output + "\nSomething went wrong. Please double-check your " + this.getType(this.state.item) +  + " and try again...\n"
                   });
                }
            })
            .catch((error) => {
                var output = this.state.output;
                this.setState({
                    finished: true,
                    output: output + "\nSomething went wrong. Please double-check your " + this.getType(this.state.item) + " and try again.\n"
                });
            });
        ;
    };

    handleRunningOutput = async(response) => {
        var taskRuns = response.results;
        var runUuid = response.pipeline_run_uuid;
        var output = this.state.fixedOutput;
        var tasksInProgress = 0;
        var tasksFailed = 0;
        var tasksChecked = 0;

        /*taskRuns = [
            {
                "id": 10537,
                "state": "executed",
                "task_uuid": "ee7b60b5-505c-45d5-92d3-6f201b94bc88",
                "task_run_uuid": "00000000-0000-0000-0000-000000000000",
                "pipeline_run_uuid": "3e071a5f-ffee-4100-8047-e43b6663c473",
                "is_successful": true,
                "output": "W3siYW1vdW50IjoxMC40NCwiY3JlYXRlZF9hdCI6IjIwMjItMDItMDUgMjI6NDY6MDgiLCJpZCI6NSwib3JkZXJfaWQiOjQsInN0YXR1cyI6ImZhaWxlZCJ9XQ==",
                "errors": [],
                "created_at": "",
                "updated_at": ""
            },
            {
                "id": 10538,
                "state": "executed",
                "task_uuid": "d8466c05-8a1e-47e9-b6b2-ae8b1dceda05",
                "task_run_uuid": "00000000-0000-0000-0000-000000000000",
                "pipeline_run_uuid": "3e071a5f-ffee-4100-8047-e43b6663c473",
                "is_successful": true,
                "output": "eyJyZXN1bHQiOnRydWUsIm9yaWdpbmFsX2lucHV0IjoiVzNzaVlXMXZkVzUwSWpveE1DNDBOQ3dpWTNKbFlYUmxaRjloZENJNklqSXdNakl0TURJdE1EVWdNakk2TkRZNk1EZ2lMQ0pwWkNJNk5Td2liM0prWlhKZmFXUWlPalFzSW5OMFlYUjFjeUk2SW1aaGFXeGxaQ0o5WFE9PSJ9",
                "errors": [],
                "created_at": "",
                "updated_at": ""
            },
            {
                "id": 10540,
                "state": "executed",
                "task_uuid": "fc29e1d5-9089-4d2d-a423-a7b03fef20de",
                "task_run_uuid": "00000000-0000-0000-0000-000000000000",
                "pipeline_run_uuid": "3e071a5f-ffee-4100-8047-e43b6663c473",
                "is_successful": false,
                "output": "",
                "errors": [
                    {
                        "id": 987,
                        "task_run_result_id": 10540,
                        "code": 10100,
                        "severity": "error",
                        "message": "invalid_auth"
                    }
                ],
                "created_at": "",
                "updated_at": ""
            }
        ];*/

        for(var i = 0; i < taskRuns.length; i++) {
            if (taskRuns[i].pipeline_run_uuid === runUuid) {
                var taskName = this.getTaskNameByUuid(taskRuns[i].task_uuid);
                if (taskRuns[i].state === RUN_STATE_PENDING) {
                    output = output + 'Task ' + taskName + ' is in progress\n';
                    tasksInProgress++;
                } else if (taskRuns[i].state === RUN_STATE_EXECUTED) {
                    if (taskRuns[i].is_successful === true) {
                        output = output + 'Task ' + taskName + ' successfully executed\n';
                        if (taskRuns[i].output !== "") {
                            let decodedOutput = atob(taskRuns[i].output);
                            output = output + decodeOutput(decodedOutput);
                        }
                    } else {
                        output = output + 'Task ' + taskName + ' failed\n';
                        output = output + "Errors:\n" + handleErrors(taskRuns[i].errors) + "\n";
                        tasksFailed++;
                    }
                }
                output = output + '\n';
                tasksChecked++;
            }
        }

        if (tasksChecked === 0) {
            tasksInProgress++;
        }

        await this.promisedSetState({output});
        if (tasksFailed === 0 && tasksInProgress === 0) {
            if (tasksFailed === 0) {
                output = output + this.getType(this.state.item) + " successfully finished\n";
                this.setState({output});
            }
            this.setState({
                finished: true,
            });
            this.props.outputHandler(taskRuns, true, output);
        } else if (tasksFailed !== 0) {
            output = output + this.getType(this.state.item) + " failed. Please fix errors above and try again.\n";
            this.setState({
                finished: true,
                output
            });
            this.props.outputHandler(taskRuns, true, output);
        } else {
            var attempts = this.state.attempts;
            if (attempts < MAX_ATTEMPTS && this.mounted) {
                this.props.outputHandler(taskRuns, false, this.state.output);
                attempts++;
                await this.promisedSetState({attempts});
                await this.sleep(2000);
                this.handlePipelineRunResults(this.state.item.uuid);
            } else if (attempts >= MAX_ATTEMPTS) {
                output = output + "\nTimeout error. Task execution takes longer than 60 seconds.\nPlease check that your data source processes the query as expected and try again.\n";
                this.setState({
                    finished: true,
                    output
                });
                this.props.outputHandler(taskRuns, true, output, true, true);
            } else {
                output = output + "\nPipeline stopped.";
                this.setState({
                    finished: true,
                    output
                });
                this.props.outputHandler(taskRuns, true, output);
            }
        }
    };

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    };

    onChangeActiveTask = async(task) => {
        this.setState({activeTask: task})
    };

    render() {
        const { finished, output, item } = this.state;

        return (
            <>
            {item.type === PIPELINE_TYPE_METRIC &&
                <div className="CLIImitation">
                    {output}
                    {
                        finished === null
                        && this.props.showLogWithoutRunning !== true
                            && <div className="text-center mt-3"><Spinner animation="grow" className="spinner-secondary"/></div> 
                    }
                </div>
            }
            </>
        );
    }
}

export default PipelineRun;
