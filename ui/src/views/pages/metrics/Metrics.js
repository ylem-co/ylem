import React from 'react';
import { Navigate, useParams } from 'react-router-dom';
import { connect } from "react-redux";
import { Fade } from "react-awesome-reveal";
import { DateTime } from "luxon";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";

import Tooltip from '@mui/material/Tooltip';
import Circle from '@mui/icons-material/Circle';

import Edit from '@mui/icons-material/Edit';
import DeleteOutlined from '@mui/icons-material/DeleteOutlined';
import BarChartRounded from '@mui/icons-material/BarChartRounded';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import { MetricsInfo } from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import FullScreenModal from "../../../components/modals/fullScreenModal.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component";
import PipelineStatistic from "../../../components/pipelines/statistic/pipelineStatistic.component";

import MetricForm from "../../../components/forms/metricForm.component";
import {TimeAgo} from "../../../components/timeAgo.component";

import PipelineService, { PIPELINE_TYPE_METRIC } from "../../../services/pipeline.service";

import { PERMISSION_LOGGED_IN, validatePermissions } from "../../../actions/pipeline";

const { SearchBar } = Search;

function withParams(Component) {
  return props => <Component {...props} params={useParams()} />;
}

class Metrics extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            isInfoOpen: false,
            isFormOpen: false,
            isTerminationModalOpen: false,
            isStatisticOpen: false,
            statItem: null,
            items: null,
            terminationName: null,
            terminationUuid: null,
            activeItem: null,
            dateFrom: this.props.params.dateFrom || DateTime.now().plus({ days: -7 }).toSQLDate() + " 00:00:00",
            dateTo: this.props.params.dateTo || DateTime.now().toSQLDate() + " 23:59:59",
        };
    }

    componentDidMount = async() => {
        document.title = 'Metrics';

        await this.handleGetMetrics(
            this.state.organization.uuid
        );
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    handleGetMetrics = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = PipelineService.getPipelinesByOrganization(uuid);

            await Promise.resolve(items)
                .then(async(items) => {
                    if (items.data && items.data.items  && items.data.length !== 0) {
                        var data = items.data.items.filter(k => k.type === PIPELINE_TYPE_METRIC);
                        
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

    rowEvents = () => {
        return {
            onClick: (e, row, rowIndex) => {
                this.toogleForm(row);
            },
        };
    }

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

        await this.handleGetMetrics(
            this.state.organization.uuid
        );
    };

    openStatistic = (item) => {
        this.toogleStatistic(item);
    };

    toogleStatistic = async(item = null) => {
        await this.promisedSetState({
            isStatisticOpen: !this.state.isStatisticOpen,
            statItem: item,
        });
    };

    closeStatistic = () => {
        this.setState({
            isStatisticOpen: false,
            statItem: null,
        });
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

        await PipelineService.deletePipeline(uuid);
        await this.promisedSetState({items: null});
        await this.handleGetMetrics(
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

    changeDatesHandler = async(dateFrom, dateTo) => {
        await this.promisedSetState({dateTo, dateFrom});
    };

    render() {
        const { 
            isInfoOpen,
            isFormOpen,
            activeItem,
            isTerminationModalOpen,
            terminationName,
            statItem,
            isStatisticOpen,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (
            !validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)
        ) {
            return <Navigate to="/login" />;
        }

        const columns = [
        {
            dataField: 'schedule',
            text: '',
            sort: false,
            searchable: false,
            headerStyle: { width: '50px' },
            formatter: (cellContent, row) => (
                <>
                    <Tooltip title={row.schedule !== ""  ? "Metric is scheduled" : "Metric is not scheduled"} placement="right">
                        <Circle 
                            className={row.schedule !== "" ? "icon_online" : "icon_new"}
                            alt={row.schedule !== "" ? "Metric is scheduled" : "Metric is not scheduled"}
                        />
                    </Tooltip>
                </>
            )
        },
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
            dataField: 'updated_at',
            text: 'Last time updated at',
            sort: true,
            sortFunc: (a, b, order, dataField, rowA, rowB) => {
                let a1 = rowA.updated_at !== "" 
                    ? DateTime.fromSQL(rowA.updated_at, { zone: 'utc'})
                    : DateTime.fromISO(rowA.created_at, { zone: 'utc'});
                
                let b1 = rowB.updated_at !== "" 
                    ? DateTime.fromSQL(rowB.updated_at, { zone: 'utc'})
                    : DateTime.fromISO(rowB.created_at, { zone: 'utc'});

                if (order === 'asc') {
                  return a1.ts - b1.ts;
                }
                return b1.ts - a1.ts;
              },
            headerStyle: { width: '300px' },
            searchable: true,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {row.updated_at !== "" 
                            ? TimeAgo(DateTime.fromSQL(row.updated_at, { zone: 'utc'})) 
                            : row.created_at !== "" 
                                ? TimeAgo(DateTime.fromISO(row.created_at, { zone: 'utc'}))
                                : "New metric"
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
                    <Tooltip title="Edit" placement="left">
                        <Edit 
                            className="iconDelete"
                            onClick={() => this.toogleForm(row)}
                        />
                    </Tooltip>
                    <Tooltip title="Statistics" placement="top">
                        <BarChartRounded 
                            className="iconDelete"
                            onClick={() => this.openStatistic(row)}
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
          dataField: 'updated_at',
          order: 'desc'
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
                                <h1>Metrics</h1>
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
                                        Add metric
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
                                                                rowEvents={ this.rowEvents() }
                                                            />
                                                        </Card.Body>
                                                    </>
                                                }
                                            </Card>
                                        )
                                        }
                                    </ToolkitProvider>
                                : <div className="text-center note mb-3">
                                    You didn't create any metric yet
                                  </div>
                             : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                            }
                            </div>
                        </Row>
                    </div>

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={MetricsInfo.title}
                    content={MetricsInfo.content}
                />

                <FullScreenModal
                    show={isFormOpen}
                    content={
                      <MetricForm 
                        item={activeItem}
                        successHandler={this.closeForm}
                      />
                    }
                    item={activeItem}
                    title={activeItem === null ? "Create metric" : "Edit metric"}
                    onHide={this.closeForm}
                  />

                  <FullScreenModal
                    show={isStatisticOpen}
                    onHide={this.closeStatistic}
                    title={statItem !== null && statItem.name ? statItem.name + ". Statistic" : "Metric statistic"}
                    content={
                        <>
                            <PipelineStatistic 
                                item={statItem}
                                type={PIPELINE_TYPE_METRIC}
                                changeDatesHandler={this.changeDatesHandler}
                            />
                        </>
                    }
                    altButton={false}
                    item={statItem}
                    confirmButton={false}
                  />

                  <ConfirmationModal
                    show={isTerminationModalOpen}
                    title="Delete metric"
                    body={'Are you sure you want to delete metric "' + terminationName + '"?'}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />
            </Fade>
        )
    }
}

function MapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    return {
        isLoggedIn,
        user,
    };
}

export default connect(MapStateToProps)(withParams(Metrics));
