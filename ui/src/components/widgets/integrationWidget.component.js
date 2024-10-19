import React from 'react';
import { useNavigate } from 'react-router-dom';

import Circle from '@mui/icons-material/Circle';
import Tooltip from '@mui/material/Tooltip';

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Spinner from "react-bootstrap/Spinner";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Button from "react-bootstrap/Button";

import IntegrationService, { INTEGRATION_TYPE_SQL } from "../../services/integration.service";

function withParams(Component) {
    return props => <Component {...props} history={useNavigate()} />;
}

class IntegrationWidget extends React.Component {
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
                this.props.history('/integrations/details/' + row.uuid);
            },
        }
    }

    handleGetItems = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = IntegrationService.getIntegrationsByOrganization(uuid);

            Promise.resolve(items)
                .then(items => {
                    if (items.data) {
                        this.setState({items: items.data.items});
                    } else {
                        this.setState({items: []});
                    }
                });
        }
    };

    render() {
        const columns = [
        {
            dataField: 'name',
            text: '',
            sort: false,
            formatter: (cellContent, row) => (
                <>
                    <div className={"dropdownItemWithBg dropdownItemWithBg-" + (row.type === INTEGRATION_TYPE_SQL ? row.value : row.type) + " bgInTheList float-left"}></div>
                    {row.name}
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
                        <Tooltip title={row.status} placement="right">
                            <Circle 
                                className={"icon_" + row.status}
                                alt={row.status === "offline" ? "Connection doesn't work" : "Connection works"}
                            />
                        </Tooltip>
                    </div>
                </>
            )
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
                    : <div className="text-center note mb-3">You didn't create any integration yet</div>
                 : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                    <div className="text-right">
                        <ButtonGroup className="mr-4">
                            <Button
                                color="primary"
                                className="mx-0"
                                href="/integrations"
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

export default withParams(IntegrationWidget);
