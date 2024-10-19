import React from 'react';
import { useNavigate } from 'react-router-dom';

import Circle from '@mui/icons-material/Circle';
import Tooltip from '@mui/material/Tooltip';

import { DateTime } from "luxon";

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Spinner from "react-bootstrap/Spinner";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Button from "react-bootstrap/Button";

import PipelineService, {PIPELINE_TYPE_GENERIC} from "../../services/pipeline.service";

import {TimeAgo} from "../../components/timeAgo.component";

function withParams(Component) {
  return props => <Component {...props} history={useNavigate()} />;
}

class PipelineWidget extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            items: null,
        };
    }

    componentDidMount() {
        this.handleGetItems(this.state.organization.uuid);
    };

    getRowEvents() {
        return {
            onClick: (e, row, rowIndex) => {
                this.props.history(
                    row.type === PIPELINE_TYPE_GENERIC
                        ? (row.folder_uuid === ""
                            ? '/pipelines/folder/root/details/' + row.uuid
                            : '/pipelines/folder/' + row.folder_uuid + '/details/' + row.uuid)
                        : (row.folder_uuid === ""
                            ? '/metrics/folder/root/details/' + row.uuid
                            : '/metrics/folder/' + row.folder_uuid + '/details/' + row.uuid)
                );
            },
        }
    }

    handleGetItems = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = PipelineService.getPipelinesByOrganization(uuid);

            Promise.resolve(items)
                .then(items => {
                    if (items.data && items.data.items && items.data.items !== null) {
                        items = items.data.items.filter(k => k.type === this.props.type);

                        this.setState({items});
                    } else {
                        this.setState({items: []});
                    }
                });
        }
    };

    render() {
        const { type } = this.props;

        const columns = [
        {
            dataField: 'name',
            text: '',
            sort: false,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {row.name}
                    </div>
                </>
            )
        },
        {
            dataField: '',
            text: '',
            sort: false,
            style: { width: '20%' },
            formatter: (cellContent, row) => (
                <>
                    <div>
                        <span className="note">{row.updated_at !== "" ? TimeAgo(DateTime.fromSQL(row.updated_at, { zone: 'utc'})) : "New pipeline"}</span>
                    </div>
                </>
            )
        },
        {
            dataField: '',
            text: '',
            sort: false,
            style: { width: '10%' },
            formatter: (cellContent, row) => (
                <>
                    <div>
                        <Tooltip title={row.schedule !== ""  ? "Is scheduled" : "Is not scheduled"} placement="right">
                            <Circle 
                                className={row.schedule !== "" ? "icon_online" : "icon_new"}
                                alt={row.schedule !== "" ? "Is scheduled" : "Is not scheduled"}
                            />
                        </Tooltip>
                    </div>
                </>
            ),
        },
    ];

    const paginationOptions = {
        sizePerPage: 5,
        totalSize: 5,
        hideSizePerPage: true,
        hidePageListOnlyOnePage: true
    };

        return (
            <>
                <div>
                {this.state.items !== null ?
                    this.state.items.length > 0 ?
                        <ToolkitProvider
                            bootstrap4
                            keyField="name"
                            data={ this.state.items }
                            columns={ columns }
                        >
                            {
                            props => (
                                <div>
                                    {this.state.items !== null &&
                                        <div>
                                            <BootstrapTable
                                                {...props.baseProps}
                                                keyField="name"
                                                headerClasses="hiddenHeader"
                                                bordered={false}
                                                hover
                                                pagination={paginationFactory(paginationOptions)}
                                                rowClasses="detailedTableRow"
                                                rowEvents={ this.getRowEvents() }
                                            />
                                        </div>
                                    }
                                </div>
                            )
                            }
                        </ToolkitProvider>
                    : <div className="text-center note mb-3">
                        {type === PIPELINE_TYPE_GENERIC 
                            ? "You didn't create any pipeline yet" 
                            : "You didn't create any metric yet"
                        }
                      </div>
                 : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                    <div className="text-right">
                        <ButtonGroup className="mr-4">
                            <Button
                                color="primary"
                                className="mx-0"
                                href={type === PIPELINE_TYPE_GENERIC ? "/pipelines" : "/metrics"}
                            >
                                Manage
                            </Button>
                        </ButtonGroup>
                    </div>
                </div>
            </>
        );
    }
}

export default withParams(PipelineWidget);
