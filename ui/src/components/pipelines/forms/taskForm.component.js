import React from 'react';

import Spinner from "react-bootstrap/Spinner";

import TaskService, {
    TASK_TYPE_QUERY,
    TASK_TYPE_CONDITION,
    TASK_TYPE_AGGREGATOR,
    TASK_TYPE_TRANSFORMER,
    TASK_TYPE_NOTIFICATION,
    TASK_TYPE_API_CALL,
    TASK_TYPE_FOR_EACH,
    TASK_TYPE_MERGE,
    TASK_TYPE_FILTER,
    TASK_TYPE_EXTERNAL_TRIGGER,
    TASK_TYPE_RUN_PIPELINE, 
    TASK_TYPE_CODE, 
    TASK_TYPE_PYTHON,
    TASK_TYPE_GPT,
    TASK_TYPE_PROCESSOR,
} from "../../../services/task.service";

import ConditionForm from "./conditionForm.component";
import AggregatorForm from "./aggregatorForm.component";
import QueryForm from "./queryForm.component";
import ApiCallForm from "./apiCallForm.component";
import NotificationForm from "./notificationForm.component";
import TransformerForm from "./transformerForm.component";
import ForEachForm from "./forEachForm.component";
import MergeForm from "./mergeForm.component";
import FilterForm from "./filterForm.component";
import ExternalTriggerForm from "./externalTriggerForm.component";
import PipelineRunForm from "./pipelineRunForm.component";
import CodeForm from "./codeForm.component";
import GptForm from "./gptForm.component";
import ProcessorForm from "./processorForm.component";

class TaskForm extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            item: this.props.item,
            pipeline: this.props.pipeline,
            task: null,
            trigger: null,
        };
    }

    componentDidMount = async() => {
        if (this.props.item !== null && this.props.pipeline !== null) {
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                expression: this.props.item.expression || "",
            });

            var task = await TaskService.getTask(this.props.item.id, this.props.pipeline.uuid);
            if (task && task.data) {
                await this.promisedSetState({task: task.data});
            }
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    render() {
        const { task, trigger } = this.state;

        return (
            <>
                {
                    task !== null ?
                    (
                        task.type === TASK_TYPE_CONDITION
                        ? <ConditionForm item={task} successHandler={this.props.successHandler}/>
                        : task.type === TASK_TYPE_AGGREGATOR
                            ? <AggregatorForm item={task} successHandler={this.props.successHandler}/>
                            : task.type === TASK_TYPE_QUERY
                                ? <QueryForm item={task} trigger={trigger} successHandler={this.props.successHandler}/>
                                : task.type === TASK_TYPE_API_CALL
                                    ? <ApiCallForm item={task} successHandler={this.props.successHandler}/>
                                    : task.type === TASK_TYPE_NOTIFICATION
                                        ? <NotificationForm item={task} successHandler={this.props.successHandler}/>
                                        : task.type === TASK_TYPE_TRANSFORMER
                                            ? <TransformerForm item={task} successHandler={this.props.successHandler}/>
                                            : task.type === TASK_TYPE_FOR_EACH
                                                ? <ForEachForm item={task} successHandler={this.props.successHandler}/>
                                                : task.type === TASK_TYPE_MERGE
                                                    ? <MergeForm item={task} successHandler={this.props.successHandler}/>
                                                    : task.type === TASK_TYPE_FILTER
                                                        ? <FilterForm item={task} successHandler={this.props.successHandler}/>
                                                            : task.type === TASK_TYPE_RUN_PIPELINE
                                                                ? <PipelineRunForm item={task} successHandler={this.props.successHandler}/>
                                                                : task.type === TASK_TYPE_EXTERNAL_TRIGGER
                                                                    ? <ExternalTriggerForm item={task} successHandler={this.props.successHandler}/>
                                                                    : (task.type === TASK_TYPE_CODE || task.type === TASK_TYPE_PYTHON)
                                                                        ? <CodeForm item={task} successHandler={this.props.successHandler}/>
                                                                        : task.type === TASK_TYPE_GPT
                                                                            ? <GptForm item={task} successHandler={this.props.successHandler}/>
                                                                            : task.type === TASK_TYPE_PROCESSOR
                                                                                ? <ProcessorForm item={task} successHandler={this.props.successHandler}/>
                                                                                : <div></div>
                    ) 
                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
            </>
        );
    }
}

export default TaskForm;
