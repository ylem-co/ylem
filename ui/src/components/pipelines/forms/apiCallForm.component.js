import React, { Component } from "react";
import { Navigate } from 'react-router-dom';

import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";

import FloatingLabel from "react-bootstrap/FloatingLabel";
import Spinner from "react-bootstrap/Spinner";

import CodeEditor from '@uiw/react-textarea-code-editor';
import rehypePrism from 'rehype-prism-plus';

import { connect } from "react-redux";
import { updateTask } from "../../../actions/tasks";

import { clearMessage } from "../../../actions/message";

import Input from "../../formControls/input.component";
import { required } from "../../formControls/validations";
import { TextareaEditor } from "../../formControls/textareaEditor.component";

import IntegrationService, { INTEGRATION_TYPE_API } from "../../../services/integration.service";

import { TASK_SEVERITY_MEDIUM } from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

class ApiCallForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeIntegration = this.onChangeIntegration.bind(this);
        this.onChangePayload = this.onChangePayload.bind(this);
        this.onChangePayloadFromOutside = this.onChangePayloadFromOutside.bind(this);
        this.onChangeQueryString = this.onChangeQueryString.bind(this);
        this.onChangeQueryStringFromOutside = this.onChangeQueryStringFromOutside.bind(this);
        this.onChangeHeaders = this.onChangeHeaders.bind(this);
        this.onChangeHeadersFromOutside = this.onChangeHeadersFromOutside.bind(this);
        this.onChangeAttachedFileName = this.onChangeAttachedFileName.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            severity: TASK_SEVERITY_MEDIUM,
            payload: "",
            queryString: "",
            headers: "",
            attachedFileName: "",
            item: this.props.item,
            task: null,
            loading: false,
            successful: false,
            integrations: null,
            integrationUuid: null,
            integrationName: null,
            integrationType: null,
        };
    }

    componentDidMount() {
        this.props.dispatch(clearMessage());

        if (this.props.item !== null) {
            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: this.props.item.severity || TASK_SEVERITY_MEDIUM,
                payload: this.props.item.implementation.payload || "",
                queryString: this.props.item.implementation.query_string || "",
                headers: this.props.item.implementation.headers || "",
                attachedFileName: this.props.item.implementation.attached_file_name || "",
            });
            this.handleGetIntegrations(this.state.organization.uuid);
        }
    };

    handleGetIntegrations = async(uuid) => {
        let integrations = this.state.integrations;

        if (
            integrations === null
            || integrations.length === 0
        ) {
            integrations = IntegrationService.getIntegrationsByOrganization(uuid);

            Promise.resolve(integrations)
                .then(integrations => {
                    if (integrations.data) {
                        let ds = [];
                        for(var i = 0; i < integrations.data.items.length; i++){
                            if (integrations.data.items[i].type === INTEGRATION_TYPE_API) {
                                ds.push(integrations.data.items[i]);
                            }
                        }

                        this.setState({integrations: ds});

                        if (ds.length > 0) {
                            if (this.props.item.implementation.destination_uuid === "") {
                                this.setState({
                                    integrationUuid: ds[0].uuid,
                                    integrationName: ds[0].name,
                                    integrationType: ds[0].type,
                                })
                            } else {
                                let integration = ds.find(o => o.uuid === this.props.item.implementation.destination_uuid);
                                if (integration) {
                                    this.setState({
                                        integrationUuid: integration.uuid,
                                        integrationName: integration.name,
                                        integrationType: integration.type,
                                    })
                                }
                            }
                        }
                    } else {
                        this.setState({integrations: []});
                    }
                });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangePayload(e) {
        this.setState({
            payload: e.target.value,
        });
    }

    onChangePayloadFromOutside(el) {
        this.setState({
            payload: el.value,
        });
    }

    onChangeQueryString(e) {
        this.setState({
            queryString: e.target.value,
        });
    }

    onChangeQueryStringFromOutside(el) {
        this.setState({
            queryString: el.value,
        });
    }

    onChangeHeaders(e) {
        this.setState({
            headers: e.target.value,
        });
    }

    onChangeHeadersFromOutside(el) {
        this.setState({
            headers: el.value,
        });
    }

    onChangeAttachedFileName(e) {
        this.setState({
            attachedFileName: e.target.value,
        });
    }

    onChangeIntegration(integrationUuid) {
        let integration = this.state.integrations.find(o => o.uuid === integrationUuid);
        if (integration) {
            this.setState({
                integrationUuid: integration.uuid,
                integrationName: integration.name,
                integrationType: integration.type,
            })
        }
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
                        "type": INTEGRATION_TYPE_API,
                        "destination_uuid": this.state.integrationUuid,
                        "payload": this.state.payload,
                        "headers": this.state.headers,
                        "query_string": this.state.queryString,
                        "attached_file_name": this.state.attachedFileName,
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

                    <span className="inputLabel">Integration</span><br/>
                    {this.state.integrations !== null ?
                        <Dropdown size="lg" className="mb-4">
                            <Dropdown.Toggle 
                                className={"dropdownItemWithBg dropdownItemWithBg-" + this.state.integrationType}
                                variant="light" 
                                id="dropdown-basic"
                            >
                                {this.state.integrationName}
                            </Dropdown.Toggle>

                            <Dropdown.Menu>
                                {this.state.integrations.map(value => (
                                    <Dropdown.Item
                                        className={"dropdownItemWithBg dropdownItemWithBg-" + value.type}
                                        value={value.name}
                                        key={value.uuid}
                                        active={value.uuid === this.state.integrationUuid}
                                        onClick={(e) => this.onChangeIntegration(value.uuid)}
                                    >
                                        {value.name}
                                    </Dropdown.Item>
                                ))}
                            </Dropdown.Menu>
                        </Dropdown>
                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                    }

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <span className="inputLabel">Payload (JSON)</span>
                            <CodeEditor
                                    className="form-control form-control-lg codeEditor"
                                    type="text"
                                    language="json"
                                    id="floatingPayload"
                                    minHeight={200}
                                    autoComplete="payload"
                                    name="payload"
                                    value={this.state.payload}
                                    onChange={this.onChangePayload}
                                    rehypePlugins={[
                                        [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                    ]}
                                    style={{
                                        fontSize: 14,
                                        fontFamily: 'Source Code Pro, monospace',
                                    }}
                                />
                            <div className="inputTip">
                                Body payload for POST requests.
                                <br/><br/>Use brackets &#123;&#123; place_it_here &#125;&#125; to use input data, aggregation functions or environment variables. 
                                <br/>E.g &#123;&#123; field_name &#125;&#125;, &#123;&#123; SUM(amount) &#125;&#125;, or &#123;&#123; ENV_variable_name &#125;&#125;
                            </div>
                            <TextareaEditor 
                                txtId="floatingPayload"
                                callback={this.onChangePayloadFromOutside}
                                brackets={true}
                            />
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <span className="inputLabel">Query string</span>
                            <CodeEditor
                                    className="form-control form-control-lg codeEditor"
                                    type="text"
                                    language="json"
                                    id="floatingQueryString"
                                    minHeight={200}
                                    autoComplete="queryString"
                                    name="queryString"
                                    value={this.state.queryString}
                                    onChange={this.onChangeQueryString}
                                    rehypePlugins={[
                                        [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                    ]}
                                    style={{
                                        fontSize: 14,
                                        fontFamily: 'Source Code Pro, monospace',
                                    }}
                                />
                            <div className="inputTip">
                                Query string for GET requests.
                                <br/><br/>Use brackets &#123;&#123; place_it_here &#125;&#125; to use input data, aggregation functions or environment variables. 
                                <br/>E.g &#123;&#123; field_name &#125;&#125;, &#123;&#123; SUM(amount) &#125;&#125;, or &#123;&#123; ENV_variable_name &#125;&#125;
                            </div>
                            <TextareaEditor 
                                txtId="floatingQueryString"
                                callback={this.onChangeQueryStringFromOutside} 
                                brackets={true}
                            />
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <span className="inputLabel">Request headers (optional)</span>
                            <CodeEditor
                                    className="form-control form-control-lg codeEditor"
                                    type="text"
                                    language="json"
                                    id="floatingHeaders"
                                    minHeight={200}
                                    autoComplete="headers"
                                    name="headers"
                                    value={this.state.headers}
                                    onChange={this.onChangeHeaders}
                                    rehypePlugins={[
                                        [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                    ]}
                                    style={{
                                        fontSize: 14,
                                        fontFamily: 'Source Code Pro, monospace',
                                    }}
                                />
                            <div className="inputTip">
                                Request headers in JSON.
                                <br/><br/>Use brackets &#123;&#123; place_it_here &#125;&#125; to use input data, aggregation functions or environment variables. 
                                <br/>E.g &#123;&#123; field_name &#125;&#125;, &#123;&#123; SUM(amount) &#125;&#125;, or &#123;&#123; ENV_variable_name &#125;&#125;
                            </div>
                            <TextareaEditor 
                                txtId="floatingHeaders" 
                                callback={this.onChangeHeadersFromOutside}
                                brackets={true}
                            />
                        </div>
                    </InputGroup>

                    <InputGroup className="mb-4">
                        <div className="registrationFormControl">
                            <FloatingLabel controlId="floatingAttachedFileName" label="Attached filename">
                                <Input
                                    className="form-control form-control-lg"
                                    type="text"
                                    id="floatingAttachedFileName"
                                    placeholder="Attached filename"
                                    autoComplete="attachedFileName"
                                    name="attachedFileName"
                                    value={this.state.attachedFileName}
                                    onChange={this.onChangeAttachedFileName}
                                />
                                <div className="inputTip">In case you want to attach a filename, give it a name</div>
                            </FloatingLabel>
                        </div>
                    </InputGroup>

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

export default connect(mapStateToProps)(ApiCallForm);
