import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Cron from "react-js-cron";
import 'react-js-cron/dist/styles.css';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import Spinner from "react-bootstrap/Spinner";
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Tooltip from '@mui/material/Tooltip';
import Edit from '@mui/icons-material/Edit';
import DeleteOutlined from '@mui/icons-material/DeleteOutlined';
import ClearRounded from '@mui/icons-material/ClearRounded';

import CodeEditor from '@uiw/react-textarea-code-editor';
import rehypePrism from 'rehype-prism-plus';

import { connect } from "react-redux";
import { addPipeline, updatePipeline } from "../../actions/pipelines";
import { addTask, updateTask } from "../../actions/tasks";
import { addTaskTrigger } from "../../actions/taskTriggers";

import PipelineRun from '../pipelines/pipelineRun.component';
import RightModal from "../modals/rightModal.component";
import ConfirmationModal from "../modals/confirmationModal.component";

import Input from "../formControls/input.component";
import QueryUI from "../formControls/queryUI.component";
import { TextareaEditor } from "../formControls/textareaEditor.component";
import { required, isCron } from "../formControls/validations";

import { clearMessage } from "../../actions/message";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {
    Button,
    Col,
    InputGroup,
    Card,
    Dropdown,
    Row
} from 'react-bootstrap'

import IntegrationService , { INTEGRATION_TYPE_SQL } from "../../services/integration.service";

import TasksService, {
    TASK_TYPE_QUERY,
    TASK_TYPE_CONDITION,
    TASK_TYPE_AGGREGATOR,
    TASK_TYPE_RUN_PIPELINE,
    TASK_SEVERITY_MEDIUM,
} from "../../services/task.service";

import TasksTriggerService, {
    TRIGGER_TYPE_OUTPUT,
    TRIGGER_TYPE_CONDITION_TRUE,
} from "../../services/taskTrigger.service";

import PipelineService, { PIPELINE_TYPE_METRIC } from "../../services/pipeline.service";

class MetricForm extends Component {
    constructor(props) {
        super(props);
        this.handleCreate = this.handleCreate.bind(this);
        this.handleUpdate = this.handleUpdate.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.handleGetTasks = this.handleGetTasks.bind(this);
        this.handleGetPipelines = this.handleGetPipelines.bind(this);
        this.handlePrepareThresholds = this.handlePrepareThresholds.bind(this);
        this.onChangeSchedule = this.onChangeSchedule.bind(this);
        this.onChangeCrontab = this.onChangeCrontab.bind(this);
        this.enableScheduleEditor = this.enableScheduleEditor.bind(this);
        this.onChangeSource = this.onChangeSource.bind(this);
        this.onChangeSQLQuery = this.onChangeSQLQuery.bind(this);
        this.onChangeExpression = this.onChangeExpression.bind(this);
        this.onChangeExpressionFromOutside = this.onChangeExpressionFromOutside.bind(this);
        this.closeRunningForm = this.closeRunningForm.bind(this);
        this.openRunningForm = this.openRunningForm.bind(this);
        this.removeThresholds = this.removeThresholds.bind(this);
        this.handleThresholds = this.handleThresholds.bind(this);
        this.saveThreshold = this.saveThreshold.bind(this);
        this.onChangeActiveThresholdPipeline = this.onChangeActiveThresholdPipeline.bind(this);
        this.onChangeActiveThresholdExpression = this.onChangeActiveThresholdExpression.bind(this);
        this.handleOpenAddThresholdForm = this.handleOpenAddThresholdForm.bind(this);
        this.handleRunningOutput = this.handleRunningOutput.bind(this);
        this.onChangeActiveThresholdExpressionFromOutside = this.onChangeActiveThresholdExpressionFromOutside.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            schedule: "",
            scheduleEditorEnabled: false,
            item: this.props.item,
            sourceUuid: null,
            sourceName: null,
            sourceType: null,
            sourceValue: null,
            SQLQuery: "",
            sources: null,
            tasks: null,
            taskTriggers: null,
            queryTask: null,
            queryAggregatorTaskTrigger: null,
            aggregatorTask: null,
            isRunningFormOpen: false,
            isTerminationModalOpen: false,
            expression: "",
            pipelines: null,
            thresholds: null,
            loading: false,
            successful: false,
            activePipelineUuid: "",
            activePipelineName: "",
            terminationThreshold: null,
            thresholdsToRemove: [],
            thresholdsToUpdate: [],
            thresholdsToAdd: [],
            isEditThresholdFormOpen: false,
            updateThreshold: null,
            thresholdExpression: "value > 0",
            activeThresholdPipelineUuid: "",
            activeThresholdPipelineName: "",
            thresholdAdditionCounter: 0,
            folderUuid: this.props.folderUuid,
        };
    }

    componentDidMount = async() => {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            await this.handleGetTasks(this.props.item.uuid);
            await this.handleGetTaskTriggers(this.props.item.uuid);
            await this.promisedSetState({
                item: this.props.item,
                name: this.props.item.name || "",
                folderUuid: this.props.folderUuid,
                schedule: this.props.item.schedule || "",
                scheduleEditorEnabled: this.props.item.schedule !== null,
                activePipelineUuid: "",
                activePipelineName: "",
                thresholds: null,
            });
        }
        await this.handleGetSources(this.state.organization.uuid);
        await this.handleGetPipelines(this.state.organization.uuid);
        await this.handlePrepareThresholds();
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handlePrepareThresholds = async() => {
       let tasks = this.state.tasks;

        if (
            tasks !== null
            && tasks.length > 0
        ) {
            var conditionTasks = tasks.filter(o => o.type === TASK_TYPE_CONDITION);
            if (conditionTasks) {
                let thresholds = [];
                let threshold = [];
                let condition = undefined;
                let conditionTrigger = undefined;
                let pipeline = undefined;
                let pipelineTrigger = undefined;
                let pipelineName = "";
                for(var i = 0; i < conditionTasks.length; i++) {
                    threshold = [];
                    condition = conditionTasks[i];
                    conditionTrigger = this.state.taskTriggers.find(o => o.triggered_task_uuid === condition.uuid);
                    pipelineTrigger = this.state.taskTriggers.find(o => o.trigger_task_uuid === condition.uuid);
                    if (pipelineTrigger) {
                        pipeline = await this.state.tasks.find(o => o.uuid === pipelineTrigger.triggered_task_uuid);
                        pipelineName = this.state.pipelines.find(o => o.uuid === pipeline.implementation.pipeline_uuid);
                    }
                    threshold.condition = condition;
                    threshold.conditionTrigger = conditionTrigger;
                    threshold.pipeline = pipeline;
                    threshold.pipelineTrigger = pipelineTrigger;
                    threshold.pipelineName = pipelineName;
                    thresholds.push(threshold);
                }
                this.setState({thresholds});
            } else {
                this.setState({thresholds: []});
            }
        
        } else {
            this.setState({thresholds: []});
        }
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
                        
                        if (
                            items !== null
                            && items.length > 0
                        ) {
                            await this.promisedSetState({tasks: items});

                            var queryTask = items.find(o => o.type === TASK_TYPE_QUERY);
                            if (queryTask) {
                                await this.promisedSetState({
                                    queryTask,
                                    SQLQuery: queryTask.implementation.sql_query || "",
                                });
                            }

                            var aggregatorTask = items.find(o => o.type === TASK_TYPE_AGGREGATOR);
                            if (aggregatorTask) {
                                await this.promisedSetState({
                                    aggregatorTask,
                                    expression: aggregatorTask.implementation.expression || "",
                                });
                            }
                        } else {
                            await this.promisedSetState({tasks: []});
                        }
                    } else {
                        await this.promisedSetState({tasks: []});
                    }
                });
        }
    };

    handleGetTaskTriggers = async(uuid) => {
        let taskTriggers = this.state.taskTriggers;

        if (
            taskTriggers === null
            || taskTriggers.length === 0
        ) {
            taskTriggers = TasksTriggerService.getTaskTriggersByPipeline(uuid);

            await Promise.resolve(taskTriggers)
                .then(async(taskTriggers) => {
                    if (taskTriggers.data && taskTriggers.data.items !== null) {
                        var items = taskTriggers.data.items;
                        
                        if (
                            items !== null
                            && items.length > 0
                        ) {
                            await this.promisedSetState({taskTriggers: items});
                            if (
                                this.state.queryTask !== null
                                && this.state.aggregatorTask !== null
                            ) {
                                var queryAggregatorTaskTrigger = items.find(
                                    o => o.trigger_task_uuid === this.state.queryTask.uuid 
                                    && o.triggered_task_uuid === this.state.aggregatorTask.uuid
                                );
                                if (queryAggregatorTaskTrigger !== undefined) {
                                    await this.promisedSetState({
                                        queryAggregatorTaskTrigger,
                                    });
                                }
                            }
                        } else {
                            await this.promisedSetState({taskTriggers: []});
                        }
                    } else {
                        await this.promisedSetState({taskTriggers: []});
                    }
                });
        }
    };

    handleGetSources = async(uuid) => {
        let sources = this.state.sources;

        if (
            sources === null
            || sources.length === 0
        ) {
            sources = IntegrationService.getIntegrationsByOrganization(uuid, INTEGRATION_TYPE_SQL);

            Promise.resolve(sources)
                .then(async(sources) => {
                    if (sources.data) {
                        this.setState({sources: sources.data.items});
                        if (sources.data.items.length > 0) {
                            if (this.state.queryTask === null || this.state.queryTask.implementation.source_uuid === "") {
                                await this.promisedSetState({
                                    sourceUuid: sources.data.items[0].uuid,
                                    sourceName: sources.data.items[0].name,
                                    sourceType: sources.data.items[0].type,
                                    sourceValue: sources.data.items[0].value,
                                })
                            } else {
                                let source = sources.data.items.find(o => o.uuid === this.state.queryTask.implementation.source_uuid);
                                if (source) {
                                    await this.promisedSetState({
                                        sourceUuid: source.uuid,
                                        sourceName: source.name,
                                        sourceType: source.type,
                                        sourceValue: source.value,
                                    })
                                }
                            }
                        }
                    } else {
                        await this.promisedSetState({sources: []});
                    }
                });
        }
    };

    handleGetPipelines = async(uuid) => {
        let pipelines = this.state.pipelines;

        if (
            pipelines === null
            || pipelines.length === 0
        ) {
            pipelines = PipelineService.getPipelinesByOrganization(uuid);

            await Promise.resolve(pipelines)
                .then(async(pipelines) => {
                    if (pipelines.data && pipelines.data.items && pipelines.data.items !== null) {
                        var items = pipelines.data.items.filter(
                            k => k.type !== PIPELINE_TYPE_METRIC
                        );
                        
                        if (
                            items.length > 0
                        ) {
                            this.setState({
                                activeThresholdPipelineUuid: items[0].uuid,
                                activeThresholdPipelineName: items[0].name,
                            });
                        }

                        this.setState({pipelines: items});
                    } else {
                        this.setState({pipelines: []});
                    }
                });
        }
    };

    openRunningForm = async() => {
        await this.promisedSetState({
            isRunningFormOpen: true,
        });
    };

    closeRunningForm = () => {
        this.setState({
            isRunningFormOpen: false,
        });
    };

    onChangeExpression(e) {
        this.setState({
            expression: e.target.value,
        });
    }

    onChangeExpressionFromOutside(el) {
        this.setState({
            expression: el.value,
        });
    }

    onChangeSource(sourceUuid) {
        let source = this.state.sources.find(o => o.uuid === sourceUuid);
        if (source) {
            this.setState({
                sourceUuid: source.uuid,
                sourceName: source.name,
                sourceType: source.type,
                sourceValue: source.value,
            })
        }
    }

    onChangeSQLQuery(query) {
        this.setState({
            SQLQuery: query,
        });
    }

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeSchedule(schedule) {
        this.setState({schedule});
    }

    onChangeSchedule(schedule) {
        this.setState({
            schedule,
            scheduleEditorEnabled: schedule !== "",
        });
    }

    onChangeActiveThresholdPipeline(pipeline) {
        this.setState({
            activeThresholdPipelineUuid: pipeline.uuid,
            activeThresholdPipelineName: pipeline.name,
        });
    }

    onChangeActiveThresholdExpression(e) {
        this.setState({
            thresholdExpression: e.target.value,
        });
    }

    onChangeActiveThresholdExpressionFromOutside(el) {
        this.setState({
            thresholdExpression: el.value,
        });
    }

    enableScheduleEditor() {
        this.setState({
            scheduleEditorEnabled: true,
        });
    }

    onChangeCrontab(e) {
        this.setState({
            schedule: e.target.value,
            scheduleEditorEnabled: e.target.value !== "",
        });
    }

    handleUpdate(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                updatePipeline(
                    this.state.item.uuid,
                    this.state.name,
                    this.state.item.folder_uuid,
                    [],
                    this.state.schedule
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        "item": this.props.message.item
                    });
                }

                if (this.state.queryTask === null) {
                    this.handleCreateQuery(e);
                } else {
                    this.handleUpdateQuery(e);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleCreate(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                addPipeline(
                    this.state.name,
                    this.state.folderUuid,
                    this.state.organization.uuid,
                    [],
                    PIPELINE_TYPE_METRIC,
                    this.state.schedule
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        "item": this.props.message.item
                    });
                }

                if (this.state.queryTask === null) {
                    this.handleCreateQuery(e);
                } else {
                    this.handleUpdateQuery(e);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleUpdateQuery(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                updateTask(
                    this.state.queryTask.uuid, 
                    this.state.queryTask.pipeline_uuid,
                    TASK_TYPE_QUERY,
                    TASK_SEVERITY_MEDIUM,
                    TASK_TYPE_QUERY,
                    {
                        "sql_query": this.state.SQLQuery,
                        "source_uuid": this.state.sourceUuid,
                    }
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        "queryTask": this.props.message.item
                    });
                }

                if (this.state.aggregatorTask === null) {
                    this.handleCreateAggregator(e);
                } else {
                    this.handleUpdateAggregator(e);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleCreateQuery = async(e) => {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                addTask(
                    this.state.item.uuid, 
                    TASK_TYPE_QUERY,
                    TASK_TYPE_QUERY
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        queryTask: this.props.message.item
                    });
                    this.handleUpdateQuery(e);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleUpdateAggregator(e) {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                updateTask(
                    this.state.aggregatorTask.uuid, 
                    this.state.aggregatorTask.pipeline_uuid,
                    TASK_TYPE_AGGREGATOR,
                    TASK_SEVERITY_MEDIUM,
                    TASK_TYPE_AGGREGATOR,
                    {
                        "expression": this.state.expression,
                        "variable_name": "value", 
                    }
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    loading: false,
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        "aggregatorTask": this.props.message.item
                    });
                }

                await this.handleThresholds();

                if (
                    this.state.queryTask !== null
                    && this.state.aggregatorTask !== null
                ) {
                    if (this.state.taskTriggers !== null) {
                        var queryAggregatorTaskTrigger = this.state.taskTriggers.find(
                            o => o.trigger_task_uuid === this.state.queryTask.uuid 
                            && o.triggered_task_uuid === this.state.aggregatorTask.uuid
                        );
                        if (queryAggregatorTaskTrigger === undefined) {
                            this.handleTaskTriggerAddition(
                                this.state.queryTask.uuid,
                                this.state.aggregatorTask.uuid
                            );
                        }
                    } else {
                        this.handleTaskTriggerAddition(
                            this.state.queryTask.uuid,
                            this.state.aggregatorTask.uuid
                        );
                    }
                }

                setTimeout(() => {
                    this.props.successHandler(this.state.aggregatorTask.pipeline_uuid);
                }, 1000);
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleCreateAggregator = async(e) => {
        e.preventDefault();

        this.setState({
            loading: true,
            successful: false,
        });

        this.form.validateAll();

        const { dispatch } = this.props;

        if (this.checkBtn.context._errors.length === 0) {
            dispatch(
                addTask(
                    this.state.item.uuid, 
                    TASK_TYPE_AGGREGATOR,
                    TASK_TYPE_AGGREGATOR
                )
            )
            .then(async() => {
                await this.promisedSetState({
                    successful: true,
                });

                if (this.props.message.item) {
                    await this.promisedSetState({
                        aggregatorTask: this.props.message.item
                    });
                    this.handleUpdateAggregator(e);
                }
            })
            .catch(() => {
                this.setState({
                    loading: false,
                    successful: false,
                });
            });
        } else {
            this.setState({
                loading: false,
                successful: false,
            });
        }
    }

    handleTaskTriggerAddition = async(triggerTaskUuid, triggeredTaskUuid, triggerType = TRIGGER_TYPE_OUTPUT) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
        });

        return this.props
            .dispatch(
                addTaskTrigger(
                    this.state.item.uuid,
                    triggerTaskUuid, 
                    triggeredTaskUuid, 
                    triggerType
                )
            )
            .then(async() => {
                if (this.props.message.item) {
                    await this.promisedSetState({
                        "queryAggregatorTaskTrigger": this.props.message.item
                    });
                }
                await this.promisedSetState({
                    successful: true,
                    loading: false,
                });
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handleThresholds = async() => {
        await this.updateThresholds();
        await this.addThresholds();
        await this.removeThresholds();
    }

    removeThresholds = async() => {
        var thresholds = this.state.thresholdsToRemove;
        var pipeline = this.state.item;
        let threshold = null;
        for(var i = 0; i < thresholds.length; i++) {
            threshold = thresholds[i];

            await TasksTriggerService.deleteTaskTrigger(threshold.conditionTrigger.uuid, pipeline.uuid);
            await TasksTriggerService.deleteTaskTrigger(threshold.pipelineTrigger.uuid, pipeline.uuid);
            await TasksService.deleteTask(threshold.condition.uuid, pipeline.uuid);
            await TasksService.deleteTask(threshold.pipeline.uuid, pipeline.uuid);
        }

        await this.promisedSetState({thresholdsToRemove: []});
    }

    updateThresholds = async() => {
        var thresholds = this.state.thresholdsToUpdate;
        var pipeline = this.state.item;
        let threshold = null;
        for(var i = 0; i < thresholds.length; i++) {
            threshold = thresholds[i];

            await TasksService.updateTask(
                threshold.condition.uuid, 
                pipeline.uuid, 
                threshold.condition.name, 
                threshold.condition.severity, 
                threshold.condition.type,
                {
                    "expression": threshold.condition.implementation.expression,
                } 
            );

            await TasksService.updateTask(
                threshold.pipeline.uuid, 
                pipeline.uuid, 
                threshold.pipeline.name, 
                threshold.pipeline.severity, 
                threshold.pipeline.type,
                {
                    "pipeline_uuid": threshold.pipeline.implementation.pipeline_uuid,
                } 
            );
        }

        await this.promisedSetState({thresholdsToUpdate: []});
    }

    addThresholds = async() => {
        var thresholds = this.state.thresholdsToAdd;
        var pipeline = this.state.item;
        let threshold = null;
        for(var i = 0; i < thresholds.length; i++) {
            threshold = thresholds[i];

            let condition = TasksService.addTask(pipeline.uuid, TASK_TYPE_CONDITION, TASK_TYPE_CONDITION);
            await Promise.resolve(condition)
                .then(async(condition) => {
                    if (condition && condition.uuid) {
                        let wfRun = TasksService.addTask(pipeline.uuid, TASK_TYPE_RUN_PIPELINE, TASK_TYPE_RUN_PIPELINE);
                        await Promise.resolve(wfRun)
                            .then(async(wfRun) => {
                                if (wfRun && wfRun.uuid) {

                                    await TasksService.updateTask(
                                        condition.uuid, 
                                        pipeline.uuid, 
                                        condition.name, 
                                        TASK_SEVERITY_MEDIUM, 
                                        condition.type,
                                        {
                                            "expression": threshold.condition.implementation.expression,
                                        } 
                                    );

                                    await TasksService.updateTask(
                                        wfRun.uuid, 
                                        pipeline.uuid, 
                                        wfRun.name, 
                                        TASK_SEVERITY_MEDIUM, 
                                        wfRun.type,
                                        {
                                            "pipeline_uuid": threshold.pipeline.implementation.pipeline_uuid,
                                        } 
                                    );

                                    await TasksTriggerService.addTaskTrigger(
                                        pipeline.uuid, 
                                        this.state.aggregatorTask.uuid, 
                                        condition.uuid, 
                                        TRIGGER_TYPE_OUTPUT
                                    );

                                    await TasksTriggerService.addTaskTrigger(
                                        pipeline.uuid, 
                                        condition.uuid, 
                                        wfRun.uuid, 
                                        TRIGGER_TYPE_CONDITION_TRUE
                                    );
                                }
                            });
                    }
                });
        }

        await this.promisedSetState({thresholdsToUpdate: []});
    }

    handleOpenTerminationModal = (threshold) => {
        this.setState({
            isTerminationModalOpen: true,
            terminationThreshold: threshold,
        });
    };

    handleConfirmTermination = async() => {
        var pipeline = this.state.item;
        var threshold = this.state.terminationThreshold;

        if (threshold.condition.uuid === null) {
            var thresholdsToAdd = this.state.thresholdsToAdd.filter(
                o => o.additionNum !== threshold.additionNum
            );
            await this.promisedSetState({thresholdsToAdd});
        } else {
            if (pipeline !== null) {
                let thresholdsToRemove = this.state.thresholdsToRemove;
                thresholdsToRemove.push(threshold);
                await this.promisedSetState({thresholdsToRemove});
            }
        }

        if (threshold.condition.uuid !== null) {
            var thresholds = this.state.thresholds.filter(
                o => o.condition.uuid !== threshold.condition.uuid
            );
        } else {
            thresholds = this.state.thresholds.filter(
                o => o.additionNum !== threshold.additionNum
            );
        }

        await this.promisedSetState({thresholds});

        this.handleCloseTerminationModal();
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
            terminationThreshold: null,
        });
    };

    handleOpenEditThresholdForm = async(threshold) => {
        await this.promisedSetState({
            isEditThresholdFormOpen: true,
            updateThreshold: threshold,
            thresholdExpression: threshold.condition.implementation.expression,
        });

        let pipelines = this.state.pipelines;
        let activePipeline = pipelines.find(
            o => o.uuid === threshold.pipeline.implementation.pipeline_uuid
        );
        if (activePipeline) {
            await this.promisedSetState({
                activeThresholdPipelineUuid: activePipeline.uuid,
                activeThresholdPipelineName: activePipeline.name,
            });
        }
    };

    handleOpenAddThresholdForm = async() => {
        let threshold = [];
        let thresholdAdditionCounter = this.state.thresholdAdditionCounter;
        threshold.additionNum = thresholdAdditionCounter;
        threshold.condition = {
            "uuid": null,
            "implementation": [],
        };
        threshold.conditionTrigger = [];
        threshold.pipeline = {
            "uuid": null,
            "implementation": [],
        };
        threshold.pipelineTrigger = [];
        threshold.pipelineName = [];

        await this.promisedSetState({
            isEditThresholdFormOpen: true,
            updateThreshold: threshold,
            thresholdAdditionCounter: thresholdAdditionCounter + 1,
        });
    };

    handleCloseEditThresholdForm = () => {
        this.setState({
            isEditThresholdFormOpen: false,
            updateThreshold: null,
            thresholdExpression: "value > 0",
        });
    };

    handleRunningOutput = async(taskRuns, finished = false, oldOutput = "") => {

    }

    saveThreshold = async() => {
        this.form.validateAll();

        if (
            this.state.thresholdExpression !== ""
            && this.state.updateThreshold !== null
        ) {
            var threshold = this.state.updateThreshold;

            var thresholds = threshold.condition.uuid === null 
                ? this.state.thresholdsToAdd
                : this.state.thresholdsToUpdate
            ;

            var thresholdsToSave = threshold.condition.uuid !== null
                ?   thresholds.filter(
                        o => o.condition.uuid !== threshold.condition.uuid
                    )
                :   thresholds.filter(
                        o => o.additionNum !== threshold.additionNum
                    );

            threshold.condition.implementation.expression = this.state.thresholdExpression;
            threshold.pipeline.implementation.pipeline_uuid = this.state.activeThresholdPipelineUuid;
            threshold.pipelineName.name = this.state.activeThresholdPipelineName;

            thresholdsToSave.push(threshold);

            var thresholdsToView = this.state.thresholds;
            if (threshold.condition.uuid !== null) {
                await this.promisedSetState({thresholdsToUpdate: thresholdsToSave});
                for(var i = 0; i < thresholdsToView.length; i++) {
                    if (thresholdsToView[i].condition.uuid === threshold.condition.uuid) {
                        thresholdsToView[i].condition.implementation.expression = this.state.thresholdExpression;
                        thresholdsToView[i].pipeline.implementation.pipeline_uuid = this.state.activeThresholdPipelineUuid;
                        thresholdsToView[i].pipelineName.name = this.state.activeThresholdPipelineName;
                    }
                }
            } else {
                await this.promisedSetState({thresholdsToAdd: thresholdsToSave});
                let found = false;
                for(i = 0; i < thresholdsToView.length; i++) {
                    if (thresholdsToView[i].additionNum === threshold.additionNum) {
                        thresholdsToView[i].condition.implementation.expression = this.state.thresholdExpression;
                        thresholdsToView[i].pipeline.implementation.pipeline_uuid = this.state.activeThresholdPipelineUuid;
                        thresholdsToView[i].pipelineName.name = this.state.activeThresholdPipelineName;
                        found = true;
                    }
                }
                if (found === false) {
                    await thresholdsToView.push(threshold);
                }
            }

            await this.promisedSetState({thresholds: thresholdsToView});

            this.handleCloseEditThresholdForm();
        }
    }

    render() {
        const { isLoggedIn, user, message, item } = this.props;

        const { 
            isRunningFormOpen, 
            tasks, 
            isTerminationModalOpen, 
            isEditThresholdFormOpen,
            updateThreshold,
            thresholdExpression,
            activeThresholdPipelineUuid,
            activeThresholdPipelineName,
        } = this.state;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div>
                <div>
                    <Row>
                        <Col>
                            <Form
                                onSubmit={
                                    item === null
                                        ? this.handleCreate
                                        : this.handleUpdate
                                }
                                ref={(c) => {
                                    this.form = c;
                                }}
                            >
                                {message && this.state.loading === false && (
                                    <div className="form-group">
                                        <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                            {
                                                message.item 
                                                    ? 
                                                        (
                                                            item === null
                                                                ? "Metric successfully created"
                                                                : "Metric successfully updated"  
                                                        )
                                                    : message
                                            }
                                        </div>
                                    </div>
                                )}
                                <Card className="onboardingCard noBorder mb-5">
                                    <Card.Body className="p-4">
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                                <FloatingLabel controlId="floatingName" label="Name">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        type="text"
                                                        id="floatingName"
                                                        placeholder="Name"
                                                        autoComplete="name"
                                                        name="name"
                                                        value={this.state.name}
                                                        onChange={this.onChangeName}
                                                        validations={[required]}
                                                    />
                                                </FloatingLabel>
                                            </div>
                                        </InputGroup>
                                    </Card.Body>
                                </Card>

                                <Card className="onboardingCard noBorder mb-5">
                                    <Card.Body className="p-4">
                                        <h4>How to retrieve data for your metric</h4>
                                        {this.state.sources !== null ?
                                            <Dropdown size="lg" className="mb-4" variant="dark">
                                                <Dropdown.Toggle 
                                                    className={"dropdownItemWithBg dropdownItemWithBg-" + (this.state.sourceType === INTEGRATION_TYPE_SQL ? this.state.sourceValue : this.state.sourceType)}
                                                    variant="light" 
                                                    id="dropdown-basic"
                                                >
                                                    {this.state.sourceName}
                                                </Dropdown.Toggle>

                                                <Dropdown.Menu>
                                                    {this.state.sources.map(value => (
                                                        <Dropdown.Item
                                                            className={"dropdownItemWithBg dropdownItemWithBg-" + (value.type === INTEGRATION_TYPE_SQL ? value.value : value.type)}
                                                            value={value.name}
                                                            key={value.uuid}
                                                            active={value.uuid === this.state.sourceUuid}
                                                            onClick={(e) => this.onChangeSource(value.uuid)}
                                                        >
                                                            {value.name}
                                                        </Dropdown.Item>
                                                    ))}
                                                </Dropdown.Menu>
                                            </Dropdown>
                                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                        }

                                        { this.state.sources !== null
                                            && this.state.sourceUuid !== null
                                            &&
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <QueryUI
                                                        query={this.state.SQLQuery}
                                                        onChangeQuery={this.onChangeSQLQuery}
                                                        integrationUuid={this.state.sourceUuid}
                                                    />
                                                </div>
                                            </InputGroup>
                                        }
                                    </Card.Body>
                                </Card>

                                <Card className="onboardingCard noBorder mb-5">
                                    <Card.Body className="p-4">
                                        <h4>How to calculate your metric</h4>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                                    <CodeEditor
                                                        className="form-control form-control-lg codeEditor"
                                                        type="text"
                                                        language="xls"
                                                        id="floatingExpression"
                                                        minHeight={200}
                                                        autoComplete="Expression"
                                                        name="expression"
                                                        value={this.state.expression}
                                                        onChange={this.onChangeExpression}
                                                        validations={[required]}
                                                        rehypePlugins={[
                                                            [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                                        ]}
                                                        style={{
                                                            fontSize: 14,
                                                            fontFamily: 'Source Code Pro, monospace',
                                                        }}
                                                    />
                                                    <TextareaEditor 
                                                        txtId="floatingExpression" 
                                                        callback={this.onChangeExpressionFromOutside}
                                                        PipelineType={PIPELINE_TYPE_METRIC}
                                                    />
                                            </div>
                                        </InputGroup>
                                    </Card.Body>
                                </Card>

                                <Card className="onboardingCard noBorder mb-5">
                                    <Card.Body className="p-4">
                                        <h4>When or how often to run your metric</h4>
                                        <InputGroup className="mb-4">
                                            <div className="registrationFormControl">
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <div className="p-3" onClick={this.enableScheduleEditor}>
                                                            <Cron 
                                                                value={this.state.schedule}
                                                                setValue={this.onChangeSchedule}
                                                                humanizeValue={false}
                                                                humanizeLabels={true}
                                                                leadingZero
                                                                defaultPeriod={'minute'}
                                                                clearButtonAction={'empty'}
                                                                disabled={!this.state.scheduleEditorEnabled}
                                                                allowEmpty={'always'}
                                                            />
                                                        </div>
                                                        <span className="inputLabel">Or write it in a cronjob format</span><br/>
                                                        <Input
                                                            className="form-control"
                                                            type="text"
                                                            id="floatingSchedule"
                                                            placeholder=""
                                                            autoComplete="schedule"
                                                            name="schedule"
                                                            value={this.state.schedule}
                                                            onChange={this.onChangeCrontab}
                                                            validations={[isCron]}
                                                        />
                                                    </div>
                                                </InputGroup>
                                            </div>
                                        </InputGroup>
                                    </Card.Body>
                                </Card>

                                <Card className="onboardingCard noBorder mb-5">
                                    <Card.Body className="p-4">
                                        <Row className="mb-3">
                                            <Col sm={8}>
                                                <h4>How to react on the value of your metric</h4>
                                            </Col>
                                            <Col sm={4} className="text-right">
                                                <Button
                                                    variant="secondary"
                                                    className="mx-0"
                                                    onClick={() => this.handleOpenAddThresholdForm()}
                                                >
                                                    Add threshold
                                                </Button>
                                            </Col>
                                        </Row>
                                        {
                                            isEditThresholdFormOpen === false 
                                            && updateThreshold === null
                                            &&
                                            (
                                            this.state.thresholds !== null ?
                                                this.state.thresholds.length > 0 ?
                                                    <div>
                                                        <Row>
                                                            <Col sm={6} className="p-4 bottomBorder"><strong>Condition</strong></Col>
                                                            <Col sm={4} className="p-4 bottomBorder"><strong>Pipeline</strong></Col>
                                                            <Col sm={2} className="p-4 bottomBorder"><strong>Actions</strong></Col>
                                                        </Row>
                                                        {this.state.thresholds.map(value => (
                                                            <Row>
                                                                <Col sm={6} className="p-4 bottomBorder">{value.condition.implementation.expression}</Col>
                                                                <Col sm={4} className="p-4 bottomBorder">
                                                                    {value.pipelineName.name}
                                                                </Col>
                                                                <Col sm={2} className="p-4 bottomBorder">
                                                                    <Tooltip title="Edit" placement="left">
                                                                        <Edit 
                                                                            className="iconDelete"
                                                                            onClick={() => this.handleOpenEditThresholdForm(value)}
                                                                        />
                                                                    </Tooltip>
                                                                    <Tooltip title="Delete" placement="right">
                                                                        <DeleteOutlined 
                                                                            className="iconDelete"
                                                                            onClick={() => this.handleOpenTerminationModal(value)}
                                                                        />
                                                                    </Tooltip>
                                                                </Col>
                                                            </Row>
                                                        ))}
                                                    </div>
                                                : <div className="text-center pt-4 pb-4">You didn't create any threshold yet</div>
                                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                            )
                                        }
                                        {
                                            isEditThresholdFormOpen === true 
                                            && updateThreshold !== null
                                            &&
                                            <div className="filterForm">
                                                <ClearRounded onClick={() => this.handleCloseEditThresholdForm()} className="closeFilterFormIcon"/>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                            <CodeEditor
                                                                className="form-control form-control-lg codeEditor"
                                                                type="text"
                                                                language="xls"
                                                                id="floatingCondition"
                                                                minHeight={100}
                                                                autoComplete="Condition"
                                                                name="condition"
                                                                value={thresholdExpression}
                                                                onChange={this.onChangeActiveThresholdExpression}
                                                                validations={[required]}
                                                                rehypePlugins={[
                                                                    [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                                                ]}
                                                                style={{
                                                                    fontSize: 14,
                                                                    fontFamily: 'Source Code Pro, monospace',
                                                                }}
                                                            />

                                                            <TextareaEditor 
                                                                txtId="floatingCondition" 
                                                                callback={this.onChangeActiveThresholdExpressionFromOutside}
                                                                pipelineType={PIPELINE_TYPE_METRIC}
                                                            />
                                                    </div>
                                                </InputGroup>
                                                <span className="inputLabel">Which pipeline to run if the condition is true</span><br/>
                                                {this.state.pipelines !== null &&
                                                    this.state.pipelines.length > 0 ? 
                                                    <>
                                                        <Dropdown size="lg" className="mb-4">
                                                            <Dropdown.Toggle 
                                                                variant="light" 
                                                                id="dropdown-basic"
                                                                className="mt-2"
                                                            >
                                                            { activeThresholdPipelineUuid !== "" &&
                                                                <>
                                                                    {activeThresholdPipelineName}
                                                                </>
                                                            }
                                                            </Dropdown.Toggle>

                                                            <Dropdown.Menu className="tasks">
                                                                {this.state.pipelines.map(value => (
                                                                    <Dropdown.Item
                                                                        value={value.name}
                                                                        key={value.uuid}
                                                                        active={this.state.activeThresholdPipelineUuid !== "" && value.uuid === this.state.activeThresholdPipelineUuid}
                                                                        onClick={(e) => this.onChangeActiveThresholdPipeline(value)}
                                                                    >
                                                                        {value.name}
                                                                    </Dropdown.Item>
                                                                ))}
                                                            </Dropdown.Menu>
                                                        </Dropdown>
                                                    </>
                                                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                                }
                                                <Button
                                                    className="px-3 mt-3 block btn btn-secondary"
                                                    onClick={() => this.saveThreshold()}
                                                >
                                                    <span>Save threshold</span>
                                                </Button>
                                            </div>
                                        }
                                    </Card.Body>
                                </Card>

                                        <Row>
                                            <Col className="text-right">
                                                <Button
                                                    className="px-4 mx-3 btn btn-primary"
                                                    disabled={this.state.loading}
                                                    type="submit"
                                                >
                                                    {this.state.loading && (
                                                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                    )}
                                                    <span>Save</span>
                                                </Button>
                                                <Button 
                                                    variant="secondary" 
                                                    onClick={this.openRunningForm}
                                                >
                                                    Test run
                                                </Button>
                                            </Col>
                                        </Row>
                                        {message && this.state.loading === false && (
                                            <div className="form-group">
                                                <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                                    {
                                                        message.item 
                                                            ? 
                                                                (
                                                                    item === null
                                                                    ? "Metric successfully created"
                                                                    : "Metric successfully updated"  
                                                                )
                                                            : message
                                                    }
                                                </div>
                                            </div>
                                        )}
                                        <CheckButton
                                            style={{ display: "none" }}
                                                ref={(c) => {
                                                    this.checkBtn = c;
                                                }}
                                        />
                                    </Form>
                                
                        </Col>
                    </Row>
                </div>

                <RightModal
                    show={isRunningFormOpen}
                    content={
                      <PipelineRun 
                        item={item} 
                        elements={tasks}
                        outputHandler={this.handleRunningOutput}
                        showLogWithoutRunning={false}
                        oldOutput={""}
                      />
                    }
                    item={item}
                    title={item !== null ? "Metric '" + item.name + "' is running" : "Metric is running"}
                    onHide={this.closeRunningForm}
                />

                <ConfirmationModal
                    show={isTerminationModalOpen}
                    title="Delete threshold"
                    body={'Are you sure you want to delete this threshold?'}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />                
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { message } = state.message;
    const { user } = state.auth;
    return {
        isLoggedIn,
        message,
        user
    };
}

export default connect(mapStateToProps)(MetricForm);
