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
import QueryUI from "../../formControls/queryUI.component";
import { required } from "../../formControls/validations";

import IntegrationService, { INTEGRATION_TYPE_SQL } from "../../../services/integration.service";

import { TASK_SEVERITY_MEDIUM } from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

class QueryForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeSource = this.onChangeSource.bind(this);
        this.onChangeSQLQuery = this.onChangeSQLQuery.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            severity: TASK_SEVERITY_MEDIUM,
            SQLQuery: "",
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            sources: null,
            sourceUuid: null,
            sourceName: null,
            sourceType: null,
            sourceValue: null,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
                SQLQuery: this.props.item.implementation.sql_query || "",
                sourceUuid: this.props.item.implementation.source_uuid || "",
            });
            this.handleGetSources(this.state.organization.uuid);
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
                .then(sources => {
                    if (sources.data) {
                        this.setState({sources: sources.data.items});
                        if (sources.data.items.length > 0) {
                            if (this.props.item.implementation.source_uuid === "") {
                                this.setState({
                                    sourceUuid: sources.data.items[0].uuid,
                                    sourceName: sources.data.items[0].name,
                                    sourceType: sources.data.items[0].type,
                                    sourceValue: sources.data.items[0].value,
                                })
                            } else {
                                let source = sources.data.items.find(o => o.uuid === this.props.item.implementation.source_uuid);
                                if (source) {
                                    this.setState({
                                        sourceUuid: source.uuid,
                                        sourceName: source.name,
                                        sourceType: source.type,
                                        sourceValue: source.value,
                                    })
                                }
                            }
                        }
                    } else {
                        this.setState({sources: []});
                    }
                });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
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

    onChangeSeverity(severity) {
        this.setState({ severity })
    }

    onChangeSQLQuery(query) {
        this.setState({
            SQLQuery: query,
        });
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
                        "sql_query": this.state.SQLQuery,
                        "source_uuid": this.state.sourceUuid,
                    }
                )
            )
            .then(() => {
                var item = this.state.item;
                item.name = this.state.name;
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

                    <span className="inputLabel">Source</span><br/>
                    {this.state.sources !== null ?
                        <Dropdown size="lg" className="mb-4">
                            <Dropdown.Toggle 
                                className={"dropdownItemWithBg dropdownItemWithBg-" + (this.state.sourceType === INTEGRATION_TYPE_SQL ? this.state.sourceValue : this.state.sourceType) }
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

                    {this.state.sources !== null
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
                                {typeof message === 'string' && message}
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

export default connect(mapStateToProps)(QueryForm);
