import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";
import Spinner from "react-bootstrap/Spinner";

import { connect } from "react-redux";
import { updateTask } from "../../../actions/tasks";

import { clearMessage } from "../../../actions/message";

import Input from "../../formControls/input.component";
import { required } from "../../formControls/validations";

import { TASK_SEVERITY_MEDIUM } from "../../../services/task.service";
import PipelineService, { PIPELINE_TYPE_GENERIC } from "../../../services/pipeline.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

class PipelineRunForm extends Component {
    constructor(props) {
        super(props);
        this.handleGetPipelines = this.handleGetPipelines.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeActivePipeline = this.onChangeActivePipeline.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: this.props.item.name,
            severity: TASK_SEVERITY_MEDIUM,
            pipelines: null,
            activePipelineUuid: "",
            activePipelineName: "",
        };
    }

    componentDidMount = async() => {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            await this.promisedSetState({
                item: this.props.item,
                name: this.props.item.name || "",
                pipelines: null,
                activePipelineUuid: this.props.item.implementation.pipeline_uuid || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
            });
            this.handleGetPipelines(this.state.organization.uuid);
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

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
                            k => k.type === PIPELINE_TYPE_GENERIC
                                && k.uuid !== this.state.item.pipeline_uuid
                        );
                        
                        if (
                            items.length > 0
                            && this.state.activePipelineUuid !== ""
                        ) {
                            var wf = pipelines.data.items.find(k => k.uuid === this.state.activePipelineUuid);
                            if (wf !== undefined) {
                               this.setState({
                                    activePipelineUuid: wf.uuid,
                                    activePipelineName: wf.name,
                               }); 
                            }
                        } else if (
                            items.length > 0
                            && this.state.activePipelineUuid === ""
                        ) {
                            this.setState({
                                activePipelineUuid: items[0].uuid,
                                activePipelineName: items[0].name,
                            });
                        }

                        this.setState({pipelines: items});
                    } else {
                        this.setState({pipelines: []});
                    }
                });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeActivePipeline(pipeline) {
        this.setState({
            activePipelineUuid: pipeline.uuid,
            activePipelineName: pipeline.name,
        });
    }

    onChangeSeverity(severity) {
        this.setState({ severity })
    }

    handleSubmit(e) {
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
                    this.state.item.uuid, 
                    this.state.item.pipeline_uuid,
                    this.state.name,
                    this.state.severity,
                    this.state.item.type,
                    {
                        'pipeline_uuid': this.state.activePipelineUuid,
                    }
                )
            )
            .then(() => {
                var item = this.state.item;
                item.name = this.state.name;
                item.implementation.pipeline_uuid = this.state.activePipelineUuid;

                this.setState({
                    loading: false,
                    successful: true,
                    item,
                });

                this.props.successHandler(item);
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

    render() {
        const { isLoggedIn, user, message } = this.props;

        const { activePipelineUuid, activePipelineName } = this.state;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div>
                <Form
                    onSubmit={this.handleSubmit}
                    ref={(c) => {
                        this.form = c;
                    }}
                >
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

                    <span className="inputLabel">Run pipeline</span><br/>
                    {this.state.pipelines !== null &&
                        this.state.pipelines.length > 0 ? 
                        <>
                            <Dropdown size="lg" className="mb-4">
                                <Dropdown.Toggle 
                                    variant="light" 
                                    id="dropdown-basic"
                                    className="mt-2"
                                >
                                { activePipelineUuid !== "" &&
                                    <>
                                        {activePipelineName}
                                    </>
                                }
                                </Dropdown.Toggle>

                                <Dropdown.Menu className="tasks">
                                    {this.state.pipelines.map(value => (
                                        <Dropdown.Item
                                            value={value.name}
                                            key={value.uuid}
                                            active={this.state.activePipelineUuid !== "" && value.uuid === this.state.activePipelineUuid}
                                            onClick={(e) => this.onChangeActivePipeline(value)}
                                        >
                                            {value.name}
                                        </Dropdown.Item>
                                    ))}
                                </Dropdown.Menu>
                            </Dropdown>
                        </>
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }

                    <Row>
                        <Col xs="6">
                            <Button
                                className="px-4 btn btn-primary"
                                disabled={this.state.loading}
                                type="submit"
                            >
                                {this.state.loading && (
                                    <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                )}
                                <span>Save</span>
                            </Button>
                        </Col>
                    </Row>
                    {message && (
                        <div className="form-group">
                            <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                {message}
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

export default connect(mapStateToProps)(PipelineRunForm);
