import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate } from 'react-router-dom';

import Spinner from "react-bootstrap/Spinner";
import Alert from "react-bootstrap/Alert";
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../formControls/input.component";
import Textarea from "../formControls/textarea.component";
import { required, isTrialHost } from "../formControls/validations";

import VisibilityOutlined from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlined from '@mui/icons-material/VisibilityOffOutlined';

import Tooltip from '@mui/material/Tooltip';

import { connect } from "react-redux";

import { addIntegration, updateIntegration, testNewIntegration, testExistingIntegration } from "../../actions/integrations";

import {
    Col,
    InputGroup,
    Row,
    Dropdown,
} from 'react-bootstrap'

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import IntegrationService, {
    INTEGRATION_TYPE_SQL,
    INTEGRATION_TYPE_SNOWFLAKE,
    INTEGRATION_TYPE_GOOGLE_BIG_QUERY,
    INTEGRATION_TYPE_POSTGRESQL,
    INTEGRATION_TYPE_IMMUTA,
    INTEGRATION_TYPE_REDSHIFT,
    INTEGRATION_TYPE_ELASTICSEARCH,
    SQL_INTEGRATION_CONNECTION_TYPES,
    SQL_INTEGRATION_CONNECTION_TYPE_DIRECT,
    SQL_INTEGRATION_SSL_MODES,
    SQL_INTEGRATION_TYPE_PORT_MAP,
    ELASTICSEARCH_VERSIONS,
    INTEGRATION_IO_TYPE_READ_WRITE,
} from "../../services/integration.service";

class SQLIntegrationForm extends Component {
    constructor(props) {
        super(props);
        this.handleCreateSQLIntegration = this.handleCreateSQLIntegration.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeType = this.onChangeType.bind(this);
        this.onChangeHost = this.onChangeHost.bind(this);
        this.onChangePort = this.onChangePort.bind(this);
        this.onChangeUser = this.onChangeUser.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
        this.onChangeDatabase = this.onChangeDatabase.bind(this);
        this.onChangeConnectionType = this.onChangeConnectionType.bind(this);
        this.onChangeSshHost = this.onChangeSshHost.bind(this);
        this.onChangeSshPort = this.onChangeSshPort.bind(this);
        this.onChangeSshUser = this.onChangeSshUser.bind(this);
        this.onChangeSslMode = this.onChangeSslMode.bind(this);
        this.onChangeProjectId = this.onChangeProjectId.bind(this);
        this.onChangeCredentials = this.onChangeCredentials.bind(this);
        this.onChangeEsVersion = this.onChangeEsVersion.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            item: this.props.item,
            name: null,
            type: 
                this.props.item
                    ? this.props.item.type
                    : this.props.type,
            ioType: this.props.ioType || INTEGRATION_IO_TYPE_READ_WRITE,
            host: "",
            port: 3306,
            user: "",
            password: "",
            passwordType: "password",
            database: "",
            connectionType: SQL_INTEGRATION_CONNECTION_TYPE_DIRECT,
            sshHost: "",
            sshPort: 22,
            sshUser: "",
            projectId: "",
            credentials: "",
            esVersion: null,
            sslMode: false,
            successful: false,
            loading: false,
            isInProgress: false,
        };
    };

    componentDidMount() {
        if (this.props.item !== null) {
            this.handleGetItem(
                this.props.item.uuid,
                this.props.item.type
            );
        } else {
            this.setState({
                type: this.props.type
            });
        }
    };

    UNSAFE_componentWillReceiveProps(props) {
        if (
            props.item === null
        ) {
            this.setState({
                ioType: props.ioType,
                type: props.type,
            });
        } else {
            this.setState({
                ioType: props.ioType,
                type: props.item.type !== INTEGRATION_TYPE_SQL ? props.item.type : props.type,
            });
        }
    };

    handleGetItem(uuid, type) {
        var item = IntegrationService.getIntegration(uuid, type);

        Promise.resolve(item)
            .then(item => {
                this.mapItemToForm(item);
            })
            .catch(() => {
            });
    };

    mapItemToForm(item) {
        this.setState({
            name: item.data.integration.name,
            type: item.data.type,
            ioType: item.data.integration.io_type,
            host: item.data.host,
            port: item.data.port,
            user: item.data.user,
            database: item.data.database,
            connectionType: item.data.connection_type,
            sshHost: item.data.ssh_host,
            sshPort: item.data.ssh_port,
            sshUser: item.data.ssh_user,
            sslMode: item.data.ssl_enabled || false,
            projectId: item.data.project_id || "",
            credentials: item.data.credentials || "",
            esVersion: item.data.es_version || null,
        });
    };

    mapFormToItem(operation = "save") {
        let data = {
            "organization_uuid": this.state.organization.uuid,
            "type": this.state.type,
        };

        if (operation === "save") {
            data.name = this.state.name;

            if (this.state.type === INTEGRATION_TYPE_ELASTICSEARCH) {
                data.es_version = this.state.esVersion
            }
        }

        if (this.state.type === INTEGRATION_TYPE_GOOGLE_BIG_QUERY) {
            data.project_id = this.state.projectId;
            data.credentials = this.state.credentials;
        } else {
            data.host = this.state.host;
            data.port = parseInt(this.state.port);
            data.user = this.state.user;
            data.password = this.state.password;
            data.database = this.state.database;
            data.connection_type = this.state.connectionType;
            data.ssh_host = this.state.sshHost;
            data.ssh_port = this.state.sshPort !== "" ? parseInt(this.state.sshPort) : "";
            data.ssh_user = this.state.sshUser;
            data.ssl_enabled = this.state.sslMode;
        }

        return data;
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeProjectId(e) {
        this.setState({
            projectId: e.target.value,
        });
    }

    onChangeCredentials(e) {
        this.setState({
            credentials: e.target.value,
        });
    }

    onChangeEsVersion(esVersion) {
        this.setState({esVersion});
    }

    onChangeHost(e) {
        this.setState({
            host: e.target.value,
        });
    }

    onChangePort(e) {
        this.setState({
            port: e.target.value,
        });
    }

    onChangeUser(e) {
        this.setState({
            user: e.target.value,
        });
    }

    onChangePassword(e) {
        this.setState({
            password: e.target.value,
        });
    }

    onChangeDatabase(e) {
        this.setState({
            database: e.target.value,
        });
    }

    onChangeType(type) {
        this.setState({type});

        if (this.state.item !== null) {
            return
        }

        this.setState({
            port: SQL_INTEGRATION_TYPE_PORT_MAP[type]
        })
    }

    onChangeConnectionType(connectionType) {
        this.setState({connectionType});
    }

    onChangeSslMode(sslMode) {
        this.setState({sslMode});
    }

    onChangeSshHost(e) {
        this.setState({
            sshHost: e.target.value,
        });
    }

    onChangeSshPort(e) {
        this.setState({
            sshPort: e.target.value,
        });
    }

    onChangeSshUser(e) {
        this.setState({
            sshUser: e.target.value,
        });
    }

    handleEyeClick = () => this.setState(({passwordType}) => ({
        passwordType: passwordType === 'text' ? 'password' : 'text'
    }));

    handleCreateSQLIntegration(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            var data = this.mapFormToItem();

            if (this.state.item === null) {
                this.props
                    .dispatch(
                        addIntegration(
                            INTEGRATION_TYPE_SQL,
                            data,
                            this.state.type
                        )
                    )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: false
                    });
                    setTimeout(() => {
                        this.props.successHandler()
                    },2000);
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: false,
                    });
                });
            } else {
                this.props
                    .dispatch(
                        updateIntegration(
                            this.state.item.uuid, 
                            INTEGRATION_TYPE_SQL,  
                            data,
                            this.state.type
                        )
                    )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: false,
                    });
                    setTimeout(() => {
                        this.props.successHandler()
                    },2000);
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: false,
                    });
                });
            }
        } else {
            this.setState({
                loading: false,
            });
        }
    }

    handleTestConnection(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            var data = this.mapFormToItem("test");

            if (this.state.item === null) {
                this.props
                    .dispatch(
                        testNewIntegration(
                            this.state.type, 
                            data
                        )
                    )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: false
                    });
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: false,
                    });
                });
            } else {
                this.props
                    .dispatch(
                        testExistingIntegration(
                            this.state.item.uuid,  
                            this.state.type, 
                            data
                        )
                    )
                .then(() => {
                    this.setState({
                        successful: true,
                        loading: false,
                    });
                })
                .catch(() => {
                    this.setState({
                        successful: false,
                        loading: false,
                    });
                });
            }
        } else {
            this.setState({
                loading: false,
            });
        }
    };

    render() {
        const { isLoggedIn, message, user } = this.props;

        const { isInProgress, type } = this.state;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div>
                        {
                            (this.state.item === null
                                || (this.state.item !== null && this.state.name !== null)
                            )
                            ?
                            <div>
                                    <Form
                                        onSubmit={this.handleCreateSQLIntegration}
                                        ref={(c) => {
                                            this.form = c;
                                        }}
                                    >

                                    <InputGroup className="mb-4">
                                            <div className="formControl100">
                                            <FloatingLabel controlId="floatingName" label="Name">
                                            <Input
                                                className="form-control form-control-lg"
                                                id="floatingName"
                                                type="text"
                                                placeholder="Name"
                                                autoComplete="name"
                                                name="name"
                                                value={this.state.name}
                                                onChange={this.onChangeName}
                                                autoFocus
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>


                                    { type !== INTEGRATION_TYPE_GOOGLE_BIG_QUERY ?
                                    <>
                                        <Row>
                                            <Col xs="9">
                                                <div className="mb-4">
                                                    <div className="formControl100">
                                                    <FloatingLabel controlId="floatingHost" label={this.state.type === INTEGRATION_TYPE_SNOWFLAKE ? "Account identifier" : "Host"}>
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingHost"
                                                        type="text"
                                                        placeholder={this.state.type === INTEGRATION_TYPE_SNOWFLAKE ? "Account identifier" : "Host"}
                                                        autoComplete="host"
                                                        name="host"
                                                        value={this.state.host}
                                                        onChange={this.onChangeHost}
                                                    />
                                                    </FloatingLabel>
                                                    </div>
                                                </div>
                                            </Col>
                                            <Col xs="3">
                                                <InputGroup className="mb-4">
                                                    <div className="formControl100">
                                                    <FloatingLabel controlId="floatingPort" label="Port">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingPort"
                                                        type="text"
                                                        placeholder="Port"
                                                        autoComplete="port"
                                                        name="port"
                                                        value={this.state.port}
                                                        onChange={this.onChangePort}
                                                    />
                                                    </FloatingLabel>
                                                    </div>
                                                </InputGroup>
                                            </Col>
                                        </Row>
                                        <Row>
                                            <Col xs="6">
                                                <div className="mb-4">
                                                    <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingUser" label="User">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingUser"
                                                        type="text"
                                                        placeholder="User"
                                                        autoComplete="user"
                                                        name="user"
                                                        value={this.state.user}
                                                        onChange={this.onChangeUser}
                                                    />
                                                    </FloatingLabel>
                                                    </div>
                                                </div>
                                            </Col>
                                            <Col xs="6">
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingPassword" label="Password">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingPassword"
                                                        type={this.state.passwordType}
                                                        placeholder="Password"
                                                        autoComplete="password"
                                                        name="password"
                                                        value={this.state.password}
                                                        isCreation={this.state.item === null}
                                                        onChange={this.onChangePassword}
                                                    />
                                                    </FloatingLabel>
                                                    </div>
                                                    <span
                                                        onClick={this.handleEyeClick}
                                                        className="eye"
                                                    >
                                                        {
                                                            this.state.passwordType === 'text' 
                                                            ? <Tooltip title="Hide" placement="right"><VisibilityOffOutlined/></Tooltip>
                                                            : <Tooltip title="Show" placement="right"><VisibilityOutlined/></Tooltip>
                                                        }
                                                    </span>
                                                </InputGroup>
                                            </Col>
                                        </Row>

                                        <InputGroup className="mb-4">
                                            {(this.state.type === INTEGRATION_TYPE_POSTGRESQL || this.state.type === INTEGRATION_TYPE_REDSHIFT) &&
                                                <div className="formControl100">
                                                <FloatingLabel controlId="floatingDatabase" label="Database">
                                                <Input
                                                    className="form-control form-control-lg"
                                                    id="floatingDatabase"
                                                    type="text"
                                                    placeholder="Database"
                                                    autoComplete="database"
                                                    name="database"
                                                    value={this.state.database}
                                                    onChange={this.onChangeDatabase}
                                                    validations={[required]}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            }
                                            {(this.state.type !== INTEGRATION_TYPE_POSTGRESQL && this.state.type !== INTEGRATION_TYPE_REDSHIFT) &&
                                                <div className="formControl100">
                                                <FloatingLabel controlId="floatingDatabase" label="Database (optional)">
                                                <Input
                                                    className="form-control form-control-lg"
                                                    id="floatingDatabase"
                                                    type="text"
                                                    placeholder="Database (optional)"
                                                    autoComplete="database"
                                                    name="database"
                                                    value={this.state.database}
                                                    onChange={this.onChangeDatabase}
                                                />
                                                </FloatingLabel>
                                                </div>
                                            }
                                        </InputGroup>

                                        {(this.state.type === INTEGRATION_TYPE_POSTGRESQL || this.state.type === INTEGRATION_TYPE_IMMUTA) &&
                                            <>
                                        <span className="inputLabel">SSL mode</span><br/>
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle
                                                variant="light" 
                                                id="dropdown-basic"
                                            >
                                                {Object.keys(SQL_INTEGRATION_SSL_MODES).find(key => SQL_INTEGRATION_SSL_MODES[key] === this.state.sslMode)}
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu>
                                                {Object.keys(SQL_INTEGRATION_SSL_MODES).map((key, index) => (
                                                    <Dropdown.Item
                                                        value={SQL_INTEGRATION_SSL_MODES[key]}
                                                        key={key}
                                                        active={SQL_INTEGRATION_SSL_MODES[key] === this.state.sslMode}
                                                        onClick={(e) => this.onChangeSslMode(SQL_INTEGRATION_SSL_MODES[key])}
                                                    >
                                                        {key}
                                                    </Dropdown.Item>
                                                ))}
                                            </Dropdown.Menu>
                                        </Dropdown>
                                        </>
                                        }

                                        {(this.state.type === INTEGRATION_TYPE_ELASTICSEARCH && this.state.item !== null) &&
                                            <>
                                        <span className="inputLabel">Elasticsearch version</span><br/>
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle
                                                variant="light"
                                                id="dropdown-basic"
                                            >
                                                {this.state.esVersion}
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu>
                                                {ELASTICSEARCH_VERSIONS.map((key, value) => (
                                                    <Dropdown.Item
                                                        value={key}
                                                        key={key}
                                                        active={value === this.state.sslMode}
                                                        onClick={(e) => this.onChangeEsVersion(key)}
                                                    >
                                                        {key}
                                                    </Dropdown.Item>
                                                ))}
                                            </Dropdown.Menu>
                                        </Dropdown>
                                        </>
                                        }

                                        {(this.state.type !== INTEGRATION_TYPE_ELASTICSEARCH && this.state.type !== INTEGRATION_TYPE_REDSHIFT) &&
                                            <div>
                                        <span className="inputLabel">Connection type</span><br/>
                                            <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle
                                            variant="light"
                                            id="dropdown-basic"
                                            >
                                            {this.state.connectionType.toUpperCase()}
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu>
                                            {SQL_INTEGRATION_CONNECTION_TYPES.map(type => (
                                                <Dropdown.Item
                                                    value={type}
                                                    key={type}
                                                    active={type === this.state.type}
                                                    onClick={(e) => this.onChangeConnectionType(e.target.getAttribute('value'))}
                                                >
                                                    {type.toUpperCase()}
                                                </Dropdown.Item>
                                            ))}
                                            </Dropdown.Menu>
                                            </Dropdown>
                                            </div>
                                        }

                                        {(this.state.type === INTEGRATION_TYPE_REDSHIFT) &&
                                        <Alert variant="info" className="mt-4">
                                            Contact us if you need SSH for Redshift.
                                        </Alert>
                                        }

                                        {this.state.connectionType === "ssh" &&
                                            <div>
                                                <Row>
                                                    <Col xs="9">
                                                        <div className="mb-4">
                                                            <div className="formControl100">
                                                                <FloatingLabel controlId="floatingSshHost" label="SSH host">
                                                                    <Input
                                                                        className="form-control form-control-lg"
                                                                        id="floatingSshHost"
                                                                        type="text"
                                                                        placeholder="SSH host"
                                                                        autoComplete="sshHost"
                                                                        name="sshHost"
                                                                        value={this.state.sshHost}
                                                                        onChange={this.onChangeSshHost}
                                                                    />
                                                                </FloatingLabel>
                                                            </div>
                                                        </div>
                                                    </Col>
                                                    <Col xs="3">
                                                        <InputGroup className="mb-4">
                                                            <div className="formControl100">
                                                                <FloatingLabel controlId="floatingSshPort" label="SSH port">
                                                                    <Input
                                                                        className="form-control form-control-lg"
                                                                        id="floatingSshPort"
                                                                        type="text"
                                                                        placeholder="SSH port"
                                                                        autoComplete="sshPort"
                                                                        name="sshPort"
                                                                        value={this.state.sshPort}
                                                                        onChange={this.onChangeSshPort}
                                                                    />
                                                                </FloatingLabel>
                                                            </div>
                                                        </InputGroup>
                                                    </Col>
                                                </Row>
                                                <InputGroup className="mb-4">
                                                    <div className="formControl100">
                                                        <FloatingLabel controlId="floatingSshUser" label="SSH user">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingSshUser"
                                                                type="text"
                                                                placeholder="SSH user"
                                                                autoComplete="sshUser"
                                                                name="sshUser"
                                                                value={this.state.sshUser}
                                                                onChange={this.onChangeSshUser}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                            </div>
                                        }
                                    </>
                                    :
                                    <>
                                        <InputGroup className="mb-4">
                                            <div className="formControl100">
                                                <FloatingLabel controlId="floatingProjectId" label="Project ID (Optional)">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingProjectId"
                                                        type="text"
                                                        placeholder="Project ID (Optional)"
                                                        autoComplete="projectID"
                                                        name="projectID"
                                                        value={this.state.projectId}
                                                        onChange={this.onChangeProjectId}
                                                    />
                                                </FloatingLabel>
                                            </div>
                                        </InputGroup>

                                        <InputGroup className="mb-4">
                                            <div className="formControl100">
                                                <FloatingLabel controlId="floatingCredentials" label="Credentials (JSON)">
                                                    <Textarea
                                                        className="form-control form-control-lg codeEditor"
                                                        id="floatingCredentials"
                                                        type="textarea"
                                                        placeholder="Credentials (JSON)"
                                                        autoComplete="credentials"
                                                        name="credentials"
                                                        value={this.state.credentials}
                                                        onChange={this.onChangeCredentials}
                                                    />
                                                </FloatingLabel>
                                                <div className="inputTip clearfix">
                                                    <a href="https://cloud.google.com/docs/authentication/getting-started#creating_a_service_account" target="_blank" rel="noreferrer">
                                                        How to create a Google service account
                                                    </a>. You need to grant it "BigQuery Data Viewer" and "BigQuery Job User" permissions. 
                                                </div>
                                            </div>
                                        </InputGroup>
                                    </>
                                    }

                                        <Row className="pt-3">
                                            { !isTrialHost(this.state.host) &&
                                            <Col xs="3">
                                                <button
                                                    className="px-4 btn btn-primary"
                                                    disabled={this.state.loading}
                                                    type="submit"
                                                >
                                                    {this.state.loading && (
                                                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                    )}
                                                    <span>{this.state.item !== null ? "Save" : "Create"}</span>
                                                </button>
                                            </Col>
                                            }
                                            <Col xs="9">
                                                <button
                                                    className="px-4 btn btn-secondary"
                                                    disabled={this.state.loading}
                                                    onClick={(e) => this.handleTestConnection(e)}
                                                    type="button"
                                                >
                                                    {this.state.loading && (
                                                        <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                                    )}
                                                    <span>Test connection</span>
                                                </button>
                                            </Col>
                                        </Row>
                                        {
                                            message && isInProgress && (
                                            <div className="form-group">
                                                <div className={ this.state.successful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                                    {   
                                                        typeof message === 'object' 
                                                        && message !== null
                                                        && 'item' in message
                                                        ? 
                                                            (
                                                                this.state.item === null
                                                                    ? "Integration successfully created"
                                                                    : "Integration successfully updated"
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
                            </div>
                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                        }
            </div>
        );
    }
}

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

export default connect(mapStateToProps)(SQLIntegrationForm);
