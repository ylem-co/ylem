import React from 'react';
import { Navigate } from 'react-router-dom';
import { connect } from "react-redux";
import { Fade } from "react-awesome-reveal";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Dropdown from "react-bootstrap/Dropdown";

import Tooltip from '@mui/material/Tooltip';
import ContentCopy from '@mui/icons-material/ContentCopy';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';

import { OAuthInfo } from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import RightModal from "../../../components/modals/rightModal.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component";

import OAuthClientForm from "../../../components/forms/OAuthClientForm.component";

import OAuthService from "../../../services/oauth.service";

import { PERMISSION_LOGGED_IN, validatePermissions } from "../../../actions/pipeline";
import { ROLE_ORGANIZATION_ADMIN } from "../../../actions/roles";

function copyLink(copiedLink){
  navigator.clipboard.writeText(copiedLink);
}

class APIClients extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            isFormOpen: false,
            clients: null,
            activeItem: null,
            isTerminationModalOpen: false,
            clientToRemove: null,
        };
    }

    componentDidMount = async() => {
        document.title = 'API OAuth Clients';

        await this.handleGetClients(
            this.state.organization.uuid
        );
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handleGetClients = async(uuid) => {
        let clients = this.state.clients;

        if (
            clients === null
            || clients.length === 0
        ) {
            clients = OAuthService.getClients(uuid);

            await Promise.resolve(clients)
                .then(async(clients) => {
                    if (clients.data && clients.data.length !== 0) {
                        var items = clients.data;
                        
                        await this.promisedSetState({clients: items});
                    } else {
                        await this.promisedSetState({clients: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({clients: []});
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
        await this.promisedSetState({
            isFormOpen: false,
            activeItem: null,
            clients: null,
        });

        await this.handleGetClients(
            this.state.organization.uuid
        );
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
            clientToRemove: null,
        });
    };

    handleOpenTerminationModal = (client) => {
        this.setState({
            isTerminationModalOpen: true,
            clientToRemove: client,
        });
    };

    handleConfirmTermination = async() => {
        var uuid = this.state.clientToRemove.uuid;

        await OAuthService.deleteClient(uuid);
        await this.promisedSetState({
            clients: null,
        });
        await this.handleGetClients(
            this.state.organization.uuid
        );
        await this.handleCloseTerminationModal();
    };

    render() {
        const { 
            isInfoOpen,
            isFormOpen,
            activeItem,
            isTerminationModalOpen,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (
            !validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)
            || !user.roles.includes(ROLE_ORGANIZATION_ADMIN)
        ) {
            return <Navigate to="/login" />;
        }

        return (
            <Fade>
                    <div>
                        <Row className="mb-5">
                            <Col sm="9">
                                <h1>API OAuth Clients</h1>
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
                                        Add OAuth Client
                                    </Button>
                                </ButtonGroup>
                            </Col>
                        </Row>
                        <Row>
                        {this.state.clients !== null ?
                            this.state.clients.length > 0 ?
                                this.state.clients.map(value => (
                                    <Col className="col-3 mb-4" key={value.uuid}>
                                        <Card className="withHeader">
                                            <Card.Header>
                                                <Row>
                                                    <Col className="col-10">
                                                        {value.name}
                                                    </Col>
                                                </Row>
                                            </Card.Header> 
                                            <Card.Body>
                                                <div className="code">
                                                    <Row>
                                                        <Col xs={10}>
                                                            <strong>Client ID:</strong> {value.uuid}
                                                        </Col>
                                                        <Col xs={2} className="text-right">
                                                            <Tooltip title="Click to copy to clipboard" placement="left">
                                                                <ContentCopy className="pointer"
                                                                    onClick={() => copyLink(value.uuid)}
                                                                />
                                                            </Tooltip>
                                                        </Col>
                                                    </Row>
                                                </div>
                                                <div className="text-right">
                                                    <Dropdown as={ButtonGroup}>
                                                        <Dropdown.Toggle split size="sm" variant="transparent" id="dropdown-split-basic" />

                                                        <Dropdown.Menu>
                                                            <Dropdown.Item onClick={() => this.handleOpenTerminationModal(value)} className="dangerAction">
                                                                <DeleteOutlineRounded className="dropdownIcon" /> Delete
                                                            </Dropdown.Item>
                                                        </Dropdown.Menu>
                                                    </Dropdown>
                                                </div> 
                                            </Card.Body>
                                        </Card>
                                    </Col>
                                )) 
                                : <div className="text-center">You didn't create any OAuth Client yet</div>
                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                        }
                        </Row>
                    </div>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={OAuthInfo.title}
                    content={OAuthInfo.content}
                />

                <ConfirmationModal
                    show={isTerminationModalOpen}
                    title={"Delete OAuth client"}
                    body={"Are you sure you want to delete this OAuth client? Carefully double-check if any of your API clients still use it."}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />

                <RightModal
                    show={isFormOpen}
                    content={
                      <OAuthClientForm 
                        item={activeItem}
                        successHandler={this.handleFormSuccess}
                      />
                    }
                    item={activeItem}
                    title={activeItem === null ? "Create client" : "Edit client"}
                    onHide={this.closeForm}
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

export default connect(mapStateToProps)(APIClients);
