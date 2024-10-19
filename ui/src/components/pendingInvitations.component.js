import React from 'react';

import { DateTime } from "luxon";

import Tooltip from '@mui/material/Tooltip';

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Spinner from "react-bootstrap/Spinner";

import ContentCopy from '@mui/icons-material/ContentCopy';

import Avatar from "../components/avatar.component";
import {TimeAgo} from "../components/timeAgo.component";

import InvitationService from "../services/invitation.service";

class PendingInvitations extends React.Component {
    constructor(props) {
        super(props);

        this.copyLink = this.copyLink.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            items: null,
        };
    }

    componentDidMount() {
        this.handleGetInvitations(this.state.organization.uuid);
    };

    copyLink = async(copiedLink) => {
        navigator.clipboard.writeText(copiedLink);
    }

    handleGetInvitations = async(uuid) => {
        let items = this.state.items;

        if (
            items === null
            || items.length === 0
        ) {
            items = InvitationService.getPendingByOrganization(uuid);

            Promise.resolve(items)
                .then(items => {
                    if (items.data && items.data.items && items.data.items !== null) {
                        this.setState({items: items.data.items});
                    } else {
                        this.setState({items: []});
                    }
                })
                .catch(async() => {
                    this.setState({items: []});
                });
        }
    };

    render() {
        const columns = [
        {
            dataField: '',
            text: '',
            sort: false,
            headerStyle: { width: '100px' },
            formatter: (cellContent, row) => (
                <>
                    <div className={row.is_active === 0 && "fadedOut"}>
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
            dataField: 'email',
            text: 'Email',
            sort: true,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {row.email}
                    </div>
                </>
            )
        },
        {
            dataField: 'created_at',
            text: 'Sent at',
            sort: true,
            formatter: (cellContent, row) => (
                <>
                    <div>
                        {TimeAgo(DateTime.fromISO(row.created_at, { zone: 'utc'}))}
                    </div>
                </>
            )
        },
        {
            dataField: 'invitation_code',
            text: 'Copy invitation link',
            sort: true,
            headerStyle: { width: '250px' },
            formatter: (cellContent, row) => (
                <>
                    <div>
                        <Tooltip title="Click to copy link to clipboard" placement="right">
                            <ContentCopy
                                onClick={() => this.copyLink("https://app.ylem.co/invitation/" + row.invitation_code, row.uuid)}
                            />
                        </Tooltip>
                    </div>
                </>
            )
        },
    ];

    const defaultSorted = [{
      dataField: 'email',
      order: 'desc'
    }];

    const paginationOptions = {
        sizePerPage: 5,
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
                    keyField="email"
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
                                keyField="email"
                                bordered={false}
                                defaultSorted={defaultSorted}
                                pagination={paginationFactory(paginationOptions)}
                                rowClasses="detailedTableRow"
                                hover
                            />
                            </div>
                            }
                        </div>
                    )
                    }
                </ToolkitProvider>
                    : <div className="text-center note">Currently you don't have pending invitations</div>
                 : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                </div>
            </>
        );
    }
}

export default PendingInvitations;
