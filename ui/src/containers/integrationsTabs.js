import React from 'react';

import Nav from 'react-bootstrap/Nav';

class IntegrationsTabs extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            page: this.props.page,
        };
    }

    render() {
        return (
            <>
                <Nav variant="tabs" className="settingTabs mb-4">
                    <Nav.Item>
                        <Nav.Link
                            href="/slack-authorizations"
                            className="navLinkSlack"
                            active={this.state.page === "slack"}
                        >
                            Slack Authorizations
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            href="/jira-authorizations"
                            className="navLinkJira"
                            active={this.state.page === "jira"}
                        >
                            Jira Cloud Authorizations
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            href="/hubspot-authorizations"
                            className="navLinkHubspot"
                            active={this.state.page === "hubspot"}
                        >
                            Hubspot Authorizations
                        </Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link
                            href="/salesforce-authorizations"
                            className="navLinkSalesforce"
                            active={this.state.page === "salesforce"}
                        >
                            Salesforce Authorizations
                        </Nav.Link>
                    </Nav.Item>
                </Nav>
            </>
        );
    }
}

export default IntegrationsTabs;
