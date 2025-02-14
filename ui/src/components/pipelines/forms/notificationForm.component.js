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
import { TextareaEditor } from "../../formControls/textareaEditor.component";
import { required, integrationBody } from "../../formControls/validations";

import IntegrationService, { 
    INTEGRATION_TYPE_HUBSPOT,
    INTEGRATION_TYPE_INCIDENT_IO, 
    INTEGRATION_TYPE_JENKINS,
    INTEGRATION_TYPE_OPSGENIE,
    INTEGRATION_TYPE_SALESFORCE,
    INTEGRATION_TYPE_WHATSAPP,
    INTEGRATION_IO_TYPE_WRITE,
} from "../../../services/integration.service";

import {
    TASK_SEVERITY_CRITICAL,
    TASK_SEVERITY_HIGH,
    TASK_SEVERITY_LOW,
    TASK_SEVERITY_LOWEST,
    TASK_SEVERITY_MEDIUM
} from "../../../services/task.service";

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../../actions/pipeline";

import {
    Button,
    Col,
    Dropdown,
    InputGroup,
    Row
} from 'react-bootstrap'

const DEFAULT_SEVERITIES = [
    {
        "id": TASK_SEVERITY_LOWEST,
        "name": "Lowest"
    },
    {
        "id": TASK_SEVERITY_LOW,
        "name": "Low"
    },
    {
        "id": TASK_SEVERITY_MEDIUM,
        "name": "Medium"
    },
    {
        "id": TASK_SEVERITY_HIGH,
        "name": "High"
    },
    {
        "id": TASK_SEVERITY_CRITICAL,
        "name": "Critical"
    },
];

const HUBSPOT_SEVERITIES = [
    {
        "id": "LOW",
        "name": "Low"
    },
    {
        "id": "MEDIUM",
        "name": "Medium"
    },
    {
        "id": "HIGH",
        "name": "High"
    },
];

const SALESFORCE_SEVERITIES = [
    {
        "id": "Low",
        "name": "Low"
    },
    {
        "id": "Medium",
        "name": "Medium"
    },
    {
        "id": "High",
        "name": "High"
    },
];

const OPSGENIE_SEVERITIES = [
    {
        "id": "P5",
        "name": "Informational"
    },
    {
        "id": "P4",
        "name": "Low"
    },
    {
        "id": "P3",
        "name": "Moderate"
    },
    {
        "id": "P2",
        "name": "High"
    },
    {
        "id": "P1",
        "name": "Critical"
    },
];

class NotificationForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeIntegration = this.onChangeIntegration.bind(this);
        this.onChangeBody = this.onChangeBody.bind(this);
        this.onChangeBodyFromOutside = this.onChangeBodyFromOutside.bind(this);
        this.onChangeAttachedFileName = this.onChangeAttachedFileName.bind(this);
        this.onChangeSeverity = this.onChangeSeverity.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            name: "",
            severity: TASK_SEVERITY_MEDIUM,
            severityName: "Medium",
            severities: DEFAULT_SEVERITIES,
            body: "",
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
            let svr = this.props.item.severity || TASK_SEVERITY_MEDIUM
            let svrObj = DEFAULT_SEVERITIES.find(x => x.id === svr);
            let svrName = "Select"
            if (svrObj !== undefined) {
                svrName = svrObj.name
            }

            this.setState({
                item: this.props.item,
                name: this.props.item.name || "",
                severity: svr,
                severityName: svrName,
                body: this.props.item.implementation.body || "",
                attachedFileName: this.props.item.implementation.attached_file_name || "",
            }, async () => {
                return this.handleGetIntegrations(this.state.organization.uuid).then(ds => {
                    let type = this.props.item.implementation.type
                    let uuid = this.props.item.implementation.destination_uuid

                    if (ds.length > 0) {
                        if (type === "") {
                            type = ds[0].type
                        }

                        if (uuid === "" ) {
                            uuid = ds[0].uuid
                        }
                    }

                    if (type === INTEGRATION_TYPE_INCIDENT_IO) {
                        this.handleGetIncidentIoSeverities(uuid)
                    }

                    if (type === INTEGRATION_TYPE_HUBSPOT) {
                        this.handleGetHubspotSeverities()
                    }

                    if (type === INTEGRATION_TYPE_SALESFORCE) {
                        this.handleGetSalesforceSeverities()
                    }

                    if (type === INTEGRATION_TYPE_OPSGENIE) {
                        this.handleGetOpsgenieSeverities()
                    }
                });
            });
        }
    };

    handleGetIncidentIoSeverities = async(uuid) => {
        let severities = IntegrationService.getIncidentIoIntegrationSeverities(uuid);
        let svr = this.state.severity
        this.setState({
            severityName: "Loading..."
        })

        Promise.resolve(severities)
            .then(severities => {
                if (severities.data && severities.data.items) {
                    let svs = [];
                    for (let i = 0; i < severities.data.items.length; i++) {
                        svs.push(severities.data.items[i]);
                    }

                    let svrObj = svs.find(x => x.id === svr);
                    let svrName = "Select"
                    if (svrObj !== undefined) {
                        svrName = svrObj.name
                        svr = svrObj.id
                    } else {
                        svr = ""
                    }

                    this.setState({
                        severities: svs,
                        severity: svr,
                        severityName: svrName
                    })
                } else {
                    console.error("did not get severities")
                }
            });
    }

    handleGetHubspotSeverities() {
        let svs = HUBSPOT_SEVERITIES
        let svr = this.state.severity
        let svrObj = svs.find(x => x.id === svr);
        let svrName = "Select"
        if (svrObj !== undefined) {
            svrName = svrObj.name
            svr = svrObj.id
        } else {
            svr = ""
        }

        this.setState({
            severities: svs,
            severity: svr,
            severityName: svrName
        })
    }

    handleGetSalesforceSeverities() {
        let svs = SALESFORCE_SEVERITIES
        let svr = this.state.severity
        let svrObj = svs.find(x => x.id === svr);
        let svrName = "Select"
        if (svrObj !== undefined) {
            svrName = svrObj.name
            svr = svrObj.id
        } else {
            svr = ""
        }

        this.setState({
            severities: svs,
            severity: svr,
            severityName: svrName
        })
    }

    handleGetOpsgenieSeverities() {
        let svs = OPSGENIE_SEVERITIES
        let svr = this.state.severity
        let svrObj = svs.find(x => x.id === svr);
        let svrName = "Select"
        if (svrObj !== undefined) {
            svrName = svrObj.name
            svr = svrObj.id
        } else {
            svr = ""
        }

        this.setState({
            severities: svs,
            severity: svr,
            severityName: svrName
        })
    }

    handleGetIntegrations = async(uuid) => {
        let integrations = this.state.integrations;

        if (
            integrations === null
            || integrations.length === 0
        ) {
            integrations = IntegrationService.getIntegrationsByOrganization(uuid, INTEGRATION_IO_TYPE_WRITE);

            return Promise.resolve(integrations)
                .then(integrations => {
                    let ds = []

                    if (integrations.data) {
                        ds = integrations.data.items

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

                    return ds
                });
        }
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeBody(e) {
        this.setState({
            body: e.target.value,
        });
    }

    onChangeBodyFromOutside(el) {
        this.setState({
            body: el.value,
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

            if (integration.type === INTEGRATION_TYPE_INCIDENT_IO) {
                this.handleGetIncidentIoSeverities(integration.uuid)
            } else if (integration.type === INTEGRATION_TYPE_OPSGENIE) {
                this.handleGetOpsgenieSeverities()
            } else if (integration.type === INTEGRATION_TYPE_JENKINS) {
                this.setState({
                    severityName: "Medium",
                })
            } else {
                let severities = DEFAULT_SEVERITIES
                if (integration.type === INTEGRATION_TYPE_HUBSPOT) {
                    severities = HUBSPOT_SEVERITIES
                }
                if (integration.type === INTEGRATION_TYPE_SALESFORCE) {
                    severities = SALESFORCE_SEVERITIES
                }

                let svr = this.state.severity
                let svrObj = severities.find(x => x.id === svr);
                let svrName = "Select"
                if (svrObj !== undefined) {
                    svrName = svrObj.name
                    svr = svrObj.id
                } else {
                    svr = ""
                }

                this.setState({
                    severities: severities,
                    severity: svr,
                    severityName: svrName
                })
            }
        }
    }

    onChangeSeverity(severity) {
        let svrObj = this.state.severities.find(x => x.id === severity);
        this.setState({
            severity,
            severityName: svrObj.name
        })
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
                        "type": this.state.integrationType,
                        "destination_uuid": this.state.integrationUuid,
                        "body": this.state.body,
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

                    <Row>
                        <Col xs="5">
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
                        </Col>
                        <Col xs="7">

                    {this.state.integrationType !== "jenkins" &&
                    <>
                        <span className="inputLabel">Severity</span><br/>
                        <Dropdown size="lg" className="mb-4">
                            <Dropdown.Toggle
                                variant="light"
                                id="dropdown-basic"
                            >
                                <div className={"dropdownSeverity dropdownSeverity_" + this.state.severity + " dropdownSeverity_" + this.state.severityName.toLowerCase()}></div>
                                {this.state.severityName}
                            </Dropdown.Toggle>

                            <Dropdown.Menu>
                                {this.state.severities.map(value => (
                                    <Dropdown.Item
                                        value={value.name}
                                        key={value.id}
                                        active={value.id.toLowerCase() === this.state.severity.toLowerCase()}
                                        onClick={(e) => this.onChangeSeverity(value.id)}
                                    >
                                        <div
                                            className={"dropdownSeverity dropdownSeverity_" + value.id + " dropdownSeverity_" + value.name.toLowerCase()}></div>
                                        {value.name}
                                    </Dropdown.Item>
                                ))}
                            </Dropdown.Menu>
                        </Dropdown>
                    </>
                    }

                        </Col>
                    </Row>

                    <div className={this.state.integrationType === "jenkins" && "hiddenHeader"}>
                        <InputGroup className="mb-4">
                            <div className="registrationFormControl">
                                <label className="nonFloatingLabel">{this.state.integrationType !== INTEGRATION_TYPE_WHATSAPP ? "Message" : "Content (JSON)" }</label>
                                    <CodeEditor
                                        className="form-control form-control-lg codeEditor"
                                        type="text"
                                        language="xls"
                                        id="floatingBody"
                                        minHeight={200}
                                        autoComplete="message"
                                        name="message"
                                        value={this.state.body}
                                        onChange={this.onChangeBody}
                                        integration={this.state.integrationType}
                                        validations={[integrationBody]}
                                        rehypePlugins={[
                                            [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                                        ]}
                                        style={{
                                            fontSize: 14,
                                            fontFamily: 'Source Code Pro, monospace',
                                        }}
                                    />
                                    {
                                        this.state.integrationType === INTEGRATION_TYPE_WHATSAPP
                                        &&
                                        <div className="inputTip">
                                            Place here JSON in the following format:<br/>
                                            {'{'}<br/>
                                            &nbsp;&nbsp;&nbsp;&nbsp;"phoneTo": &#123;&#123; phone &#125;&#125;,<br/>
                                            &nbsp;&nbsp;&nbsp;&nbsp;"contentVariables": "&#123;\"1\":\"12/1\",\"2\":\"3pm\"&#125;"<br/>
                                            {'}'}
                                        </div>
                                    }
                                    <div className="inputTip">
                                        Use brackets &#123;&#123; place_it_here &#125;&#125; to use input data,
                                        aggregation functions or environment variables.
                                        <br/>E.g &#123;&#123; field_name &#125;&#125;, &#123;&#123; SUM(amount) &#125;&#125;,
                                        or &#123;&#123; ENV_variable_name &#125;&#125;
                                    </div>
                                    <TextareaEditor 
                                        txtId="floatingBody"
                                        callback={this.onChangeBodyFromOutside} 
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
                                    <div className="inputTip">In case you want to attach a filename, give it a name
                                    </div>
                                </FloatingLabel>
                            </div>
                        </InputGroup>
                    </div>

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

export default connect(mapStateToProps)(NotificationForm);
