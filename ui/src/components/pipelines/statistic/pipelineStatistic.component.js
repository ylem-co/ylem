import React from 'react';

import PipelineStatisticSummary from "./pipelineStatisticSummary.component";
import Datepicker from "../../datepicker.component";

import Spinner from "react-bootstrap/Spinner";
import Dropdown from "react-bootstrap/Dropdown";

import TasksService from "../../../services/task.service";
import { PIPELINE_TYPE_GENERIC, PIPELINE_TYPE_METRIC } from "../../../services/pipeline.service";

class PipelineStatistic extends React.Component {
    constructor(props) {
        super(props);
        this.handleGetTasks = this.handleGetTasks.bind(this);
        this.onChangeActiveTask = this.onChangeActiveTask.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            item: this.props.item,
            dateFrom: this.props.dateFrom,
            dateTo: this.props.dateTo,
            type: this.props.type || PIPELINE_TYPE_GENERIC,
            tasks: null,
            activeTask: null,
        };
    }

    componentDidMount() {
        this.handleGetTasks(this.props.item.uuid);
    };

    UNSAFE_componentWillReceiveProps(props) {
        this.setState({
            item: props.item,
            dateFrom: props.dateFrom,
            dateTo: props.dateTo,
        });
    }

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
                        
                        this.setState({tasks: items});
                        if (items.length > 0) {
                            this.setState({activeTask: items[0]});
                        }
                    } else {
                        this.setState({tasks: []});
                    }
                });
        }
    };

    onChangeActiveTask = async(task) => {
        this.setState({activeTask: task})
    };

    render() {
        const { item } = this.props;

        const { dateFrom, dateTo, showDatePicker, activeTask } = this.state;

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
                <PipelineStatisticSummary 
                    item={item} 
                    type="pipeline"
                    dateFrom={dateFrom}
                    dateTo={dateTo}
                    changeDatesHandler={this.props.changeDatesHandler}
                />

                { this.props.type !== PIPELINE_TYPE_METRIC &&
                <>
                <h2 className="float-left px-3 tasksDropdownPreTitle">Tasks:</h2> 
                {this.state.tasks !== null ?
                        <Dropdown size="lg" className="mb-4 float-left">
                            <Dropdown.Toggle 
                                variant="light" 
                                id="dropdown-basic"
                                className="tasksDropdown"
                            >
                                { activeTask !== null &&
                                    <>
                                        <div className={"dropdownTaskType " + activeTask.type + "Node draggableNode"}>
                                            {activeTask.type[0].charAt(0).toUpperCase()}
                                        </div>
                                        <div className="dropdownTaskTitle float-left">{activeTask.name}</div>
                                    </>
                                }
                            </Dropdown.Toggle>

                            <Dropdown.Menu className="tasks">
                                {this.state.tasks.map(value => (
                                    <Dropdown.Item
                                        value={value.name}
                                        key={value.uuid}
                                        active={this.state.activeTask !== null && value.uuid === this.state.activeTask.uuid}
                                        onClick={(e) => this.onChangeActiveTask(value)}
                                    >
                                        <div className={"dropdownTaskType " + value.type + "Node draggableNode"}>
                                            {value.type[0].charAt(0).toUpperCase()}
                                        </div>
                                        <div className="dropdownTaskTitle">{value.name}</div>
                                    </Dropdown.Item>
                                ))}
                            </Dropdown.Menu>
                        </Dropdown>
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }
                <div className="clearfix"></div>
                { activeTask !== null &&
                    <PipelineStatisticSummary 
                        item={activeTask} 
                        type="task"
                        dateFrom={dateFrom}
                        dateTo={dateTo}
                        changeDatesHandler={this.props.changeDatesHandler}
                    />
                }
                </>
                }
            </>
        );
    }
}

export default PipelineStatistic;
