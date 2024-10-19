import React, { Component } from 'react';
import { DateTime } from "luxon";

import html2canvas from 'html2canvas'

import ReactFlow, {
  removeElements,
  addEdge,
  MiniMap,
  Controls,
  Background,
} from 'react-flow-renderer';

import WarningAmberOutlined from '@mui/icons-material/WarningAmberOutlined';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup"; 

import Input from "../formControls/input.component";
import { required } from "../formControls/validations";

import { connect } from "react-redux";

import PipelineSidebar from './pipelineSidebar.component';
import PipelineRun, { decodeOutput, handleErrors } from './pipelineRun.component';
import PipelineRunOutput from './pipelineRunOutput.component';

import {TimeAgo} from "../timeAgo.component";
import TaskForm from "./forms/taskForm.component"
import ScheduleForm from "../forms/scheduleForm.component";

import { AggregatorNodeComponent } from "./nodes/aggregatorNode.component";
import { ConditionNodeComponent, CONDITION_CONNECTOR_TRUE, CONDITION_CONNECTOR_FALSE } from "./nodes/conditionNode.component";
import { QueryNodeComponent } from "./nodes/queryNode.component";
import { NotificationNodeComponent } from "./nodes/notificationNode.component";
import { TransformerNodeComponent } from "./nodes/transformerNode.component";
import { APICallNodeComponent } from "./nodes/APICallNode.component";
import { ForEachNodeComponent } from "./nodes/forEachNode.component";
import { MergeNodeComponent } from "./nodes/mergeNode.component";
import { FilterNodeComponent } from "./nodes/filterNode.component";
import { PipelineRunNodeComponent } from "./nodes/pipelineRunNode.component";
import { ExternalTriggerNodeComponent } from "./nodes/externalTriggerNode.component";
import { ProcessorNodeComponent } from "./nodes/processorNode.component";
import { CodeNodeComponent } from "./nodes/codeNode.component";
import { GptNodeComponent } from "./nodes/gptNode.component";

import RightModal from "../modals/rightModal.component";
import TaskModal from "../modals/taskModal.component";

import { addPipeline, updatePipeline } from "../../actions/pipelines";
import { addTask, deleteTask } from "../../actions/tasks";
import { addTaskTrigger, updateTaskTrigger, deleteTaskTrigger } from "../../actions/taskTriggers";
import { clearMessage } from "../../actions/message";

import { 
  TRIGGER_TYPE_SCHEDULE,
  TRIGGER_TYPE_CONDITION_TRUE,
  TRIGGER_TYPE_CONDITION_FALSE,
  TRIGGER_TYPE_OUTPUT
} from "../../services/taskTrigger.service";

import PipelineService from "../../services/pipeline.service"; 

import {
  TASK_TYPE_QUERY,
  TASK_TYPE_CONDITION,
  TASK_TYPE_RUN_PIPELINE,
} from "../../services/task.service";

export const nodeTypes = {
  aggregator: AggregatorNodeComponent,
  condition: ConditionNodeComponent,
  query: QueryNodeComponent,
  notification: NotificationNodeComponent,
  api_call: APICallNodeComponent,
  transformer: TransformerNodeComponent,
  for_each: ForEachNodeComponent,
  merge: MergeNodeComponent,
  filter: FilterNodeComponent,
  external_trigger: ExternalTriggerNodeComponent,
  code: CodeNodeComponent,
  python: CodeNodeComponent,
  gpt: GptNodeComponent,
  run_pipeline: PipelineRunNodeComponent,
  processor: ProcessorNodeComponent,
};

export const nodeEdgeTypes = {
  schedule: "",
  output: "",
  condition_true: "",
  condition_false: "",
};

export const nodeColors = {
  query: '#52cc7f',
  condition: '#FAB910',
  processor: '#35b0bb',
  aggregator: '#F26023',
  transformer: '#9769FF',
  notification: '#02A5FF',
  for_each: '#FE7188',
  api_call: '#aadb40',
  filter: '#4264fc',
  merge: '#ba4eb6',
  external_trigger: '#9C7178',
  run_pipeline: '#e88f09',
  code: '#666600',
  gpt: '#598094',
};

const RUN_STATE_PENDING = "pending"
const RUN_STATE_EXECUTED = "executed"

const TASK_PAUSED = "taskPaused";
const TASK_FAILED = "taskFailed";
const TASK_SUCCESS = "taskSuccess";
const TASK_RUNNING = "taskRunning";

const TEXT_TASK_PAUSED = "Task paused";
const TEXT_TASK_FAILED = "Task failed";
const TEXT_TASK_EXECUTED = "Task executed";
const TEXT_TASK_IN_PROGRESS = "Task is in progress";
const TEXT_TASK_FAILED_TIMEOUT = "Timeout error.<br/><br/>Task execution takes longer than 60 seconds.<br/><br/>Please check that your data source processes the query as expected and try again later.";
const TEXT_TASK_FAILED_PIPELINE_CANNOT_RUN = "Something went wrong.<br/><br/>Please make sure you have access permisions to it, and the pipeline is configured correctly, and try again.<br/><br/>If you are not an administrator of the organization, please ask an administrator for a help";

class Pipeline extends Component {
    constructor(props) {
      super(props);
      this.handlePipeline = this.handlePipeline.bind(this);
      this.handleTaskAddition = this.handleTaskAddition.bind(this);
      this.handleTaskDeletion = this.handleTaskDeletion.bind(this);
      this.handleTaskTriggerAddition = this.handleTaskTriggerAddition.bind(this);
      this.handleTaskTriggerDeletion = this.handleTaskTriggerDeletion.bind(this);
      this.handleTaskTriggerUpdate = this.handleTaskTriggerUpdate.bind(this);
      this.handleForm = this.handleForm.bind(this);
      this.onChangeName = this.onChangeName.bind(this);
      this.onChangeItem = this.onChangeItem.bind(this);
      this.onLoad = this.onLoad.bind(this);
      this.onElementsRemove = this.onElementsRemove.bind(this);
      this.onConnect = this.onConnect.bind(this);
      this.onUpdateScheduledConnection = this.onUpdateScheduledConnection.bind(this);
      this.onNodeDragStop = this.onNodeDragStop.bind(this);
      this.onDragOver = this.onDragOver.bind(this);
      this.onDrop = this.onDrop.bind(this);
      this.onAdd = this.onAdd.bind(this);
      this.onNodeDoubleClick = this.onNodeDoubleClick.bind(this);
      this.addNode = this.addNode.bind(this);
      this.setElements = this.setElements.bind(this);
      this.isTask = this.isTask.bind(this);
      this.isTaskTrigger = this.isTaskTrigger.bind(this);
      this.isTaskCondition = this.isTaskCondition.bind(this);
      this.defineTaskTriggerType = this.defineTaskTriggerType.bind(this);
      this.captureCanvas = this.captureCanvas.bind(this);
      this.generateUpdatedAt = this.generateUpdatedAt.bind(this);
      this.onConnectStart = this.onConnectStart.bind(this);
      this.openTaskForm = this.openTaskForm.bind(this);
      this.closeTaskForm = this.closeTaskForm.bind(this);
      this.closeRunningForm = this.closeRunningForm.bind(this);
      this.closeRunningOutput = this.closeRunningOutput.bind(this);
      this.openRunningForm = this.openRunningForm.bind(this);
      this.openRunningOutput = this.openRunningOutput.bind(this);
      this.closeScheduleForm = this.closeScheduleForm.bind(this);
      this.openScheduleForm = this.openScheduleForm.bind(this);
      this.handlePipelineRunningSuccess = this.handlePipelineRunningSuccess.bind(this);
      this.handleTaskFormSuccess = this.handleTaskFormSuccess.bind(this);
      this.isTriggerExist = this.isTriggerExist.bind(this);
      this.isReverseTriggerExist = this.isReverseTriggerExist.bind(this);
      this.handleRunningOutput = this.handleRunningOutput.bind(this);
      this.cleanOutput = this.cleanOutput.bind(this);
      this.changeRunningClass = this.changeRunningClass.bind(this);

      this.reactFlowWrapper = React.createRef();

      this.state = {
        organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
        elements: JSON.parse(this.props.elements),
        name: this.props.newTitle,
        preview: "",
        item: this.props.item,
        folderUuid: this.props.folderUuid,
        activeTaskItem: null,
        activeTaskTrigger: null,
        successful: false,
        loading: false,
        isInProgress: false,
        isTaskFormOpen: false,
        isRunningFormOpen: false,
        isScheduleFormOpen: false,
        isOutputVisible: false,
        isPipelineRunning: false,
        isPipelineFailed: false,
        isRunningOutputOpen: false,
        showLogWithoutRunning: false,
        oldOutput: "",
        reactFlowInstance: null,
      };
    }

    componentDidMount() {
        this.mounted = true;
        if (this.props.item !== null) {
            this.setState({
                name: this.props.item.name,
                folderUuid: this.props.folderUuid,
            });
        }
    };

    componentWillUnmount() {
        this.mounted = false;
        this.setState({
            isRunningFormOpen: false,
            isPipelineRunning: false,
            isRunningOutputOpen: false,
        });
    }

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    isTask = (item) => {
        return (item.type in nodeTypes)
    };

    isTaskTrigger = (item) => {
        var element = this.state.elements.find(o => o.id === item.id);

        if (element) {
            return (element.trigger_type in nodeEdgeTypes);
        }

        return false;
    };

    isTaskCondition = (id) => {
        var element = this.state.elements.find(o => o.id === id);

        return this.isTask(element) && element.type === TASK_TYPE_CONDITION;
    };

    isTriggerExist = (target, source, trigger_type) => {
        var element = this.state.elements.find(o => o.trigger_type === trigger_type && o.target === target && o.source === source);

        return element !== undefined;
    };

    isReverseTriggerExist = (target, source, trigger_type) => {
        var element = this.state.elements.find(
            o => o.trigger_type === trigger_type 
            && o.source === target 
            && o.target === source
        );

        return element !== undefined;
    };

    onLoad = (reactFlowInstance) => {
        reactFlowInstance.fitView();
        this.setState({reactFlowInstance});
    };

    onNodeDoubleClick = (e, node) => {
        this.openTaskForm(node);
    };

    onChangeItem(item) {
        this.setState({item});

        this.props.handleSetActiveItem(item);
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    cleanOutput(e) {
        document.querySelectorAll(".runningInfo").forEach(el => el.remove());
        document.querySelectorAll(".runningText").forEach(el => el.remove());
        this.setState({
            isOutputVisible: false,
            oldOutput: "",
        });
    }

    setElements(elements) {
      var item = this.state.item;
      if (item !== null) {
        item.elements_layout = JSON.stringify(elements);
      }
      this.setState({elements, item});
    };

    onElementsRemove = async(elementsToRemove) => {
        for(var i = 0; i < elementsToRemove.length; i++){
            if (this.isTask(elementsToRemove[i])) {
              await this.handleTaskDeletion(elementsToRemove[i].id);
            } else if (this.isTaskTrigger(elementsToRemove[i])) {
              await this.handleTaskTriggerDeletion(elementsToRemove[i].id);
            }
        }

        var newElements = removeElements(elementsToRemove, this.state.elements);
        await this.setElements(newElements);
        this.handlePipeline();
        this.captureCanvas();
    };

    onNodeDragStop = async(event, node) => {
        var elements = this.state.elements;

        for(var i = 0; i < elements.length; i++){
            if (elements[i].id === node.id) {
                elements[i] = node;
            }
        }

        await this.setElements(elements);
        let pipeline = await this.handlePipeline();
        if(pipeline !== undefined) {
            this.captureCanvas();
        }
    };

    defineTaskTriggerType = (triggerTaskUuid, triggeredTaskUuid, targetHandle, schedule = "") => {
      if (triggerTaskUuid !== 0) {  
        var element = this.state.elements.find(o => o.id === triggerTaskUuid);

        if (element) {
            if (element.type === TASK_TYPE_CONDITION) {
                /*var edges = this.state.elements.filter(
                  o => o.target === triggerTaskUuid && o.targetHandle === targetHandle
                );
                if (edges.length === 0) {*/
                    if (targetHandle === CONDITION_CONNECTOR_TRUE) {
                        return TRIGGER_TYPE_CONDITION_TRUE;
                    } else if (targetHandle === CONDITION_CONNECTOR_FALSE) {
                        return TRIGGER_TYPE_CONDITION_FALSE;
                    } else {
                        return null;
                    }
                /*} else {
                    return null;
                }*/
            }
        }
      } else if (triggeredTaskUuid !== 0) {  
        element = this.state.elements.find(o => o.id === triggeredTaskUuid);

        if (element) {
            if (element.type === TASK_TYPE_QUERY && schedule !== "") {
                return TRIGGER_TYPE_SCHEDULE;
            }
        }
      }

      return TRIGGER_TYPE_OUTPUT;
    };

    onConnect = async(params) => {
        params.trigger_type = this.defineTaskTriggerType(params.target, params.source, params.targetHandle, params.schedule);

        if (
            params.trigger_type === null
            || params.target === params.source
            || this.isTriggerExist(
                params.target, 
                params.source,
                params.trigger_type
            )
            || this.isReverseTriggerExist(
                params.target, 
                params.source,
                params.trigger_type
            )
        ) {
            return;
        }

        if (
            params.trigger_type !== TRIGGER_TYPE_SCHEDULE
        ) {
            params.schedule = "";
        }

        if (
            params.trigger_type === TRIGGER_TYPE_CONDITION_TRUE
        ) {
            params.label = "true"
        } else if (
            params.trigger_type === TRIGGER_TYPE_CONDITION_FALSE
        ) {
            params.label = "false"
        }

        var trigger = await this.handleTaskTriggerAddition(
            params.target !== 0 ? params.target : "", 
            params.source, 
            params.trigger_type, 
            params.schedule
        );

        if (trigger && trigger.uuid) {
            params.id = trigger.uuid;

            if (params.trigger_type !== TRIGGER_TYPE_SCHEDULE){
                var newElements = addEdge(params, this.state.elements);
            } else {
                var elements = this.state.elements;
                newElements = elements.concat(params);
            }
            await this.setElements(newElements);
            this.handlePipeline();
        }
    };

    onUpdateScheduledConnection = async(trigger, schedule) => {
        var update = await this.handleTaskTriggerUpdate(
            trigger.id,
            TRIGGER_TYPE_SCHEDULE,
            schedule
        );

        if (update === true) {
            var elements = this.state.elements;
            
            for (var i = 0; i < elements.length; i++) {
                if (elements[i].id === trigger.id) {
                    elements[i].schedule = schedule;
                }
            }

            await this.setElements(elements);
            this.handlePipeline();
        }
    };

    onDragOver = (e) => {
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
    };

    onConnectStart = (event, { nodeId, handleType }) => {
        if (this.isTaskCondition(nodeId)) {

        }
    };

    onDrop = async(e) => {
        e.preventDefault();
        var reactFlowInstance = this.state.reactFlowInstance;

        const reactFlowBounds = this.reactFlowWrapper.current.getBoundingClientRect();
        const type = e.dataTransfer.getData('application/reactflow');
        const position = reactFlowInstance.project({
          x: e.clientX - reactFlowBounds.left,
          y: e.clientY - reactFlowBounds.top,
        });
        
        this.addNode(type, position);
    };

    addNode = async(type, position) => {
        var name = type.charAt(0).toUpperCase() + type.slice(1);

        var node = await this.handleTaskAddition(name, type);

        if (node && node.uuid) {
          const newNode = {
            id: node.uuid,
            type,
            position,
            data: { name },
          };

          var elements = this.state.elements;
          var newElements = elements.concat(newNode);
          await this.setElements(newElements);
          let pipeline = await this.handlePipeline();
          if(pipeline !== undefined) {
            this.captureCanvas();
            return newNode;
          } else {
            return null;
          }
        }
    };

    handleRunningOutput = async(taskRuns, finished = false, oldOutput = "", timeoutError = false, isExternalFailure = false) => {
        //console.log(taskRuns);

        if (this.mounted !== true) {
            return;
        }

        var elementText = "...";
        var elementTitle = "Task is pending";
        var elementClass = "";

        if (finished !== true) {
            await this.promisedSetState({isPipelineRunning: true});
            Array.from(document.getElementsByClassName('node')).forEach(function(element){
                element.classList.add('nodeTBD');
                let elementUuid = element.parentElement.getAttribute("data-id");
                let infoElements = element.querySelectorAll('.runningInfo');

                if (infoElements.length === 0){
                    // Icon element
                    let div = document.createElement('div');
                    div.classList.add('runningInfo');

                    // Info element under icon
                    let div2 = document.createElement('div');
                    div2.classList.add('runningText');
                    let newId = "text-for-" + elementUuid;
                    div2.setAttribute("id", newId);

                    let div3 = document.createElement('div');
                    div3.classList.add('runningTextTitle');
                    const newContent3 = document.createTextNode(elementTitle);
                    div3.appendChild(newContent3);
                    div2.appendChild(div3);

                    let div5 = document.createElement('div');
                    div5.classList.add('runningInfoClose');
                    const newContent5 = document.createTextNode("X");
                    div5.appendChild(newContent5);
                    div2.appendChild(div5);

                    let div4 = document.createElement('div');
                    div4.classList.add('runningTextInside');
                    const newContent4 = document.createTextNode(elementText);
                    div4.appendChild(newContent4);
                    div2.appendChild(div4);

                    element.appendChild(div);
                    element.appendChild(div2);

                    div.onclick = function() {
                        if (
                            !this.classList.contains(TASK_PAUSED)
                            && !this.classList.contains(TASK_FAILED)
                            && !this.classList.contains(TASK_SUCCESS)
                            && !this.classList.contains(TASK_RUNNING)
                        ) {
                            return
                        }
                        let div = this.parentElement.querySelector(".runningText");
                        if (div) {
                            let hasClass = div.classList.contains('runningTextVisible');
                            div.removeAttribute('class');
                            div.classList.add('runningText');

                            if (!hasClass) {
                                div.classList.add('runningTextVisible');
                            }
                        } 
                    };

                    div5.onclick = function() {
                        let infoElement = this.parentElement.parentElement.querySelector(".runningInfo");
                        
                        if (
                            !infoElement.classList.contains(TASK_PAUSED)
                            && !infoElement.classList.contains(TASK_FAILED)
                            && !infoElement.classList.contains(TASK_SUCCESS)
                            && !infoElement.classList.contains(TASK_RUNNING)
                        ) {
                            return
                        }
                        let div = infoElement.parentElement.querySelector(".runningText");
                        if (div) {
                            let hasClass = div.classList.contains('runningTextVisible');
                            div.removeAttribute('class');
                            div.classList.add('runningText');

                            if (!hasClass) {
                                div.classList.add('runningTextVisible');
                            }
                        } 
                    };
                }
            });
        } else {
            await this.promisedSetState({isPipelineRunning: false});
            this.closeRunningForm(timeoutError, isExternalFailure);
            Array.from(document.getElementsByClassName('node')).forEach(element => (
                element.classList.remove('nodeTBD')
            ));
        }

        var element = null;
        for(var i = 0; i < taskRuns.length; i++) {
            element = document.querySelector('[data-id="' + taskRuns[i].task_uuid + '"]');
            if (
                taskRuns[i].state === RUN_STATE_PENDING
                || taskRuns[i].state === RUN_STATE_EXECUTED
            ) {
                element.firstChild.classList.remove('nodeTBD');
            }

            if (taskRuns[i].state === RUN_STATE_PENDING) {
                elementTitle = TEXT_TASK_IN_PROGRESS;
                elementText = "...";
                elementClass = TASK_RUNNING;
            } else {
                if (taskRuns[i].state === RUN_STATE_EXECUTED) {
                    if (taskRuns[i].is_successful === true) {
                        elementClass = TASK_SUCCESS;
                        elementTitle = TEXT_TASK_EXECUTED;
                        elementText = decodeOutput(atob(taskRuns[i].output), true, false);
                    } else {
                        elementClass = TASK_FAILED;
                        elementTitle = TEXT_TASK_FAILED;
                        elementText = "Errors:\n" + handleErrors(taskRuns[i].errors);
                    }
                }
            }

            element.firstChild.querySelector(".runningTextTitle").textContent = elementTitle;
            element.firstChild.querySelector(".runningTextInside").innerText = elementText;
            this.changeRunningClass(element.firstChild.querySelector(".runningInfo"), elementClass);

            elementText = "...";
            elementTitle = "Task is pending";
            elementClass = "";
        }

        this.setState({isOutputVisible: true, oldOutput});
    };

    changeRunningClass = (e, c) => {
        e.removeAttribute('class');
        e.classList.add('runningInfo');
        e.classList.add(c);
    }

    onAdd = async(e, type) => {
        /*const viewBox = document.getElementsByClassName('react-flow__minimap')[0].getAttribute("viewBox");
        const dimensions = viewBox.split(" ");
        console.log(
          viewBox,
          parseInt(dimensions[0]) + parseInt(dimensions[2])/3,
          parseInt(dimensions[1]) + parseInt(dimensions[3])/3
        );
        this.addNode(
          type, 
          {
            x: parseInt(dimensions[0]) + parseInt(dimensions[2])/3,
            y: parseInt(dimensions[1]) + parseInt(dimensions[3])/3,
          }
        );*/
    }

    handleForm = async(e) => {
        e.preventDefault();
        let pipeline = await this.handlePipeline();
        if(pipeline !== undefined) {
            this.captureCanvas();
        }
    };

    handleTaskTriggerAddition = async(triggerTaskUuid, triggeredTaskUuid, triggerType, schedule) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        var pipeline = this.state.item;

        return this.props
            .dispatch(
                addTaskTrigger(
                    pipeline.uuid,
                    triggerTaskUuid, 
                    triggeredTaskUuid, 
                    triggerType, 
                    schedule
                )
            )
            .then(() => {
                this.setState({
                    successful: true,
                    loading: false
                });

                let item = this.state.item;
                item.updated_at = this.generateUpdatedAt();
                this.setState({
                    item, 
                    isInProgress: false,
                });

                if (this.props.message.item) {
                    return this.props.message.item
                }
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handleTaskTriggerUpdate = async(uuid, triggerType, schedule) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        var pipeline = this.state.item;

        return this.props
            .dispatch(
                updateTaskTrigger(
                    uuid,
                    pipeline.uuid, 
                    triggerType, 
                    schedule
                )
            )
            .then(() => {
                this.setState({
                    successful: true,
                    loading: false
                });

                let item = this.state.item;
                item.updated_at = this.generateUpdatedAt();
                this.setState({
                    item,
                    isInProgress: false,
                });

                return true;
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handleTaskTriggerDeletion = async(uuid) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        var pipeline = this.state.item;

        return this.props
            .dispatch(
                deleteTaskTrigger(
                    uuid,
                    pipeline.uuid
                )
            )
            .then(() => {
                this.setState({
                    successful: true,
                    loading: false
                });

                let item = this.state.item;
                item.updated_at = this.generateUpdatedAt();
                this.setState({
                    item,
                    isInProgress: false,
                });

                if (this.props.message) {
                    return this.props.message
                }
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handleTaskAddition = async(name, type) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        var pipeline = this.state.item;

        if (pipeline === null) {
            pipeline = await this.handlePipeline();
            if(pipeline !== undefined) {
                await this.captureCanvas();
            } else {
                return null;
            }
        }

        return this.props
            .dispatch(
                addTask(
                    pipeline.uuid,
                    name, 
                    type
                )
            )
            .then(() => {
                this.setState({
                    successful: true,
                    loading: false
                });

                let item = this.state.item;
                item.updated_at = this.generateUpdatedAt();
                this.setState({
                    item,
                    isInProgress: false,
                });

                if (this.props.message.item) {
                    return this.props.message.item
                }
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handleTaskDeletion = async(uuid) => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        var pipeline = this.state.item;

        return this.props
            .dispatch(
                deleteTask(
                    uuid,
                    pipeline.uuid
                )
            )
            .then(() => {
                this.setState({
                    successful: true,
                    loading: false
                });

                let item = this.state.item;
                item.updated_at = this.generateUpdatedAt();
                this.setState({
                    item,
                    isInProgress: false,
                });

                if (this.props.message) {
                    return this.props.message
                }
            })
            .catch(() => {
                this.setState({
                    successful: false,
                    loading: false,
                    isInProgress: false,
                });
            });
    }

    handlePipeline = async() => {
        this.props.dispatch(clearMessage());

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            if (this.state.item === null) {
                return this.props
                    .dispatch(
                        addPipeline(
                          this.state.name,
                          this.state.folderUuid,
                          this.state.organization.uuid,
                          this.state.elements
                        )
                    )
                    .then(() => {
                        this.setState({
                            successful: true,
                            loading: false,
                        });

                        var item = this.props.message.item;

                        if (item) {
                            item.updated_at = this.generateUpdatedAt();
                            this.setState({
                                item,
                                isInProgress: false,
                            });

                            this.props.handleSetActiveItem(item);
                        }

                        this.props.dispatch(clearMessage());

                        return item;
                    })
                    .catch(() => {
                        this.setState({
                            successful: false,
                            loading: false,
                            isInProgress: false,
                        });
                    });
                } else {
                  return this.props
                        .dispatch(
                            updatePipeline(
                                this.state.item.uuid,
                                this.state.name,
                                this.state.item.folder_uuid,
                                this.state.elements,
                                this.state.item.schedule
                            )
                        )
                    .then(() => {
                        this.setState({
                            successful: true,
                            loading: false,
                        });

                        var item = this.state.item;
                        item.name = this.state.name;
                        item.elements = this.state.elements;
                        item.preview = this.state.preview;
                        item.folder_uuid = this.state.folderUuid;
                        item.elements_layout = JSON.stringify(this.state.elements);
                        item.updated_at = this.generateUpdatedAt();
                        
                        this.setState({
                            item,
                            isInProgress: false,
                        });

                        this.props.handleSetActiveItem(item);

                        this.props.dispatch(clearMessage());

                        return item;
                    })
                    .catch(() => {
                        this.setState({
                            successful: false,
                            loading: false,
                            isInProgress: false,
                        });
                    });
                }
            } else {
                this.setState({
                    loading: false,
                    isInProgress: false,
                });
            }
    }

    generateUpdatedAt = () => {
        return DateTime.local({ zone: "utc" }).toSQL({ includeOffset: false });
    };

    captureCanvas = async () => {
        if (
            this.state.isPipelineRunning === true || 
            this.state.isOutputVisible === true
        ) {
            return;
        }

        let base64 = await html2canvas(
            document.getElementsByClassName('react-flow__nodes')[0],
            {backgroundColor:null}
        )
            .then(canvas => 
                {
                  return canvas.toDataURL('image/png', 0.1)
                }
            )

        await Promise.resolve(base64)
            .then(async(base64) => {
                const base64Response = await fetch(base64);
                const blob = await base64Response.blob();
                PipelineService.updatePipelinePreview(this.state.item.uuid, blob);
            });
    }; 

    openTaskForm = async(item) => {
        let taskTrigger = null;
        let elements = this.state.elements;

        for (var i = 0; i < elements.length; i++) {
            if (
              elements[i].source === item.id
              && elements[i].trigger_type === TRIGGER_TYPE_SCHEDULE
            ) {
                taskTrigger = elements[i];
            }
        }

        await this.promisedSetState({
            isTaskFormOpen: true,
            activeTaskItem: item,
            activeTaskTrigger: taskTrigger,
        });
    };

    closeTaskForm = () => {
        this.setState({
            isTaskFormOpen: false,
            activeTaskItem: null,
            activeTaskTrigger: null,
        });
    };

    openRunningForm = async(showLogWithoutRunning = false) => {
        await this.promisedSetState({
            isRunningFormOpen: true,
            showLogWithoutRunning: showLogWithoutRunning,
        });

        Array.from(document.querySelectorAll(".runningInfo")).forEach(function(e){
            e.removeAttribute('class');
            e.classList.add('runningInfo');
        });

        document.querySelectorAll(".runningTextVisible").forEach(function(e){
            e.classList.remove('runningTextVisible');
        });
    };

    closeRunningForm = async(timeoutError = false, isExternalFailure = false) => {
        await this.promisedSetState({
            isRunningFormOpen: false,
        });

        if (isExternalFailure === true) {
            this.setState({isPipelineFailed: true});
        }

        let runningTasks = document.querySelectorAll("." + TASK_RUNNING);

        if (runningTasks.length > 0) {
            for(var i = 0; i < runningTasks.length; i++) {
                if (isExternalFailure === true) {
                    this.changeRunningClass(runningTasks[i], TASK_FAILED);
                    runningTasks[i].parentElement.querySelector(".runningText .runningTextTitle").textContent = TEXT_TASK_FAILED;
                    if (timeoutError === true) {
                        runningTasks[i].parentElement.querySelector(".runningText .runningTextInside").innerHTML = TEXT_TASK_FAILED_TIMEOUT;
                    } else {
                        runningTasks[i].parentElement.querySelector(".runningText .runningTextInside").innerHTML = TEXT_TASK_FAILED_PIPELINE_CANNOT_RUN;
                    }
                } else {
                    this.changeRunningClass(runningTasks[i], TASK_PAUSED);
                    runningTasks[i].parentElement.querySelector(".runningText .runningTextTitle").textContent = TEXT_TASK_PAUSED; 
                }
            }
        }
    };

    openRunningOutput = async() => {
        await this.promisedSetState({
            isRunningOutputOpen: true,
        });
    };

    closeRunningOutput = () => {
        this.setState({
            isRunningOutputOpen: false,
        });
    };

    openScheduleForm = async() => {
        await this.promisedSetState({
            isScheduleFormOpen: true,
        });
    };

    closeScheduleForm = () => {
        this.setState({
            isScheduleFormOpen: false,
        });
    };

    handlePipelineRunningSuccess = async() => {
        this.closeRunningForm()
    };

    handleScheduleFormSuccess = async(newSchedule) => {
        let item = this.state.item;
        item.schedule = newSchedule;

        await this.promisedSetState({item});

        this.closeScheduleForm()
    };

    handleTaskFormSuccess = async(node) => {
        if (node && node.uuid) {
          var elements = this.state.elements;

          for(var i = 0; i < elements.length; i++){
              if (elements[i].id === node.uuid) {
                  elements[i].data.name = node.name;
              }
          }

          await this.setElements(elements);
          await this.handlePipeline();
          this.captureCanvas();

          /*if (node.type === TASK_TYPE_QUERY) {
              let taskTrigger = null;
              for (var i = 0; i < elements.length; i++) {
                  if (
                    elements[i].source === node.uuid
                    && elements[i].trigger_type === TRIGGER_TYPE_SCHEDULE
                  ) {
                      taskTrigger = elements[i];
                  }
              }

              if (taskTrigger === null) {
                  // create scheduled trigger
                  this.onConnect({
                      target: 0, 
                      source: node.uuid, 
                      schedule: node.schedule,
                  });
              } else {
                  // update scheduled trigger
                  this.onUpdateScheduledConnection(
                      taskTrigger,
                      node.schedule,
                  );
              }
          }*/
        }
        this.closeTaskForm()
    };

    render() {
        const { message } = this.props;
        const { 
            elements, 
            isInProgress, 
            isPipelineFailed, 
            isRunningOutputOpen, 
            isPipelineRunning, 
            oldOutput, 
            showLogWithoutRunning, 
            isOutputVisible, 
            loading, 
            item, 
            isScheduleFormOpen, 
            isRunningFormOpen, 
            isTaskFormOpen, 
            activeTaskItem, 
            activeTaskTrigger 
        } = this.state;

      return (
        <ReactFlow
          elements={elements}
          nodeTypes={nodeTypes}
          onElementsRemove={this.onElementsRemove}
          onConnect={this.onConnect}
          onNodeDragStop={this.onNodeDragStop}
          onLoad={this.onLoad}
          onDrop={this.onDrop}
          onNodeDoubleClick={this.onNodeDoubleClick}
          onDragOver={this.onDragOver}
          onConnectStart={this.onConnectStart}
          snapToGrid={true}
          snapGrid={[15, 15]}
          defaultZoom={5}
          ref={this.reactFlowWrapper}
        >
          <div className="modal-title pipelineTitle">
            <Form
              onSubmit={this.handleForm}
              ref={(c) => {
                this.form = c;
              }}
            >
            <div className="float-left">
            <Input
              className="form-control form-control-lg"
              type="text"
              placeholder={this.state.name}
              autoComplete="name"
              name="name"
              value={this.state.name}
              onChange={this.onChangeName}
              onBlur={this.handleForm}
              validations={[required]}
            />
            </div>
            {message 
                && isInProgress === false 
                && this.state.successful === false 
                && (
              <div className="form-group float-left">
                <div className={ this.state.successful ? "alert alert-success" : "alert alert-danger" } role="alert">
                  {message.item ? "Pipeline successfully created" : message}
                </div>
              </div>
            )}
            
              {loading 
                ? <div className="pipelineNameIcon float-left"><span className="spinner-border float-left spinner-border-sm spinner-primary"></span></div>
                : item !== null &&
                  <>
                  <span className="note clearfix px-3 pt-1 float-left">
                    {item.updated_at !== "" &&
                        "Last modified: " + TimeAgo(DateTime.fromSQL(item.updated_at, { zone: 'utc'}))
                    }
                  </span>
                  </>
              }
            
            <div className="clearfix"></div>
            <CheckButton
              style={{ display: "none" }}
              ref={(c) => {
                this.checkBtn = c;
              }}
            />
          </Form>
          </div>
          <Controls showInteractive={false}/>
          <MiniMap 
              nodeStrokeColor={(n) => {
                if (nodeColors[n.type]) return nodeColors[n.type];

                return '#eee';
              }}
              nodeColor={(n) => {
                if (nodeColors[n.type]) return nodeColors[n.type];

                return '#fff';
              }}
              nodeBorderRadius={4}
          />
          <Background color="#aaa" gap={8} />
          <PipelineSidebar 
            onAdd={this.onAdd}
            isDraggingAllowed={!isInProgress}
          />

          <TaskModal
            show={isTaskFormOpen}
            content={
              <TaskForm 
                item={activeTaskItem}
                trigger={activeTaskTrigger}  
                pipeline={item}
                successHandler={this.handleTaskFormSuccess}
              />
            }
            item={activeTaskItem}
            title={activeTaskItem !== null 
                ? activeTaskItem.type !== TASK_TYPE_RUN_PIPELINE
                    ? "Edit " + activeTaskItem.type
                    : "Edit run_pipeline" 
                : "Edit task"}
            onHide={this.closeTaskForm}
          />

          <RightModal
            show={isRunningOutputOpen}
            content={
              <PipelineRunOutput
                output={oldOutput}
                finished={!isPipelineRunning}
              />
            }
            item={item}
            title={item !== null ? "Pipeline '" + item.name + "' is running" : "Pipeline is running"}
            onHide={this.closeRunningOutput}
          />

          {
            isRunningFormOpen
            && <PipelineRun 
                item={item} 
                elements={elements}
                outputHandler={this.handleRunningOutput}
                showLogWithoutRunning={showLogWithoutRunning}
                oldOutput={oldOutput}
              />
          }

          <TaskModal
            show={isScheduleFormOpen}
            content={
                <ScheduleForm 
                    item={item}
                    successHandler={this.handleScheduleFormSuccess}
                />
            }
            item={item}
            title={"Edit pipeline running schedule"}
            onHide={this.closeScheduleForm}
          />

          { item !== null && elements.length > 0 &&
            <ButtonGroup className="mr-4 runningButton">
                {
                    isOutputVisible === false
                    && isPipelineRunning === false
                    && 
                    <Button 
                    variant="light" 
                    onClick={this.openScheduleForm}
                    >
                        Schedule pipeline
                    </Button>
                }
                { isOutputVisible === true
                    && isPipelineRunning === false 
                    &&
                    <Button 
                        variant="secondary" 
                        onClick={this.cleanOutput}
                    >
                        Clear output
                    </Button>
                }
                { isPipelineRunning === true 
                    &&
                    <Button 
                        variant="secondary" 
                        onClick={() => this.closeRunningForm()}
                    >
                        Stop pipeline
                    </Button>
                }
                { isOutputVisible === true 
                    &&
                    <Button 
                        variant="light" 
                        onClick={() => this.openRunningOutput()}
                    >
                        {
                            isPipelineFailed === true
                            && <WarningAmberOutlined className="warningLogIcon" fontSize="small" sx={{ color: "#fff" }}/>
                        }
                        Show output log
                    </Button>
                }
                <Button 
                    variant="primary" 
                    onClick={() => this.openRunningForm(false)}
                    disabled={isPipelineRunning}
                >
                    {isPipelineRunning && (
                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                    )}
                    Run pipeline
                </Button>
            </ButtonGroup>
          }
          { /*<div className="inputTip pipelineTip">Double click on the task to edit it</div> */ }
        </ReactFlow>
      );
  }
};

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    const { message } = state.message;
    return {
        isLoggedIn,
        message,
        user,
    };
}

export default connect(mapStateToProps)(Pipeline);
