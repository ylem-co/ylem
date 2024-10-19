import React from 'react';
import { Navigate } from 'react-router-dom';
import {connect} from "react-redux";
import { Fade } from "react-awesome-reveal";

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Card from 'react-bootstrap/Card';

import {UsersInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import EmailChips from "../../../components/emailChips.component";

import ClearOutlined from '@mui/icons-material/ClearOutlined';
import VerticalAlignTop from '@mui/icons-material/VerticalAlignTop';
import VerticalAlignBottom from '@mui/icons-material/VerticalAlignBottom';
import CheckCircleOutline from '@mui/icons-material/CheckCircleOutline';

import Tooltip from '@mui/material/Tooltip';

import Spinner from "react-bootstrap/Spinner";

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";

import { ROLE_ORGANIZATION_ADMIN, ROLE_TEAM_MEMBER, showUserFriendlyRoles } from "../../../actions/roles";

import UserService from "../../../services/user.service";
import InvitationService from "../../../services/invitation.service";

import Avatar from "../../../components/avatar.component";
import PendingInvitations from "../../../components/pendingInvitations.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component"

const { SearchBar } = Search;

class Users extends React.Component {
    constructor(props) {
        super(props);

        this.inviteItemsChangesHandler = this.inviteItemsChangesHandler.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            items: null,
            isInfoOpen: false,
            isTerminationModalOpen: false,
            terminationUuid: null,
            terminationName: null,
            isActivationModalOpen: false,
            activationUuid: null,
            activationName: null,
            isUpgradeModalOpen: false, 
            upgradeUuid: null,
            upgradeName: null,
            isDowngradeModalOpen: false, 
            downgradeUuid: null,
            downgradeName: null,
            isPendingInvitationsOpen: false,
            isFormOpen: false,
            inviteItems: null,
            submitButtonLoading: false,
            submissionMessage: null,
            submissionSuccessful: null,
        };
    }

    componentDidMount = async() => {
        document.title = 'Users'
        
        await this.handleGetUsers(this.state.organization.uuid);
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    inviteItemsChangesHandler = (inviteItems) => {
        this.setState({inviteItems});
    };

    handleGetUsers = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = UserService.getUsersByOrganization(uuid);

            Promise.resolve(items)
                .then(async(items) => {
                    if (items.data) {
                        await this.promisedSetState({items: items.data.items});
                    } else {
                        await this.promisedSetState({items: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({items: []});
                });
        }
    };

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    };

    tooglePendingInvitations = async() => {
        await this.promisedSetState({
            isPendingInvitationsOpen: !this.state.isPendingInvitationsOpen,
        });
    };

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    };

    closePendingInvitations = () => {
        this.setState({isPendingInvitationsOpen: false});
    };

    handleConfirmTermination = async() => {
        var uuid = this.state.terminationUuid;

        await UserService.terminateUser(uuid);
        await this.promisedSetState({items: null});
        await this.handleGetUsers(this.state.organization.uuid);
        this.handleCloseTerminationModal();
    };

    handleConfirmActivation = async() => {
        var uuid = this.state.activationUuid;

        await UserService.activateUser(uuid);
        await this.promisedSetState({items: null});
        await this.handleGetUsers(this.state.organization.uuid);
        this.handleCloseActivationModal();
    };

    handleConfirmDowngrade = async() => {
        var uuid = this.state.downgradeUuid;

        await UserService.assignRoleToUser(uuid, ROLE_TEAM_MEMBER);
        await this.promisedSetState({items: null});
        await this.handleGetUsers(this.state.organization.uuid);
        this.handleCloseDowngradeModal();
    };

    handleConfirmUpgrade = async() => {
        var uuid = this.state.upgradeUuid;

        await UserService.assignRoleToUser(uuid, ROLE_ORGANIZATION_ADMIN);
        await this.promisedSetState({items: null});
        await this.handleGetUsers(this.state.organization.uuid);
        this.handleCloseUpgradeModal();
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
            terminationName: null,
            terminationUuid: null,
        });
    };

    handleOpenTerminationModal = (name, uuid) => {
        this.setState({
            isTerminationModalOpen: true,
            terminationName: name,
            terminationUuid: uuid,
        });
    };

    handleCloseActivationModal = () => {
        this.setState({
            isActivationModalOpen: false,
            activationName: null,
            activationUuid: null,
        });
    };

    handleOpenActivationModal = (name, uuid) => {
        this.setState({
            isActivationModalOpen: true,
            activationName: name,
            activationUuid: uuid,
        });
    };

    handleCloseUpgradeModal = () => {
        this.setState({
            isUpgradeModalOpen: false,
            upgradeName: null,
            upgradeUuid: null,
        });
    };

    handleOpenUpgradeModal = (name, uuid) => {
        this.setState({
            isUpgradeModalOpen: true,
            upgradeName: name,
            upgradeUuid: uuid,
        });
    };

    handleCloseDowngradeModal = () => {
        this.setState({
            isDowngradeModalOpen: false,
            downgradeName: null,
            downgradeUuid: null,
        });
    };

    handleOpenDowngradeModal = (name, uuid) => {
        this.setState({
            isDowngradeModalOpen: true,
            downgradeName: name,
            downgradeUuid: uuid,
        });
    };

    toogleForm = async() => {
        await this.promisedSetState({
            isFormOpen: !this.state.isFormOpen,
        });
    };

    closeForm = () => {
        this.setState({
            isFormOpen: false,
            submitButtonLoading: false,
            submissionMessage: null,
            submissionSuccessful: null,
        });
    };

    submitForm = async() => {
        var inviteItems = this.state.inviteItems;

        if (inviteItems !== null && inviteItems.length > 0) {
            await this.promisedSetState({
                submitButtonLoading: true,
                submissionMessage: null,
                submissionSuccessful: null,
            });

            var response = InvitationService.sendInvitations(this.state.organization.uuid, inviteItems.join(","));

            Promise.resolve(response)
                .then(response => {
                    this.setState({
                        inviteItems: null, 
                        submitButtonLoading: false,
                        submissionMessage: "Invitations have been successfully sent",
                        submissionSuccessful: true,
                    });

                    setTimeout(() => {
                        this.closeForm();
                    },1000);
                })
                .catch(error => {
                    this.setState({
                        submitButtonLoading: false,
                        submissionMessage: "Something went wrong, please try again",
                        submissionSuccessful: false,
                    });
                });
        }
    };

    render() {
        const { 
            isInfoOpen, 
            isTerminationModalOpen, 
            terminationName,
            isActivationModalOpen, 
            activationName,
            isUpgradeModalOpen, 
            upgradeName,
            isDowngradeModalOpen, 
            downgradeName, 
            isPendingInvitationsOpen,
            isFormOpen, 
            submitButtonLoading,
            submissionMessage,
            submissionSuccessful,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (
            !validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)
            || !user.roles.includes(ROLE_ORGANIZATION_ADMIN)
        ) {
            return <Navigate to="/login" />;
        }

    const columns = [
        {
            dataField: '',
            text: '',
            sort: false,
            headerStyle: { width: '100px' },
            formatter: (cellContent, row) => (
                <>
                    <div className={row.is_active === 0 ? "fadedOut" : undefined}>
                        <Avatar
                            email={row.email}
                            avatar_url={null}
                            size={24}
                        /> 
                    </div>
                </>
            )
        },
        {
          dataField: 'first_name',
          text: 'Name',
          sort: true,
          formatter: (cellContent, row) => (
              <div className={row.is_active === 0 && "fadedOut"}>
                {row.first_name + " " + row.last_name}
              </div>
          )
        },
        {
            dataField: 'last_name',
            hidden: true
        },
        {
            dataField: 'email',
            text: 'Email',
            sort: true,
            formatter: (cellContent, row) => (
                <>
                    <div className={row.is_active === 0 && "fadedOut"}>
                        {row.email}
                    </div>
                </>
            )
        },
        {
            dataField: 'roles',
            text: 'Roles',
            sort: true,
            searchable: true,
            formatter: (cellContent, row) => (
                <>
                    <div className={row.is_active === 0 && "fadedOut"}>
                        {showUserFriendlyRoles(row.roles)}
                    </div>
                </>
            )
        },
        {
            dataField: 'is_active',
            text: 'Status',
            sort: true,
            searchable: false,
            headerStyle: { width: '140px' },
            formatter: (cellContent, row) => (
                <>
                    <div className={row.is_active === 0 && "fadedOut"}>
                        {
                            row.is_active === 1 ? "Active" : "Inactive"
                        }
                    </div>
                </>
            )
        },
        {
            dataField: '',
            text: 'Actions',
            sort: false,
            headerStyle: { width: '100px' },
            formatter: (cellContent, row) => (
                <>
                    {
                        row.roles.includes(ROLE_ORGANIZATION_ADMIN) 
                        ?
                            <Tooltip title="Downgrade to team member" placement="left">
                            <VerticalAlignBottom 
                                className="iconDelete"
                                onClick={() => this.handleOpenDowngradeModal(
                                    row.first_name + " " + row.last_name,
                                    row.uuid
                                )}
                            />
                            </Tooltip>
                        :
                            <Tooltip title="Upgrade to administrator" placement="left">
                            <VerticalAlignTop 
                                className="iconUpgrade"
                                onClick={() => this.handleOpenUpgradeModal(
                                    row.first_name + " " + row.last_name,
                                    row.uuid
                                )}
                            />
                            </Tooltip>
                    }

                    {
                        row.is_active === 1 
                        ?
                            <Tooltip title="Deactivate" placement="left">
                            <ClearOutlined 
                                className="iconDelete"
                                onClick={() => this.handleOpenTerminationModal(
                                    row.first_name + " " + row.last_name,
                                    row.uuid
                                )}
                            />
                            </Tooltip>
                        :
                            <Tooltip title="Activate" placement="left">
                            <CheckCircleOutline 
                                className="iconUpgrade"
                                onClick={() => this.handleOpenActivationModal(
                                    row.first_name + " " + row.last_name,
                                    row.uuid
                                )}
                            />
                            </Tooltip>
                    }
                </>
            )
        },
    ];

    const defaultSorted = [{
      dataField: 'last_name',
      order: 'desc'
    }];

    const paginationOptions = {
        sizePerPage: 20,
        hideSizePerPage: true,
        hidePageListOnlyOnePage: true
    };

        return (
            <Fade>
                <div>
                <Row className="mb-3">
                    <Col sm="6">
                        <h1>Users</h1>
                        <Tooltip title="Info" placement="right">
                            <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                        </Tooltip>
                    </Col>
                    <Col sm="6" className="text-right">
                    <ButtonGroup>
                        <Button
                            variant="secondary"
                            className="mx-0"
                            onClick={() => this.tooglePendingInvitations()}
                        >
                            Pending invitations
                        </Button>
                        <Button
                            variant="primary"
                            className="mx-0"
                            onClick={() => this.toogleForm()}
                        >
                            Invite users
                        </Button>
                        </ButtonGroup>
                    </Col>
                </Row>
                
                <div className="pt-4 clearfix">
                {this.state.items !== null ?
                <ToolkitProvider
                    bootstrap4
                    keyField="email"
                    data={ this.state.items }
                    columns={ columns }
                    search
                >
                    {
                    props => (
                        <Card className="withHeader">
                            {this.state.items !== null &&
                            <>
                            <Card.Header className="noBottom">
                                <SearchBar {...props.searchProps} srText="" />
                            </Card.Header>
                            <Card.Body>
                            <BootstrapTable
                                {...props.baseProps}
                                keyField="email"
                                bordered={false}
                                defaultSorted={defaultSorted}
                                pagination={paginationFactory(paginationOptions)}
                                rowClasses="detailedTableRow"
                                hover
                            />
                            </Card.Body>
                            </>
                            }
                        </Card>
                    )
                    }
                </ToolkitProvider>
                 : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                </div>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={UsersInfo.title}
                    content={UsersInfo.content}
                />

                <InfoModal
                    show={isPendingInvitationsOpen}
                    onHide={this.closePendingInvitations}
                    title="Pending invitations"
                    content={<><PendingInvitations/></>}
                />

                <ConfirmationModal
                    show={isTerminationModalOpen}
                    title="Deactivate user"
                    body={"Are you sure you want to deactivate user " + terminationName + "?"}
                    confirmText="Deactivate"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />

                <ConfirmationModal
                    show={isActivationModalOpen}
                    title="Activate user"
                    body={"Are you sure you want to activate user " + activationName + " again?"}
                    confirmText="Activate"
                    onCancel={this.handleCloseActivationModal}
                    onHide={this.handleCloseActivationModal}
                    onConfirm={this.handleConfirmActivation}
                />

                <ConfirmationModal
                    show={isUpgradeModalOpen}
                    title="Upgrade user to administrator"
                    body={"Are you sure you want to upgrade user " + upgradeName + " to administrator?"}
                    confirmText="Upgrade"
                    onCancel={this.handleCloseUpgradeModal}
                    onHide={this.handleCloseUpgradeModal}
                    onConfirm={this.handleConfirmUpgrade}
                />

                <ConfirmationModal
                    show={isDowngradeModalOpen}
                    title="Downgrade administrator to user"
                    body={"Are you sure you want to downgrade administrator " + downgradeName + " to team member?"}
                    confirmText="Downgrade"
                    onCancel={this.handleCloseDowngradeModal}
                    onHide={this.handleCloseDowngradeModal}
                    onConfirm={this.handleConfirmDowngrade}
                />

                <InfoModal
                    show={isFormOpen}
                    onHide={this.closeForm}
                    title="Invite users"
                    content={
                        <>
                            <div className="mb-5">
                                <EmailChips
                                    changesHandler={this.inviteItemsChangesHandler}
                                />
                            </div>
                            {submissionMessage && (
                                <div className="form-group">
                                    <div className={submissionSuccessful ? "alert alert-success mt-3" : "alert alert-danger mt-3" } role="alert">
                                        {submissionMessage}
                                    </div>
                                </div>
                            )}
                        </>
                    }
                    confirmText="Invite"
                    confirmButton={true}
                    onConfirm={this.submitForm}
                    submitButtonLoading={submitButtonLoading}
                />
                </div>
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

export default connect(mapStateToProps)(Users);
