import React from 'react';
import { Navigate, useParams, useNavigate } from 'react-router-dom';
import {connect} from "react-redux";
import { DateTime } from "luxon";
import { Fade } from "react-awesome-reveal";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";
import Button from "react-bootstrap/Button";
import Dropdown from "react-bootstrap/Dropdown";
import ButtonGroup from "react-bootstrap/ButtonGroup";

import Circle from '@mui/icons-material/Circle';

import Tooltip from '@mui/material/Tooltip';
import IntegrationService, {INTEGRATION_TYPE_SQL, INTEGRATION_IO_TYPES_FORM} from "../../../services/integration.service";
import {IntegrationsInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import FullScreenModal from "../../../components/modals/fullScreenModal.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component"
import {TimeAgo} from "../../../components/timeAgo.component";
import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";

import IntegrationForm from "../../../components/forms/integrationForm.component"

function withParams(Component) {
  return props => <Component {...props} params={useParams()} history={useNavigate()} />;
}

class Integrations extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            isFormOpen: false,
            integrations: null,
            activeItem: null,
            itemUuid: this.props.params.page === "details"
                ? (this.props.params.itemUuid || null)
                : null,
            page: this.props.params.page || "type",
            type: this.props.params.page === "type"
                ? (this.props.params.type || "all")
                : "all",
        };
    }

    componentDidMount = async() => {
        document.title = 'Integrations'
        await this.handleGetIntegrations(this.state.organization.uuid);

        if (
            this.state.page === "details"
            && this.state.itemUuid !== null 
            && this.state.integrations !== null
        ) {
            var element = this.state.integrations.find(o => o.uuid === this.state.itemUuid);

            if (element) {
                this.toogleForm(element);
            }
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    UNSAFE_componentWillReceiveProps = async(props) => {
        await this.promisedSetState({
            itemUuid: props.params.page === "details"
                ? (props.params.itemUuid || null)
                : null,
            page: props.params.page || "type",
            type: props.params.page === "type"
                ? (props.params.type || "all")
                : "all",
            integrations: null,
        })
        await this.handleGetIntegrations(this.state.organization.uuid);

        if (
            this.state.page === "details"
            && this.state.itemUuid !== null 
            && this.state.integrations !== null
        ) {
            var element = this.state.integrations.find(o => o.uuid === this.state.itemUuid);

            if (element) {
                this.toogleForm(element);
            }
        }
    }

    handleGetIntegrations = async(uuid) => {
        let integrations = this.state.integrations;

        if (
            integrations === null
            || integrations.length === 0
        ) {
            integrations = IntegrationService.getIntegrationsByOrganization(uuid, this.state.type);

            await Promise.resolve(integrations)
                .then(async(integrations) => {
                    if (integrations.data) {
                        await this.promisedSetState({integrations: integrations.data.items});
                    } else {
                        await this.promisedSetState({integrations: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({integrations: []});
                });
        }
    };

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    };

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    };

    toogleForm = async(item = null) => {
        await this.promisedSetState({
            isFormOpen: !this.state.isFormOpen,
            activeItem: item,
        });
    };

    closeForm = async() => {
        let item = this.state.activeItem;

        await this.promisedSetState({
            isFormOpen: false,
            activeItem: null,
        });

        if (item !== null) {
            this.props.history('/integrations/');
        }
    };

    openForm = (item) => {
        this.props.history('/integrations/details/' + item.uuid);
    };

    onChangeFilterType = (type) => {
        this.props.history('/integrations/type/' + type);
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
        });
    };

    handleOpenTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: true,
        });
    };

    handleConfirmTermination = async() => {
        var uuid = this.state.activeItem.uuid;

        await IntegrationService.deleteIntegration(uuid);
        await this.promisedSetState({
            integrations: null,
        });
        await this.handleGetIntegrations(this.state.organization.uuid);
        await this.handleCloseTerminationModal();
        await this.closeForm();

        window.location.reload();
    };

    handleCloseFormAfterSuccess = async() => {
        await this.promisedSetState({
            integrations: null,
        });
        await this.handleGetIntegrations(this.state.organization.uuid);
        await this.closeForm();
        
        window.location.reload();
    };

    render() {
        const { 
            isInfoOpen, 
            isFormOpen, 
            isTerminationModalOpen, 
            activeItem,
            type,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

        return (
            <>
                <Dropdown size="lg" className="mb-5">
                    <Dropdown.Toggle
                        variant="light" 
                        id="dropdown-basic"
                    >
                        {type.charAt(0).toUpperCase() + type.slice(1) + " integrations"}
                    </Dropdown.Toggle>
                    <Dropdown.Menu>
                        {INTEGRATION_IO_TYPES_FORM.map(type => (
                            <Dropdown.Item
                                value={type}
                                key={"io_dropdown"+type}
                                active={type === this.state.type}
                                onClick={() => this.onChangeFilterType(type)}
                            >
                                {type.charAt(0).toUpperCase() + type.slice(1)}
                            </Dropdown.Item>
                        ))}
                    </Dropdown.Menu>
                </Dropdown>
            <Fade>
                <Row className="mb-5">
                    <Col sm="9">
                        <h1>{type.charAt(0).toUpperCase() + type.slice(1) + " integrations"}</h1>
                        <Tooltip title="Info" placement="right">
                            <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                        </Tooltip>
                    </Col>
                    <Col sm="3" className="text-right">
                        <ButtonGroup className="mr-4">
                            <Button
                                color="primary"
                                className="mx-0"
                                onClick={() => this.toogleForm()}
                            >
                                Add integration
                            </Button>
                        </ButtonGroup>
                    </Col>
                </Row>
                    
                <Row>

                {this.state.integrations !== null ?
                    this.state.integrations.length > 0 ?
                        this.state.integrations.map(value => (
                            <Col className="col-3 mb-4" key={value.uuid}>
                                <Card className="withHeader editableCard" onClick={() => this.openForm(value)}>
                                    <Card.Header>
                                        <Row>
                                            <Col className="col-10">
                                                {value.name}
                                            </Col>
                                            <Col className="col-2 text-right">
                                                <Tooltip title={value.status} placement="right">
                                                    <Circle 
                                                        className={"icon_" + value.status}
                                                        alt={value.status === "offline" ? "Connection doesn't work" : "Connection works"}
                                                    />
                                                </Tooltip>
                                            </Col>
                                        </Row>
                                    </Card.Header> 
                                    <Card.Body className={"cardWithBg cardWithBg-" + (value.type === INTEGRATION_TYPE_SQL ? value.value : value.type) }>
                                        {value.type !== INTEGRATION_TYPE_SQL &&
                                            <span>{value.value}</span>
                                        }
                                        <br/>
                                        <span className="note">Last modified: {TimeAgo(DateTime.fromISO(value.user_updated_at, { zone: 'utc'}))}</span>
                                    </Card.Body>    
                                </Card>
                            </Col>
                        ))
                    : <div className="text-center">You didn't create any integration yet</div>
                : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }

                </Row>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={IntegrationsInfo.title}
                    content={IntegrationsInfo.content}
                />

                <FullScreenModal
                    show={isFormOpen}
                    onHide={this.closeForm}
                    title={activeItem !== null && activeItem.name ? activeItem.name : "Add integration"}
                    content={
                        <>
                            <IntegrationForm
                                item={activeItem}
                                integrationType={type}
                                successHandler={this.handleCloseFormAfterSuccess}
                            />
                        </>
                    }
                    onAltButtonClick={this.handleOpenTerminationModal}
                    altButtonText="Delete integration"
                    altButton={activeItem !== null && activeItem.name ? true : false}
                    item={activeItem}
                    confirmButton={false}
                />

                <ConfirmationModal
                    show={isTerminationModalOpen}
                    title="Delete integration"
                    body={"Are you sure you want to delete this integration?"}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />
            </Fade>
        </>
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

export default connect(mapStateToProps)(withParams(Integrations));
