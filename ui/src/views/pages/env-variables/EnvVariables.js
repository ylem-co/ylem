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

import Tooltip from '@mui/material/Tooltip';

import Edit from '@mui/icons-material/Edit';
import DeleteOutlined from '@mui/icons-material/DeleteOutlined';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import { EnvVariablesInfo } from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import RightModal from "../../../components/modals/rightModal.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component";

import EnvVariableForm from "../../../components/forms/envVariableForm.component";

import EnvVariableService from "../../../services/envVariable.service";

import { PERMISSION_LOGGED_IN, validatePermissions } from "../../../actions/pipeline";

const { SearchBar } = Search;

class EnvVariables extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            isFormOpen: false,
            isTerminationModalOpen: false,
            items: null,
            terminationName: null,
            terminationUuid: null,
            activeItem: null,
        };
    }

    componentDidMount = async() => {
        document.title = 'Environment variables';

        await this.handleGetVariables(
            this.state.organization.uuid
        );
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handleGetVariables = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = EnvVariableService.getVariablesByOrganization(uuid);

            await Promise.resolve(items)
                .then(async(items) => {
                    if (items.data && items.data.items  && items.data.length !== 0) {
                        var data = items.data.items;
                        
                        await this.promisedSetState({items: data});
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
            items: null,
        });

        await this.handleGetVariables(
            this.state.organization.uuid
        );
    };

    handleOpenTerminationModal = (name, uuid) => {
        this.setState({
            isTerminationModalOpen: true,
            terminationName: name,
            terminationUuid: uuid,
        });
    };

    handleConfirmTermination = async() => {
        var uuid = this.state.terminationUuid;

        await EnvVariableService.deleteVariable(uuid);
        await this.promisedSetState({items: null});
        await this.handleGetVariables(
            this.state.organization.uuid
        );
        this.handleCloseTerminationModal();
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
            terminationName: null,
            terminationUuid: null,
        });
    };

    render() {
        const { 
            isInfoOpen,
            isFormOpen,
            activeItem,
            isTerminationModalOpen,
            terminationName,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (
            !validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)
        ) {
            return <Navigate to="/login" />;
        }

        const columns = [
        {
            dataField: 'name',
            text: 'Name',
            sort: true,
            searchable: true,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {row.name}
                    </div>
                </>
            )
        },
        {
            dataField: 'value',
            text: 'Value',
            sort: true,
            searchable: true,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {row.value}
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
                    <Tooltip title="Edit" placement="left">
                        <Edit 
                            className="iconEdit"
                            onClick={() => this.toogleForm(row)}
                        />
                    </Tooltip>
                    <Tooltip title="Delete" placement="right">
                        <DeleteOutlined 
                            className="iconDelete"
                            onClick={() => this.handleOpenTerminationModal(row.name, row.uuid)}
                        />
                    </Tooltip>
                </>
            )
        },
        ];

        const defaultSorted = [{
          dataField: 'name',
          order: 'asc'
        }];

        const paginationOptions = {
            sizePerPage: 25,
            hidePageListOnlyOnePage: true
        };

        return (
            <Fade>
                    <div>
                        <Row className="mb-5">
                            <Col sm="9">
                                <h1>Environment Variables</h1>
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
                                        Add variable
                                    </Button>
                                </ButtonGroup>
                            </Col>
                        </Row>
                        <Row>
                        <div>
                            {this.state.items !== null ?
                                this.state.items.length > 0 ?
                                    <ToolkitProvider
                                        bootstrap4
                                        keyField="name"
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
                                                                keyField="name"
                                                                bordered={false}
                                                                hover
                                                                defaultSorted={defaultSorted}
                                                                pagination={paginationFactory(paginationOptions)}
                                                                rowClasses="detailedTableRow"
                                                            />
                                                        </Card.Body>
                                                    </>
                                                }
                                            </Card>
                                        )
                                        }
                                    </ToolkitProvider>
                                : <div className="text-center mb-3">
                                    You didn't create any environment variable yet
                                  </div>
                             : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                            }
                            </div>
                        </Row>
                    </div>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={EnvVariablesInfo.title}
                    content={EnvVariablesInfo.content}
                />

                <RightModal
                    show={isFormOpen}
                    content={
                      <EnvVariableForm 
                        item={activeItem}
                        successHandler={this.closeForm}
                      />
                    }
                    item={activeItem}
                    title={activeItem === null ? "Create environment variable" : "Edit environment variable"}
                    onHide={this.closeForm}
                  />

                  <ConfirmationModal
                    show={isTerminationModalOpen}
                    title="Delete environment variable"
                    body={'Are you sure you want to delete environment variable "' + terminationName + '"?'}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
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

export default connect(mapStateToProps)(EnvVariables);
