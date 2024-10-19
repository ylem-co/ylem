import React from 'react';
import { Navigate, useParams, useLocation, useNavigate } from 'react-router-dom';
import {connect} from "react-redux";
import { Fade } from "react-awesome-reveal";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";

import Tooltip from '@mui/material/Tooltip';
import Circle from '@mui/icons-material/Circle';

import {JiraAuthorizationInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";

import IntegrationService from "../../../services/integration.service";

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";
import log from "loglevel";
import FullScreenModal from "../../../components/modals/fullScreenModal.component";
import JiraAuthorizationForm from "../../../components/forms/jiraAuthorizationForm.component";

const JIRA_JUST_CONNECTED_QUERY_STRING = "?justConnected";

function withParams(Component) {
  return props => <Component {...props} params={useParams()} location={useLocation()} history={useNavigate()} />;
}

class JiraAuthorizations extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            authorizations: null,
            authorizationLink: null,
            loading: false,
            activeItem: null,
            itemUuid: this.props.params.itemUuid || null,
            isFormOpen: false,
        };
    }

    componentDidMount = async() => {
        log.debug("JiraAuthorizations component did mount")

        document.title = 'Jira Cloud Authorizations'
        await this.handleGetAuthorizations(this.state.organization.uuid);

        if (this.state.itemUuid !== null) {
            log.info("Looking for an element with uuid", this.state.itemUuid)
            await this.handleGetAuthorization(this.state.itemUuid);

            if (this.state.activeItem !== null) {
                    log.info("Found an element", this.state.activeItem)
                    this.toogleForm(this.state.activeItem);
            } else {
                log.debug("Authorization was not found")
            }
        } else {
            log.debug("No uuid presented")
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handleCreateAuthorization = async(uuid) => {
        let authorizationLink = this.state.authorizationLink;

        if (authorizationLink === null) {
            authorizationLink = IntegrationService.createJiraAuthorization(uuid);

            return Promise.resolve(authorizationLink)
                .then(async(authorizationLink) => {
                    if (authorizationLink.data) {
                        await this.promisedSetState({authorizationLink: authorizationLink.data.url});
                        return authorizationLink.data.url;
                    } else {
                        return null;
                    }
                })
                .catch(async() => {
                    return null;
                });
        } else {
            return authorizationLink;
        }
    };

    handleGetAuthorizations = async(uuid) => {
        log.debug("handle get authorizations");
        let authorizations = this.state.authorizations;

        if (
            authorizations === null
            || authorizations.length === 0
        ) {
            authorizations = IntegrationService.getJiraAuthorizations(uuid);

            await Promise.resolve(authorizations)
                .then(async(authorizations) => {
                    if (authorizations.data) {
                        await this.promisedSetState({authorizations: authorizations.data.items});
                    } else {
                        await this.promisedSetState({authorizations: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({authorizations: []});
                });
        }
    };

    handleGetAuthorization = async(uuid) => {
        log.debug("handle get authorization");
            let authorization = IntegrationService.getJiraAuthorization(uuid);

            await Promise.resolve(authorization)
                .then(async(result) => {
                    if (result.data) {
                        await this.promisedSetState({activeItem: result.data});
                    } else {
                        await this.promisedSetState({activeItem: null});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({activeItem: null});
                });
    };

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    };

    toogleAuthorization = async() => {
        await this.promisedSetState({
            loading: true,
        });

        var link = await this.handleCreateAuthorization(this.state.organization.uuid);

        if (link !== null) {
            window.location.href = link;
        }

        await this.promisedSetState({
            loading: false,
        });
    };

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    };

    toogleForm = async(item = null) => {
        log.debug("toggle form")
        await this.promisedSetState({
            isFormOpen: !this.state.isFormOpen,
            activeItem: item,
        });
    }

    closeForm = async() => {
        log.debug("close form")
        let item = this.state.activeItem;

        await this.promisedSetState({
            isFormOpen: false,
            activeItem: null,
        });

        if (item !== null) {
            this.props.history('/jira-authorizations/');
        }
    }

    openForm = (item) => {
        log.debug("open form")
        this.props.history('/jira-authorizations/' + item.uuid);
    }

    handleCloseFormAfterSuccess = async() => {
        log.debug("handle close form after success")
        await this.promisedSetState({
            authorizations: null,
        });
        await this.handleGetAuthorizations(this.state.organization.uuid);
        await this.closeForm();
    }

    render() {
        const { 
            isInfoOpen,
            isFormOpen,
            activeItem,
        } = this.state;

        const { isLoggedIn, user, location } = this.props;

        const { search } = location;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

        return (
            <Fade>
                        <Row className="mb-3">
                            <Col sm="9">
                                <h1>Jira Cloud Authorizations</h1>
                                <Tooltip title="Info" placement="right">
                                    <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                                </Tooltip>
                            </Col>
                            <Col sm="3" className="text-right">
                                <ButtonGroup className="mr-4">
                                    <Button
                                        color="primary"
                                        className="mx-0"
                                        onClick={() => this.toogleAuthorization()}
                                        disabled={this.state.loading}
                                    >
                                        {this.state.loading && (
                                            <span className="spinner-border spinner-border-sm spinner-primary"></span>
                                        )}
                                        <span>Add Jira authorization</span>
                                    </Button>
                                </ButtonGroup>
                            </Col>
                        </Row>
                        <Row>

                {this.state.authorizations !== null ?
                    this.state.authorizations.length > 0 ?
                        this.state.authorizations.map(value => (
                            <Col className="col-3 mb-4" key={value.uuid}>
                                <Card className="withHeader editableCard" onClick={() => this.openForm(value)}>
                                    <Card.Header className="noBottom">
                                        <Row>
                                            <Col className="col-10">
                                                {value.name}
                                            </Col>
                                            <Col className="col-2 text-right">
                                                <Tooltip title={value.is_active === false ? "Authorization is not active" : "Authorization is active"} placement="right">
                                                    <Circle 
                                                        className={value.is_active === false ? "icon_offline" : "icon_online"}
                                                        alt={value.is_active === false ? "Authorization is not active" : "Authorization is active"}
                                                    />
                                                </Tooltip>
                                            </Col>
                                        </Row>
                                    </Card.Header>
                                </Card>
                            </Col>
                        )) 
                        : <div className="text-center">You didn't create any jira authorization yet</div>
                    : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                        </Row>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={JiraAuthorizationInfo.title}
                    content={JiraAuthorizationInfo.content}
                />

                <FullScreenModal
                    show={isFormOpen}
                    onHide={this.closeForm}
                    title="Edit Jira Authorization"
                    content={
                        <>
                            <JiraAuthorizationForm
                                item={activeItem}
                                successHandler={this.handleCloseFormAfterSuccess}
                                justConnected={search === JIRA_JUST_CONNECTED_QUERY_STRING}
                            />
                        </>
                    }
                    item={activeItem}
                    confirmButton={false}
                />
            </Fade>
        )
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    return {
        isLoggedIn,
        user,
    };
}

export default connect(mapStateToProps)(withParams(JiraAuthorizations));
