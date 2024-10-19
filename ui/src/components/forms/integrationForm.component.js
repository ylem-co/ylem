import React, { Component } from "react";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import { Navigate } from 'react-router-dom';

import Spinner from "react-bootstrap/Spinner";
import Alert from "react-bootstrap/Alert";
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "../formControls/input.component";
import {required, requiredForCreation} from "../formControls/validations";

import VisibilityOutlined from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlined from '@mui/icons-material/VisibilityOffOutlined';

import Tooltip from '@mui/material/Tooltip';

import { connect } from "react-redux";
import { addIntegration, updateIntegration, confirmIntegration } from "../../actions/integrations";

import SQLIntegrationForm from "./sqlIntegrationForm.component";

import {
    Card,
    Col,
    Container,
    InputGroup,
    Row,
    Dropdown,
} from 'react-bootstrap'

import { validatePermissions, PERMISSION_LOGGED_OUT } from "../../actions/pipeline";

import {
    INTEGRATION_TYPES_PER_IO_FORM,
    INTEGRATION_TYPES_PER_IO,
    INTEGRATION_IO_TYPES,
    INTEGRATION_IO_TYPE_READ_WRITE,
    INTEGRATION_TYPE_API,
    INTEGRATION_TYPE_SMS,
    INTEGRATION_TYPE_EMAIL,
    INTEGRATION_TYPE_SLACK,
    INTEGRATION_TYPE_API_AUTH_TYPES,
    INTEGRATION_TYPE_API_AUTH_TYPE_BASIC,
    INTEGRATION_TYPE_API_AUTH_TYPE_BEARER,
    INTEGRATION_TYPE_API_AUTH_TYPE_HEADER,
    INTEGRATION_TYPE_API_AUTH_TYPE_NONE,
    INTEGRATION_TYPE_JIRA,
    INTEGRATION_TYPE_INCIDENT_IO,
    INTEGRATION_TYPE_TO_HUMAN,
    INTEGRATION_API_METHODS,
    INTEGRATION_API_METHOD_POST,
    INTEGRATION_TYPE_TABLEAU,
    INTEGRATION_TYPE_HUBSPOT,
    INTEGRATION_TYPE_SALESFORCE,
    INTEGRATION_TYPE_GOOGLE_SHEETS,
    INTEGRATION_TYPE_OPSGENIE, 
    INTEGRATION_TYPE_JENKINS,
    INTEGRATION_TYPE_SQL,
    SQL_INTEGRATION_TYPES,
} from "../../services/integration.service";

import IntegrationService from "../../services/integration.service";
import Textarea from "../formControls/textarea.component";

// noinspection ES6ConvertVarToLetConst
class IntegrationForm extends Component {
    constructor(props) {
        super(props);
        this.handleRegister = this.handleRegister.bind(this);
        this.handleSQLFormAfterSuccess = this.handleSQLFormAfterSuccess.bind(this);
        this.onChangeName = this.onChangeName.bind(this);
        this.onChangeType = this.onChangeType.bind(this);
        this.onChangeIoType = this.onChangeIoType.bind(this);
        this.onChangeApiMethod = this.onChangeApiMethod.bind(this);
        this.onChangeAuthType = this.onChangeAuthType.bind(this);
        this.onChangeAuthBearerToken = this.onChangeAuthBearerToken.bind(this);
        this.onChangeAuthHeaderName = this.onChangeAuthHeaderName.bind(this);
        this.onChangeAuthHeaderValue = this.onChangeAuthHeaderValue.bind(this);
        this.onChangeAuthBasicUserName = this.onChangeAuthBasicUserName.bind(this);
        this.onChangeAuthBasicPassword = this.onChangeAuthBasicPassword.bind(this);
        this.onChangeUrl = this.onChangeUrl.bind(this);
        this.onChangeNumber = this.onChangeNumber.bind(this);
        this.onChangeEmail = this.onChangeEmail.bind(this);
        this.onChangeCode = this.onChangeCode.bind(this);
        this.onChangeChannel = this.onChangeChannel.bind(this);
        this.onChangeSlackAuthorization = this.onChangeSlackAuthorization.bind(this);
        this.onChangeJiraAuthorization = this.onChangeJiraAuthorization.bind(this);
        this.onChangeJiraIssueType = this.onChangeJiraIssueType.bind(this);
        this.onChangeJiraProjectKey = this.onChangeJiraProjectKey.bind(this);
        this.onChangeIncidentIoApiKey = this.onChangeIncidentIoApiKey.bind(this);
        this.onChangeOpsgenieApiKey = this.onChangeOpsgenieApiKey.bind(this);
        this.onChangeIncidentIoMode = this.onChangeIncidentIoMode.bind(this);
        this.onChangeIncidentIoVisibility = this.onChangeIncidentIoVisibility.bind(this);
        this.onChangeTableauServer = this.onChangeTableauServer.bind(this);
        this.onChangeTableauUsername = this.onChangeTableauUsername.bind(this);
        this.onChangeTableauPassword = this.onChangeTableauPassword.bind(this);
        this.onChangeTableauSiteName = this.onChangeTableauSiteName.bind(this);
        this.onChangeTableauProjectName = this.onChangeTableauProjectName.bind(this);
        this.onChangeTableauDatasourceName = this.onChangeTableauDatasourceName.bind(this);
        this.onChangeTableauMode = this.onChangeTableauMode.bind(this);
        this.onChangeHubspotAuthorization = this.onChangeHubspotAuthorization.bind(this);
        this.onChangeSalesforceAuthorization = this.onChangeSalesforceAuthorization.bind(this);
        this.onChangeHubspotPipelineCode = this.onChangeHubspotPipelineCode.bind(this);
        this.onChangeHubspotPipelineStageCode = this.onChangeHubspotPipelineStageCode.bind(this);
        this.onChangeHubspotOwnerCode = this.onChangeHubspotOwnerCode.bind(this);
        this.onChangeGoogleSheetsCredentials = this.onChangeGoogleSheetsCredentials.bind(this);
        this.onChangeGoogleSheetsSpreadsheetId = this.onChangeGoogleSheetsSpreadsheetId.bind(this);
        this.onChangeGoogleSheetsSheetId = this.onChangeGoogleSheetsSheetId.bind(this);
        this.onChangeGoogleSheetsMode = this.onChangeGoogleSheetsMode.bind(this);
        this.onChangeGoogleSheetsWriteHeader = this.onChangeGoogleSheetsWriteHeader.bind(this);
        this.onChangeJenkinsBaseUrl = this.onChangeJenkinsBaseUrl.bind(this);
        this.onChangeJenkinsProjectName = this.onChangeJenkinsProjectName.bind(this);
        this.onChangeJenkinsToken = this.onChangeJenkinsToken.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            item: this.props.item,
            name: null,
            url: "",
            uuid: null,
            type: this.props.integrationType && this.props.integrationType in INTEGRATION_TYPES_PER_IO
                ? INTEGRATION_TYPES_PER_IO_FORM[this.props.integrationType][0]
                : INTEGRATION_TYPES_PER_IO_FORM[INTEGRATION_IO_TYPE_READ_WRITE][0],
            ioType: this.props.integrationType && this.props.integrationType in INTEGRATION_TYPES_PER_IO
                ? this.props.integrationType
                : INTEGRATION_IO_TYPE_READ_WRITE,
            apiMethod: INTEGRATION_API_METHOD_POST,
            number: "",
            incidentIoApiKey: "",
            incidentIoVisibility: "public",
            incidentIoMode: "test",
            opsgenieApiKey: "",
            email: "",
            code: "",

            channel: "",
            slackAuthorizationUuid: "",
            slackAuthorizationName: "",
            slackAuthorizations: null,

            jenkinsBaseUrl: "",
            jenkinsToken: "",
            jenkinsProjectName: "",

            jiraAuthorizationUuid: "",
            jiraAuthorizationName: "",
            jiraProjectKey: "",
            jiraIssueType: "",
            jiraAuthorizations: null,

            tableauServer: "",
            tableauUsername: "",
            tableauPassword: "",
            tableauSiteName: "",
            tableauProjectName: "",
            tableauDatasourceName: "",
            tableauMode: "overwrite",

            hubspotAuthorizationUuid: "",
            hubspotAuthorizationName: "",
            hubspotPipelineCode: "",
            hubspotPipelineStageCode: "",
            hubspotOwnerCode: "",
            hubspotAuthorizations: null,

            salesforceAuthorizationUuid: "",
            salesforceAuthorizationName: "",
            salesforceAuthorizations: null,

            googleSheetsCredentials: null,
            googleSheetsMode: "append",
            googleSheetsSpreadsheetId: null,
            googleSheetsSheetId: null,
            googleSheetsWriteHeader: "yes",

            authType: INTEGRATION_TYPE_API_AUTH_TYPE_NONE,
            authBearerToken: "",
            authHeaderName: "",
            authHeaderValue: "",
            authBasicUserName: "",
            authBasicPassword: "",
            authBasicPasswordType: "password",
            successful: false,
            loading: false,
            isInProgress: false,
            isWaitingForConfirmation: false,
        };
    };

    componentDidMount() {
        if (this.props.item !== null) {
            this.handleGetItem(
                this.props.item.uuid,
                this.props.item.type,
            );
        }

        this.handleGetSlackAuthorizations(this.state.organization.uuid);
        this.handleGetJiraAuthorizations(this.state.organization.uuid);
        this.handleGetHubspotAuthorizations(this.state.organization.uuid);
        this.handleGetSalesforceAuthorizations(this.state.organization.uuid);
    };

    handleGetSlackAuthorizations(uuid) {
        var slackAuthorizations = IntegrationService.getSlackAuthorizations(uuid);

        Promise.resolve(slackAuthorizations)
            .then(slackAuthorizations => {
                var auths = slackAuthorizations.data.items.filter(k => k.is_active === true);

                this.setState({slackAuthorizations: auths});

                if (
                    auths.length > 0
                    && this.state.item === null
                ) {
                    this.setState({
                        slackAuthorizationUuid: auths[0].uuid,
                        slackAuthorizationName: auths[0].name
                    });
                } else if (this.state.item !== null) {
                    
                } else {
                    this.setState({
                        slackAuthorizationUuid: null,
                        slackAuthorizationName: "Please add Slack authorization"
                    });
                }
            })
            .catch(() => {
            });
    };

    handleGetJiraAuthorizations(uuid) {
        var jiraAuthorizations = IntegrationService.getJiraAuthorizations(uuid);

        Promise.resolve(jiraAuthorizations)
            .then(jiraAuthorizations => {
                var auths = jiraAuthorizations.data.items.filter(k => k.is_active === true);

                this.setState({jiraAuthorizations: auths});

                if (
                    auths.length > 0
                    && this.state.item === null
                ) {
                    this.setState({
                        jiraAuthorizationUuid: auths[0].uuid,
                        jiraAuthorizationName: auths[0].name
                    });
                } else if (this.state.item !== null) {

                } else {
                    this.setState({
                        jiraAuthorizationUuid: null,
                        jiraAuthorizationName: "Please add Jira Cloud authorization"
                    });
                }
            })
            .catch(() => {
            });
    };

    handleGetHubspotAuthorizations(uuid) {
        var hubspotAuthorizations = IntegrationService.getHubspotAuthorizations(uuid);

        Promise.resolve(hubspotAuthorizations)
            .then(hubspotAuthorizations => {
                var auths = hubspotAuthorizations.data.items.filter(k => k.is_active === true);

                this.setState({hubspotAuthorizations: auths});

                if (
                    auths.length > 0
                    && this.state.item === null
                ) {
                    this.setState({
                        hubspotAuthorizationUuid: auths[0].uuid,
                        hubspotAuthorizationName: auths[0].name
                    });
                } else if (this.state.item !== null) {

                } else {
                    this.setState({
                        hubspotAuthorizationUuid: null,
                        hubspotAuthorizationName: "Please add Hubspot authorization"
                    });
                }
            })
            .catch(() => {
            });
    };

    handleGetSalesforceAuthorizations(uuid) {
        var salesforceAuthorizations = IntegrationService.getSalesforceAuthorizations(uuid);

        Promise.resolve(salesforceAuthorizations)
            .then(salesforceAuthorizations => {
                var auths = salesforceAuthorizations.data.items.filter(k => k.is_active === true);

                this.setState({salesforceAuthorizations: auths});

                if (
                    auths.length > 0
                    && this.state.item === null
                ) {
                    this.setState({
                        salesforceAuthorizationUuid: auths[0].uuid,
                        salesforceAuthorizationName: auths[0].name
                    });
                } else if (this.state.item !== null) {

                } else {
                    this.setState({
                        salesforceAuthorizationUuid: null,
                        salesforceAuthorizationName: "Please add Salesforce authorization"
                    });
                }
            })
            .catch(() => {
            });
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
        if (item.data.integration.type === INTEGRATION_TYPE_API) {
            this.setState({
                url: item.data.integration.value,
                apiMethod: item.data.method !== "" ? item.data.method : INTEGRATION_API_METHOD_POST,
                authType: item.data.auth_type,
                authBearerToken: item.data.auth_bearer_token,
                authHeaderName: item.data.auth_header_name,
                authHeaderValue: item.data.auth_header_value,
                authBasicUserName: item.data.auth_basic_user_name,
                authBasicPassword: item.data.auth_basic_password,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_SMS) {
            this.setState({
                number: item.data.integration.value,
                isWaitingForConfirmation: !item.data.is_confirmed,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_EMAIL) {
            this.setState({
                email: item.data.integration.value,
                isWaitingForConfirmation: !item.data.is_confirmed,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_SLACK) {
            this.setState({
                channel: item.data.integration.value,
                slackAuthorizationUuid: item.data.authorization.uuid,
                slackAuthorizationName: item.data.authorization.name,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_JIRA) {
            this.setState({
                jiraProjectKey: item.data.integration.value,
                jiraIssueType: item.data.issue_type,
                jiraAuthorizationUuid: item.data.authorization.uuid,
                jiraAuthorizationName: item.data.authorization.name,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_INCIDENT_IO) {
            this.setState({
                incidentIoApiKey: item.data.api_key,
                incidentIoMode: item.data.integration.value,
                incidentIoVisibility: item.data.visibility,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_OPSGENIE) {
            this.setState({
                opsgenieApiKey: item.data.api_key,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_TABLEAU) {
            this.setState({
                tableauServer: item.data.integration.value,
                tableauUsername: item.data.username,
                tableauSiteName: item.data.site_name,
                tableauProjectName: item.data.project_name,
                tableauDatasourceName: item.data.datasource_name,
                tableauMode: item.data.mode,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_HUBSPOT) {
            this.setState({
                hubspotPipelineCode: item.data.integration.value,
                hubspotPipelineStageCode: item.data.pipeline_stage_code,
                hubspotOwnerCode: item.data.owner_code,
                hubspotAuthorizationUuid: item.data.authorization.uuid,
                hubspotAuthorizationName: item.data.authorization.name,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_SALESFORCE) {
            this.setState({
                salesforceAuthorizationUuid: item.data.authorization.uuid,
                salesforceAuthorizationName: item.data.authorization.name,
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_GOOGLE_SHEETS) {
            this.setState({
                googleSheetsCredentials: item.data.credentials,
                googleSheetsMode: item.data.mode,
                googleSheetsSpreadsheetId: item.data.spreadsheet_id,
                googleSheetsSheetId: item.data.sheet_id.toString(),
                googleSheetsWriteHeader: item.data.write_header ? "yes" : "no",
            });
        } else if (item.data.integration.type === INTEGRATION_TYPE_JENKINS) {
            this.setState({
                jenkinsProjectName: item.data.integration.value,
                jenkinsBaseUrl: item.data.base_url,
                jenkinsToken: item.data.token,
            });
        }

        this.setState({
            name: item.data.integration.name,
            type: item.data.integration.type === INTEGRATION_TYPE_SQL
                ? item.data.type
                : item.data.integration.type
            ,
            uuid: item.data.integration.uuid,
        });
    };

    mapFormToItem() {
        let data = {
            "name": this.state.name,
        };

        if (
            this.state.item === null
            && this.state.type !== INTEGRATION_TYPE_SLACK
            && this.state.type !== INTEGRATION_TYPE_JIRA
            && this.state.type !== INTEGRATION_TYPE_HUBSPOT
            && this.state.type !== INTEGRATION_TYPE_SALESFORCE
        ) {
            data.organization_uuid = this.state.organization.uuid;
        }

        if (this.state.type === INTEGRATION_TYPE_API) {
            data.url = this.state.url;
            data.auth_type = this.state.authType;
            data.method = this.state.apiMethod;

            if (this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_BEARER) {
                data.auth_bearer_token = this.state.authBearerToken;
            } else if (this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_HEADER) {
                data.auth_header_name = this.state.authHeaderName;
                data.auth_header_value = this.state.authHeaderValue;
            } else if (this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_BASIC) {
                data.auth_basic_user_name = this.state.authBasicUserName;
                data.auth_basic_password = this.state.authBasicPassword;
            }
        } else if (this.state.type === INTEGRATION_TYPE_SMS) {
            data.number = this.state.number;
        } else if (this.state.type === INTEGRATION_TYPE_EMAIL) {
            data.email = this.state.email;
        } else if (this.state.type === INTEGRATION_TYPE_SLACK) {
            data.authorization_uuid = this.state.slackAuthorizationUuid;
            data.channel = this.state.channel;
        } else if (this.state.type === INTEGRATION_TYPE_JIRA) {
            data.authorization_uuid = this.state.jiraAuthorizationUuid;
            data.issue_type = this.state.jiraIssueType;
            data.project_key = this.state.jiraProjectKey;
        } else if (this.state.type === INTEGRATION_TYPE_INCIDENT_IO) {
            data.api_key = this.state.incidentIoApiKey;
            data.mode = this.state.incidentIoMode;
            data.visibility = this.state.incidentIoVisibility;
        } else if (this.state.type === INTEGRATION_TYPE_OPSGENIE) {
            data.api_key = this.state.opsgenieApiKey;
        } else if (this.state.type === INTEGRATION_TYPE_TABLEAU) {
            data.server = this.state.tableauServer;
            data.username = this.state.tableauUsername;
            data.password = this.state.tableauPassword;
            data.site_name = this.state.tableauSiteName;
            data.project_name = this.state.tableauProjectName;
            data.datasource_name = this.state.tableauDatasourceName;
            data.mode = this.state.tableauMode;
        } else if (this.state.type === INTEGRATION_TYPE_HUBSPOT) {
            data.authorization_uuid = this.state.hubspotAuthorizationUuid;
            data.pipeline_code = this.state.hubspotPipelineCode;
            data.pipeline_stage_code = this.state.hubspotPipelineStageCode;
            data.owner_code = this.state.hubspotOwnerCode;
        } else if (this.state.type === INTEGRATION_TYPE_SALESFORCE) {
            data.authorization_uuid = this.state.salesforceAuthorizationUuid;
        } else if (this.state.type === INTEGRATION_TYPE_GOOGLE_SHEETS) {
            data.credentials = this.state.googleSheetsCredentials;
            data.mode = this.state.googleSheetsMode;
            data.spreadsheet_id = this.state.googleSheetsSpreadsheetId;
            data.sheet_id = parseInt(this.state.googleSheetsSheetId);
            data.write_header = this.state.googleSheetsWriteHeader === "yes";
        } else if (this.state.type === INTEGRATION_TYPE_JENKINS) {
            data.token = this.state.jenkinsToken;
            data.base_url = this.state.jenkinsBaseUrl;
            data.project_name = this.state.jenkinsProjectName;
        }

        return data;
    };

    onChangeName(e) {
        this.setState({
            name: e.target.value,
        });
    }

    onChangeChannel(e) {
        this.setState({
            channel: e.target.value,
        });
    }

    onChangeJiraProjectKey(e) {
        this.setState({
            jiraProjectKey: e.target.value,
        });
    }

    onChangeJiraIssueType(e) {
        this.setState({
            jiraIssueType: e.target.value,
        });
    }

    onChangeCode(e) {
        this.setState({
            code: e.target.value,
        });
    }

    onChangeNumber(e) {
        this.setState({
            number: e.target.value,
        });
    }

    onChangeIncidentIoApiKey(e) {
        this.setState({
            incidentIoApiKey: e.target.value,
        });
    }

    onChangeOpsgenieApiKey(e) {
        this.setState({
            opsgenieApiKey: e.target.value,
        });
    }

    onChangeJenkinsToken(e) {
        this.setState({
            jenkinsToken: e.target.value,
        });
    }

    onChangeJenkinsBaseUrl(e) {
        this.setState({
            jenkinsBaseUrl: e.target.value,
        });
    }

    onChangeJenkinsProjectName(e) {
        this.setState({
            jenkinsProjectName: e.target.value,
        });
    }

    onChangeIncidentIoMode(mode) {
        this.setState({
            incidentIoMode: mode,
        });
    }

    onChangeIncidentIoVisibility(visibility) {
        this.setState({
            incidentIoVisibility: visibility,
        });
    }

    onChangeEmail(e) {
        this.setState({
            email: e.target.value,
        });
    }

    onChangeTableauServer(e) {
        this.setState({
            tableauServer: e.target.value,
        });
    }

    onChangeTableauUsername(e) {
        this.setState({
            tableauUsername: e.target.value,
        });
    }

    onChangeTableauPassword(e) {
        this.setState({
            tableauPassword: e.target.value,
        });
    }

    onChangeTableauSiteName(e) {
        this.setState({
            tableauSiteName: e.target.value,
        });
    }

    onChangeTableauProjectName(e) {
        this.setState({
            tableauProjectName: e.target.value,
        });
    }

    onChangeTableauDatasourceName(e) {
        this.setState({
            tableauDatasourceName: e.target.value,
        });
    }

    onChangeTableauMode(tableauMode) {
        this.setState({
            tableauMode
        });
    }

    onChangeAuthBasicUserName(e) {
        this.setState({
            authBasicUserName: e.target.value,
        });
    }

    onChangeAuthBasicPassword(e) {
        this.setState({
            authBasicPassword: e.target.value,
        });
    }

    onChangeUrl(e) {
        this.setState({
            url: e.target.value,
        });
    }

    onChangeType(type, ioType) {
        this.setState({type, ioType});
    }

    onChangeIoType(ioType) {
        this.setState({
            ioType,
            type: INTEGRATION_TYPES_PER_IO[ioType][0],
        });
    }

    onChangeApiMethod(apiMethod) {
        this.setState({apiMethod});
    }

    onChangeSlackAuthorization(uuid, name) {
        this.setState({
            slackAuthorizationUuid: uuid,
            slackAuthorizationName: name,
        });
    }

    onChangeJiraAuthorization(uuid, name) {
        this.setState({
            jiraAuthorizationUuid: uuid,
            jiraAuthorizationName: name,
        });
    }

    onChangeHubspotAuthorization(uuid, name) {
        this.setState({
            hubspotAuthorizationUuid: uuid,
            hubspotAuthorizationName: name,
        });
    }

    onChangeSalesforceAuthorization(uuid, name) {
        this.setState({
            salesforceAuthorizationUuid: uuid,
            salesforceAuthorizationName: name,
        });
    }

    onChangeHubspotPipelineCode(e) {
        this.setState({
            hubspotPipelineCode: e.target.value,
        });
    }

    onChangeHubspotPipelineStageCode(e) {
        this.setState({
            hubspotPipelineStageCode: e.target.value,
        });
    }

    onChangeHubspotOwnerCode(e) {
        this.setState({
            hubspotOwnerCode: e.target.value,
        });
    }

    onChangeGoogleSheetsCredentials(e) {
        this.setState({
            googleSheetsCredentials: e.target.value,
        });
    }

    onChangeGoogleSheetsSpreadsheetId(e) {
        this.setState({
            googleSheetsSpreadsheetId: e.target.value,
        });
    }

    onChangeGoogleSheetsSheetId(e) {
        this.setState({
            googleSheetsSheetId: e.target.value,
        });
    }

    onChangeGoogleSheetsMode(mode) {
        this.setState({
            googleSheetsMode: mode,
        });
    }

    onChangeGoogleSheetsWriteHeader(header) {
        this.setState({
            googleSheetsWriteHeader: header,
        });
    }

    onChangeAuthBearerToken(e) {
        this.setState({
            authBearerToken: e.target.value,
        });
    }

    onChangeAuthHeaderName(e) {
        this.setState({
            authHeaderName: e.target.value,
        });
    }

    onChangeAuthHeaderValue(e) {
        this.setState({
            authHeaderValue: e.target.value,
        });
    }

    onChangeAuthType(authType) {
        this.setState({authType});
    }

    handleEyeClick = () => this.setState(({authBasicPasswordType}) => ({
        authBasicPasswordType: authBasicPasswordType === 'text' ? 'password' : 'text'
    }));

    handleSQLFormAfterSuccess() {
        this.props.successHandler();
    }

    handleRegister(e) {
        e.preventDefault();

        this.setState({
            successful: false,
            loading: true,
            isInProgress: true,
        });

        if (this.state.isWaitingForConfirmation === true) {
            if (this.state.code !== "") {
                this.props
                    .dispatch(
                        confirmIntegration(
                            this.state.type, 
                            this.state.code,
                            this.state.uuid
                        )
                    )
                    .then(() => {
                        if (this.state.type === INTEGRATION_TYPE_EMAIL) {
                            let user = JSON.parse(localStorage.getItem('user'));
                            user.is_email_confirmed = "1";
                            localStorage.setItem("user", JSON.stringify(user));
                        }

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
                this.setState({
                    loading: false,
                });
            }
        } else {
            this.form.validateAll();

            if (this.checkBtn.context._errors.length === 0) {
                var data = this.mapFormToItem();

                if (this.state.item === null) {
                    this.props
                        .dispatch(
                            addIntegration(
                                this.state.type, 
                                data
                            )
                        )
                    .then(() => {
                        this.setState({
                            successful: true,
                            loading: false
                        });

                        if (this.state.type !== INTEGRATION_TYPE_SMS && this.state.type !== INTEGRATION_TYPE_EMAIL) {
                            setTimeout(() => {
                                this.props.successHandler()
                            },2000);
                        } else {
                            if (this.props.message.item && this.props.message.item.is_confirmed === false) {
                                this.setState({
                                    "isWaitingForConfirmation": true,
                                    "item": this.props.message.item.integration,
                                    "uuid": this.props.message.item.integration.uuid,
                                });
                            } else {
                                setTimeout(() => {
                                    this.props.successHandler()
                                },2000);
                            }
                        }
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
                                this.state.type, 
                                data
                            )
                        )
                    .then(() => {
                        this.setState({
                            successful: true,
                            loading: false,
                        });
                        if (this.state.type !== INTEGRATION_TYPE_SMS && this.state.type !== INTEGRATION_TYPE_EMAIL) {
                            setTimeout(() => {
                                this.props.successHandler()
                            },2000);
                        } else {
                            if (this.props.message.item && this.props.message.item.is_confirmed === false) {
                                this.setState({
                                    "isWaitingForConfirmation": true,
                                    "item": this.props.message.item.integration,
                                    "uuid": this.props.message.item.integration.uuid,
                                });
                            } else {
                                setTimeout(() => {
                                    this.props.successHandler()
                                },2000);
                            }
                        }
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
    }

    render() {
        const { isLoggedIn, message, user } = this.props;

        const { isInProgress } = this.state;

        if (validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_OUT)) {
            return <Navigate to="/login" />;
        }

        return (
            <div className="align-items-center">
                <Container>
                    <Row className="justify-content-center">
                        <Col md="10" lg="9" xl="8">
                        {
                            (this.state.item === null
                                || (this.state.item !== null && this.state.name !== null)
                            )
                            ?
                            <Card className="onboardingCard noBorder mb-5">
                                <Card.Body className="p-4">
                                    <Row>
                                        <Col xs="5">
                                        { this.state.item === null ?
                                            <div>
                                                <span className="inputLabel">Integration type</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle 
                                                        className={"dropdownItemWithBg dropdownItemWithBg-" + this.state.type}
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                    >
                                                        {INTEGRATION_TYPE_TO_HUMAN[this.state.type]}
                                                    </Dropdown.Toggle>
                                                    <Dropdown.Menu>

                                                        {INTEGRATION_IO_TYPES.map(io_type => (
                                                            <div>
                                                                <div className="pt-2 pb-3">
                                                                    <span className="inputLabel">
                                                                        {io_type.charAt(0).toUpperCase() + io_type.slice(1)} integrations
                                                                    </span>
                                                                </div>
                                                                {INTEGRATION_TYPES_PER_IO_FORM[io_type].map(type => (
                                                                    <Dropdown.Item
                                                                        className={"dropdownItemWithBg dropdownItemWithBg-" + type} 
                                                                        value={type}
                                                                        key={type}
                                                                        active={type === this.state.type}
                                                                        onClick={() => this.onChangeType(type, io_type)}
                                                                    >
                                                                        {INTEGRATION_TYPE_TO_HUMAN[type]}
                                                                    </Dropdown.Item>
                                                                ))}
                                                            </div>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                            </div>
                                            :
                                            <div>
                                                <span className="inputLabel">Integration type</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle 
                                                        className={"dropdownItemWithBg queryFilter dropdownItemWithBg-" + this.state.type}
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                        disabled={true}
                                                    >
                                                        {INTEGRATION_TYPE_TO_HUMAN[this.state.type]}
                                                    </Dropdown.Toggle>
                                                </Dropdown>
                                            </div>
                                        }
                                        </Col>
                                        <Col xs="7">
                                            <div>
                                                <span className="inputLabel">Input/Output type</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle 
                                                        className={"queryFilter"}
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                        disabled={true}
                                                    >
                                                        { this.state.item === null 
                                                            ? this.state.ioType.charAt(0).toUpperCase() + this.state.ioType.slice(1)
                                                            : this.state.item.io_type.charAt(0).toUpperCase() + this.state.item.io_type.slice(1)
                                                        }
                                                    </Dropdown.Toggle>
                                                </Dropdown>
                                            </div>
                                        </Col>
                                    </Row>

                                    {
                                        SQL_INTEGRATION_TYPES.includes(this.state.type)
                                        ?
                                            <SQLIntegrationForm
                                                item={this.state.item}
                                                successHandler={this.handleSQLFormAfterSuccess}
                                                ioType={this.state.ioType}
                                                type={this.state.type}
                                            />
                                        :

                                    <div>
                                        <Form
                                            onSubmit={this.handleRegister}
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
                                                autoFocus
                                                readOnly={this.state.isWaitingForConfirmation}
                                                validations={[required]}
                                            />
                                            </FloatingLabel>
                                            </div>
                                        </InputGroup>

                                        { this.state.type === INTEGRATION_TYPE_SMS &&
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingNumber" label="Phone number">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingNumber"
                                                        type="text"
                                                        placeholder="Phone number"
                                                        autoComplete="number"
                                                        name="number"
                                                        value={this.state.number}
                                                        onChange={this.onChangeNumber}
                                                        readOnly={this.state.isWaitingForConfirmation}
                                                    />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        }

                                        {this.state.type === INTEGRATION_TYPE_INCIDENT_IO &&
                                            <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingIncidentIoApiKey" label="API Key">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingIncidentIoApiKey"
                                                            type="text"
                                                            placeholder="API Key"
                                                            autoComplete="number"
                                                            name="api_key"
                                                            value={this.state.incidentIoApiKey}
                                                            onChange={this.onChangeIncidentIoApiKey}
                                                            readOnly={this.state.isWaitingForConfirmation}
                                                        />
                                                        </FloatingLabel>
                                                        <div className="inputTip clearfix">
                                                            You can find or create your API Key in the <a href="https://app.incident.io/settings/api-keys" rel="noreferrer" target="_blank">Incident.io settings</a>.
                                                        </div>
                                                    </div>
                                                </InputGroup>

                                                <span className="inputLabel">Mode</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.incidentIoMode}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {["real", "test", "tutorial"].map(value => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={value}
                                                                active={value === this.state.incidentIoMode}
                                                                onClick={(e) => this.onChangeIncidentIoMode(e.target.getAttribute('value'))}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                    <div className="inputTip clearfix">
                                                        Indicates whether incidents created via Ylem are real, or you test so far. Or maybe you educate somebody.
                                                    </div>
                                                </Dropdown>

                                                <span className="inputLabel">Visibility</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.incidentIoVisibility}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {["public", "private"].map(value => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={value}
                                                                active={value === this.state.incidentIoVisibility}
                                                                onClick={(e) => this.onChangeIncidentIoVisibility(e.target.getAttribute('value'))}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                    <div className="inputTip clearfix">
                                                        Whether the incident should be open to anyone in your Slack workspace (public), or invite-only (private). For more information on Private Incidents see <a href="https://help.incident.io/en/articles/5947963-can-we-mark-incidents-as-sensitive-and-restrict-access" rel="noreferrer" target="_blank">Incident.io help centre</a>.
                                                    </div>
                                                </Dropdown>
                                            </div>
                                        }

                                        {this.state.type === INTEGRATION_TYPE_OPSGENIE &&
                                        <div>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingIncidentIoApiKey" label="API Key">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingIncidentIoApiKey"
                                                            type="text"
                                                            placeholder="API Key"
                                                            autoComplete="number"
                                                            name="api_key"
                                                            onChange={this.onChangeOpsgenieApiKey}
                                                            readOnly={this.state.isWaitingForConfirmation}
                                                        />
                                                    </FloatingLabel>
                                                    <div className="inputTip clearfix">
                                                        How to create a GeniKey, you can check out <a href="https://support.atlassian.com/opsgenie/docs/create-a-default-api-integration/" rel="noreferrer" target="_blank">in the Atlassian's documentation</a>.
                                                    </div>
                                                </div>
                                            </InputGroup>
                                        </div>
                                        }

                                        {this.state.type === INTEGRATION_TYPE_JENKINS &&
                                        <div>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingJenkinsBaseUrl" label="Base Url">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingJenkinsBaseUrl"
                                                            type="text"
                                                            placeholder="Base Url"
                                                            name="base_url"
                                                            onChange={this.onChangeJenkinsBaseUrl}
                                                            readOnly={this.state.isWaitingForConfirmation}
                                                            value={this.state.jenkinsBaseUrl}
                                                        />
                                                    </FloatingLabel>
                                                    <div className="inputTip clearfix">
                                                        Your jenkins base url. E.g., <strong>https://jenkins-ci.example.com</strong>. Please whitelist our IP for Ylem to be able to reach your server. @todo: write a doc in the docs on how to whitelist an IP
                                                    </div>
                                                </div>
                                            </InputGroup>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingJenkinsProjectName" label="Project Name">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingJenkinsProjectName"
                                                            type="text"
                                                            placeholder="Project Name"
                                                            name="project_name"
                                                            onChange={this.onChangeJenkinsProjectName}
                                                            readOnly={this.state.isWaitingForConfirmation}
                                                            value={this.state.jenkinsProjectName}
                                                        />
                                                    </FloatingLabel>
                                                    <div className="inputTip clearfix">
                                                        A build of what project needs to be triggered?
                                                    </div>
                                                </div>
                                            </InputGroup>
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingJenkinsToken" label="CSRF Token">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingJenkinsToken"
                                                            type="text"
                                                            placeholder="CSRF Token"
                                                            name="token"
                                                            onChange={this.onChangeJenkinsToken}
                                                            readOnly={this.state.isWaitingForConfirmation}
                                                            value={this.state.jenkinsToken}
                                                        />
                                                    </FloatingLabel>
                                                    <div className="inputTip clearfix">
                                                        This token is used to <a href="https://www.jenkins.io/doc/book/security/csrf-protection/" rel="noreferrer" target="_blank">protect the queries to your API</a>.
                                                    </div>
                                                </div>
                                            </InputGroup>
                                        </div>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_EMAIL &&
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingEmail" label="Email">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingEmail"
                                                        type="text"
                                                        placeholder="Email"
                                                        autoComplete="email"
                                                        name="email"
                                                        value={this.state.email}
                                                        onChange={this.onChangeEmail}
                                                        readOnly={this.state.isWaitingForConfirmation}
                                                    />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_TABLEAU &&
                                            <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingServer" label="Server">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingServer"
                                                                type="text"
                                                                placeholder="Server"
                                                                name="server"
                                                                value={this.state.tableauServer}
                                                                onChange={this.onChangeTableauServer}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingUsername" label="Username">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingUsername"
                                                                type="text"
                                                                placeholder="Username"
                                                                name="username"
                                                                value={this.state.tableauUsername}
                                                                onChange={this.onChangeTableauUsername}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingPassword" label="Password">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingPassword"
                                                                type="password"
                                                                placeholder="Password"
                                                                name="password"
                                                                value={this.state.tableauPassword}
                                                                onChange={this.onChangeTableauPassword}
                                                                iscreation={this.state.item === null ? "true" : "false"}
                                                                validations={[requiredForCreation]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingSiteName" label="Site Name">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingSiteName"
                                                                type="text"
                                                                placeholder="Site Name"
                                                                name="siteName"
                                                                value={this.state.tableauSiteName}
                                                                onChange={this.onChangeTableauSiteName}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingProjectName" label="Project Name">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingProjectName"
                                                                type="text"
                                                                placeholder="Project Name"
                                                                name="projectName"
                                                                value={this.state.tableauProjectName}
                                                                onChange={this.onChangeTableauProjectName}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingDatasourceName" label="Datasource Name">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingProjectName"
                                                                type="text"
                                                                placeholder="Datasource Name"
                                                                name="datasourceName"
                                                                value={this.state.tableauDatasourceName}
                                                                onChange={this.onChangeTableauDatasourceName}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <span className="inputLabel">Mode</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.tableauMode}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {["overwrite", "append"].map(value => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={value}
                                                                active={value === this.state.tableauMode}
                                                                onClick={(e) => this.onChangeTableauMode(e.target.getAttribute('value'))}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                            </div>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_GOOGLE_SHEETS &&
                                            <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingServer" label="Credentials">
                                                            <Textarea
                                                                className="form-control form-control-lg codeEditor"
                                                                id="floatingCredentials"
                                                                type="textarea"
                                                                placeholder="Credentials (JSON)"
                                                                autoComplete="credentials"
                                                                name="credentials"
                                                                value={this.state.googleSheetsCredentials}
                                                                onChange={this.onChangeGoogleSheetsCredentials}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingUsername" label="Spreadsheet Id">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingSpreadsheetId"
                                                                type="text"
                                                                placeholder="Spreadsheet Id"
                                                                name="spreadsheetId"
                                                                value={this.state.googleSheetsSpreadsheetId}
                                                                onChange={this.onChangeGoogleSheetsSpreadsheetId}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingPassword" label="Sheet Id">
                                                            <Input
                                                                className="form-control form-control-lg"
                                                                id="floatingSheetId"
                                                                type="text"
                                                                placeholder="Sheet Id"
                                                                name="sheetId"
                                                                value={this.state.googleSheetsSheetId}
                                                                onChange={this.onChangeGoogleSheetsSheetId}
                                                                validations={[required]}
                                                            />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <span className="inputLabel">Mode</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.googleSheetsMode}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {["overwrite", "append"].map(value => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={value}
                                                                active={value === this.state.googleSheetsMode}
                                                                onClick={(e) => this.onChangeGoogleSheetsMode(e.target.getAttribute('value'))}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>

                                                <span className="inputLabel">Write header</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.googleSheetsWriteHeader}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {["yes", "no"].map(value => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={value}
                                                                active={value === this.state.googleSheetsWriteHeader}
                                                                onClick={(e) => this.onChangeGoogleSheetsWriteHeader(e.target.getAttribute('value'))}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                            </div>
                                        }


                                        { this.state.type === INTEGRATION_TYPE_HUBSPOT &&
                                        <div>
                                            { this.state.hubspotAuthorizations !== null ?
                                                <div className="mb-4">
                                                    <div className="float-left">
                                                        <span className="inputLabel">Hubspot authorization</span><br/>
                                                        <Dropdown size="lg">
                                                            <Dropdown.Toggle
                                                                variant="light"
                                                                id="dropdown-basic"
                                                            >
                                                                {this.state.hubspotAuthorizationName}
                                                            </Dropdown.Toggle>

                                                            <Dropdown.Menu>
                                                                {this.state.hubspotAuthorizations.map(value => (
                                                                    value.is_active === true &&
                                                                    (
                                                                        <Dropdown.Item
                                                                            value={value.name}
                                                                            key={value.uuid}
                                                                            active={value.uuid === this.state.hubspotAuthorizationUuid}
                                                                            onClick={(e) => this.onChangeHubspotAuthorization(value.uuid, value.name)}
                                                                        >
                                                                            {value.name}
                                                                        </Dropdown.Item>
                                                                    )
                                                                ))}
                                                            </Dropdown.Menu>
                                                        </Dropdown>
                                                    </div>
                                                    <div className="pt-4 mt-2 px-3 float-left">
                                                        <a href="/hubspot-authorizations" rel="noreferrer" target="_blank">Manage Hubspot authorizations</a>
                                                    </div>
                                                    <div className="clearfix"></div>
                                                </div>
                                                : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                            }

                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingPipelineCode" label="Pipeline Id">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingPipelineCode"
                                                            type="text"
                                                            placeholder="Pipeline Id. E.g., 0"
                                                            name="pipeline_code"
                                                            value={this.state.hubspotPipelineCode}
                                                            onChange={this.onChangeHubspotPipelineCode}
                                                            validations={[required]}
                                                        />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>

                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingPipelineStageCode" label="Pipeline Stage Id">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingPipelineStageCode"
                                                            type="text"
                                                            placeholder="Pipeline Stage Id. E.g., 1"
                                                            name="pipeline_stage_code"
                                                            value={this.state.hubspotPipelineStageCode}
                                                            onChange={this.onChangeHubspotPipelineStageCode}
                                                            validations={[required]}
                                                        />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>

                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingOwnerCode" label="Owner Id">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingOwnerCode"
                                                            type="text"
                                                            placeholder="Pipeline Owner Id. E.g., 421381513"
                                                            name="pipeline_ownerCode"
                                                            value={this.state.hubspotOwnerCode}
                                                            onChange={this.onChangeHubspotOwnerCode}
                                                            validations={[required]}
                                                        />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        </div>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_SALESFORCE &&
                                        <div>
                                            { this.state.hubspotAuthorizations !== null ?
                                                <div className="mb-4">
                                                    <div className="float-left">
                                                        <span className="inputLabel">Salesforce authorization</span><br/>
                                                        <Dropdown size="lg">
                                                            <Dropdown.Toggle
                                                                variant="light"
                                                                id="dropdown-basic"
                                                            >
                                                                {this.state.salesforceAuthorizationName}
                                                            </Dropdown.Toggle>

                                                            <Dropdown.Menu>
                                                                {this.state.salesforceAuthorizations.map(value => (
                                                                    value.is_active === true &&
                                                                    (
                                                                        <Dropdown.Item
                                                                            value={value.name}
                                                                            key={value.uuid}
                                                                            active={value.uuid === this.state.salesforceAuthorizationUuid}
                                                                            onClick={(e) => this.onChangeSalesforceAuthorization(value.uuid, value.name)}
                                                                        >
                                                                            {value.name}
                                                                        </Dropdown.Item>
                                                                    )
                                                                ))}
                                                            </Dropdown.Menu>
                                                        </Dropdown>
                                                    </div>
                                                    <div className="pt-4 mt-2 px-3 float-left">
                                                        <a href="/salesforce-authorizations" rel="noreferrer" target="_blank">Manage Salesforce authorizations</a>
                                                    </div>
                                                    <div className="clearfix"></div>
                                                </div>
                                                : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                            }
                                        </div>
                                        }

                                        { (this.state.type === INTEGRATION_TYPE_EMAIL
                                            || this.state.type === INTEGRATION_TYPE_SMS)
                                            && this.state.isWaitingForConfirmation === true
                                            &&
                                            <InputGroup className="mb-4">
                                                <Alert variant="info" className="mt-4">
                                                    We just sent you a confirmation code. Please enter it here and click on "Save".
                                                </Alert>
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingCode" label="Confirmation code">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingCode"
                                                            type="text"
                                                            placeholder="Confirmation code"
                                                            autoComplete="code"
                                                            name="code"
                                                            value={this.state.code}
                                                            onChange={this.onChangeCode}
                                                        />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_SLACK &&
                                        <div>
                                            { this.state.slackAuthorizations !== null ?
                                            <div className="mb-4">
                                                <div className="float-left">
                                                <span className="inputLabel">Slack authorization</span><br/>
                                                <Dropdown size="lg">
                                                    <Dropdown.Toggle
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                    >
                                                        {this.state.slackAuthorizationName}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {this.state.slackAuthorizations.map(value => (
                                                            value.is_active === true &&
                                                            (
                                                            <Dropdown.Item
                                                                value={value.name}
                                                                key={value.uuid}
                                                                active={value.uuid === this.state.slackAuthorizationUuid}
                                                                onClick={(e) => this.onChangeSlackAuthorization(value.uuid, value.name)}
                                                            >
                                                                {value.name}
                                                            </Dropdown.Item>
                                                            )
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                                </div>
                                                <div className="pt-4 mt-2 px-3 float-left">
                                                    <a href="/slack-authorizations" rel="noreferrer" target="_blank">Manage Slack authorizations</a>
                                                </div>
                                                <div className="clearfix"></div>
                                            </div>
                                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                            }
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingChannel" label="Slack channel">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingChannel"
                                                        type="text"
                                                        placeholder="Slack channel"
                                                        autoComplete="channel"
                                                        name="channel"
                                                        value={this.state.channel}
                                                        onChange={this.onChangeChannel}
                                                        validations={[required]}
                                                    />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        </div>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_JIRA &&
                                        <div>
                                            { this.state.jiraAuthorizations !== null ?
                                            <div className="mb-4">
                                                <div className="float-left">
                                                <span className="inputLabel">Jira Cloud authorization</span><br/>
                                                <Dropdown size="lg">
                                                    <Dropdown.Toggle
                                                        variant="light"
                                                        id="dropdown-basic"
                                                    >
                                                        {this.state.jiraAuthorizationName}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {this.state.jiraAuthorizations.map(value => (
                                                            value.is_active === true &&
                                                            (
                                                            <Dropdown.Item
                                                                value={value.name}
                                                                key={value.uuid}
                                                                active={value.uuid === this.state.jiraAuthorizationUuid}
                                                                onClick={(e) => this.onChangeJiraAuthorization(value.uuid, value.name)}
                                                            >
                                                                {value.name}
                                                            </Dropdown.Item>
                                                            )
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                                </div>
                                                <div className="pt-4 mt-2 px-3 float-left">
                                                    <a href="/jira-authorizations" rel="noreferrer" target="_blank">Manage Jira Cloud authorizations</a>
                                                </div>
                                                <div className="clearfix"></div>
                                            </div>
                                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                            }
                                            <InputGroup className="mb-4">
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingProjectKey" label="Project Key">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingProjectKey"
                                                        type="text"
                                                        placeholder="Project Key, e.g., DEMO"
                                                        autoComplete="project_Key"
                                                        name="project_key"
                                                        value={this.state.jiraProjectKey}
                                                        onChange={this.onChangeJiraProjectKey}
                                                        validations={[required]}
                                                    />
                                                    </FloatingLabel>
                                                </div>
                                                <div className="registrationFormControl">
                                                    <FloatingLabel controlId="floatingIssueType" label="Issue Type">
                                                    <Input
                                                        className="form-control form-control-lg"
                                                        id="floatingIssueType"
                                                        type="text"
                                                        placeholder="Issue Type, e.g., Task"
                                                        autoComplete="issue_type"
                                                        name="issue_type"
                                                        value={this.state.jiraIssueType}
                                                        onChange={this.onChangeJiraIssueType}
                                                        validations={[required]}
                                                    />
                                                    </FloatingLabel>
                                                </div>
                                            </InputGroup>
                                        </div>
                                        }

                                        { this.state.type === INTEGRATION_TYPE_API &&
                                            <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingUrl" label="URL">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingUrl"
                                                            type="text"
                                                            placeholder="URL"
                                                            autoComplete="url"
                                                            name="url"
                                                            value={this.state.url}
                                                            onChange={this.onChangeUrl}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>

                                                <span className="inputLabel">Method</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.apiMethod.toUpperCase()}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {INTEGRATION_API_METHODS.map(method => (
                                                            <Dropdown.Item 
                                                                value={method}
                                                                key={method}
                                                                active={method === this.state.apiMethod}
                                                                onClick={(e) => this.onChangeApiMethod(method)}
                                                            >
                                                                {method.toUpperCase()}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>

                                                <span className="inputLabel">Authorization type</span><br/>
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle variant="light" id="dropdown-basic">
                                                        {this.state.authType}
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {INTEGRATION_TYPE_API_AUTH_TYPES.map(type => (
                                                            <Dropdown.Item 
                                                                value={type}
                                                                key={type}
                                                                active={type === this.state.authType}
                                                                onClick={(e) => this.onChangeAuthType(e.target.getAttribute('value'))}
                                                            >
                                                                {type}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown>

                                                {this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_BEARER &&
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingToken" label="Token">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingToken"
                                                            type="text"
                                                            placeholder="Token"
                                                            autoComplete="authBearerToken"
                                                            name="authBearerToken"
                                                            value={this.state.authBearerToken}
                                                            onChange={this.onChangeAuthBearerToken}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>
                                                }

                                                {this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_HEADER &&
                                                    <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingHeader" label="Header">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingHeader"
                                                            type="text"
                                                            placeholder="Header"
                                                            autoComplete="authHeaderName"
                                                            name="authHeaderName"
                                                            value={this.state.authHeaderName}
                                                            onChange={this.onChangeAuthHeaderName}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingValue" label="Value">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingValue"
                                                            type="text"
                                                            placeholder="Value"
                                                            autoComplete="authHeaderValue"
                                                            name="authHeaderValue"
                                                            value={this.state.authHeaderValue}
                                                            onChange={this.onChangeAuthHeaderValue}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>
                                                    </div>
                                                }

                                                {this.state.authType === INTEGRATION_TYPE_API_AUTH_TYPE_BASIC &&
                                                    <div>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingUsername" label="Username">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingUsername"
                                                            type="text"
                                                            placeholder="Username"
                                                            autoComplete="authBasicUserName"
                                                            name="authBasicUserName"
                                                            value={this.state.authBasicUserName}
                                                            onChange={this.onChangeAuthBasicUserName}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                </InputGroup>
                                                <InputGroup className="mb-4">
                                                    <div className="registrationFormControl">
                                                        <FloatingLabel controlId="floatingPassword" label="Password">
                                                        <Input
                                                            className="form-control form-control-lg"
                                                            id="floatingPassword"
                                                            type={this.state.authBasicPasswordType}
                                                            placeholder="Password"
                                                            autoComplete="authBasicPassword"
                                                            name="authBasicPassword"
                                                            value={this.state.authBasicPassword}
                                                            onChange={this.onChangeAuthBasicPassword}
                                                        />
                                                        </FloatingLabel>
                                                    </div>
                                                    <span
                                                        onClick={this.handleEyeClick}
                                                        className="eye"
                                                    >
                                                        {
                                                            this.state.authBasicPasswordType === 'text' 
                                                            ? <Tooltip title="Hide" placement="right"><VisibilityOffOutlined/></Tooltip>
                                                            : <Tooltip title="Show" placement="right"><VisibilityOutlined/></Tooltip>
                                                        }
                                                    </span>
                                                </InputGroup>
                                                    </div>
                                                }

                                            </div>
                                        }
                                        
                                            <div>
                                                <Row className="pt-3">
                                                    <Col xs="12">
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
                                                </Row>
                                                {message && isInProgress && (
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
                                            </div>
                                        </Form>
                                        </div>
                                    }

                                </Card.Body>
                            </Card>
                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                        }
                        </Col>
                    </Row>
                </Container>
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

export default connect(mapStateToProps)(IntegrationForm);
